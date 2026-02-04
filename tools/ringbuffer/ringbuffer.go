// Package ringbuffer 提供线程安全的智能环形缓冲区
package ringbuffer

import (
	"errors"
	"github.com/helays/utils/v2/tools/ringbuffer/impl"
	"reflect"
	"sync"
	"sync/atomic"
)

// RingBuffer 线程安全的环形缓冲区实现
type RingBuffer[T any] struct {
	impl  impl.Buffer[T] // 底层实现
	size  int            // 容量
	mu    sync.RWMutex   // 读写锁
	count int32          // 原子计数器
}

// New 创建新的环形缓冲区
func New[T any](size int) (*RingBuffer[T], error) {
	if size <= 0 {
		return nil, errors.New("缓冲区大小必须为正整数")
	}

	rb := &RingBuffer[T]{
		size: size,
	}

	// 智能选择实现方式
	if shouldUseLinkedList[T](size) {
		rb.impl = impl.NewLinkedListBuffer[T](size)
	} else {
		rb.impl = impl.NewArrayBuffer[T](size)
	}

	return rb, nil
}

// shouldUseLinkedList 判断是否使用链表实现
func shouldUseLinkedList[T any](size int) bool {
	const (
		largeElementSize = 16 // 大元素阈值(字节)
		smallCapacity    = 32 // 小容量阈值
	)

	elementSize := estimateElementSize[T]()
	return elementSize > largeElementSize && size <= smallCapacity
}

// estimateElementSize 估算泛型类型大小
func estimateElementSize[T any]() int {
	var zero T
	t := reflect.TypeOf(zero)

	switch t.Kind() {
	case reflect.Bool, reflect.Int8, reflect.Uint8:
		return 1
	case reflect.Int16, reflect.Uint16:
		return 2
	case reflect.Int32, reflect.Uint32, reflect.Float32:
		return 4
	case reflect.Int64, reflect.Uint64, reflect.Float64, reflect.Complex64:
		return 8
	case reflect.Complex128:
		return 16
	case reflect.String:
		return 32
	case reflect.Struct:
		return int(t.Size())
	case reflect.Ptr, reflect.Slice, reflect.Map, reflect.Interface:
		return 16
	default:
		return 32
	}
}

// Push 添加单个元素
func (rb *RingBuffer[T]) Push(item T) {
	rb.mu.Lock()
	defer rb.mu.Unlock()

	old := rb.impl.Len()
	rb.impl.Push(item)
	if rb.impl.Len() != old {
		atomic.StoreInt32(&rb.count, int32(rb.impl.Len()))
	}
}

// PushAll 批量添加元素
func (rb *RingBuffer[T]) PushAll(items []T) {
	if len(items) == 0 {
		return
	}

	rb.mu.Lock()
	defer rb.mu.Unlock()

	old := rb.impl.Len()
	rb.impl.PushAll(items)
	if rb.impl.Len() != old {
		atomic.StoreInt32(&rb.count, int32(rb.impl.Len()))
	}
}

// GetAll 获取所有元素
func (rb *RingBuffer[T]) GetAll() []T {
	rb.mu.RLock()
	defer rb.mu.RUnlock()

	return rb.impl.GetAll()
}

// GetLast 获取最后n个元素
func (rb *RingBuffer[T]) GetLast(n int) []T {
	rb.mu.RLock()
	defer rb.mu.RUnlock()

	return rb.impl.GetLast(n)
}

// Len 当前元素数量
func (rb *RingBuffer[T]) Len() int {
	return int(atomic.LoadInt32(&rb.count))
}

// Cap 缓冲区容量
func (rb *RingBuffer[T]) Cap() int {
	return rb.size
}

// IsFull 是否已满
func (rb *RingBuffer[T]) IsFull() bool {
	return rb.Len() == rb.size
}

// IsEmpty 是否为空
func (rb *RingBuffer[T]) IsEmpty() bool {
	return rb.Len() == 0
}

// Clear 清空缓冲区
func (rb *RingBuffer[T]) Clear() {
	rb.mu.Lock()
	defer rb.mu.Unlock()

	rb.impl.Clear()
	atomic.StoreInt32(&rb.count, 0)
}

// Iterator 创建元素迭代器
func (rb *RingBuffer[T]) Iterator() <-chan T {
	ch := make(chan T, rb.Len())

	go func() {
		rb.mu.RLock()
		defer rb.mu.RUnlock()

		for item := range rb.impl.Iterator() {
			ch <- item
		}
		close(ch)
	}()

	return ch
}

// MustNew 创建环形缓冲区(失败时panic)
func MustNew[T any](size int) *RingBuffer[T] {
	rb, err := New[T](size)
	if err != nil {
		panic(err)
	}
	return rb
}
