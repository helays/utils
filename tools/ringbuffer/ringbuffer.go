// Package ringbuffer 固定大小的滑动窗口
package ringbuffer

import "sync"

// RingBuffer 是一个线程安全的泛型环形缓冲区(循环队列)实现
// 使用读写锁(sync.RWMutex)保证并发安全，适用于多读少写的场景
// 当缓冲区满时，新元素会自动覆盖最旧的元素(FIFO淘汰策略)
// 类型参数T表示缓冲区中存储的元素类型
type RingBuffer[T any] struct {
	buffer []T          // 底层数组，存储实际元素
	size   int          // 缓冲区总容量(固定大小)
	head   int          // 指向队列头部(最旧元素)的索引
	tail   int          // 指向队列尾部(下一个写入位置)的索引
	count  int          // 当前缓冲区中的元素数量
	mu     sync.RWMutex // 读写锁，保证并发安全
}

// New 创建并返回一个新的RingBuffer实例
// 参数:
//   - size: 缓冲区容量，必须为正整数
//
// 返回值:
//   - *RingBuffer[T]: 初始化后的环形缓冲区指针
//
// 注意:
//   - 创建后缓冲区容量不可变
//   - 如果size<=0会导致panic
func New[T any](size int) *RingBuffer[T] {
	return &RingBuffer[T]{
		buffer: make([]T, size),
		size:   size,
	}
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
	if rb.count < rb.size {
		rb.count++ // 缓冲区未满，元素数量增加
	} else {
		rb.head = (rb.head + 1) % rb.size // 缓冲区已满，移动head指针
	}
}

// GetAll 获取缓冲区中所有元素的副本(按插入顺序从旧到新)
// 线程安全：使用读锁(sync.RLock)保证并发安全
// 返回值:
//   - []T: 包含所有元素的切片(可能为空)
//
// 性能:
//   - 需要分配新切片并复制元素，时间复杂度O(n)
func (rb *RingBuffer[T]) GetAll() []T {
	rb.mu.RLock()
	defer rb.mu.RUnlock()

	// 创建结果切片
	result := make([]T, rb.count)

	// 按顺序复制元素(考虑环形布局)
	for i := 0; i < rb.count; i++ {
		result[i] = rb.buffer[(rb.head+i)%rb.size]
	}
	return result
}

// Len 返回当前缓冲区中的元素数量
// 线程安全：使用读锁(sync.RLock)保证并发安全
// 返回值:
//   - int: 当前元素数量(0 <= count <= size)
func (rb *RingBuffer[T]) Len() int {
	rb.mu.RLock()
	defer rb.mu.RUnlock()
	return rb.count
}

// Clear 清空缓冲区(重置所有指针和计数器)
// 线程安全：使用互斥锁(sync.Mutex)保证并发安全
// 注意:
//   - 不会实际删除元素，只是重置索引
//   - 已存在的元素会被后续Push操作覆盖
func (rb *RingBuffer[T]) Clear() {
	rb.mu.Lock()
	defer rb.mu.Unlock()

	rb.head = 0
	rb.tail = 0
	rb.count = 0
}
