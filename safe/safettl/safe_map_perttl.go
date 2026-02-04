package safettl

import (
	"strings"
	"sync"
	"time"

	"helay.net/go/utils/v3/tools"
)

// PerKeyTTLMap 是支持每个键单独 TTL 的泛型安全 Map
// Deprecated: 弃用，最新采用 safe.Map
type PerKeyTTLMap[K comparable, V any] struct {
	mu              sync.Map
	closeChan       chan struct{}
	cleanupInterval time.Duration // 定义清理间隔
}

// NewPerKeyTTLMap 创建一个支持单独 TTL 的安全 Map
// Deprecated: 弃用，最新采用 safe.Map
func NewPerKeyTTLMap[K comparable, V any]() *PerKeyTTLMap[K, V] {
	m := &PerKeyTTLMap[K, V]{
		closeChan:       make(chan struct{}),
		cleanupInterval: time.Second,
	}

	// 启动后台清理协程
	go m.cleanupExpired()
	return m
}

// NewPerKeyTTLMapWithInterval 创建一个支持单独 TTL 的安全 Map，可指定清理间隔
// Deprecated: 弃用，最新采用 safe.Map
func NewPerKeyTTLMapWithInterval[K comparable, V any](cleanupInterval time.Duration) *PerKeyTTLMap[K, V] {
	// 确保清理间隔在合理范围内
	cleanupInterval = tools.AutoTimeDuration(cleanupInterval, time.Second, time.Second) // 绝对最小1秒

	m := &PerKeyTTLMap[K, V]{
		closeChan:       make(chan struct{}),
		cleanupInterval: cleanupInterval,
	}

	// 启动后台清理协程
	go m.cleanupExpired()

	return m
}

// LoadAndDelete 删除键的值，返回之前的值（如果存在且未过期）
func (m *PerKeyTTLMap[K, V]) LoadAndDelete(key K) (value V, loaded bool) {
	if val, ok := m.mu.LoadAndDelete(key); ok {
		item := val.(*itemWithExpiry[V])

		// 检查是否过期
		if item.isExpired() {
			var zeroV V
			return zeroV, false
		}

		return item.value, true
	}
	var zeroV V
	return zeroV, false
}

