// Package impl 实现环形缓冲区的数组存储版本
package impl

import (
	"github.com/helays/utils/tools/ringbuffer/calculator"
)

// arrayBuffer 基于数组的环形缓冲区实现
type arrayBuffer[T any] struct {
	buffer []T
	size   int
	head   int
	tail   int
	count  int
	calc   calculator.IndexCalculator
}

// NewArrayBuffer 创建数组实现的环形缓冲区
func NewArrayBuffer[T any](size int) Buffer[T] {
	return &arrayBuffer[T]{
		buffer: make([]T, size),
		size:   size,
		calc:   calculator.NewSmartCalculator(size),
	}
}

func (a *arrayBuffer[T]) Push(item T) {
	a.buffer[a.tail] = item
	a.tail = a.calc.Next(a.tail)

	if a.count < a.size {
		a.count++
	} else {
		a.head = a.calc.Next(a.head)
	}
}

func (a *arrayBuffer[T]) PushAll(items []T) {
	if len(items) == 0 {
		return
	}

	remaining := a.size - a.tail
	if len(items) <= remaining {
		copy(a.buffer[a.tail:], items)
	} else {
		copy(a.buffer[a.tail:], items[:remaining])
		copy(a.buffer, items[remaining:])
	}

	newCount := a.count + len(items)
	if newCount <= a.size {
		a.count = newCount
	} else {
		a.count = a.size
		a.head = a.calc.Sub(a.head, newCount-a.size)
	}

	a.tail = (a.tail + len(items)) % a.size
}

func (a *arrayBuffer[T]) GetAll() []T {
	if a.count == 0 {
		return make([]T, 0)
	}

	result := make([]T, a.count)
	if a.head < a.tail || (a.head+a.count) <= a.size {
		copy(result, a.buffer[a.head:a.head+a.count])
	} else {
		firstPart := a.size - a.head
		copy(result, a.buffer[a.head:])
		copy(result[firstPart:], a.buffer[:a.count-firstPart])
	}
	return result
}

func (a *arrayBuffer[T]) GetLast(n int) []T {
	if n <= 0 {
		return nil
	}

	if n > a.count {
		n = a.count
	}

	result := make([]T, n)
	start := a.calc.Sub(a.tail, n)

	if start+n <= a.size {
		copy(result, a.buffer[start:start+n])
	} else {
		firstPart := a.size - start
		copy(result, a.buffer[start:])
		copy(result[firstPart:], a.buffer[:n-firstPart])
	}
	return result
}

func (a *arrayBuffer[T]) Clear() {
	a.head = 0
	a.tail = 0
	a.count = 0
}

func (a *arrayBuffer[T]) Len() int {
	return a.count
}

func (a *arrayBuffer[T]) Cap() int {
	return a.size
}

func (a *arrayBuffer[T]) IsFull() bool {
	return a.count == a.size
}

func (a *arrayBuffer[T]) IsEmpty() bool {
	return a.count == 0
}

func (a *arrayBuffer[T]) Iterator() <-chan T {
	ch := make(chan T, a.count)

	go func() {
		items := a.GetAll()
		for _, item := range items {
			ch <- item
		}
		close(ch)
	}()

	return ch
}
