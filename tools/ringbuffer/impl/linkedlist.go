// Package impl 实现环形缓冲区的链表存储版本
package impl

// linkedNode 链表节点
type linkedNode[T any] struct {
	value T
	next  *linkedNode[T]
	prev  *linkedNode[T]
}

// linkedListBuffer 基于链表的环形缓冲区实现
type linkedListBuffer[T any] struct {
	head  *linkedNode[T]
	tail  *linkedNode[T]
	size  int
	count int
}

// NewLinkedListBuffer 创建链表实现的环形缓冲区
func NewLinkedListBuffer[T any](size int) Buffer[T] {
	return &linkedListBuffer[T]{
		size: size,
	}
}

func (l *linkedListBuffer[T]) Push(item T) {
	node := &linkedNode[T]{value: item}

	if l.head == nil {
		l.head = node
		l.tail = node
		l.count = 1
		return
	}

	node.prev = l.tail
	l.tail.next = node
	l.tail = node

	if l.count < l.size {
		l.count++
	} else {
		l.head = l.head.next
		l.head.prev = nil
	}
}

func (l *linkedListBuffer[T]) PushAll(items []T) {
	for _, item := range items {
		l.Push(item)
	}
}

func (l *linkedListBuffer[T]) GetAll() []T {
	result := make([]T, 0, l.count)
	for node := l.head; node != nil; node = node.next {
		result = append(result, node.value)
	}
	return result
}

func (l *linkedListBuffer[T]) GetLast(n int) []T {
	if n <= 0 || l.count == 0 {
		return nil
	}

	if n > l.count {
		n = l.count
	}

	result := make([]T, 0, n)
	for node := l.tail; n > 0 && node != nil; node = node.prev {
		result = append(result, node.value)
		n--
	}

	// 反转结果保持顺序
	for i, j := 0, len(result)-1; i < j; i, j = i+1, j-1 {
		result[i], result[j] = result[j], result[i]
	}

	return result
}

func (l *linkedListBuffer[T]) Clear() {
	l.head = nil
	l.tail = nil
	l.count = 0
}

func (l *linkedListBuffer[T]) Len() int {
	return l.count
}

func (l *linkedListBuffer[T]) Cap() int {
	return l.size
}

func (l *linkedListBuffer[T]) IsFull() bool {
	return l.count == l.size
}

func (l *linkedListBuffer[T]) IsEmpty() bool {
	return l.count == 0
}

func (l *linkedListBuffer[T]) Iterator() <-chan T {
	ch := make(chan T, l.count)

	go func() {
		for node := l.head; node != nil; node = node.next {
			ch <- node.value
		}
		close(ch)
	}()

	return ch
}