// LoadAndDeleteIf 条件性加载并删除
func (m *PerKeyTTLMap[K, V]) LoadAndDeleteIf(key K, condition func(value V) bool) (value V, deleted bool) {
	if val, ok := m.mu.Load(key); ok {
		item := val.(*itemWithExpiry[V])

		// 检查是否过期
		if item.isExpired() {
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
func (m *PerKeyTTLMap[K, V]) LoadAndDeleteExpired(key K) (value V, expired bool) {
	if val, ok := m.mu.Load(key); ok {
		item := val.(*itemWithExpiry[V])

		// 检查是否过期
		if item.isExpired() {
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
func (m *PerKeyTTLMap[K, V]) DeleteAndGetCount(keys ...K) int {
	count := 0
	for _, key := range keys {
		if _, loaded := m.LoadAndDelete(key); loaded {
			count++
		}
	}
	return count
}

// Store 设置键的值，永不过期
func (m *PerKeyTTLMap[K, V]) Store(key K, value V) {
	item := &itemWithExpiry[V]{
		value:      value,
		expiryTime: time.Time{}, // 零值表示永不过期
		ttl:        0,
	}
	m.mu.Store(key, item)
}

// StoreWithTTL 使用自定义 TTL 设置键的值，ttl=0 表示永不过期
func (m *PerKeyTTLMap[K, V]) StoreWithTTL(key K, value V, ttl time.Duration) {
	var expiryTime time.Time
	if ttl > 0 {
		expiryTime = time.Now().Add(ttl)
	}
	// 如果 ttl <= 0，expiryTime 保持为零值，表示永不过期

	item := &itemWithExpiry[V]{
		value:      value,
		expiryTime: expiryTime,
		ttl:        ttl,
	}
	m.mu.Store(key, item)
}

// Load 返回存储在 map 中给定键的值（如果未过期）
func (m *PerKeyTTLMap[K, V]) Load(key K) (V, bool) {
	if val, ok := m.mu.Load(key); ok {
		item := val.(*itemWithExpiry[V])

		// 检查是否过期
		if item.isExpired() {
			m.mu.Delete(key)
			var zeroV V
			return zeroV, false
		}

		return item.value, true
	}
	var zeroV V
	return zeroV, false
}

// LoadWithExpiry 返回存储在 map 中给定键的值和过期时间（如果未过期）
func (m *PerKeyTTLMap[K, V]) LoadWithExpiry(key K) (V, time.Time, bool) {
	if val, ok := m.mu.Load(key); ok {
		item := val.(*itemWithExpiry[V])

		// 检查是否过期
		if item.isExpired() {
			m.mu.Delete(key)
			var zeroV V
			return zeroV, time.Time{}, false
		}

		return item.value, item.getExpiry(), true
	}
	var zeroV V
	return zeroV, time.Time{}, false
}

// LoadAndRefresh 加载值并刷新过期时间（只对有 TTL 的键有效）
func (m *PerKeyTTLMap[K, V]) LoadAndRefresh(key K, ttl time.Duration) (V, bool) {
	if val, ok := m.mu.Load(key); ok {
		item := val.(*itemWithExpiry[V])

		// 检查是否过期
		if item.isExpired() {
			m.mu.Delete(key)
			var zeroV V
			return zeroV, false
		}

		// 只有原键有 TTL 时才刷新（expiryTime 不是零值）
		if !item.getExpiry().IsZero() && ttl > 0 {
			item.setExpiry(ttl)
		}

		return item.value, true
	}
	var zeroV V
	return zeroV, false
}

// Delete 移除键的值
func (m *PerKeyTTLMap[K, V]) Delete(key K) {
	m.mu.Delete(key)
}

// Clear 清空所有键值对
func (m *PerKeyTTLMap[K, V]) Clear() {
	m.mu.Range(func(k, v interface{}) bool {
		m.mu.Delete(k)
		return true
	})
}

// DeletePrefix 删除具有指定前缀的键
func (m *PerKeyTTLMap[K, V]) DeletePrefix(prefix string) {
	m.mu.Range(func(k, v interface{}) bool {
		if key, ok := k.(string); ok && strings.HasPrefix(key, prefix) {
			m.mu.Delete(k)
		}
		return true
	})
}

// DeleteSuffix 删除具有指定后缀的键
func (m *PerKeyTTLMap[K, V]) DeleteSuffix(suffix string) {
	m.mu.Range(func(k, v interface{}) bool {
		if key, ok := k.(string); ok && strings.HasSuffix(key, suffix) {
			m.mu.Delete(k)
		}
		return true
	})
}

// Range 遍历所有未过期的键值对
func (m *PerKeyTTLMap[K, V]) Range(f func(key K, value V) bool) {
	m.mu.Range(func(k, v interface{}) bool {
		item := v.(*itemWithExpiry[V])

		// 跳过过期的项
		if item.isExpired() {
			m.mu.Delete(k)
			return true
		}

		return f(k.(K), item.value)
	})
}

// LoadOrStore 如果键存在且未过期则返回值，否则存储并返回值（永不过期）
func (m *PerKeyTTLMap[K, V]) LoadOrStore(key K, value V) (actual V, loaded bool) {
	return m.LoadOrStoreWithTTL(key, value, 0)
}

// LoadOrStoreWithTTL 如果键存在且未过期则返回值，否则存储并返回值（可设置 TTL）
func (m *PerKeyTTLMap[K, V]) LoadOrStoreWithTTL(key K, value V, ttl time.Duration) (actual V, loaded bool) {
	// 先尝试加载
	if v, ok := m.Load(key); ok {
		return v, true
	}

	// 不存在或已过期，存储新值
	var expiryTime time.Time
	if ttl > 0 {
		expiryTime = time.Now().Add(ttl)
	}

	item := &itemWithExpiry[V]{
		value:      value,
		expiryTime: expiryTime,
	}

	existing, loaded := m.mu.LoadOrStore(key, item)
	if loaded {
		// 其他协程在我们之前存储了值
		existingItem := existing.(*itemWithExpiry[V])

		// 检查是否过期
		if existingItem.isExpired() {
			// 已过期，用新值替换
			m.mu.Store(key, item)
			return value, false
		}

		return existingItem.value, true
	}

	return value, false
}

// GetTTL 获取键的剩余存活时间，返回剩余时间和是否存在（永不过期的键返回 0, true）
func (m *PerKeyTTLMap[K, V]) GetTTL(key K) (time.Duration, bool) {
	if val, ok := m.mu.Load(key); ok {
		item := val.(*itemWithExpiry[V])

		// 永不过期的键
		if item.getExpiry().IsZero() {
			return 0, true
		}

		remaining := time.Until(item.getExpiry())
		if remaining > 0 {
			return remaining, true
		}

		// 已过期，删除
		m.mu.Delete(key)
	}
	return 0, false
}

// Refresh 刷新键的过期时间（只对有 TTL 的键有效）
func (m *PerKeyTTLMap[K, V]) Refresh(key K, ttl time.Duration) bool {
	if val, ok := m.mu.Load(key); ok {
		item := val.(*itemWithExpiry[V])

		// 检查是否已过期
		if item.isExpired() {
			m.mu.Delete(key)
			return false
		}

		// 只有原键有 TTL 时才刷新，并且新的 ttl 必须 > 0
		if !item.getExpiry().IsZero() && ttl > 0 {
			item.setExpiry(ttl)
		}

		return true
	}
	return false
}

// cleanupExpired 定期清理过期的键值对
func (m *PerKeyTTLMap[K, V]) cleanupExpired() {
	ticker := time.NewTicker(time.Minute) // 默认每分钟清理一次
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

func (m *PerKeyTTLMap[K, V]) cleanupBatch(maxClean int) int {
	count := 0

	m.mu.Range(func(k, v interface{}) bool {
		item := v.(*itemWithExpiry[V])
		if item.isExpired() {
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
func (m *PerKeyTTLMap[K, V]) Close() {
	close(m.closeChan)
}

// Size 返回 map 中未过期项的数量
func (m *PerKeyTTLMap[K, V]) Size() int {
	count := 0

	m.mu.Range(func(k, v interface{}) bool {
		item := v.(*itemWithExpiry[V])
		if !item.isExpired() {
			count++
		} else {
			m.mu.Delete(k)
		}
		return true
	})

	return count
}

// IsExpired 检查键是否已过期
func (m *PerKeyTTLMap[K, V]) IsExpired(key K) bool {
	if val, ok := m.mu.Load(key); ok {
		item := val.(*itemWithExpiry[V])
		return item.isExpired()
	}
	return true
}

// HasTTL 检查键是否有 TTL（不是永不过期）
func (m *PerKeyTTLMap[K, V]) HasTTL(key K) bool {
	if val, ok := m.mu.Load(key); ok {
		item := val.(*itemWithExpiry[V])
		return !item.getExpiry().IsZero()
	}
	return false
}
