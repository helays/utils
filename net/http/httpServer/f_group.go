package httpServer

import (
	"github.com/helays/utils/v2/net/http/httpServer/httpmethod"
	"golang.org/x/net/websocket"
	"net/http"
	"strings"
)

type RouterGroup struct {
	service     *HttpServer
	prefix      string           // 分组前缀
	middlewares []MiddlewareFunc // 中间件链
	parent      *RouterGroup     // 父级分组
}

func (h *HttpServer) Group(g string) *RouterGroup {
	return &RouterGroup{service: h, prefix: strings.TrimSpace(g)}
}

// Group 创建新分组
func (g *RouterGroup) Group(prefix string) *RouterGroup {
	return &RouterGroup{
		service:     g.service,
		prefix:      g.calculatePrefix(strings.TrimSpace(prefix)),
		parent:      g,
		middlewares: g.middlewares[:], // 继承父级中间件
	}
}

// 计算完整路径前缀
func (g *RouterGroup) calculatePrefix(p string) string {
	prefix := "/" + strings.Trim(g.prefix, "/")
	if p == "" {
		return prefix
	} else if prefix == "/" {
		prefix = ""
	}
	p = strings.TrimLeft(p, "/")
	return prefix + "/" + p
}

func (g *RouterGroup) addRoute(method, p string, handler http.Handler) {
	fullPath := g.calculatePrefix(strings.TrimSpace(p))
	var finalHandler = handler
	if len(g.middlewares) > 0 {
		finalHandler = Chain(g.middlewares...)(finalHandler)
	}
	var cb []MiddlewareFunc
	if method != "" {
		vm := func(next http.Handler) http.Handler {
			return httpmethod.Method(method, next)
		}
		cb = append(cb, vm)
	}
	g.service.middleware(fullPath, finalHandler, cb...) // 注册路由
}

// Ws 添加 WebSocket 支持
func (g *RouterGroup) Ws(p string, handler func(ws *websocket.Conn)) {
	fullPath := g.calculatePrefix(p)
	g.service.socketMiddleware(fullPath, handler)
}

// Use 支持链式调用
func (g *RouterGroup) Use(middleware ...MiddlewareFunc) *RouterGroup {
	g.middlewares = append(g.middlewares, middleware...)
	return g
}

func (g *RouterGroup) Get(p string, handler http.HandlerFunc) {
	g.addRoute(http.MethodGet, p, handler)
}
func (g *RouterGroup) GetHandler(p string, handler http.Handler) {
	g.addRoute(http.MethodGet, p, handler)
}

func (g *RouterGroup) Head(p string, handler http.HandlerFunc) {
	g.addRoute(http.MethodHead, p, handler)
}

func (g *RouterGroup) HeadHandler(p string, handler http.Handler) {
	g.addRoute(http.MethodPost, p, handler)
}

func (g *RouterGroup) Post(p string, handler http.HandlerFunc) {
	g.addRoute(http.MethodPost, p, handler)
}

func (g *RouterGroup) PostHandler(p string, handler http.Handler) {
	g.addRoute(http.MethodPost, p, handler)
}

func (g *RouterGroup) Put(p string, handler http.HandlerFunc) {
	g.addRoute(http.MethodPut, p, handler)
}
func (g *RouterGroup) PutHandler(p string, handler http.Handler) {
	g.addRoute(http.MethodPut, p, handler)
}

func (g *RouterGroup) Patch(p string, handler http.HandlerFunc) {
	g.addRoute(http.MethodPatch, p, handler)
}
func (g *RouterGroup) PatchHandler(p string, handler http.Handler) {
	g.addRoute(http.MethodPatch, p, handler)
}

func (g *RouterGroup) Delete(p string, handler http.HandlerFunc) {
	g.addRoute(http.MethodDelete, p, handler)
}
func (g *RouterGroup) DeleteHandler(p string, handler http.Handler) {
	g.addRoute(http.MethodDelete, p, handler)
}

func (g *RouterGroup) Connect(p string, handler http.HandlerFunc) {
	g.addRoute(http.MethodConnect, p, handler)
}
func (g *RouterGroup) ConnectHandler(p string, handler http.Handler) {
	g.addRoute(http.MethodConnect, p, handler)
}

func (g *RouterGroup) Options(p string, handler http.HandlerFunc) {
	g.addRoute(http.MethodOptions, p, handler)
}
func (g *RouterGroup) OptionsHandler(p string, handler http.Handler) {
	g.addRoute(http.MethodOptions, p, handler)
}

func (g *RouterGroup) Trace(p string, handler http.HandlerFunc) {
	g.addRoute(http.MethodTrace, p, handler)
}
func (g *RouterGroup) TraceHandler(p string, handler http.Handler) {
	g.addRoute(http.MethodTrace, p, handler)
}

func (g *RouterGroup) Any(p string, handler http.HandlerFunc) {
	g.addRoute("", p, handler)
}
func (g *RouterGroup) AnyHandler(p string, handler http.Handler) {
	g.addRoute("", p, handler)
}
