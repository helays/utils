package server

import (
	"context"
	"crypto/tls"
	"errors"
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
	"github.com/quic-go/quic-go"
	quicHttp3 "github.com/quic-go/quic-go/http3"
)

// noinspection all
func New(ctx context.Context, cfg *Config) (*Server[any], error) {
	return NewGeneric[any](ctx, cfg)
}

func NewGeneric[T any](ctx context.Context, cfg *Config) (*Server[T], error) {
	s := &Server[T]{
		opt:      cfg,
		ipAccess: middleware.NewIPAccessMiddleware(),
	}
	if err := s.validParam(); err != nil {
		return nil, err
	}

	s.enhancedWriter = middleware.NewResponseProcessor()
	s.enhancedWriter.SetCompressionConfig(s.opt.Compression)

	if len(s.opt.Logger.LogLevelConfigs) > 0 {
		logger, err := middleware.NewZapLogger(&s.opt.Logger)
		if err != nil {
			return nil, fmt.Errorf("日志模块初始化失败 %v", err)
		}
		s.enhancedWriter.AddLogHandler(logger)
	}

	if err := s.ipAccessInit(ctx); err != nil {
		return nil, err
	}
	s.httpServer()
	s.quicHttp3Server()
	if err := s.tls(); err != nil {
		return nil, err
	}
	s.routes = make(map[string]*routerRule[T])
	s.route = route.New(&s.opt.Route) // 系统 通用路由

	mime.InitMimeTypes()
	return s, nil
}

// http1.1 http2 配置处理
func (s *Server[T]) httpServer() {
	s.server = &http.Server{
		Addr:                         s.opt.Addr,
		DisableGeneralOptionsHandler: s.opt.DisableGeneralOptionsHandler,
		ReadTimeout:                  tools.AutoTimeDuration(s.opt.ReadTimeout, time.Second),
		ReadHeaderTimeout:            tools.AutoTimeDuration(s.opt.ReadHeaderTimeout, time.Second),
		WriteTimeout:                 tools.AutoTimeDuration(s.opt.WriteTimeout, time.Second),
		IdleTimeout:                  tools.AutoTimeDuration(s.opt.IdleTimeout, time.Second),
		MaxHeaderBytes:               s.opt.MaxHeaderBytes,
	}
}

// http3 配置处理
func (s *Server[T]) quicHttp3Server() {
	if !s.opt.EnableQuickH3 {
		return
	}
	s.quicH3Server = &quicHttp3.Server{
		Addr: s.opt.Addr,
		Port: s.opt.Port,
		QUICConfig: &quic.Config{
			HandshakeIdleTimeout:             tools.AutoTimeDuration(s.opt.QUICConfig.HandshakeIdleTimeout, time.Second),
			MaxIdleTimeout:                   tools.AutoTimeDuration(s.opt.QUICConfig.MaxIdleTimeout, time.Second),
			InitialStreamReceiveWindow:       s.opt.QUICConfig.InitialStreamReceiveWindow,
			MaxStreamReceiveWindow:           s.opt.QUICConfig.MaxStreamReceiveWindow,
			InitialConnectionReceiveWindow:   s.opt.QUICConfig.InitialConnectionReceiveWindow,
			MaxConnectionReceiveWindow:       s.opt.QUICConfig.MaxConnectionReceiveWindow,
			MaxIncomingStreams:               s.opt.QUICConfig.MaxIncomingStreams,
			MaxIncomingUniStreams:            s.opt.QUICConfig.MaxIncomingUniStreams,
			KeepAlivePeriod:                  s.opt.QUICConfig.KeepAlivePeriod,
			InitialPacketSize:                s.opt.QUICConfig.InitialPacketSize,
			DisablePathMTUDiscovery:          s.opt.QUICConfig.DisablePathMTUDiscovery,
			Allow0RTT:                        s.opt.QUICConfig.Allow0RTT,
			EnableDatagrams:                  s.opt.QUICConfig.EnableDatagrams,
			EnableStreamResetPartialDelivery: s.opt.QUICConfig.EnableStreamResetPartialDelivery,
		},
		EnableDatagrams:    s.opt.EnableDatagrams,
		MaxHeaderBytes:     s.opt.MaxHeaderBytes,
		AdditionalSettings: s.opt.AdditionalSettings,

		IdleTimeout: tools.AutoTimeDuration(s.opt.IdleTimeout, time.Second),
	}

}

