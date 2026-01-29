package server

import (
	"context"
	"net/http"
	"net/http/pprof"
	"time"

	"github.com/helays/utils/v2/config"
	"github.com/helays/utils/v2/logger/ulogs"
	"github.com/helays/utils/v2/net/http/route/middleware"
	"github.com/helays/utils/v2/tools"
	"golang.org/x/net/websocket"
)

func (s *Server[T]) setRoutes() {
	now := time.Now()
	ulogs.Infof("开始http router初始化 ")
	defer ulogs.Infof("完成http router初始化，耗时:%v", time.Since(now))
	s.setDebugRoutes()
	for _, route := range s.routes {
		switch route.routeType {
		case RouteTypeHTTP:
			s.middleware(route)
		case RouteTypeWebSocket:
			s.wsMiddleware(route)
		}
	}
}

func (s *Server[T]) setDebugRoutes() {
	if !config.Dbg {
		return
	}
	groupName := tools.Ternary(s.opt.Security.DebugPath == "", "debug", s.opt.Security.DebugPath)
	group := s.Group(groupName)
	group.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := context.WithValue(r.Context(), middleware.DebugCtxKey, true)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	})
	group.Get("/pprof/", pprof.Index)
	group.Get("/pprof/cmdline", pprof.Cmdline)
	group.Get("/pprof/profile", pprof.Profile)
	group.Get("/pprof/symbol", pprof.Symbol)
	group.Get("/pprof/trace", pprof.Trace)
}

func (s *Server[T]) AddRouteFunc(method, path string, handle http.HandlerFunc, cb ...Middleware) {
	s.AddRoute(method, path, handle, cb...)
}

func (s *Server[T]) AddRoute(method, path string, handle http.Handler, cb ...Middleware) {
	if !s.checkUniq(method, path) {
		s.routes = append(s.routes, &routerRule[T]{
			routeType:   RouteTypeHTTP,
			method:      method,
			path:        path,
			handle:      handle,
			middlewares: cb,
		})
	}

}

func (s *Server[T]) AddRouteFuncWithDescription(method, path string, handle http.HandlerFunc, description Description[T], cb ...Middleware) {
	s.AddRouteWithDescription(method, path, handle, description, cb...)
}

func (s *Server[T]) AddRouteWithDescription(method, path string, handle http.Handler, description Description[T], cb ...Middleware) {
	if !s.checkUniq(method, path) {
		s.routes = append(s.routes, &routerRule[T]{
			routeType:   RouteTypeHTTP,
			method:      method,
			path:        path,
			handle:      handle,
			middlewares: cb,
			description: description,
		})
	}
}

func (s *Server[T]) AddWebsocketRoute(path string, handle websocket.Handler) {
	if !s.checkUniq("ws", path) {
		s.routes = append(s.routes, &routerRule[T]{
			routeType: RouteTypeWebSocket,
			path:      path,
			wsHandle:  handle,
		})
	}
}

func (s *Server[T]) AddWebsocketRouteWithDescription(path string, handle websocket.Handler, description Description[T]) {
	if !s.checkUniq("ws", path) {
		s.routes = append(s.routes, &routerRule[T]{
			routeType:   RouteTypeWebSocket,
			path:        path,
			wsHandle:    handle,
			description: description,
		})
	}
}

func (s *Server[T]) checkUniq(method, path string) bool {
	uk := method + path
	if _, ok := s.routesMap[uk]; ok {
		return true
	}
	s.routesMap[uk] = struct{}{}
	return false
}
