package server

import (
	"net/http"
	"time"

	"github.com/helays/utils/v2/logger/zaploger"
	"github.com/helays/utils/v2/net/http/route/middleware"
	"github.com/helays/utils/v2/net/ipmatch"
	"github.com/helays/utils/v2/security/cors"
	"golang.org/x/net/websocket"
)

const Version = "vs/2.1"

type Config struct {
	Addr                         string        `json:"addr" yaml:"addr"`                                                       // 监听地址
	DisableGeneralOptionsHandler bool          `json:"disable_general_options_handler" yaml:"disable_general_options_handler"` // 如果为 true，将 "OPTIONS *" 请求传递给 Handler；否则自动响应 200 O
	ReadTimeout                  time.Duration `json:"read_timeout" yaml:"read_timeout"`                                       // 读取超时，零或负值表示无超时
	ReadHeaderTimeout            time.Duration `json:"read_header_timeout" yaml:"read_header_timeout"`                         // 读取请求头超时，零或负值表示无超时
	WriteTimeout                 time.Duration `json:"write_timeout" yaml:"write_timeout"`                                     // 写入超时，零或负值表示无超时
	IdleTimeout                  time.Duration `json:"idle_timeout" yaml:"idle_timeout"`                                       // 在 keep-alive 启用时等待下一个请求的最大时间，如果为零，使用 ReadTimeout 的值
	MaxHeaderBytes               int           `json:"max_header_bytes" yaml:"max_header_bytes"`                               // 服务器解析请求头时读取的最大字节数（包括请求行），如果为零，使用 DefaultMaxHeaderBytes（1MB）
	ServerName                   []string      `json:"server_name" yaml:"server_name"`                                         // 绑定域名

	TLS         TLSConfig         `json:"tls" yaml:"tls"`                          // TLS 配置
	Security    SecurityConfig    `ini:"security" json:"security" yaml:"security"` // 安全配置
	Compression CompressionConfig `json:"compression" yaml:"compression"`          // 压缩配置
	Logger      zaploger.Config   `json:"logger" yaml:"logger"`                    // 日志配置
}

type SecurityConfig struct {
	CORS      cors.Config    `ini:"cors" json:"cors" yaml:"cors"`                   // 跨域配置
	IPAccess  IPAccessConfig `ini:"ip_access" json:"ip_access" yaml:"ip_access"`    // IP访问控制
	DebugPath string         `ini:"debug_path" json:"debug_path" yaml:"debug_path"` // 调试路径,默认debug
}

type IPAccessConfig struct {
	Enable bool            `ini:"enable" json:"enable" yaml:"enable"`
	Allow  *ipmatch.Config `ini:"allow" json:"allow" yaml:"allow"` // 允许的IP
	Deny   *ipmatch.Config `ini:"deny" json:"deny" yaml:"deny"`    // 屏蔽的IP
	Debug  *ipmatch.Config `ini:"debug" json:"debug" yaml:"debug"` // 调试允许的IP
}

type CompressionConfig struct {
	Enabled      bool     `json:"enabled" yaml:"enabled" ini:"enabled"`                   // 是否启用压缩
	Algorithm    string   `json:"algorithm" yaml:"algorithm" ini:"algorithm"`             // 压缩算法：gzip, deflate, br 等
	Level        int      `json:"level" yaml:"level" ini:"level"`                         // 压缩级别
	MinSize      int      `json:"min_size" yaml:"min_size" ini:"min_size"`                // 最小压缩大小（字节），小于此值不压缩
	ContentTypes []string `json:"content_types" yaml:"content_types" ini:"content_types"` // 需要压缩的 MIME 类型
	ExcludePaths []string `json:"exclude_paths" yaml:"exclude_paths" ini:"exclude_paths"` // 排除的路径
}

type Server[T any] struct {
	opt         *Config
	serverNames map[string]struct{}       // 绑定的域名
	routes      map[string]*routerRule[T] // 路由集合

	logger   *middleware.LoggerMiddleware   // 日志中间件
	ipAccess *middleware.IPAccessMiddleware // IP访问控制中间件

	mux    *http.ServeMux
	server *http.Server
}

type routerRule[T any] struct {
	routeType   RouteType // 新增：路由类型
	method      string    // 请求方法
	path        string    // 路由
	handle      http.Handler
	wsHandle    websocket.Handler
	middlewares []Middleware   // 中间件
	description Description[T] // 描述
}

type RouteType int

const (
	RouteTypeHTTP RouteType = iota
	RouteTypeWebSocket
)

// Middleware 修改为支持 http.Handler
type Middleware func(next http.Handler) http.Handler

type RouteCode string

func (c RouteCode) String() string {
	return string(c)
}

func NewRouteCode(code string) RouteCode {
	return RouteCode(code)
}

type Description[T any] struct {
	Name     string    // 路由名称
	Version  string    // 路由版本
	CodeKey  string    // 路由编码字段
	Code     RouteCode // 路由编码
	Metadata T         // 自定义描述结构
}
