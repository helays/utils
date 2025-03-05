package httpClient

import (
	"bufio"
	"compress/gzip"
	"context"
	"crypto/tls"
	"errors"
	"github.com/helays/utils/close/gzipClose"
	"github.com/helays/utils/close/httpClose"
	"github.com/helays/utils/http/httpTools"
	"gopkg.in/iconv.v1"
	"io"
	"net"
	"net/http"
	"net/http/cookiejar"
	url2 "net/url"
	"regexp"
	"strings"
	"time"
	"unicode/utf8"
)

// Init 初始化
func (c *Curl) Init() {
	if c.Sleep < 1 {
		c.Sleep = 2
	}
	// 异常重新尝试次数
	if c.Retry < 1 {
		c.Retry = 10
	}
	// 请求
	if c.Connection == "" {
		c.Connection = "close"
	}
	if c.Maxbody < 1 {
		c.Maxbody = 3 * 1024 * 1024 // 最大body允许3M
	} else {
		c.Maxbody *= 1024 * 1024
	}
	if c.AcceptEncoding == "" {
		c.AcceptEncoding = "gzip"
	}
	c.acceptList = strings.Split(c.Accept, ",")
	c.initClient()
}

// 在启动一个CURL的时候，就实例化一个client
func (c *Curl) initClient() {
	c.Client = &http.Client{
		Transport: &http.Transport{
			// 表示连接池对所有host的最大链接数量，host也即dest-ip，默认为无穷大（0），
			// 但是通常情况下为了性能考虑都要严格限制该数目（实际使用中通常利用压测二分得到该参数的最佳近似值）。
			// 太大容易导致客户端和服务端的socket数量剧增，导致内存吃满，文件描述符不足等问题；
			// 太小则限制了连接池的socket数量，资源利用率较低。
			MaxIdleConns: c.Curlcommon.MaxIdleConns,
			// 表示连接池对每个host的最大链接数量，从字面意思也可以看出：
			MaxIdleConnsPerHost: c.Maxidleconnsperhost,
			// 空闲timeout设置，也即socket在该时间内没有交互则自动关闭连接,该参数通常设置为分钟级别，例如：90秒。
			IdleConnTimeout: c.Idleconntimeout * time.Second,
			// DisableCompression: true,
			// 使用短链接
			// DisableKeepAlives: true,
			// request header 超时
			ResponseHeaderTimeout: c.ResponseHeaderTimeout * time.Second,
			ExpectContinueTimeout: c.ExpectContinueTimeout * time.Second,

			DialContext: func(ctx context.Context, network, addr string) (conn net.Conn, e error) {
				// 确定，这是建立连接的时间
				// 这里主要是控制解析域名时间,TCP 建立连接超时时长设置
				conn, err := net.DialTimeout(network, addr, c.Tcpconnecttimeout*time.Second)
				if err != nil {
					return nil, err
				}
				return conn, nil
			},
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true, // client 不对https 证书进行校验
			},
		},
		Jar: c.jar,
		// time.Duration不是一个函数，只是将数据显示转为 time.Duration这个类型
		// 这个超时是总的超时时间，如果不清楚Transport里面的设置可以设置这个超时
		// 由这个超时来设置总的超时时长
		Timeout: c.Timeouttotal * time.Second,
		// 检查首次状态值
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			if c.Allowredict {
				return nil
			}
			// 这里对于抓取非严格模式，允许 www.ba.com/a 到www.ba.com/a/ 这种类型的301 或者302
			if (req.Response.StatusCode == 301 || req.Response.StatusCode == 302) && c.Strict == false {
				rl, err := req.Response.Location()
				if err != nil {
					return http.ErrUseLastResponse
				}
				ori := req.Response.Request.URL.String()
				ori = ori + "/"
				if ori == rl.String() {
					// 如果跳转后的链接 就仅仅是最后面添加了一个 / 就进行保留。
					return nil
				}
			}
			// 直接返回最后一次的状态，30x也不会进行跳转
			return http.ErrUseLastResponse
		},
	}
}

func (c *Curl) SetCookie() {
	c.jar, _ = cookiejar.New(nil)
}

// SetStrict
// 设置严格模式或者非严格模式
func (c *Curl) SetStrict(s bool) {
	c.Strict = s
}

func (c *Curl) Setrefer(refer string) {
	c.Refer = refer
}

