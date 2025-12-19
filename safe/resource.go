package safe

import (
	"sync/atomic"
)

type Resource[T any] struct {
	_        noCopy
	resource atomic.Pointer[T]
}

// NewResource 创建一个原子操作的资源
// 注意这个资源不是线程安全的
// 虽然原子操作不会出现数据竞争，但是数据读取出来后的比较不是线程安全的。可能真实值已经被其他线程修改了
func NewResource[T any](resource *T) *Resource[T] {
	s := &Resource[T]{}
	if resource == nil {
		var zero T
		s.resource.Store(&zero)
	} else {
		s.resource.Store(resource)
	}
	return s
}

func (s *Resource[T]) Read() *T {
	return s.resource.Load()
}

func (s *Resource[T]) Write(newResource *T) {
	s.resource.Store(newResource)
}
