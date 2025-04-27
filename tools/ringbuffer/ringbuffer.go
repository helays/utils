// Package ringbuffer 固定大小的滑动窗口
package ringbuffer

import (
	"errors"
	"sync"
	"sync/atomic"
)

// RingBuffer 是一个线程安全的泛型环形缓冲区(循环队列)实现
// 使用读写锁(sync.RWMutex)保证并发安全，适用于多读少写的场景
// 当缓冲区满时，新元素会自动覆盖最旧的元素(FIFO淘汰策略)
// 类型参数T表示缓冲区中存储的元素类型
type RingBuffer[T any] struct {
	buffer []T          // 底层数组，存储实际元素
	size   int          // 缓冲区总容量(固定大小)
	head   int          // 指向队列头部(最旧元素)的索引
	tail   int          // 指向队列尾部(下一个写入位置)的索引
	count  int32        // 当前缓冲区中的元素数量(使用原子操作)
	mu     sync.RWMutex // 读写锁，保证并发安全
}

// New 创建并返回一个新的RingBuffer实例
// 参数:
//   - size: 缓冲区容量，必须为正整数
//
// 返回值:
//   - *RingBuffer[T]: 初始化后的环形缓冲区指针
//   - error: 如果size<=0返回错误
func New[T any](size int) (*RingBuffer[T], error) {
	if size <= 0 {
		return nil, errors.New("size must be positive")
	}
	return &RingBuffer[T]{
		buffer: make([]T, size),
		size:   size,
	}, nil
}

// MustNew 创建并返回一个新的RingBuffer实例，如果size<=0会panic
// 参数:
//   - size: 缓冲区容量，必须为正整数
//
// 返回值:
//   - *RingBuffer[T]: 初始化后的环形缓冲区指针
func MustNew[T any](size int) *RingBuffer[T] {
	rb, err := New[T](size)
	if err != nil {
		panic(err)
	}
	return rb
}

// Push 向缓冲区添加一个元素
// 如果缓冲区已满，最旧的元素将被新元素覆盖
// 线程安全：使用互斥锁(sync.Mutex)保证并发安全
func (rb *RingBuffer[T]) Push(item T) {
	rb.mu.Lock()
	defer rb.mu.Unlock()

	// 写入元素到tail位置(下一个写入位置)
	rb.buffer[rb.tail] = item

	// 移动tail指针(循环)
	rb.tail = (rb.tail + 1) % rb.size

	// 更新计数器
	if rb.count < int32(rb.size) {
		atomic.AddInt32(&rb.count, 1) // 缓冲区未满，元素数量增加
	} else {
		rb.head = (rb.head + 1) % rb.size // 缓冲区已满，移动head指针
	}
}

// PushAll 向缓冲区批量添加元素
// 线程安全：使用互斥锁(sync.Mutex)保证并发安全
func (rb *RingBuffer[T]) PushAll(items []T) {
	if len(items) == 0 {
		return
	}

	rb.mu.Lock()
	defer rb.mu.Unlock()

	for _, item := range items {
		rb.buffer[rb.tail] = item
		rb.tail = (rb.tail + 1) % rb.size

		if rb.count < int32(rb.size) {
			atomic.AddInt32(&rb.count, 1)
		} else {
			rb.head = (rb.head + 1) % rb.size
		}
	}
}

// GetAll 获取缓冲区中所有元素的副本(按插入顺序从旧到新)
// 线程安全：使用读锁(sync.RLock)保证并发安全
// 返回值:
//   - []T: 包含所有元素的切片(可能为空)
func (rb *RingBuffer[T]) GetAll() []T {
	rb.mu.RLock()
	defer rb.mu.RUnlock()

	result := make([]T, rb.count)
	for i := 0; i < int(rb.count); i++ {
		result[i] = rb.buffer[(rb.head+i)%rb.size]
	}
	return result
}

// GetLast 获取缓冲区中最后n个元素(按插入顺序从旧到新)
// 线程安全：使用读锁(sync.RLock)保证并发安全
// 返回值:
//   - []T: 包含元素的切片，长度不超过n
func (rb *RingBuffer[T]) GetLast(n int) []T {
	rb.mu.RLock()
	defer rb.mu.RUnlock()

	if n <= 0 {
		return nil
	}

	count := int(rb.count)
	if n > count {
		n = count
	}

	result := make([]T, n)
	start := (rb.tail - n + rb.size) % rb.size
	for i := 0; i < n; i++ {
		result[i] = rb.buffer[(start+i)%rb.size]
	}
	return result
}

// Iterator 返回一个只读通道，用于遍历缓冲区中的元素
// 注意: 必须在消费完所有元素前完成遍历，否则会导致goroutine泄漏
func (rb *RingBuffer[T]) Iterator() <-chan T {
	ch := make(chan T)
	go func() {
		rb.mu.RLock()
		defer rb.mu.RUnlock()
		defer close(ch)

		for i := 0; i < int(rb.count); i++ {
			ch <- rb.buffer[(rb.head+i)%rb.size]
		}
	}()
	return ch
}

// Len 返回当前缓冲区中的元素数量
// 线程安全：使用原子操作保证并发安全
func (rb *RingBuffer[T]) Len() int {
	return int(atomic.LoadInt32(&rb.count))
}

// Cap 返回缓冲区的总容量
func (rb *RingBuffer[T]) Cap() int {
	return rb.size
}

// IsFull 检查缓冲区是否已满
func (rb *RingBuffer[T]) IsFull() bool {
	return rb.Len() == rb.size
}

// IsEmpty 检查缓冲区是否为空
func (rb *RingBuffer[T]) IsEmpty() bool {
	return rb.Len() == 0
}

// Clear 清空缓冲区(重置所有指针和计数器)
// 线程安全：使用互斥锁(sync.Mutex)保证并发安全
func (rb *RingBuffer[T]) Clear() {
	rb.mu.Lock()
	defer rb.mu.Unlock()

	rb.head = 0
	rb.tail = 0
	atomic.StoreInt32(&rb.count, 0)
}

// RingBufferPool 对象池，用于管理RingBuffer实例
type RingBufferPool[T any] struct {
	pool sync.Pool
}

// NewPool 创建一个新的RingBuffer对象池
func NewPool[T any](size int) *RingBufferPool[T] {
	return &RingBufferPool[T]{
		pool: sync.Pool{
			New: func() any {
				rb, _ := New[T](size)
				return rb
			},
		},
	}
}

// Get 从对象池获取一个RingBuffer实例
func (p *RingBufferPool[T]) Get() *RingBuffer[T] {
	return p.pool.Get().(*RingBuffer[T])
}

// Put 将RingBuffer实例放回对象池
func (p *RingBufferPool[T]) Put(rb *RingBuffer[T]) {
	rb.Clear()
	p.pool.Put(rb)
}
