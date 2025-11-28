package server

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/helays/utils/v2/close/httpClose"
	"github.com/helays/utils/v2/logger/ulogs"
	"github.com/helays/utils/v2/net/http/mime"
	"github.com/helays/utils/v2/net/http/route"
	"github.com/helays/utils/v2/net/http/route/middleware"
	"github.com/helays/utils/v2/net/ipmatch"
	"github.com/helays/utils/v2/tools"
)

// noinspection all
func New(cfg *Config) (*Server[any], error) {
	return NewGeneric[any](cfg)
}

func NewGeneric[T any](cfg *Config) (*Server[T], error) {
	s := &Server[T]{
		opt:      cfg,
		ipAccess: middleware.NewIPAccessMiddleware(),
	}
	if err := s.validParam(); err != nil {
		return nil, err
	}

	s.enhancedWriter = middleware.NewCompression()
	s.enhancedWriter.SetCompressionConfig(s.opt.Compression)
	if err := s.enhancedWriter.SetLoggerConfig(s.opt.Logger); err != nil {
		return nil, fmt.Errorf("日志模块初始化失败 %v", err)
	}

	if err := s.ipAccessInit(); err != nil {
		return nil, err
	}
	s.server = &http.Server{
		Addr:                         s.opt.Addr,
		DisableGeneralOptionsHandler: s.opt.DisableGeneralOptionsHandler,
		ReadTimeout:                  tools.AutoTimeDuration(s.opt.ReadTimeout, time.Second),
		ReadHeaderTimeout:            tools.AutoTimeDuration(s.opt.ReadHeaderTimeout, time.Second),
		WriteTimeout:                 tools.AutoTimeDuration(s.opt.WriteTimeout, time.Second),
		IdleTimeout:                  tools.AutoTimeDuration(s.opt.IdleTimeout, time.Second),
		MaxHeaderBytes:               s.opt.MaxHeaderBytes,
	}
	if err := s.tls(); err != nil {
		return nil, err
	}
	s.routes = make(map[string]*routerRule[T])
	s.route = route.New(s.opt.Route) // 系统 通用路由

	mime.InitMimeTypes()
	return s, nil
}

// 验证参数
func (s *Server[T]) validParam() error {
	s.serverNames = make(map[string]struct{})
	for _, name := range s.opt.ServerName {
		s.serverNames[strings.ToLower(name)] = struct{}{}
	}

	return nil
}

func (s *Server[T]) ipAccessInit() error {
	access := s.opt.Security.IPAccess
	if !access.Enable {
		ulogs.Infof("HTTP服务未启用 IP访问控制模块")
		return nil
	}
	now := time.Now()
	ulogs.Infof("开始初始化IP访问控制")

	if access.Allow != nil {
		if allow, err := ipmatch.NewIPMatcher(access.Allow); err != nil {
			return fmt.Errorf("HTTP服务IP访问控制模块 [白名单] 初始化失败 %v", err)
		} else {
			s.ipAccess.SetAllow(allow)
		}
	}
	if access.Deny != nil {
		if deny, err := ipmatch.NewIPMatcher(access.Deny); err != nil {
			return fmt.Errorf("HTTP服务IP访问控制模块 [黑名单] 初始化失败 %v", err)
		} else {
			s.ipAccess.SetDeny(deny)
		}
	}
	if access.Debug != nil {
		if dbg, err := ipmatch.NewIPMatcher(access.Debug); err != nil {
			return fmt.Errorf("HTTP服务IP访问控制模块 [调试] 模块初始化失败 %v", err)
		} else {
			s.ipAccess.SetDebug(dbg)
		}
	}
	ulogs.Infof("初始化IP访问控制完成，耗时：%v", time.Since(now))
	return nil
}

// TLS 配置
func (s *Server[T]) tls() error {
	tlsOpts := s.opt.TLS
	if !tlsOpts.Enable {
		ulogs.Infof("HTTP服务未启用 TLS")
		return nil
	}
	tlsConfig, err := tlsOpts.ToTLSConfig()
	if err != nil {
		return err
	}
	s.server.TLSConfig = tlsConfig
	return nil
}

func (s *Server[T]) Close() {
	s.ipAccess.Close()
	httpClose.Server(s.server)
	ulogs.Log("http server已关闭")
}

// Run 启动服务
func (s *Server[T]) Run() error {
	s.mux = http.NewServeMux()
	s.setRoutes()
	s.server.Handler = s.mux

	if s.opt.TLS.Enable {
		ulogs.Log("启动https server", s.opt.Addr)
		return s.server.ListenAndServeTLS("", "")
	}
	ulogs.Log("启动http server", s.opt.Addr)
	return s.server.ListenAndServe()
}

func (s *Server[T]) GetRouteDescriptions() []Description[T] {
	var routes = make([]Description[T], 0, len(s.routes))
	for _, r := range s.routes {
		routes = append(routes, r.description)
	}
	return routes
}

func (s *Server[T]) GetRoute() *route.Route {
	return s.route
}
