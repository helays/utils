package csrf_std

import (
	"net/http"

	"github.com/helays/utils/v2/map/safemap"
	"github.com/helays/utils/v2/security/csrf"
)

type Std struct {
	configs safemap.SyncMap[string, *csrf.Config]
}

func NewStd() *Std {
	s := &Std{}

	return s
}

// SetConfig 设置路径的CSRF配置
func (s *Std) SetConfig(pattern string, config *csrf.Config) {
	if config == nil {
		return
	}
	s.configs.Store(pattern, config)
}

func (s *Std) GetConfig(pattern string) (*csrf.Config, bool) {
	return s.configs.Load(pattern)
}

// WrapHandler 包装单个处理器
func (s *Std) WrapHandler(pattern string, handler http.HandlerFunc) http.HandlerFunc {
	config, ok := s.GetConfig(pattern)
	if !ok {
		return handler
	}
	return WrapHandler(handler, config)
}
