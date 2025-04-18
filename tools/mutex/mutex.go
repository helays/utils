package mutex

import "sync"

// SafeResourceRWMutex 使用 RWMutex 的线程安全资源
type SafeResourceRWMutex[T any] struct {
	resource T
	mu       sync.RWMutex
}

// SafeResourceMutex 使用 Mutex 的线程安全资源
type SafeResourceMutex[T any] struct {
	resource T
	mu       sync.Mutex
}

// NewSafeResourceRWMutex 创建 RWMutex 版本的安全资源
func NewSafeResourceRWMutex[T any](resource T) *SafeResourceRWMutex[T] {
	return &SafeResourceRWMutex[T]{
		resource: resource,
	}
}

// NewSafeResourceMutex 创建 Mutex 版本的安全资源
func NewSafeResourceMutex[T any](resource T) *SafeResourceMutex[T] {
	return &SafeResourceMutex[T]{
		resource: resource,
	}
}

// ----- RWMutex 版本的方法 -----

func (sr *SafeResourceRWMutex[T]) Read() *T {
	sr.mu.RLock()
	defer sr.mu.RUnlock()
	return &sr.resource
}

func (sr *SafeResourceRWMutex[T]) ReadWith(fn func(T)) {
	sr.mu.RLock()
	defer sr.mu.RUnlock()
	fn(sr.resource)
}

func (sr *SafeResourceRWMutex[T]) Write(newResource T) {
	sr.mu.Lock()
	defer sr.mu.Unlock()
	sr.resource = newResource
}

func (sr *SafeResourceRWMutex[T]) Update(fn func(T) T) {
	sr.mu.Lock()
	defer sr.mu.Unlock()
	sr.resource = fn(sr.resource)
}

// ----- Mutex 版本的方法 -----

func (sr *SafeResourceMutex[T]) Read() *T {
	sr.mu.Lock()
	defer sr.mu.Unlock()
	return &sr.resource
}

func (sr *SafeResourceMutex[T]) ReadWith(fn func(T)) {
	sr.mu.Lock()
	defer sr.mu.Unlock()
	fn(sr.resource)
}

func (sr *SafeResourceMutex[T]) Write(newResource T) {
	sr.mu.Lock()
	defer sr.mu.Unlock()
	sr.resource = newResource
}

func (sr *SafeResourceMutex[T]) Update(fn func(T) T) {
	sr.mu.Lock()
	defer sr.mu.Unlock()
	sr.resource = fn(sr.resource)
}
