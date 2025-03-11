package session

import (
	"errors"
	"github.com/helays/utils/tools"
	uuid "github.com/satori/go.uuid"
	"net/http"
	"regexp"
	"strings"
	"time"
)

//
// ━━━━━━神兽出没━━━━━━
// 　　 ┏┓     ┏┓
// 　　┏┛┻━━━━━┛┻┓
// 　　┃　　　　　 ┃
// 　　┃　　━　　　┃
// 　　┃　┳┛　┗┳  ┃
// 　　┃　　　　　 ┃
// 　　┃　　┻　　　┃
// 　　┃　　　　　 ┃
// 　　┗━┓　　　┏━┛　Code is far away from bug with the animal protecting
// 　　　 ┃　　　┃    神兽保佑,代码无bug
// 　　　　┃　　　┃
// 　　　　┃　　　┗━━━┓
// 　　　　┃　　　　　　┣┓
// 　　　　┃　　　　　　┏┛
// 　　　　┗┓┓┏━┳┓┏┛
// 　　　　 ┃┫┫ ┃┫┫
// 　　　　 ┗┻┛ ┗┻┛
//
// ━━━━━━感觉萌萌哒━━━━━━
//
//
// User helay
// Date: 2024/12/7 22:02
//

// session 支持内存存储、程序重启session失效
// 支持文件存储，可持久化存储
// 支持 redis 存储
// 支持数据库存储

const (
	Interval   = time.Hour // 默认检测频率
	CookieName = "vsclubId"
	ExpireTime = 1 * time.Hour // session默认24小时过期
)

var (
	ErrUnSupport = errors.New("不支持的session载体")
	ErrNotFound  = errors.New("session不存在")
)

type Store struct {
	*Options
	Dialector
}

// Init 初始化session
func Init(dialector Dialector, opt ...*Options) *Store {
	c := &Store{}
	if len(opt) > 0 {
		c.Options = opt[0]
		c.Options.CheckInterval = tools.AutoTimeDuration(c.Options.CheckInterval, time.Microsecond, Interval)
	} else {
		c.Options.CheckInterval = Interval
		c.Options.Carrier = "cookie"
	}
	c.Options.CookieName = tools.Ternary(c.Options.CookieName == "", CookieName, c.Options.CookieName)
	c.Options.GcProbability = tools.Ternary(c.Options.GcProbability <= 0, 0.9, c.Options.GcProbability) // 默认GC 90%
	c.Dialector = dialector
	c.Dialector.Apply(c.Options)
	return c
}

func Close(c *Store) {
	if c != nil {
		_ = c.Close()
	}
}

// 创建session ID
func newSessionId() string {
	return strings.ReplaceAll(uuid.NewV4().String(), "-", "")
}

// c0449773432e4a478d157a8a923199ac
// 用于校验session id 值是否合规
var sessionRegexp = regexp.MustCompile("^[0-9a-f]{32}$")

// GetSessionId 获取sessionId
// 如果未获取到 sessionId ，则生成一个，并写到响应信息中
func GetSessionId(w http.ResponseWriter, r *http.Request, options *Options) (string, error) {
	switch options.Carrier {
	case "", "cookie":
		cookie, err := r.Cookie(options.CookieName)
		if err != nil || !sessionRegexp.MatchString(cookie.Value) {
			sid := newSessionId()
			SetSessionId(w, sid, options)
			return sid, nil
		}
		return cookie.Value, nil
	case "header":
		sid := r.Header.Get(options.CookieName)
		if sid == "" || !sessionRegexp.MatchString(sid) {
			sid = newSessionId()
			SetSessionId(w, sid, options)
		}
		return sid, nil
	}
	return "", ErrUnSupport

}

// SetSessionId 设置sessionId
func SetSessionId(w http.ResponseWriter, sid string, options *Options) {
	expire := time.Time{}
	if options.MaxAge > 0 {
		expire = time.Now().Add(time.Duration(options.MaxAge) * time.Second)
	}
	switch options.Carrier {
	case "cookie", "":
		http.SetCookie(w, &http.Cookie{
			Name:     options.CookieName,
			Value:    sid,
			Path:     options.Path,
			Domain:   options.Domain,
			Expires:  expire,           // 可选字段，指定了一个具体的过期时间。一旦到达这个时间点，浏览器应该删除这个cookie。如果未指定，则cookie是会话cookie，当浏览器关闭时它将被删除
			MaxAge:   options.MaxAge,   // 表示cookie的最大年龄（以秒为单位）。MaxAge<0意味着立即删除cookie；MaxAge=0意味着没有指定'Max-Age'属性；MaxAge>0则表示存在'Max-Age'属性，并给出以秒计的时间
			Secure:   options.Secure,   // 如果为true，那么cookie只能通过HTTPS安全连接传输给服务器。
			HttpOnly: options.HttpOnly, // 如果为true，则JavaScript无法访问这个cookie，这有助于防止跨站脚本攻击（XSS）
			SameSite: options.SameSite, // 控制浏览器是否应该在跨站请求中包含这个cookie。它可以有以下三个值：SameSite.Lax, SameSite.Strict, 或 SameSite.None。这个属性有助于减少跨站请求伪造（CSRF）攻击的风险
		})
	case "header":
		w.Header().Set("vsclub_"+options.CookieName, sid)
	}
}

func DeleteSessionId(w http.ResponseWriter, options *Options) {
	switch options.Carrier {
	case "cookie", "":
		http.SetCookie(w, &http.Cookie{
			Name:       options.CookieName,
			Value:      "",
			Path:       options.Path,
			Domain:     options.Domain,
			Expires:    time.Unix(0, 0),
			RawExpires: "",
			MaxAge:     -1,
			Secure:     options.Secure,
		})
	case "header":
		w.Header().Del("vsclub_" + options.CookieName)
	}
}

// GetSessionName 生成session key
// 这个key 是存储在 存储系统中的
func GetSessionName(sessionId, name string) string {
	return sessionId + "_" + name
}
