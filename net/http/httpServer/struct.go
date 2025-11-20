package httpServer

import (
	"context"
	"net/http"
	"time"

	"github.com/helays/utils/v2/logger/zaploger"
	"github.com/helays/utils/v2/net/ipAccess"
	"github.com/helays/utils/v2/security/cors"
	"github.com/helays/utils/v2/tools/mutex"
	"golang.org/x/net/websocket"
)

type HttpServer struct {
	ListenAddr          string                      `ini:"listen_addr" json:"listen_addr" yaml:"listen_addr"`
	Auth                string                      `ini:"auth" json:"auth" yaml:"auth"`
	Allowip             []string                    `ini:"allowip,omitempty" json:"allowip" yaml:"allowip"`
	Denyip              []string                    `ini:"denyip,omitempty" json:"denyip" yaml:"denyip"`
	DebugAllowIp        []string                    `ini:"debug_allow_ip,omitempty" json:"debug_allow_ip" yaml:"debug_allow_ip"`
	ServerName          []string                    `ini:"server_name,omitempty" json:"server_name" yaml:"server_name"` // 绑定域名
	Ssl                 bool                        `ini:"ssl" json:"ssl" yaml:"ssl"`
	Ca                  string                      `ini:"ca" json:"ca" yaml:"ca"`
	Crt                 string                      `ini:"crt" json:"crt" yaml:"crt"`
	Key                 string                      `ini:"key" json:"key" yaml:"key"`
	SocketTimeout       time.Duration               `ini:"socket_timeout" json:"socket_timeout" yaml:"socket_timeout"`                // socket 心跳超时时间
	Hotupdate           bool                        `ini:"hotupdate" json:"hotupdate" yaml:"hotupdate"`                               // 是否启动热加载
	EnableGzip          bool                        `ini:"enable_gzip" json:"enable_gzip" yaml:"enable_gzip"`                         // 是否开启gzip
	DefaultValidFirst   bool                        `ini:"default_valid_first" json:"default_valid_first" yaml:"default_valid_first"` // 默认验证第一
	CORS                *cors.Config                `ini:"cors" json:"cors" yaml:"cors"`                                              // 跨域配置
	Route               map[string]http.HandlerFunc `yaml:"-" json:"-"`
	RouteHandle         map[string]http.Handler
	RouteSocket         map[string]func(ws *websocket.Conn)               `yaml:"-" json:"-"`
	CommonCallback      func(w http.ResponseWriter, r *http.Request) bool `yaml:"-" json:"-"`
	serverNameMap       map[string]byte                                   // 绑定的域名
	Logger              zaploger.Config                                   `json:"logger" yaml:"logger" ini:"logger" gorm:"comment:日志配置"`
	logger              *zaploger.Logger
	enableCheckIpAccess bool // 是否开启ip访问控制
	allowIpList         *ipAccess.IPList
	denyIpList          *ipAccess.IPList
	debugAllowIpList    *ipAccess.IPList

	mux    *http.ServeMux
	server *http.Server
	ctx    context.Context
	cancel context.CancelFunc
	isStop mutex.SafeResourceRWMutex[bool]
}
