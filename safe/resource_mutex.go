package safe

import "sync"

// ResourceRWMutex 使用 RWMutex 的线程安全资源
type ResourceRWMutex[T any] struct {
	resource T
	mu       sync.RWMutex
}

// NewResourceRWMutex 创建 RWMutex 版本的安全资源
func NewResourceRWMutex[T any](resource T) *ResourceRWMutex[T] {
	return &ResourceRWMutex[T]{
		resource: resource,
	}
}

// ----- RWMutex 版本的方法 -----

func (sr *ResourceRWMutex[T]) Read() T {
	sr.mu.RLock()
	defer sr.mu.RUnlock()
	return sr.resource
}

func (sr *ResourceRWMutex[T]) ReadWith(fn func(T)) {
	sr.mu.RLock()
	defer sr.mu.RUnlock()
	fn(sr.resource)
}

func (sr *ResourceRWMutex[T]) ReadWithResult(fn func(T) error) error {
	sr.mu.RLock()
	defer sr.mu.RUnlock()
	return fn(sr.resource)
}

func (sr *ResourceRWMutex[T]) Write(newResource T) {
	sr.mu.Lock()
	defer sr.mu.Unlock()
	sr.resource = newResource
}

func (sr *ResourceRWMutex[T]) Update(fn func(T) T) {
	sr.mu.Lock()
	defer sr.mu.Unlock()
	sr.resource = fn(sr.resource)
}

// ----- Mutex 版本的方法 -----

// ResourceMutex 使用 Mutex 的线程安全资源
type ResourceMutex[T any] struct {
	resource T
	mu       sync.Mutex
}

// NewResourceMutex 创建 Mutex 版本的安全资源
func NewResourceMutex[T any](resource T) *ResourceMutex[T] {
	return &ResourceMutex[T]{
		resource: resource,
	}
}
func (sr *ResourceMutex[T]) Read() T {
	sr.mu.Lock()
	defer sr.mu.Unlock()
	return sr.resource
}

func (sr *ResourceMutex[T]) ReadWith(fn func(T)) {
	sr.mu.Lock()
	defer sr.mu.Unlock()
	fn(sr.resource)
}

func (sr *ResourceMutex[T]) Write(newResource T) {
	sr.mu.Lock()
	defer sr.mu.Unlock()
	sr.resource = newResource
}

func (sr *ResourceMutex[T]) Update(fn func(T) T) {
	sr.mu.Lock()
	defer sr.mu.Unlock()
	sr.resource = fn(sr.resource)
}
