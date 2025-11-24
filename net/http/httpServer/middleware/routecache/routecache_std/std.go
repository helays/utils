package routecache_std

import (
	"context"
	"net/http"

	"github.com/helays/utils/v2/net/http/httpServer/middleware/routecache"
)

type STD[T comparable] struct {
	ctxField string // 往上下文中写入数据的字段
	cache    *routecache.RouteCache[T]
}

func New[T comparable](ctxField string) *STD[T] {
	return &STD[T]{
		ctxField: ctxField,
		cache:    routecache.New[T](),
	}
}

func (s *STD[T]) GetCache() *routecache.RouteCache[T] {
	return s.cache
}

func (s *STD[T]) Middleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return s.WrapHandler(next.ServeHTTP)
	}
}

func (s *STD[T]) WrapHandler(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		method := r.Method
		path := r.URL.Path
		if handle, ok := s.cache.Match(method, path); ok {
			ctx := r.Context()
			ctx = context.WithValue(ctx, s.ctxField, handle)
			r = r.WithContext(ctx)
		}
		handler(w, r)
	}
}
