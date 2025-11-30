package httpServer

import (
	"net/http"
	"time"

	"github.com/helays/utils/v2/logger/zaploger"
	"github.com/helays/utils/v2/net/http/route/middleware"
	"github.com/helays/utils/v2/net/ipmatch"
	"github.com/helays/utils/v2/security/cors"
	"github.com/helays/utils/v2/security/cors/cors_std"
	"golang.org/x/net/websocket"
)

type HttpServer struct {
	// 配置属性
	ListenAddr    string          `ini:"listen_addr" json:"listen_addr" yaml:"listen_addr"`
	ServerName    []string        `ini:"server_name,omitempty" json:"server_name" yaml:"server_name"` // 绑定域名
	Ssl           bool            `ini:"ssl" json:"ssl" yaml:"ssl"`
	Ca            string          `ini:"ca" json:"ca" yaml:"ca"`
	Crt           string          `ini:"crt" json:"crt" yaml:"crt"`
	Key           string          `ini:"key" json:"key" yaml:"key"`
	SocketTimeout time.Duration   `ini:"socket_timeout" json:"socket_timeout" yaml:"socket_timeout"` // socket 心跳超时时间
	Hotupdate     bool            `ini:"hotupdate" json:"hotupdate" yaml:"hotupdate"`                // 是否启动热加载
	EnableGzip    bool            `ini:"enable_gzip" json:"enable_gzip" yaml:"enable_gzip"`          // 是否开启gzip
	Security      SecurityConfig  `ini:"security" json:"security" yaml:"security"`                   // 安全配置
	Logger        zaploger.Config `json:"logger" yaml:"logger" ini:"logger" gorm:"comment:日志配置"`
	Allowip       []string        `ini:"allowip,omitempty" json:"allowip" yaml:"allowip"` // Deprecated: 请使用 addRoute
	Denyip        []string        `ini:"denyip,omitempty" json:"denyip" yaml:"denyip"`    // Deprecated: 请使用 addRoute

	// 可访问属性

	Route          map[string]http.HandlerFunc                       `yaml:"-" json:"-"` // Deprecated: 请使用 addRoute
	RouteHandle    map[string]http.Handler                           `yaml:"-" json:"-"` // Deprecated: 请使用 addRoute
	RouteSocket    map[string]func(ws *websocket.Conn)               `yaml:"-" json:"-"` // Deprecated: 请使用 addRoute
	CommonCallback func(w http.ResponseWriter, r *http.Request) bool `yaml:"-" json:"-"` // Deprecated: 请使用 addRoute

	route         map[string]*routerRule // 路由
	serverNameMap map[string]byte        // 绑定的域名
	logger        *middleware.ResponseProcessor

	allowIPMatch *ipmatch.IPMatcher
	denyIPMatch  *ipmatch.IPMatcher
	debugIPMatch *ipmatch.IPMatcher

	corsManager *cors_std.StdCORS

	mux    *http.ServeMux
	server *http.Server
}

type SecurityConfig struct {
	DefaultValidLast bool           `ini:"default_valid_last" json:"default_valid_last" yaml:"default_valid_last"` // 默认验证器是否放最后
	CORS             *cors.Config   `ini:"cors" json:"cors" yaml:"cors"`                                           // 跨域配置
	IPAccess         IPAccessConfig `ini:"ip_access" json:"ip_access" yaml:"ip_access"`                            // IP访问控制
}

type IPAccessConfig struct {
	Enable bool            `ini:"enable" json:"enable" yaml:"enable"`
	Allow  *ipmatch.Config `ini:"allow" json:"allow" yaml:"allow"` // 允许的IP
	Deny   *ipmatch.Config `ini:"deny" json:"deny" yaml:"deny"`    // 屏蔽的IP
	Debug  *ipmatch.Config `ini:"debug" json:"debug" yaml:"debug"` // 调试允许的IP
}

type RouteType int

const (
	RouteTypeHTTP RouteType = iota
	RouteTypeWebSocket
)

type routerRule struct {
	routeType RouteType // 新增：路由类型
	path      string    // 路由
	handle    http.Handler
	wsHandle  websocket.Handler
	cb        []MiddlewareFunc // 中间件
}
