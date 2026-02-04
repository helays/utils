package httpServer

import (
	"net/http"
	"net/http/pprof"
	"time"

	"helay.net/go/utils/v3/config"
	"helay.net/go/utils/v3/logger/ulogs"
	"helay.net/go/utils/v3/net/http/httpServer/router"
)

// 设置路由
func (h *HttpServer) initRouter() {
	now := time.Now()
	ulogs.Infof("开始初始化http server 路由规则")
	defer ulogs.Infof("http server 路由规则初始化完成，耗时：%v", time.Since(now))
	h.mux = http.NewServeMux()

	h.addDebugRoute()
	h.addDeprecatedRoute() // 添加兼容模式路由

	if h.route == nil {
		return
	}
	for _, r := range h.route {
		switch r.routeType {
		case RouteTypeHTTP:
			h.middleware(r.path, r.handle, r.cb...)
		case RouteTypeWebSocket:
			h.socketMiddleware(r.path, r.wsHandle)
		}
	}

}

// 添加调试路由
func (h *HttpServer) addDebugRoute() {
	debugGroup := h.Group("/debug")
	debugGroup.Get("/switch-debug", router.SwitchDebug)
	if config.Dbg {
		debugGroup.Get("/pprof/", pprof.Index)
		debugGroup.Get("/pprof/cmdline", pprof.Cmdline)
		debugGroup.Get("/pprof/profile", pprof.Profile)
		debugGroup.Get("/pprof/symbol", pprof.Symbol)
		debugGroup.Get("/pprof/trace", pprof.Trace)
	}
}

func (h *HttpServer) addDeprecatedRoute() {
	if h.Route != nil {
		for u, funcName := range h.Route {
			h.addRoute(u, funcName)
		}
	}

	if h.RouteHandle != nil {
		for u, funcName := range h.RouteHandle {
			h.addRoute(u, funcName)
		}
	}

	if h.RouteSocket != nil {
		for u, funcName := range h.RouteSocket {
			//mux.Handle(u, websocket.Handler(funcName))
			h.addWSRoute(u, funcName)
		}
	}
}
