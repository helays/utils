package server

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/helays/utils/v2/close/httpClose"
	"github.com/helays/utils/v2/logger/ulogs"
	"github.com/helays/utils/v2/logger/zaploger"
	"github.com/helays/utils/v2/net/http/mime"
	"github.com/helays/utils/v2/net/ipmatch"
)

func New(cfg *Config) (*Server[any], error) {
	return NewGeneric[any](cfg)
}

func NewGeneric[T any](cfg *Config) (*Server[T], error) {
	s := &Server[T]{
		opt: cfg,
	}
	if err := s.validParam(); err != nil {
		return nil, err
	}
	if err := s.ipAccessInit(); err != nil {
		return nil, err
	}
	s.server = &http.Server{Addr: s.opt.ListenAddr}
	if err := s.tls(); err != nil {
		return nil, err
	}
	mime.InitMimeTypes()
	return s, nil
}

// 验证参数
func (s *Server[T]) validParam() error {
	s.serverNames = make(map[string]struct{})
	for _, name := range s.opt.ServerName {
		s.serverNames[strings.ToLower(name)] = struct{}{}
	}
	if len(s.opt.Logger.LogLevelConfigs) > 0 {
		var err error
		s.logger, err = zaploger.New(&s.opt.Logger)
		if err != nil {
			return fmt.Errorf("HTTP 服务日志初始化失败 %v", err)
		}
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
	var err error
	if access.Allow != nil {
		s.allowIPMatch, err = ipmatch.NewIPMatcher(access.Allow)
		if err != nil {
			return fmt.Errorf("HTTP服务IP访问控制模块 [白名单] 初始化失败 %v", err)
		}
		s.allowIPMatch.Build()
	}
	if access.Deny != nil {
		s.denyIPMatch, err = ipmatch.NewIPMatcher(access.Deny)
		if err != nil {
			return fmt.Errorf("HTTP服务IP访问控制模块 [黑名单] 初始化失败 %v", err)
		}
		s.denyIPMatch.Build()
	}
	if access.Debug != nil {
		s.debugIPMatch, err = ipmatch.NewIPMatcher(access.Debug)
		if err != nil {
			return fmt.Errorf("HTTP服务IP访问控制模块 [调试] 模块初始化失败 %v", err)
		}
		s.debugIPMatch.Build()
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

func (s *Server[T]) Run() error {
	s.compression()
	return nil
}

func (s *Server[T]) compression() {

}

func (s *Server[T]) Close() {
	if s.denyIPMatch != nil {
		s.denyIPMatch.Close()
	}
	if s.debugIPMatch != nil {
		s.debugIPMatch.Close()
	}
	if s.allowIPMatch != nil {
		s.allowIPMatch.Close()
	}
	httpClose.Server(s.server)
	ulogs.Log("http server已关闭")
}
