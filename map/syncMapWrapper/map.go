package syncMapWrapper

import (
	"sync"
)

// SyncMapWrapper 是 sync.Map 的泛型包装器。
type SyncMapWrapper[K comparable, V any] struct {
	mu sync.Map
}

// Load 返回存储在 map 中给定键的值。
func (m *SyncMapWrapper[K, V]) Load(key K) (V, bool) {
	if val, ok := m.mu.Load(key); ok {
		return val.(V), true // 类型断言不会失败
	}
	var zeroV V
	return zeroV, false
}

// Store 设置键的值。
func (m *SyncMapWrapper[K, V]) Store(key K, value V) {
	m.mu.Store(key, value)
}

// Delete 移除键的值。
func (m *SyncMapWrapper[K, V]) Delete(key K) {
	m.mu.Delete(key)
}

// Range 依次调用 f 函数，针对 map 中的每个键和值。
// 如果 f 返回 false，则停止迭代。
func (m *SyncMapWrapper[K, V]) Range(f func(key K, value V) bool) {
	m.mu.Range(func(k, v interface{}) bool {
		return f(k.(K), v.(V)) // 类型断言不会失败
	})
}

// LoadOrStore 返回存在的键的值（如果存在）。
// 否则，它存储并返回给定的值。
// loaded 结果为 true 表示找到了该值。
func (m *SyncMapWrapper[K, V]) LoadOrStore(key K, value V) (actual V, loaded bool) {
	existing, loaded := m.mu.LoadOrStore(key, value)
	if loaded {
		actual = existing.(V) // 类型断言不会失败
	} else {
		actual = value
	}
	return actual, loaded
}

// LoadAndDelete 删除键的值，返回之前的值（如果有）。
// loaded 结果为 true 表示键存在。
func (m *SyncMapWrapper[K, V]) LoadAndDelete(key K) (value V, loaded bool) {
	existing, loaded := m.mu.LoadAndDelete(key)
	if loaded {
		value = existing.(V) // 类型断言不会失败
	}
	return value, loaded
}

// CompareAndSwap 比较键的值与提供的旧值，
// 如果相等，则将其替换为新值。
// 它返回是否成功进行了替换。
func (m *SyncMapWrapper[K, V]) CompareAndSwap(key K, oldval, newval V) bool {
	return m.mu.CompareAndSwap(key, oldval, newval)
}

// CompareAndDelete 比较键的值与提供的旧值，
// 如果相等，则删除该项。
// 它返回是否成功进行了删除。
func (m *SyncMapWrapper[K, V]) CompareAndDelete(key K, oldval V) bool {
	return m.mu.CompareAndDelete(key, oldval)
}

// Swap 将键的值与提供的新值交换，返回旧值。
func (m *SyncMapWrapper[K, V]) Swap(key K, newval V) (oldval V, swapped bool) {
	oldIface, swapped := m.mu.Swap(key, newval)
	if swapped {
		oldval = oldIface.(V) // 类型断言不会失败
	}
	return oldval, swapped
}
