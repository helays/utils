package cookiekit

import (
	"net/http"
	"time"

	"github.com/helays/utils/v2/tools"
)

type Config struct {
	Name    string    `json:"name" yaml:"name"`
	Value   string    `json:"value" yaml:"value"`
	Quoted  bool      `json:"quoted" yaml:"quoted"`   // Value 是否原本被引号包围（影响序列化时的格式）
	Path    string    `json:"path" yaml:"path"`       //  Cookie 的有效路径，指定在哪些路径下浏览器会发送该 Cookie
	Domain  string    `json:"domain" yaml:"domain"`   // Cookie 的有效域名，指定在哪些域名下浏览器会发送该 Cookie
	Expires time.Time `json:"expires" yaml:"expires"` // 过期时间
	// =0 不设置Max-Age 属性
	// <0 立即删除Cookie 等同于Max-Age:0
	// >0 设置具体的存活秒数
	MaxAge   int  `json:"max_age" yaml:"max_age"`
	Secure   bool `json:"secure" yaml:"secure"`       // 为 true 时表示 Cookie 只能通过 HTTPS 传输
	HttpOnly bool `json:"http_only" yaml:"http_only"` // 为 true 时表示 Cookie 不能被 JavaScript 访问
	// 控制跨站请求时是否发送 Cookie
	// SameSiteDefaultMode 浏览器默认行为
	// SameSiteLaxMode 宽松模式
	// SameSiteStrictMode：严格模式
	// SameSiteNoneMode：无限制
	SameSite    http.SameSite `json:"same_site" yaml:"same_site"`
	Partitioned bool          `json:"partitioned" yaml:"partitioned"` // 布尔值，表示是否为分区 Cookie（跨站上下文中的存储隔离）
}

func (c *Config) Clone() *Config {
	if c == nil {
		return nil
	}
	clone := *c // 浅拷贝基础字段
	return &clone
}

func SetCookie(w http.ResponseWriter, value *Config) {
	cookie := http.Cookie{
		Name:        value.Name,
		Value:       value.Value,
		Quoted:      value.Quoted,
		Path:        tools.Ternary(value.Path == "", "/", value.Path),
		Domain:      value.Domain,
		Expires:     value.Expires,
		MaxAge:      value.MaxAge,
		Secure:      value.Secure,
		HttpOnly:    value.HttpOnly,
		SameSite:    value.SameSite,
		Partitioned: value.Partitioned,
	}
	http.SetCookie(w, &cookie)
}

func DelCookie(w http.ResponseWriter, value *Config) {
	cookie := http.Cookie{
		Name:     value.Name,
		Value:    "",
		Path:     tools.Ternary(value.Path == "", "/", value.Path),
		Domain:   value.Domain,
		Expires:  time.Unix(0, 0),
		MaxAge:   -1,
		Secure:   value.Secure,
		HttpOnly: value.HttpOnly,
		SameSite: value.SameSite,
	}
	http.SetCookie(w, &cookie)
}
