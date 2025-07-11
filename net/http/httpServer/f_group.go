package httpServer

import (
	"github.com/helays/utils/net/http/httpServer/response"
	"github.com/helays/utils/tools"
	"golang.org/x/net/websocket"
	"net/http"
	"path"
)

type RouterGroup struct {
	service     *HttpServer
	prefix      string           // 分组前缀
	middlewares []MiddlewareFunc // 中间件链
	parent      *RouterGroup     // 父级分组
}

func (h *HttpServer) Group(g string) *RouterGroup {
	return &RouterGroup{service: h, prefix: g}
}

// Group 创建新分组
func (g *RouterGroup) Group(prefix string) *RouterGroup {
	return &RouterGroup{
		service:     g.service,
		prefix:      g.calculatePrefix(prefix),
		parent:      g,
		middlewares: g.middlewares[:], // 继承父级中间件
	}
}

// 计算完整路径前缀
func (g *RouterGroup) calculatePrefix(p string) string {
	prefix := tools.Ternary(g.prefix == "", "/", g.prefix)
	return path.Join(prefix, p)
}

// 验证请求方法是否正确
func (g *RouterGroup) methodValid(method string, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != method {
			response.MethodNotAllow(w)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func (g *RouterGroup) addRoute(method, p string, handler http.Handler) {
	fullPath := path.Join(g.prefix, p)
	var finalHandler = handler
	if len(g.middlewares) > 0 {
		finalHandler = Chain(g.middlewares...)(finalHandler)
	}
	var cb []MiddlewareFunc
	if method != "" {
		vm := func(next http.Handler) http.Handler {
			return g.methodValid(method, next)
		}
		cb = append(cb, vm)
	}
	g.service.middleware(fullPath, finalHandler, cb...) // 注册路由
}

// WS 添加 WebSocket 支持
func (g *RouterGroup) WS(p string, handler func(ws *websocket.Conn)) {
	fullPath := path.Join(g.prefix, p)
	g.service.socketMiddleware(fullPath, handler)
}

// Use 支持链式调用
func (g *RouterGroup) Use(middleware ...MiddlewareFunc) *RouterGroup {
	g.middlewares = append(g.middlewares, middleware...)
	return g
}

func (g *RouterGroup) GET(p string, handler http.HandlerFunc) {
	g.addRoute(http.MethodGet, p, handler)
}
func (g *RouterGroup) GETHandler(p string, handler http.Handler) {
	g.addRoute(http.MethodGet, p, handler)
}

func (g *RouterGroup) HEAD(p string, handler http.HandlerFunc) {
	g.addRoute(http.MethodHead, p, handler)
}

func (g *RouterGroup) HEADHandler(p string, handler http.Handler) {
	g.addRoute(http.MethodPost, p, handler)
}

func (g *RouterGroup) POST(p string, handler http.HandlerFunc) {
	g.addRoute(http.MethodPost, p, handler)
}

func (g *RouterGroup) POSTHandler(p string, handler http.Handler) {
	g.addRoute(http.MethodPost, p, handler)
}

func (g *RouterGroup) PUT(p string, handler http.HandlerFunc) {
	g.addRoute(http.MethodPut, p, handler)
}
func (g *RouterGroup) PUTHandler(p string, handler http.Handler) {
	g.addRoute(http.MethodPut, p, handler)
}

func (g *RouterGroup) PATCH(p string, handler http.HandlerFunc) {
	g.addRoute(http.MethodPatch, p, handler)
}
func (g *RouterGroup) PATCHHandler(p string, handler http.Handler) {
	g.addRoute(http.MethodPatch, p, handler)
}

func (g *RouterGroup) DELETE(p string, handler http.HandlerFunc) {
	g.addRoute(http.MethodDelete, p, handler)
}
func (g *RouterGroup) DELETEHandler(p string, handler http.Handler) {
	g.addRoute(http.MethodDelete, p, handler)
}

func (g *RouterGroup) CONNECT(p string, handler http.HandlerFunc) {
	g.addRoute(http.MethodConnect, p, handler)
}
func (g *RouterGroup) CONNECTHandler(p string, handler http.Handler) {
	g.addRoute(http.MethodConnect, p, handler)
}

func (g *RouterGroup) OPTIONS(p string, handler http.HandlerFunc) {
	g.addRoute(http.MethodOptions, p, handler)
}
func (g *RouterGroup) OPTIONSHandler(p string, handler http.Handler) {
	g.addRoute(http.MethodOptions, p, handler)
}

func (g *RouterGroup) TRACE(p string, handler http.HandlerFunc) {
	g.addRoute(http.MethodTrace, p, handler)
}
func (g *RouterGroup) TRACEHandler(p string, handler http.Handler) {
	g.addRoute(http.MethodTrace, p, handler)
}

func (g *RouterGroup) Any(p string, handler http.HandlerFunc) {
	g.addRoute("", p, handler)
}
func (g *RouterGroup) AnyHandler(p string, handler http.Handler) {
	g.addRoute("", p, handler)
}
