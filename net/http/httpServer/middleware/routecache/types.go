package routecache

import (
	"fmt"

	"github.com/go-chi/chi/v5"
	"github.com/helays/utils/v2/net/http/httpServer/router"
)

type RouteCache[T any] struct {
	staticRoute map[string]T // 静态路由缓存
	routeMap    map[string]T // 路由和配置的对应关系
	router      *chi.Mux     // 动态路由匹配器
}

func New[T any]() *RouteCache[T] {
	r := &RouteCache[T]{
		staticRoute: make(map[string]T),
		routeMap:    make(map[string]T),
		router:      chi.NewRouter(),
	}
	return r
}

// AddRoute 添加路由
func (r *RouteCache[T]) AddRoute(t router.RouteType, method string, pattern string, handle T) {
	path := method + pattern
	if t == router.RouteTypeStatic {
		r.staticRoute[path] = handle
		return
	}

	r.routeMap[path] = handle
	r.router.MethodFunc(method, pattern, nil)
}

func (r *RouteCache[T]) Match(method, path string) (T, bool) {
	if h, ok := r.staticRoute[method+path]; ok {
		return h, true
	}
	ctx := chi.NewRouteContext()
	if r.router.Match(ctx, method, path) {
		pattern := ctx.RoutePattern()
		fmt.Println("pattern:", pattern)
		h, ok := r.routeMap[method+pattern]
		return h, ok
	}
	var zero T
	return zero, false
}
