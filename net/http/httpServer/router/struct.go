package router

import (
	"embed"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/helays/utils/v2/net/http/route"
	"github.com/helays/utils/v2/net/http/session"
	"gorm.io/gorm"
)

type LoginInfo struct {
	LoginTime     time.Time // 登录时间
	IsLogin       bool      // 是否登录
	UserId        int       // 用户ID
	User          string    // 用户名
	IsManage      bool      // 是否管理员
	DemoUser      bool      // 是否演示用户
	RsaPrivateKey []byte    //ras 私钥
	RsaPublickKey []byte    // rsa 公钥
}

func (i LoginInfo) QueryIsManage() func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if !i.IsManage {
			db.Where("user_id=?", i.UserId)
		}
		return db
	}
}

type ErrPageFunc func(w http.ResponseWriter, resp *route.ErrorResp)
type Router struct {
	Default                string `ini:"default" json:"default" yaml:"default"`
	Root                   string `ini:"root" json:"root" yaml:"root"`
	HttpCache              bool   `ini:"http_cache" json:"http_cache" yaml:"http_cache"`
	HttpCacheMaxAge        string `ini:"http_cache_max_age" json:"http_cache_max_age" yaml:"http_cache_max_age"`
	UnauthorizedRespMethod int    `ini:"unauthorized_resp_method" json:"unauthorized_resp_method" yaml:"unauthorized_resp_method"` // 未登录响应方法 默认为 401，302 表示自动重定向到登录页面
	SessionLoginName       string `ini:"session_login_name" json:"session_login_name" yaml:"session_login_name"`                   // 系统中，用于在session中存储登录信息的key

	// cookie相关配置
	CookiePath     string `ini:"cookie_path" json:"cookie_path" yaml:"cookie_path"`
	CookieDomain   string `ini:"cookie_domain" json:"cookie_domain" yaml:"cookie_domain"`
	CookieSecure   bool   `ini:"cookie_secure" json:"cookie_secure" yaml:"cookie_secure"`
	CookieHttpOnly bool   `ini:"cookie_http_only" json:"cookie_http_only" yaml:"cookie_http_only"`

	Error        ErrPageFunc // 错误页面处理函数
	ErrorWithLog ErrPageFunc // 错误页面处理函数

	dev           bool                  // 开发模式
	staticEmbedFS map[string]*embedInfo // 静态文件

	IsLogin                bool             // 是否登录
	LoginPath              string           // 登录页面
	HomePage               string           //首页
	UnLoginPath            map[string]bool  // 免授权页面
	UnLoginPathRegexp      []*regexp.Regexp // 免授权页面正则
	MustLoginPath          map[string]bool  //必须登录才能访问的页面
	MustLoginPathRegexp    []*regexp.Regexp // 必须登录才能访问的页面正则
	DisableLoginPath       map[string]bool  // 登录状态下不能访问的页面
	DisableLoginPathRegexp []*regexp.Regexp // 登录状态下不能访问的页面正则
	ManagePage             map[string]bool  // 管理员访问
	ManagePageRegexp       []*regexp.Regexp

	session *session.Manager
}

type embedInfo struct {
	embedFS *embed.FS
	prefix  string
}

func (ro *Router) SetDev(dev bool) {
	ro.dev = dev
}

// SetStaticEmbedFs 设置内置 embedFS
func (ro *Router) SetStaticEmbedFs(p string, embedFS *embed.FS, prefix ...string) {
	if ro.staticEmbedFS == nil {
		ro.staticEmbedFS = make(map[string]*embedInfo)
	}
	ro.staticEmbedFS[p] = &embedInfo{embedFS: embedFS}
	if len(prefix) > 0 {
		if !strings.HasPrefix(prefix[0], "/") {
			prefix[0] = "/" + prefix[0]
		}
		ro.staticEmbedFS[p].prefix = prefix[0]
	}
}

func (ro *Router) SetSession(sessionManager *session.Manager) {
	ro.session = sessionManager
}
func (ro *Router) GetSession() *session.Manager {
	return ro.session
}
