package sessionmgr

import (
	"net/http"
	"time"
)

// Options session 配置
type Options struct {
	CookieName    string        `json:"cookie_name" yaml:"cookie_name" ini:"cookie_name"`          // 从cookie或者 header中读取 session的标识
	DisableGc     bool          `json:"disable_gc" yaml:"disable_gc" ini:"disable_gc"`             // 是否禁用session gc
	GcProbability float64       `json:"gc_probability" yaml:"gc_probability" ini:"gc_probability"` // session gc 概率
	CheckInterval time.Duration `json:"check_interval" yaml:"check_interval" ini:"check_interval"` // session 检测默认有效期
	Carrier       CookieCarrier `json:"carrier" yaml:"carrier" ini:"carrier"`                      // session 载体，默认cookie
	// cookie相关配置
	Path   string `json:"path" yaml:"path" ini:"path"`
	Domain string `json:"domain" yaml:"domain" ini:"domain"`
	// MaxAge=0 means no Max-Age attribute specified and the cookie will be
	// deleted after the browser session ends.
	// MaxAge<0 means delete cookie immediately.
	// MaxAge>0 means Max-Age attribute present and given in seconds.
	MaxAge   int           `json:"max_age" yaml:"max_age" ini:"max_age"`
	Secure   bool          `json:"secure" yaml:"secure" ini:"secure"`
	HttpOnly bool          `json:"http_only" yaml:"http_only" ini:"http_only"`
	SameSite http.SameSite `json:"same_site" yaml:"same_site" ini:"same_site"`
}
