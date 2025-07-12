package httpServer

import (
	"github.com/helays/utils/config"
	"github.com/helays/utils/net/http/httpServer/router"
	"net/http/pprof"
)

// 设置路由
func (h *HttpServer) initRouter() {
	debugGroup := h.Group("/debug")
	debugGroup.Get("/switch-debug", router.SwitchDebug)
	if config.Dbg {
		debugGroup.Get("/pprof/", pprof.Index)
		debugGroup.Get("/pprof/cmdline", pprof.Cmdline)
		debugGroup.Get("/pprof/profile", pprof.Profile)
		debugGroup.Get("/pprof/symbol", pprof.Symbol)
		debugGroup.Get("/pprof/trace", pprof.Trace)
	}
	if h.Route != nil {
		for u, funcName := range h.Route {
			h.middleware(u, funcName)
		}
	}

	if h.RouteHandle != nil {
		for u, funcName := range h.RouteHandle {
			h.middleware(u, funcName)
		}
	}

	if h.RouteSocket != nil {
		for u, funcName := range h.RouteSocket {
			//mux.Handle(u, websocket.Handler(funcName))
			h.socketMiddleware(u, funcName)
		}
	}
}
