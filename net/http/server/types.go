package server

import (
	"net/http"

	"github.com/helays/utils/v2/logger/zaploger"
	"github.com/helays/utils/v2/net/ipmatch"
	"github.com/helays/utils/v2/security/cors"
	"github.com/helays/utils/v2/security/cors/cors_std"
	"golang.org/x/net/websocket"
)

type Config struct {
	ListenAddr string   `json:"listen_addr" yaml:"listen_addr"` // 监听地址
	ServerName []string `json:"server_name" yaml:"server_name"` // 绑定域名

	TLS         TLSConfig         `json:"tls" yaml:"tls"`                          // TLS 配置
	Security    SecurityConfig    `ini:"security" json:"security" yaml:"security"` // 安全配置
	Compression CompressionConfig `json:"compression" yaml:"compression"`          // 压缩配置
	Logger      zaploger.Config   `json:"logger" yaml:"logger"`                    // 日志配置
}

type SecurityConfig struct {
	DefaultValidLast bool           `ini:"default_valid_last" json:"default_valid_last" yaml:"default_valid_last"` // 默认验证器是否放最后
	CORS             cors.Config    `ini:"cors" json:"cors" yaml:"cors"`                                           // 跨域配置
	IPAccess         IPAccessConfig `ini:"ip_access" json:"ip_access" yaml:"ip_access"`                            // IP访问控制
}

type IPAccessConfig struct {
	Enable bool            `ini:"enable" json:"enable" yaml:"enable"`
	Allow  *ipmatch.Config `ini:"allow" json:"allow" yaml:"allow"` // 允许的IP
	Deny   *ipmatch.Config `ini:"deny" json:"deny" yaml:"deny"`    // 屏蔽的IP
	Debug  *ipmatch.Config `ini:"debug" json:"debug" yaml:"debug"` // 调试允许的IP
}

type CompressionConfig struct {
	// 是否启用压缩
	Enabled bool `json:"enabled" yaml:"enabled" ini:"enabled"`

	// 压缩算法：gzip, deflate, br 等
	Algorithm string `json:"algorithm" yaml:"algorithm" ini:"algorithm"`

	// 压缩级别
	Level int `json:"level" yaml:"level" ini:"level"`

	// 最小压缩大小（字节），小于此值不压缩
	MinSize int `json:"min_size" yaml:"min_size" ini:"min_size"`

	// 需要压缩的 MIME 类型
	ContentTypes []string `json:"content_types" yaml:"content_types" ini:"content_types"`

	// 排除的路径
	ExcludePaths []string `json:"exclude_paths" yaml:"exclude_paths" ini:"exclude_paths"`
}

type Server[T any] struct {
	opt         *Config
	serverNames map[string]struct{}       // 绑定的域名
	routes      map[string]*routerRule[T] // 路由集合
	logger      *zaploger.Logger

	allowIPMatch *ipmatch.IPMatcher
	denyIPMatch  *ipmatch.IPMatcher
	debugIPMatch *ipmatch.IPMatcher

	corsManager *cors_std.StdCORS

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

type Description[T any] struct {
	Name     string // 路由名称
	Version  string // 路由版本
	Metadata T      // 自定义描述结构
}
