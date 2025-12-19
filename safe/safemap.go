package safe

import (
	"strings"
	"sync"

	"github.com/helays/utils/v2/tools"
)

type (
	value[K comparable, V any] struct {
		key K
		val V
	}
	shard[K comparable, V any] struct {
		items map[uint64]*value[K, V]
		mu    sync.RWMutex
	}
	Map[K comparable, V any] struct {
		shards    []*shard[K, V]
		shardSize uint64
		hasher    Hasher[K]
	}
)

func NewMap[K comparable, V any](hasher Hasher[K], shardSize ...uint64) *Map[K, V] {
	m := &Map[K, V]{
		hasher: hasher,
	}
	if len(shardSize) > 0 {
		if !tools.IsPowerOfTwo(shardSize[0]) {
			panic("分片数量必须是2的N次幂")
		}
		m.shardSize = shardSize[0]
	} else {
		m.shardSize = defaultShardSize
	}
	m.shards = make([]*shard[K, V], m.shardSize)
	for i := uint64(0); i < m.shardSize; i++ {
		m.shards[i] = &shard[K, V]{
			items: make(map[uint64]*value[K, V]),
		}
	}
	return m
}

// Load 获取键的值。
func (m *Map[K, V]) Load(key K) (V, bool) {
	sd, k := m.getShard(key)
	sd.mu.RLock()
	defer sd.mu.RUnlock()
	item, ok := sd.items[k]
	return item.val, ok
}

// Store 存储键的值。
func (m *Map[K, V]) Store(key K, val V) {
	sd, k := m.getShard(key)
	sd.mu.Lock()
	defer sd.mu.Unlock()
	sd.items[k] = &value[K, V]{
		key: key,
		val: val,
	}
}

// Delete 移除键的值。
func (m *Map[K, V]) Delete(key K) {
	sd, k := m.getShard(key)
	sd.mu.Lock()
	defer sd.mu.Unlock()
	delete(sd.items, k)
}

// DeleteAll 删除所有键值对
func (m *Map[K, V]) DeleteAll() {
	for _, sd := range m.shards {
		sd.mu.Lock()
		sd.items = make(map[uint64]*value[K, V])
		sd.mu.Unlock()
	}
}

func (m *Map[K, V]) Range(f func(key K, value V) bool) {
	for _, sd := range m.shards {
		sd.mu.RLock()
		for _, v := range sd.items {
			if !f(v.key, v.val) {
				sd.mu.RUnlock()
				return
			}
		}
		sd.mu.RUnlock()
	}
}

// DeletePrefix 删除指定前缀的键值对
func (m *Map[K, V]) DeletePrefix(prefix string) {
	// 如果K是string类型，才执行这个
	var zeroK K
	if _, ok := any(zeroK).(string); ok {
		for _, sd := range m.shards {
			sd.mu.Lock()
			for k, v := range sd.items {
				if strings.HasPrefix(tools.Any2string(v.key), prefix) {
					delete(sd.items, k)
				}
			}
			sd.mu.Unlock()
		}
	}
}

// DeleteSuffix 删除指定后缀的键值对
func (m *Map[K, V]) DeleteSuffix(suffix string) {
	// 如果K是string类型，才执行这个
	var zeroK K
	if _, ok := any(zeroK).(string); ok {
		for _, sd := range m.shards {
			sd.mu.Lock()
			for k, v := range sd.items {
				if strings.HasSuffix(tools.Any2string(v.key), suffix) {
					delete(sd.items, k)
				}
			}
			sd.mu.Unlock()
		}
	}
}

// LoadOrStore 获取键的值，如果没有则存储键的值。
func (m *Map[K, V]) LoadOrStore(key K, val V) (actual V, loaded bool) {
	sd, k := m.getShard(key)
	sd.mu.Lock()
	defer sd.mu.Unlock()
	item, ok := sd.items[k]
	if ok {
		return item.val, true
	}
	sd.items[k] = &value[K, V]{
		key: key,
		val: val,
	}
	return val, false
}

// LoadOrStoreFunc 获取键的值，如果没有则存储键的值。
func (m *Map[K, V]) LoadOrStoreFunc(key K, valueFunc func(k K) (V, error)) (V, bool, error) {
	sd, k := m.getShard(key)
	sd.mu.Lock()
	defer sd.mu.Unlock()
	item, ok := sd.items[k]
	if ok {
		return item.val, true, nil
	}
	val, err := valueFunc(key)
	if err != nil {
		return val, false, err
	}
	sd.items[k] = &value[K, V]{
		key: key,
		val: val,
	}
	return val, false, nil
}

// LoadAndDelete 获取键的值并删除键值对。
func (m *Map[K, V]) LoadAndDelete(key K) (value V, loaded bool) {
	sd, k := m.getShard(key)
	sd.mu.Lock()
	defer sd.mu.Unlock()
	item, ok := sd.items[k]
	if ok {
		delete(sd.items, k)
		return item.val, true
	}
	return value, false
}

// 获取分片
func (m *Map[K, V]) getShard(key K) (*shard[K, V], uint64) {
	k := m.hasher.Hash(key)
	return m.shards[k&(m.shardSize-1)], k
}
