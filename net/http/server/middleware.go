package server

import "net/http"

// Chain 中间件链（兼容 http.Handler）
func Chain(middlewares ...Middleware) Middleware {
	return func(final http.Handler) http.Handler {
		for i := len(middlewares) - 1; i >= 0; i-- {
			final = middlewares[i](final)
		}
		return final
	}
}

func (s *Server[T]) middleware(route *routerRule[T]) {

}

func (s *Server[T]) wsMiddleware(route *routerRule[T]) {
}
