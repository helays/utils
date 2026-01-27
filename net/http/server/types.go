package server

import (
	"net/http"
	"time"

	"github.com/helays/utils/v2/logger/zaploger"
	"github.com/helays/utils/v2/net/http/route"
	"github.com/helays/utils/v2/net/http/route/middleware"
	"github.com/helays/utils/v2/net/ipmatch"
	"github.com/helays/utils/v2/net/tlsconfig"
	"github.com/helays/utils/v2/security/cors"
	quicHttp3 "github.com/quic-go/quic-go/http3"
	"golang.org/x/net/websocket"
)

const Version = "vs/2.1"

type Config struct {
	Addr                         string        `json:"addr" yaml:"addr"`                                                       // 监听地址
	Port                         int           `json:"port" yaml:"port"`                                                       //用于 Alt-Svc 响应头中的端口,用于防火墙重定向等场景，允许客户端使用与服务器监听端口不同的端口
	DisableGeneralOptionsHandler bool          `json:"disable_general_options_handler" yaml:"disable_general_options_handler"` // 如果为 true，将 "OPTIONS *" 请求传递给 Handler；否则自动响应 200 O
	ReadTimeout                  time.Duration `json:"read_timeout" yaml:"read_timeout"`                                       // 读取超时，零或负值表示无超时
	ReadHeaderTimeout            time.Duration `json:"read_header_timeout" yaml:"read_header_timeout"`                         // 读取请求头超时，零或负值表示无超时
	WriteTimeout                 time.Duration `json:"write_timeout" yaml:"write_timeout"`                                     // 写入超时，零或负值表示无超时
	IdleTimeout                  time.Duration `json:"idle_timeout" yaml:"idle_timeout"`                                       // 在 keep-alive 启用时等待下一个请求的最大时间，如果为零，使用 ReadTimeout 的值
	MaxHeaderBytes               int           `json:"max_header_bytes" yaml:"max_header_bytes"`                               // 服务器解析请求头时读取的最大字节数（包括请求行），如果为零，使用 DefaultMaxHeaderBytes（1MB）
	ServerName                   []string      `json:"server_name" yaml:"server_name"`                                         // 绑定域名

	TLS tlsconfig.TLSConfig `json:"tls" yaml:"tls"` // TLS 配置

	EnableQuickH3      bool              `json:"enable_quick_h3" yaml:"enable_quick_h3"`         // 启用 QUIC HTTP/3
	EnableDatagrams    bool              `json:"enable_datagrams" yaml:"enable_datagrams"`       // 启用 HTTP/3 数据报支持（RFC 9297）
	AdditionalSettings map[uint64]uint64 `json:"additional_settings" yaml:"additional_settings"` // 额外的设置，QUIC 包配置
	QUICConfig         QUICConfig        `json:"quic_config" yaml:"quic_config"`                 // QUIC 配置

	Security    SecurityConfig               `ini:"security" json:"security" yaml:"security"` // 安全配置
	Compression middleware.CompressionConfig `json:"compression" yaml:"compression"`          // 压缩配置
	Logger      zaploger.Config              `json:"logger" yaml:"logger"`                    // 日志配置

	Route route.Config `json:"route" yaml:"route"` // 路由配置
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

type Server[T any] struct {
	opt         *Config
	serverNames map[string]struct{}       // 绑定的域名
	routes      map[string]*routerRule[T] // 路由集合
	route       *route.Route              // 系统默认路由

	enhancedWriter *middleware.ResponseProcessor  // 通用响应处理中间件
	ipAccess       *middleware.IPAccessMiddleware // IP 访问控制中间件

	mux          *http.ServeMux
	server       *http.Server
	quicH3Server *quicHttp3.Server
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

type Description[T any] struct {
	Name     string // 路由名称
	Version  string // 路由版本
	CodeKey  string // 路由编码字段
	Method   string
	Path     string
	Code     string // 路由编码
	Metadata *T     // 自定义描述结构
}

func NewDesc[T any](key, code string, meta *T) Description[T] {
	return Description[T]{
		Code:     code,
		CodeKey:  key,
		Metadata: meta,
	}
}

var HttpMethodList = []string{
	http.MethodGet,
	http.MethodHead,
	http.MethodPost,
	http.MethodPut,
	http.MethodPatch,
	http.MethodDelete,
	http.MethodConnect,
	http.MethodOptions,
	http.MethodTrace,
}
