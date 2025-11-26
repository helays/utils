package server

import (
	"net/http"
	"net/http/pprof"
	"time"

	"github.com/helays/utils/v2/config"
	"github.com/helays/utils/v2/logger/ulogs"
	"github.com/helays/utils/v2/tools"
	"golang.org/x/net/websocket"
)

func (s *Server[T]) setRoutes() {
	now := time.Now()
	ulogs.Infof("开始初始化 http router")
	defer ulogs.Infof("完成初始化 http router,耗时:%v", time.Since(now))
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
	group.Get("/pprof/", pprof.Index)
	group.Get("/pprof/cmdline", pprof.Cmdline)
	group.Get("/pprof/profile", pprof.Profile)
	group.Get("/pprof/symbol", pprof.Symbol)
	group.Get("/pprof/trace", pprof.Trace)
}

func (s *Server[T]) AddRoute(method, path string, handle http.Handler, cb ...Middleware) {
	s.routes[path] = &routerRule[T]{
		routeType:   RouteTypeHTTP,
		method:      method,
		path:        path,
		handle:      handle,
		middlewares: cb,
	}
}

func (s *Server[T]) AddRouteWithDescription(method, path string, handle http.Handler, description Description[T], cb ...Middleware) {
	s.routes[path] = &routerRule[T]{
		routeType:   RouteTypeHTTP,
		method:      method,
		path:        path,
		handle:      handle,
		middlewares: cb,
		description: description,
	}
}

func (s *Server[T]) AddWebsocketRoute(path string, handle websocket.Handler, cb ...WSMiddleware) {
	s.routes[path] = &routerRule[T]{
		routeType:     RouteTypeWebSocket,
		path:          path,
		wsHandle:      handle,
		wsMiddlewares: cb,
	}
}

func (s *Server[T]) AddWebsocketRouteWithDescription(path string, handle websocket.Handler, description Description[T], cb ...WSMiddleware) {
	s.routes[path] = &routerRule[T]{
		routeType:     RouteTypeWebSocket,
		path:          path,
		wsHandle:      handle,
		wsMiddlewares: cb,
		description:   description,
	}
}
