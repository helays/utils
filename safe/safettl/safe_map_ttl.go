package safettl

import (
	"strings"
	"sync"
	"time"
)

// Map 是带有超时功能的泛型安全 Map
// Deprecated: 弃用，最新采用 safe.Map
type Map[K comparable, V any] struct {
	mu        sync.Map
	ttl       time.Duration
	closeChan chan struct{}
}

// New 创建一个带有 TTL 的安全 Map
// Deprecated: 弃用，最新采用 safe.Map
func New[K comparable, V any](ttl time.Duration) *Map[K, V] {
	m := &Map[K, V]{
		ttl:       ttl,
		closeChan: make(chan struct{}),
	}

	// 启动后台清理协程
	go m.cleanupExpired()

	return m
}

// LoadAndDelete 删除键的值，返回之前的值（如果存在且未过期）
func (m *Map[K, V]) LoadAndDelete(key K) (value V, loaded bool) {
	if val, ok := m.mu.LoadAndDelete(key); ok {
		item := val.(*itemWithExpiry[V])

		// 检查是否过期
		if time.Now().After(item.getExpiry()) {
			var zeroV V
			return zeroV, false
		}

		return item.value, true
	}
	var zeroV V
	return zeroV, false
}

// LoadAndDeleteIf 条件性加载并删除
func (m *Map[K, V]) LoadAndDeleteIf(key K, condition func(value V) bool) (value V, deleted bool) {
	if val, ok := m.mu.Load(key); ok {
		item := val.(*itemWithExpiry[V])

		// 检查是否过期
		if time.Now().After(item.getExpiry()) {
			m.mu.Delete(key)
			var zeroV V
			return zeroV, false
		}

		// 检查条件
		if condition == nil || condition(item.value) {
			if m.mu.CompareAndDelete(key, val) {
				return item.value, true
			}
			// 如果 CAS 失败，说明值已被其他 goroutine 修改
			var zeroV V
			return zeroV, false
		}

		return item.value, false
	}
	var zeroV V
	return zeroV, false
}

// LoadAndDeleteExpired 加载并删除已过期的键（用于手动清理）
func (m *Map[K, V]) LoadAndDeleteExpired(key K) (value V, expired bool) {
	if val, ok := m.mu.Load(key); ok {
		item := val.(*itemWithExpiry[V])

		// 检查是否过期
		if time.Now().After(item.getExpiry()) {
			if m.mu.CompareAndDelete(key, val) {
				return item.value, true
			}
			// 如果 CAS 失败，说明值已被其他 goroutine 修改或删除
			var zeroV V
			return zeroV, false
		}

		// 未过期，不删除
		var zeroV V
		return zeroV, false
	}
	var zeroV V
	return zeroV, false
}

// DeleteAndGetCount 删除多个键并返回删除的数量
func (m *Map[K, V]) DeleteAndGetCount(keys ...K) int {
	count := 0
	for _, key := range keys {
		if _, loaded := m.LoadAndDelete(key); loaded {
			count++
		}
	}
	return count
}

// 其他方法保持不变...
// Store, StoreWithTTL, Load, LoadAndRefresh, Delete, Clear,
// DeletePrefix, DeleteSuffix, Range, LoadOrStore, GetTTL,
// Refresh, Close, Size, IsExpired 等方法

// Store 设置键的值，并设置过期时间
func (m *Map[K, V]) Store(key K, value V) {
	item := &itemWithExpiry[V]{
		value:      value,
		expiryTime: time.Now().Add(m.ttl),
	}
	m.mu.Store(key, item)
}

// StoreWithTTL 使用自定义 TTL 设置键的值
func (m *Map[K, V]) StoreWithTTL(key K, value V, ttl time.Duration) {
	item := &itemWithExpiry[V]{
		value:      value,
		expiryTime: time.Now().Add(ttl),
	}
	m.mu.Store(key, item)
}

// Load 返回存储在 map 中给定键的值（如果未过期）
func (m *Map[K, V]) Load(key K) (V, bool) {
	if val, ok := m.mu.Load(key); ok {
		item := val.(*itemWithExpiry[V])

		// 检查是否过期
		if time.Now().After(item.getExpiry()) {
			m.mu.Delete(key)
			var zeroV V
			return zeroV, false
		}

		return item.value, true
	}
	var zeroV V
	return zeroV, false
}

// LoadAndRefresh 加载值并刷新过期时间
func (m *Map[K, V]) LoadAndRefresh(key K) (V, bool) {
	if val, ok := m.mu.Load(key); ok {
		item := val.(*itemWithExpiry[V])

		// 检查是否过期
		if time.Now().After(item.getExpiry()) {
			m.mu.Delete(key)
			var zeroV V
			return zeroV, false
		}

		// 刷新过期时间
		item.setExpiry(m.ttl)
		return item.value, true
	}
	var zeroV V
	return zeroV, false
}

// Delete 移除键的值
func (m *Map[K, V]) Delete(key K) {
	m.mu.Delete(key)
}

// Clear 清空所有键值对
func (m *Map[K, V]) Clear() {
	m.mu.Range(func(k, v interface{}) bool {
		m.mu.Delete(k)
		return true
	})
}

