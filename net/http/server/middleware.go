package server

import (
	"context"
	"net/http"

	"github.com/helays/utils/v2/close/httpClose"
	"github.com/helays/utils/v2/close/vclose"
	"github.com/helays/utils/v2/net/http/route/middleware"
	"golang.org/x/net/websocket"
)

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
	var handler http.Handler
	mid := []Middleware{
		s.ipAccess.DenyIpAccess,              // 设置黑名单防火墙
		s.ipAccess.AllowIpAccess,             // 设置白名单防火墙
		s.ipAccess.DebugIPAccess,             // debug 模式ip限制
		middleware.Cors(s.opt.Security.CORS), // 跨域
		s.enhancedWriter.Handler,             // 多功能响应处理器
	}
	mid = append(mid, route.middlewares...)

	handler = Chain(mid...)(route.handle)

	desc := route.description
	needSetCode := desc.CodeKey != "" && desc.Code != ""
	codeKey := desc.CodeKey
	codeValue := desc.Code

	// 最终的处理函数
	finalHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer httpClose.CloseReq(r)
		if needSetCode {
			ctx := context.WithValue(r.Context(), codeKey, codeValue)
			r = r.WithContext(ctx)
		}
		w.Header().Set("server", Version)
		handler.ServeHTTP(w, r)
	})
	path := route.path
	if route.method != "" {
		path = route.method + " " + route.path
	}
	s.mux.Handle(path, finalHandler)
}

func (s *Server[T]) wsMiddleware(route *routerRule[T]) {
	wsHandler := websocket.Handler(func(conn *websocket.Conn) {
		defer vclose.Close(conn)
		route.wsHandle(conn)
	})
	// 构建中间件链
	mid := []Middleware{
		s.ipAccess.DenyIpAccess,              // 设置黑名单防火墙
		s.ipAccess.AllowIpAccess,             // 设置白名单防火墙
		middleware.Cors(s.opt.Security.CORS), // 跨域
		s.enhancedWriter.Handler,             // 多功能响应处理器
	}
	mid = append(mid, route.middlewares...)

	// 将 WebSocket 处理器包装为 http.Handler
	handler := Chain(mid...)(wsHandler)
	// 最终的处理函数
	finalHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer httpClose.CloseReq(r)
		w.Header().Set("server", Version)
		handler.ServeHTTP(w, r)
	})
	s.mux.Handle(route.path, finalHandler)

}
