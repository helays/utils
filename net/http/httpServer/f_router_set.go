package httpServer

import (
	"github.com/helays/utils/config"
	"github.com/helays/utils/net/http/httpServer/router"
	"net/http/pprof"
)

// 设置路由
func (h *HttpServer) initRouter() {
	h.Route["/debug/switch-debug"] = router.SwitchDebug
	if config.Dbg {
		h.Route["/debug/pprof/"] = pprof.Index
		h.Route["/debug/pprof/cmdline"] = pprof.Cmdline
		h.Route["/debug/pprof/profile"] = pprof.Profile
		h.Route["/debug/pprof/symbol"] = pprof.Symbol
		h.Route["/debug/pprof/trace"] = pprof.Trace
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