// 验证参数
func (s *Server[T]) validParam() error {
	s.serverNames = make(map[string]struct{})
	for _, name := range s.opt.ServerName {
		s.serverNames[strings.ToLower(name)] = struct{}{}
	}

	return nil
}

func (s *Server[T]) ipAccessInit(ctx context.Context) error {
	access := s.opt.Security.IPAccess
	if !access.Enable {
		ulogs.Infof("HTTP服务未启用 IP访问控制模块")
		return nil
	}
	now := time.Now()
	ulogs.Infof("开始初始化IP访问控制")

	if access.Allow != nil {
		if allow, err := ipmatch.NewIPMatcher(ctx, access.Allow); err != nil {
			return fmt.Errorf("HTTP服务IP访问控制模块 [白名单] 初始化失败 %v", err)
		} else {
			s.ipAccess.SetAllow(allow)
		}
	}
	if access.Deny != nil {
		if deny, err := ipmatch.NewIPMatcher(ctx, access.Deny); err != nil {
			return fmt.Errorf("HTTP服务IP访问控制模块 [黑名单] 初始化失败 %v", err)
		} else {
			s.ipAccess.SetDeny(deny)
		}
	}
	if access.Debug != nil {
		if dbg, err := ipmatch.NewIPMatcher(ctx, access.Debug); err != nil {
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

	if s.quicH3Server != nil {
		tlsConfig, err := tlsOpts.ToTLSConfig()
		if err != nil {
			return err
		}
		s.quicH3Server.TLSConfig = tlsConfig
	}
	tlsConfig, err := tlsOpts.ToTLSConfig()
	if err != nil {
		return err
	}
	s.server.TLSConfig = tlsConfig
	return nil
}

func (s *Server[T]) close() {
	httpClose.Server(s.server)
	httpClose.ServerQuick(s.quicH3Server)
	ulogs.Log("http server已关闭")
}

// Run 启动服务
func (s *Server[T]) Run(ctx context.Context) error {
	go tools.RunOnContextDone(ctx, func() { s.close() })
	s.mux = http.NewServeMux()
	s.setRoutes()
	s.server.Handler = s.mux

	if s.quicH3Server != nil {
		s.quicH3Server.Handler = s.mux
		go func() {
			ulogs.Log("启动quic http3服务", s.opt.Addr)
			err := s.quicH3Server.ListenAndServe()
			if err != nil && !errors.Is(err, http.ErrServerClosed) {
				panic(fmt.Errorf("quic http3服务启动失败：%s", err))
			}
		}()
	}

	if s.opt.TLS.Enable {
		ulogs.Log("启动https server", s.opt.Addr)
		err := s.server.ListenAndServeTLS("", "")
		if err == nil || errors.Is(err, http.ErrServerClosed) {
			return nil
		}
		return fmt.Errorf("https服务启动失败：%s", err)
	}
	ulogs.Log("启动http server", s.opt.Addr)
	err := s.server.ListenAndServe()
	if err == nil || errors.Is(err, http.ErrServerClosed) {
		return nil
	}
	return fmt.Errorf("http服务启动失败：%s", err)
}

func (s *Server[T]) GetRouteDescriptions() []Description[T] {
	var routes = make([]Description[T], 0, len(s.routes))
	for _, r := range s.routes {
		r.description.Path = r.path
		r.description.Method = r.method
		routes = append(routes, r.description)
	}
	return routes
}

func (s *Server[T]) GetRoute() *route.Route {
	return s.route
}

// AddLogHandler 添加日志处理
func (s *Server[T]) AddLogHandler(le ...middleware.Logger) {
	s.enhancedWriter.AddLogHandler(le...)
}

func (s *Server[T]) TLS() *tls.Config {
	
	return s.server.TLSConfig
}
