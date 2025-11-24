package cors_std

import (
	"net/http"

	"github.com/helays/utils/v2/map/safemap"
	"github.com/helays/utils/v2/security/cors"
)

type StdCORS struct {
	configs           safemap.SyncMap[string, *cors.Config]
	routeCodeCtxField string // 路由code字段
}

func New(routeCodeCtxField string) *StdCORS {
	return &StdCORS{routeCodeCtxField: routeCodeCtxField}
}

func (s *StdCORS) SetConfig(pattern string, config *cors.Config) {
	s.configs.Store(pattern, config)
}

func (s *StdCORS) GetConfig(pattern string) (*cors.Config, bool) {
	return s.configs.Load(pattern)
}

func (s *StdCORS) Middleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return s.WrapHandler(next.ServeHTTP)
	}
}

func (s *StdCORS) WrapHandler(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var path string
		if s.routeCodeCtxField == "" {
			path = r.URL.Path
		} else {
			path = r.Context().Value(s.routeCodeCtxField).(string)
		}
		if config, exists := s.GetConfig(path); exists {
			if r.Method == http.MethodOptions {
				config.HandlePreflight(w, r.Header.Get("Origin"))
				return
			}
			if shouldContinue := config.Apply(w, r.Header.Get("Origin")); !shouldContinue {
				// 严格模式下被拒绝
				return
			}
		}

		handler(w, r)
	}
}
