package safe

import (
	"context"
	"strings"
	"sync"
	"time"

	"github.com/helays/utils/v2/tools"
)

// 泛型 并发安全map,利用分片map提升并发性能。
// 如果

type (
	MapConfig struct {
		EnableCleanup bool          // 是否启用自动清理功能
		ClearInterval time.Duration // 清理间隔,推荐值是设置成 ttl/2 或者 ttl/3
		TTL           time.Duration // 默认TTL为0，表示不过期
		ShardSize     uint64        // 分片数量，默认为2的7次方
	}

	value[K comparable, V any] struct {
		key       K
		val       V
		expire    int64 // 过期时间，单位纳秒，0表示不过期
		ttl       time.Duration
		heartbeat time.Time // 用于记录最后一次心跳时间的。
	}
	shard[K comparable, V any] struct {
		items map[uint64]*value[K, V]
		mu    sync.RWMutex
	}
	Map[K comparable, V any] struct {
		ctx       context.Context // 上下文
		shards    []*shard[K, V]
		shardSize uint64
		hasher    Hasher[K]

		enableCleanup bool // 是否启用自动清理功能
		// 默认TTL为0，表示不过期，
		// 如果启用了enableCleanup，默认ttl和私有ttl都没设置，则不过期
		//  在启用enableCleanup后，如果有私有ttl就用私有ttl，否则用默认ttl
		defaultTTL time.Duration
		// 清理间隔,推荐值是设置成 ttl/2 或者 ttl/3
		// 当前启用自动清理后，clearInterval必须设置
		clearInterval time.Duration
		onExpired     onExpired[K] // 再过期时刻触发的回调操作。
	}
)

func NewMap[K comparable, V any](ctx context.Context, hasher Hasher[K], configs ...MapConfig) *Map[K, V] {
	m := &Map[K, V]{
		ctx:       ctx,
		hasher:    hasher,
		shardSize: defaultShardSize,
	}
	if len(configs) > 0 {
		config := configs[0]
		if !tools.IsPowerOfTwo(config.ShardSize) {
			panic("分片数量必须是2的N次幂")
		}
		m.shardSize = config.ShardSize
		m.enableCleanup = config.EnableCleanup
		if m.enableCleanup {
			m.defaultTTL = config.TTL
			m.clearInterval = tools.AutoTimeDuration(config.ClearInterval, time.Second, 30*time.Second) // 默认三十秒清理一次
		}

	}

	m.shards = make([]*shard[K, V], m.shardSize)
	for i := uint64(0); i < m.shardSize; i++ {
		m.shards[i] = &shard[K, V]{
			items: make(map[uint64]*value[K, V]),
		}
	}
	go m.cleanupDaemon() // 启动清理守护进程

	return m
}

func (m *Map[K, V]) SetOnExpired(onExpired onExpired[K]) {
	m.onExpired = onExpired
}

// Load 获取键的值。
func (m *Map[K, V]) Load(key K) (V, bool) {
	sd, k := m.getShard(key)
	sd.mu.RLock()
	defer sd.mu.RUnlock()
	if item, ok := sd.items[k]; ok {
		return item.val, true
	}
	var zero V
	return zero, false
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

// 清理进程
func (m *Map[K, V]) cleanupDaemon() {
	if !m.enableCleanup {
		return
	}
	tick := time.NewTicker(m.clearInterval)
	defer tick.Stop()
	for {
		select {
		case <-m.ctx.Done():
			return
		case <-tick.C:
			m.cleanupShards()

		}
	}
}

// 清理所有分片数据
func (m *Map[K, V]) cleanupShards() {
	for idx := range m.shards {
		select {
		case <-m.ctx.Done():
			return
		default:
			m.cleanupShard(idx)
			time.Sleep(time.Millisecond * 50) // 短暂休眠，让出CPU时间，避免清理进程占用cpu时间太长影响其他业务
		}
	}
}

// 清理单个分片数据
func (m *Map[K, V]) cleanupShard(idx int) {
	sd := m.shards[idx]
	sd.mu.Lock()
	defer sd.mu.Unlock()
	if len(sd.items) == 0 {
		return
	}
	now := time.Now().UnixNano()
	var expiredKeys []K
	for k, item := range sd.items {
		// 判断过期时间是否 > 0，并且需要时到期。
		if item.expire > 0 && now > item.expire {
			delete(sd.items, k)
			if m.onExpired != nil {
				if expiredKeys == nil {
					expiredKeys = make([]K, 0, 16) // 预先分配一个小缓冲队列
				}
				expiredKeys = append(expiredKeys, item.key)
			}
		}
	}
	if m.onExpired != nil && len(expiredKeys) > 0 {
		go m.onExpired(expiredKeys)
	}

}
