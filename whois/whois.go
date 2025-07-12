package whois

import (
	"github.com/helays/utils/v2/close/netClose"
	"github.com/helays/utils/v2/logger/ulogs"
	"github.com/helays/utils/v2/worker"
	"io/ioutil"
	"net"
	"net/url"
	"regexp"
	"strings"
	"sync"
	"time"
)

var (
	Root             = "whois.iana.org" // 查询域名后缀 whois server的服务器
	DialTimeout      = 20 * time.Second // 连接超时时间
	ReadTimeout      = 10 * time.Second // 读取数据超时
	WriteTimeout     = 10 * time.Second // 写发送数据超时
	work             = new(worker.StartWorker)
	whoisServer      sync.Map
	whoisServerRegex = regexp.MustCompile(`(?i)whois:\s+(\S+)`)
)

type Whois struct {
	retry int // 重试3次
}

type findParams struct {
	server string      // 查询服务器地址
	url    string      // 被查询的数据
	res    chan string // 查询结果信号
}

// New 初始化
func New(thred, retry int) *Whois {
	work.WorkerPool = make(chan chan worker.Job)
	work.MaxSize = thred
	work.Init()
	return &Whois{retry}
}

// Search 查询函数
// return status whoisinfo ip whoiserver
func (this *Whois) Search(u string) (int, string, string, string) {
	var ip, _server string
	host, suffix := this.parse(u)
	if host == "" || suffix == "" {
		return 0, "", ip, _server
	}
	_server = this.getWhoisServer(suffix)
	if _server == "" {
		return 0, "", ip, _server
	}
	var ch = make(chan string)
	work.Run(&worker.Job{
		Func: this.handle,
		Params: findParams{
			url:    host,
			server: _server,
			res:    ch,
		},
	})
	res := <-ch
	if res == "" {
		return 0, "", ip, _server
	}
	if addr, err := net.ResolveIPAddr("ip", host); err == nil {
		ip = addr.String()
	}
	return 1, res, ip, _server
}

// 解析url
func (this *Whois) parse(u string) (string, string) {
	u = strings.TrimSpace(u)
	preg := regexp.MustCompile(`^https?://`)
	if preg.MatchString(u) == false {
		u = "http://" + u
	}
	parse, err := url.Parse(u)
	if err != nil {
		ulogs.Error(err, "parseUrl")
		return "", ""
	}
	ur := parse.Host
	l := len(ur)
	lastPointIndex := strings.LastIndexAny(ur, ".")
	if lastPointIndex < 0 {
		return "", ""
	}
	return ur, ur[lastPointIndex:l]
}

// 根据后缀，查询对应后缀的whois server信息
func (this *Whois) getWhoisServer(suffix string) string {
	if _w, ok := whoisServer.Load(suffix); ok {
		return _w.(string)
	}
	// 现场获取
	var ch = make(chan string)
	work.Run(&worker.Job{
		Func: this.handle,
		Params: findParams{
			url:    suffix,
			server: Root,
			res:    ch,
		},
	})
	res := <-ch
	tmp := whoisServerRegex.FindAllStringSubmatch(res, -1)
	if len(tmp) == 1 && len(tmp[0]) == 2 {
		_server := strings.TrimSpace(tmp[0][1])
		whoisServer.Store(suffix, _server)
		return _server
	}
	return ""
}

// 根据whois 协议，去查询相关数据
func (this *Whois) handle(i interface{}) {
	var (
		ii     = i.(findParams)
		con    net.Conn
		err    error
		buffer []byte
	)
	defer netClose.CloseConn(con)
	for j := 0; j < this.retry; j++ {
		buffer = nil
		con, err = net.DialTimeout("tcp", net.JoinHostPort(ii.server, "43"), DialTimeout)
		if err != nil {
			netClose.CloseConn(con)
			continue
		}

		if err := con.SetWriteDeadline(time.Now().Add(WriteTimeout)); err != nil {
			netClose.CloseConn(con)
			continue
		}

		if err := con.SetReadDeadline(time.Now().Add(ReadTimeout)); err != nil {
			netClose.CloseConn(con)
			continue
		}

		_, err = con.Write([]byte(ii.url + "\r\n"))
		if err != nil {
			netClose.CloseConn(con)
			continue
		}
		buffer, err = ioutil.ReadAll(con)
		if err != nil {
			netClose.CloseConn(con)
			continue
		}
		break
	}
	ulogs.Checkerr(err, "whois", ii.url)
	ii.res <- string(buffer)
}
