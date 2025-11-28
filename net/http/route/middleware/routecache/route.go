package routecache

import (
	"github.com/helays/utils/v2/net/http/route"
	"github.com/helays/utils/v2/net/http/route/middleware/routecache/tree"
)

// AddRoute 添加路由
func (r *RouteCache[T]) AddRoute(t route.RouteType, method string, pattern string, handle T) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if t == route.Static {
		path := method + pattern
		r.staticRoute[path] = handle
		return nil
	}
	if r.radixTree[method] == nil {
		r.radixTree[method] = tree.NewRadixTree[T]()
	}
	return r.radixTree[method].AddRoute(pattern, handle)
}

func (r *RouteCache[T]) Match(method, path string) (T, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	if h, ok := r.staticRoute[method+path]; ok {
		return h, true
	}
	var zero T
	radixTree, ok := r.radixTree[method]
	if !ok {
		return zero, false
	}
	h, _, _, err := radixTree.GetValue(path)
	if err != nil {
		return zero, false
	}
	return h, true
}

// Clear 清空缓存
func (r *RouteCache[T]) Clear() {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.staticRoute = make(map[string]T)
	r.radixTree = make(map[string]*tree.RadixTreeNode[T])
}
