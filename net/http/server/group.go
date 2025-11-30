package server

import (
	"net/http"
	"strings"

	"golang.org/x/net/websocket"
)

type Group[T any] struct {
	serv        *Server[T]
	prefix      string
	middlewares []Middleware
	parent      *Group[T] // 父级分组
}

func (s *Server[T]) Group(groupName string) *Group[T] {
	return &Group[T]{
		serv:        s,
		prefix:      strings.TrimSpace(groupName),
		middlewares: make([]Middleware, 0),
	}
}

// Group 创建子分组
func (g *Group[T]) Group(groupName string) *Group[T] {
	return &Group[T]{
		serv:        g.serv,
		prefix:      g.calculatePrefix(strings.TrimSpace(groupName)),
		parent:      g,
		middlewares: g.middlewares[:],
	}
}

// 计算完整路径前缀
func (g *Group[T]) calculatePrefix(p string) string {
	prefix := "/" + strings.Trim(g.prefix, "/")
	if p == "" {
		return prefix
	} else if prefix == "/" {
		prefix = ""
	}
	p = strings.TrimLeft(p, "/")
	return prefix + "/" + p
}

func (g *Group[T]) addRoute(method, path string, handle http.Handler, descriptions ...Description[T]) {
	fullPath := g.calculatePrefix(strings.TrimSpace(path))
	if len(descriptions) > 0 {
		g.serv.AddRouteWithDescription(method, fullPath, handle, descriptions[0], g.middlewares...)
	} else {
		g.serv.AddRoute(method, fullPath, handle, g.middlewares...)
	}

}

func (g *Group[T]) Use(middleware ...Middleware) *Group[T] {
	g.middlewares = append(g.middlewares, middleware...)
	return g
}

func (g *Group[T]) Get(p string, handler http.HandlerFunc, descriptions ...Description[T]) {
	g.addRoute(http.MethodGet, p, handler, descriptions...)
}

func (g *Group[T]) GetWithHandler(p string, handler http.Handler, descriptions ...Description[T]) {
	g.addRoute(http.MethodGet, p, handler, descriptions...)
}

func (g *Group[T]) Head(p string, handler http.HandlerFunc, descriptions ...Description[T]) {
	g.addRoute(http.MethodHead, p, handler, descriptions...)
}

func (g *Group[T]) HeadWithHandler(p string, handler http.Handler, descriptions ...Description[T]) {
	g.addRoute(http.MethodHead, p, handler, descriptions...)
}

func (g *Group[T]) Post(p string, handler http.HandlerFunc, descriptions ...Description[T]) {
	g.addRoute(http.MethodPost, p, handler, descriptions...)
}

func (g *Group[T]) PostWithHandler(p string, handler http.Handler, descriptions ...Description[T]) {
	g.addRoute(http.MethodPost, p, handler, descriptions...)
}

func (g *Group[T]) Put(p string, handler http.HandlerFunc, descriptions ...Description[T]) {
	g.addRoute(http.MethodPut, p, handler, descriptions...)
}

func (g *Group[T]) PutWithHandler(p string, handler http.Handler, descriptions ...Description[T]) {
	g.addRoute(http.MethodPut, p, handler, descriptions...)
}

func (g *Group[T]) Patch(p string, handler http.HandlerFunc, descriptions ...Description[T]) {
	g.addRoute(http.MethodPatch, p, handler, descriptions...)
}

func (g *Group[T]) PatchWithHandler(p string, handler http.Handler, descriptions ...Description[T]) {
	g.addRoute(http.MethodPatch, p, handler, descriptions...)
}

func (g *Group[T]) Delete(p string, handler http.HandlerFunc, descriptions ...Description[T]) {
	g.addRoute(http.MethodDelete, p, handler, descriptions...)
}

func (g *Group[T]) DeleteWithHandler(p string, handler http.Handler, descriptions ...Description[T]) {
	g.addRoute(http.MethodDelete, p, handler, descriptions...)
}

func (g *Group[T]) Connect(p string, handler http.HandlerFunc, descriptions ...Description[T]) {
	g.addRoute(http.MethodConnect, p, handler, descriptions...)
}

func (g *Group[T]) ConnectWithHandler(p string, handler http.Handler, descriptions ...Description[T]) {
	g.addRoute(http.MethodConnect, p, handler, descriptions...)
}

func (g *Group[T]) Options(p string, handler http.HandlerFunc, descriptions ...Description[T]) {
	g.addRoute(http.MethodOptions, p, handler, descriptions...)
}

func (g *Group[T]) OptionsWithHandler(p string, handler http.Handler, descriptions ...Description[T]) {
	g.addRoute(http.MethodOptions, p, handler, descriptions...)
}

func (g *Group[T]) Trace(p string, handler http.HandlerFunc, descriptions ...Description[T]) {
	g.addRoute(http.MethodTrace, p, handler, descriptions...)
}

func (g *Group[T]) TraceWithHandler(p string, handler http.Handler, descriptions ...Description[T]) {
	g.addRoute(http.MethodTrace, p, handler, descriptions...)
}

func (g *Group[T]) Any(p string, handler http.HandlerFunc, descriptions ...Description[T]) {
	g.addRoute("", p, handler, descriptions...)
}

func (g *Group[T]) AnyWithHandler(p string, handler http.Handler, descriptions ...Description[T]) {
	g.addRoute("", p, handler, descriptions...)
}

func (g *Group[T]) WS(p string, handler func(ws *websocket.Conn)) {
	fullPath := g.calculatePrefix(strings.TrimSpace(p))
	g.serv.AddWebsocketRoute(fullPath, handler)
}

func (g *Group[T]) WSWithDescription(p string, handler func(ws *websocket.Conn), description Description[T]) {
	fullPath := g.calculatePrefix(strings.TrimSpace(p))
	g.serv.AddWebsocketRouteWithDescription(fullPath, handler, description)
}