// DeletePrefix 删除具有指定前缀的键
func (m *Map[K, V]) DeletePrefix(prefix string) {
	m.mu.Range(func(k, v interface{}) bool {
		if key, ok := k.(string); ok && strings.HasPrefix(key, prefix) {
			m.mu.Delete(k)
		}
		return true
	})
}

// DeleteSuffix 删除具有指定后缀的键
func (m *Map[K, V]) DeleteSuffix(suffix string) {
	m.mu.Range(func(k, v interface{}) bool {
		if key, ok := k.(string); ok && strings.HasSuffix(key, suffix) {
			m.mu.Delete(k)
		}
		return true
	})
}

// Range 遍历所有未过期的键值对
func (m *Map[K, V]) Range(f func(key K, value V) bool) {
	now := time.Now()
	m.mu.Range(func(k, v interface{}) bool {
		item := v.(*itemWithExpiry[V])

		// 跳过过期的项
		if now.After(item.getExpiry()) {
			m.mu.Delete(k)
			return true
		}

		return f(k.(K), item.value)
	})
}

// LoadOrStore 如果键存在且未过期则返回值，否则存储并返回值
func (m *Map[K, V]) LoadOrStore(key K, value V) (actual V, loaded bool) {
	// 先尝试加载
	if v, ok := m.Load(key); ok {
		return v, true
	}

	// 不存在或已过期，存储新值
	item := &itemWithExpiry[V]{
		value:      value,
		expiryTime: time.Now().Add(m.ttl),
	}

	existing, loaded := m.mu.LoadOrStore(key, item)
	if loaded {
		// 其他协程在我们之前存储了值
		existingItem := existing.(*itemWithExpiry[V])

		// 检查是否过期
		if time.Now().After(existingItem.getExpiry()) {
			// 已过期，用新值替换
			m.mu.Store(key, item)
			return value, false
		}

		return existingItem.value, true
	}

	return value, false
}

// GetTTL 获取键的剩余存活时间
func (m *Map[K, V]) GetTTL(key K) (time.Duration, bool) {
	if val, ok := m.mu.Load(key); ok {
		item := val.(*itemWithExpiry[V])
		remaining := time.Until(item.getExpiry())

		if remaining > 0 {
			return remaining, true
		}

		// 已过期，删除
		m.mu.Delete(key)
	}
	return 0, false
}

// Refresh 刷新键的过期时间
func (m *Map[K, V]) Refresh(key K) bool {
	if val, ok := m.mu.Load(key); ok {
		item := val.(*itemWithExpiry[V])

		// 检查是否已过期
		if time.Now().After(item.getExpiry()) {
			m.mu.Delete(key)
			return false
		}

		// 刷新过期时间
		item.setExpiry(m.ttl)
		return true
	}
	return false
}

// cleanupExpired 定期清理过期的键值对
func (m *Map[K, V]) cleanupExpired() {
	// 计算合理的清理间隔
	cleanupInterval := m.calculateCleanupInterval()
	ticker := time.NewTicker(cleanupInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			m.cleanupBatch(1000) // 批量清理，避免阻塞
		case <-m.closeChan:
			return
		}
	}
}

func (m *Map[K, V]) calculateCleanupInterval() time.Duration {
	// 基本原则：清理间隔 = TTL / 4
	interval := m.ttl / 4

	// 动态边界：确保清理间隔既不太频繁也不太稀疏
	minInterval := m.ttl / 10 // 最小为 TTL 的 1/10
	if minInterval < time.Second {
		minInterval = time.Second // 绝对最小1秒
	}

	maxInterval := m.ttl / 2 // 最大为 TTL 的 1/2
	if maxInterval > 10*time.Minute {
		maxInterval = 10 * time.Minute // 绝对最大10分钟
	}

	// 确保在合理范围内
	if interval < minInterval {
		return minInterval
	}
	if interval > maxInterval {
		return maxInterval
	}
	return interval
}

func (m *Map[K, V]) cleanupBatch(maxClean int) int {
	now := time.Now()
	count := 0

	m.mu.Range(func(k, v interface{}) bool {
		item := v.(*itemWithExpiry[V])
		if now.After(item.getExpiry()) {
			m.mu.Delete(k)
			count++
			if count >= maxClean {
				return false // 停止迭代
			}
		}
		return true
	})

	return count
}

// Close 停止后台清理协程，释放资源
func (m *Map[K, V]) Close() {
	close(m.closeChan)
}

// Size 返回 map 中未过期项的数量
func (m *Map[K, V]) Size() int {
	count := 0
	now := time.Now()

	m.mu.Range(func(k, v interface{}) bool {
		item := v.(*itemWithExpiry[V])
		if now.Before(item.getExpiry()) {
			count++
		} else {
			m.mu.Delete(k)
		}
		return true
	})

	return count
}

// IsExpired 检查键是否已过期
func (m *Map[K, V]) IsExpired(key K) bool {
	if val, ok := m.mu.Load(key); ok {
		item := val.(*itemWithExpiry[V])
		return time.Now().After(item.getExpiry())
	}
	return true
}
