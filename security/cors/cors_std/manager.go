package cors_std

import (
	"context"
	"net/http"

	"helay.net/go/utils/v3/safe"
	"helay.net/go/utils/v3/security/cors"
)

type StdCORS struct {
	configs           *safe.Map[string, *cors.Config]
	routeCodeCtxField string // 路由code字段
}

func New(ctx context.Context, routeCodeCtxField string) *StdCORS {
	s := &StdCORS{routeCodeCtxField: routeCodeCtxField}
	s.configs = safe.NewMap[string, *cors.Config](ctx, safe.StringHasher{})
	return s
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
		} else if val := r.Context().Value(s.routeCodeCtxField); val != nil {
			path = val.(string)
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
