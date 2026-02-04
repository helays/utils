package routecache

import (
	"sync"

	"helay.net/go/utils/v3/net/http/route/middleware/routecache/tree"
)

type RouteCache[T comparable] struct {
	staticRoute map[string]T // 静态路由缓存
	radixTree   map[string]*tree.RadixTreeNode[T]
	mu          sync.RWMutex
}

func New[T comparable]() *RouteCache[T] {
	r := &RouteCache[T]{
		staticRoute: make(map[string]T),
		radixTree:   make(map[string]*tree.RadixTreeNode[T]),
	}
	return r
}
