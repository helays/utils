package server

import (
	"context"
	"net/http"

	"github.com/helays/utils/v2/close/httpClose"
	"github.com/helays/utils/v2/close/vclose"
	"github.com/helays/utils/v2/security/cors/cors_std"
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
		s.ipAccess.DenyIpAccess,            // 设置黑名单防火墙
		s.ipAccess.AllowIpAccess,           // 设置白名单防火墙
		s.ipAccess.DebugIPAccess,           // debug 模式ip限制
		cors_std.Cors(s.opt.Security.CORS), // 跨域，这个是配置文件级别的跨域中间件。
		s.enhancedWriter.Handler,           // 多功能响应处理器
	}
	mid = append(mid, route.middlewares...)

	handler = Chain(mid...)(route.handle)

	// 判断当前路由 是否设置了code key 和 code。
	// 用于将路由默认的编码信息打入请求上下文中，后续在请求中就可直接使用当前code 去匹配当前路由的一些基本信息。
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
		s.ipAccess.DenyIpAccess,            // 设置黑名单防火墙
		s.ipAccess.AllowIpAccess,           // 设置白名单防火墙
		cors_std.Cors(s.opt.Security.CORS), // 跨域
		s.enhancedWriter.Handler,           // 多功能响应处理器
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