// Run
// 运行
func (c *Curl) Run(url string) (*Curlresult, error) {
	_, err := url2.Parse(url) // 判断 url 是否正常
	if err != nil {           // new request 只有这里会出错
		return nil, err
	}
	var (
		re   = new(Curlresult)
		resp = new(http.Response)
	)
	resp, err = c.startRequest(url, resp, err)
	defer httpClose.CloseResp(resp)
	if err != nil {
		return nil, err
	}
	re.Url = resp.Request.URL
	// 首先过滤资源类型的页面
	ctype := resp.Header.Get("content-type")
	if !c.filterContentType(ctype) {
		httpClose.CloseResp(resp)
		re.HttpStatus = 1
		// 如果不是允许的 头信息，返回为资源类型的页面。
		return re, errors.New("资源类型不抓取：" + ctype + " " + url)
	}
	re.HttpStatus = resp.StatusCode
	re.Header = resp.Header.Clone()
	// 然后读取 content，读取失败的就直接返回不处理
	encoding := resp.Header.Get("Content-Encoding")
	if strings.ToUpper(encoding) != "GZIP" {
		read := bufio.NewReader(resp.Body)
		body, err := c.readBody(read)
		if err != nil {
			return re, err
		}
		re.Body = body
		return re, nil
	}
	read, err := gzip.NewReader(resp.Body)
	if err != nil {
		return re, err
	}
	defer gzipClose.Close(read)
	body, err := c.readBody(read)
	if err != nil {
		return re, err
	}
	re.Body = body
	return re, nil
}

func (this *Curl) readBody(read io.Reader) ([]byte, error) {
	var body []byte
	buff := make([]byte, 1024)
	for {
		n, err := read.Read(buff)
		body = append(body, buff[:n]...)
		if err != nil {
			if err == io.EOF {
				break
			}
			return body, err
		}
		if len(body) > this.Maxbody {
			return body, errors.New("非正常HTML页面")
		}
	}
	if len(body) < 1 {
		return body, errors.New("读取到页面为空")
	}
	// 页面字符转码
	if err := this.character(&body); err != nil {
		return nil, err
	}
	return body, nil
}

// 字符串转码
func (c *Curl) character(body *[]byte) error {
	// 识别页面字符集,每次页面打开，都需要先识别一次页面字符集
	if utf8.Valid(*body) {
		return nil
	}
	searchre := httpTools.PageCharacterSetPreg.FindStringSubmatch(string(*body))
	if len(searchre) < 2 {
		return errors.New("页面字符编码识别失败 ")
	}
	charset := strings.ToLower(searchre[1])
	if charset != "utf8" && charset != "utf-8" && !utf8.Valid(*body) {
		// 转码处理
		if charset == "gbk" || charset == "gb2312" {
			charset = "gb18030"
		}
		if err := c.toUtf8(body, charset); err != nil {
			return err
		}
	}
	return nil
}

func (c *Curl) toUtf8(src *[]byte, srccode string) error {
	// 对页面数据进行转码
	icov, err := iconv.Open("UTF-8", strings.ToUpper(srccode))
	if err != nil {
		return err
	}
	defer func() {
		_ = icov.Close()
	}()
	var outbuf = make([]byte, len(*src))
	s1, _, err := icov.Conv(*src, outbuf)
	if err != nil {
		return err
	}
	*src = s1
	return nil
}

// 过滤content type 类型
func (c *Curl) filterContentType(t string) bool {
	t = strings.ToLower(t)
	for _, i := range c.acceptList {
		preg := regexp.MustCompile(i)
		if preg.MatchString(t) {
			return true
		}
	}
	return false
}

// 开始请求
func (c *Curl) startRequest(url string, resp *http.Response, err error) (*http.Response, error) {
	for i := 0; i < c.Retry; i++ {
		req := c.newRequest(url)
		// 开始请求
		resp, err = c.Client.Do(req)
		if err == nil { // 请求正常
			return resp, nil
		}
		// 失败的 也确认下，关闭系统
		if resp != nil && resp.Body != nil {
			_ = resp.Body.Close()
		}
		resp = nil
		// 等待一秒后重新开始
		time.Sleep(time.Duration(c.Sleep) * time.Second)
	}
	return nil, err
}

// 创建请求
func (c *Curl) newRequest(url string) *http.Request {
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("User-Agent", c.UserAgent)
	req.Header.Set("accept", c.Accept)
	req.Header.Set("accept-language", c.AcceptLanguage)
	req.Header.Set("cache-control", c.CacheControl)
	req.Header.Set("Connection", c.Connection)
	req.Header.Set("Accept-Encoding", c.AcceptEncoding)
	if c.Refer != "" {
		req.Header.Set("Referer", c.Refer)
	}
	// req.Close = true

	return req
}
