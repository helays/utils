package safe

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/helays/utils/v2/tools"
)

// 泛型 并发安全map,利用分片map提升并发性能。
// 支持设置 key 的过期时间
// 支持自动清理，当启用自动清理功能后，也支持设置不过期的key
// 可指定分片数量，默认为2的8次方。
// 默认以 hash 结果 uint64作为key，但也支持 key 作为map的 key。

// 使用默认 uint64 key 的安全数据量范围
// 绝对安全区：n < 100万 (1,000,000)
//碰撞概率：< 3×10⁻⁸ (< 0.000003%)
//几乎不可能发生碰撞
//适合：缓存、会话存储、配置存储
//安全区：n < 1000万 (10,000,000)
//碰撞概率：< 3×10⁻⁶ (< 0.0003%)
//实际应用中可接受
//适合：中等规模应用
//风险区：n > 1亿 (100,000,000)
//碰撞概率开始显著 (> 0.03%)
//需要哈希碰撞检测
//适合：仅限高性能计算场景

type (
	MapConfig struct {
		EnableCleanup bool          // 是否启用自动清理功能
		ClearInterval time.Duration // 清理间隔,推荐值是设置成 ttl/2 或者 ttl/3
		TTL           time.Duration // 默认TTL为0，表示不过期
		ShardSize     uint64        // 分片数量，默认为2的8次方
		UseKey        bool          // 是否使用 key 作为 hash 值
	}

	value[K comparable, V any] struct {
		key       K
		val       V
		expire    int64 // 过期时间，单位纳秒，0表示不过期
		ttl       time.Duration
		heartbeat time.Time // 用于记录最后一次心跳时间的。
	}
	shard[K comparable, V any] struct {
		hashItems map[uint64]*value[K, V]
		keyItems  map[K]*value[K, V]

		mu sync.RWMutex
	}
	Map[K comparable, V any] struct {
		ctx       context.Context // 上下文
		shards    []*shard[K, V]
		shardSize uint64
		hasher    Hasher[K]
		useKey    bool // 是否使用 key 作为 hash 值

		enableCleanup bool // 是否启用自动清理功能
		// 默认TTL为0，表示不过期，
		// 如果启用了enableCleanup，默认ttl和私有ttl都没设置，则不过期
		// 在启用enableCleanup后，如果有私有ttl就用私有ttl，否则用默认ttl
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
		if config.ShardSize > 0 {
			if !tools.IsPowerOfTwo(config.ShardSize) {
				panic(fmt.Errorf("异常的分片数量[%d]，分片数量必须是2的N次幂", config.ShardSize))
			}
			m.shardSize = config.ShardSize
		}

		m.useKey = config.UseKey
		m.enableCleanup = config.EnableCleanup
		if m.enableCleanup {
			m.defaultTTL = config.TTL
			m.clearInterval = tools.AutoTimeDuration(config.ClearInterval, time.Second, 30*time.Second) // 默认三十秒清理一次
		}

	}

	m.shards = make([]*shard[K, V], m.shardSize)
	for i := uint64(0); i < m.shardSize; i++ {
		sd := &shard[K, V]{}
		if m.useKey {
			sd.keyItems = make(map[K]*value[K, V], defaultCapacity)
		} else {
			sd.hashItems = make(map[uint64]*value[K, V], defaultCapacity) // 预分配128个空间
		}
		m.shards[i] = sd

	}
	go m.cleanupDaemon() // 启动清理守护进程

	return m
}

func (m *Map[K, V]) SetOnExpired(onExpired onExpired[K]) {
	m.onExpired = onExpired
}

// Load 获取键的值。
func (m *Map[K, V]) Load(key K) (V, bool) {
	sd, _, k := m.getShard(key)
	sd.mu.RLock()
	defer sd.mu.RUnlock()
	val, _, ok, _ := m.load(sd, key, k)
	return val, ok
}

// load 获取键的值。
// 这个函数都不会独立使用，基于上级函数，已经做了锁功能。
// 这个函数内容也已经做了过期数据判断。
// val 存储的值
// item 存储的结构，包含元数据
// ok 是否有效值
// existed 仅仅用于判断key是否存在
// 上级函数在调用这个函数的时候，如果是 RLock,如果值无效但是值存在，就不调用删除函数。
// Lock,如果不是store相关的，需要调用删除函数。
func (m *Map[K, V]) load(sd *shard[K, V], key K, k uint64) (val V, item *value[K, V], ok bool, existed bool) {

	if m.useKey {
		item, ok = sd.keyItems[key]
	} else {
		item, ok = sd.hashItems[k]
	}
	var zero V
	if !ok {
		return zero, nil, false, false
	}
	// 时间==0，表示不过期
	// 过期时间>当前时间，表示未过期
	if item.expire == 0 || item.expire > time.Now().UnixNano() {
		return item.val, item, true, true
	}
	// 这个位置 不能删除，应该由自动清理或者手动删除来操作。对于rw锁，多读是无锁的
	// 会形成竟态
	//m.delete(sd, key, k)
	return zero, nil, false, true
}

// LoadOrStore 获取键的值，如果没有则存储键的值。
func (m *Map[K, V]) LoadOrStore(key K, val V, duration ...time.Duration) (actual V, loaded bool) {
	sd, _, k := m.getShard(key)
	sd.mu.Lock()
	defer sd.mu.Unlock()
	item, _, ok, _ := m.load(sd, key, k)
	if ok {
		return item, true
	}
	m.store(sd, key, k, val, duration...)
	return val, false
}

// LoadOrStoreFunc 获取键的值，如果没有则存储键的值。
func (m *Map[K, V]) LoadOrStoreFunc(key K, valueFunc func(k K) (V, error), duration ...time.Duration) (V, bool, error) {
	sd, _, k := m.getShard(key)
	sd.mu.Lock()
	defer sd.mu.Unlock()
	item, _, ok, _ := m.load(sd, key, k)
	if ok {
		return item, true, nil
	}

	if valueFunc == nil {
		return item, false, nil
	}

	val, err := valueFunc(key)
	if err != nil {
		return val, false, err
	}
	m.store(sd, key, k, val, duration...)

	return val, false, nil
}

// LoadAndDelete 获取键的值并删除键值对。
func (m *Map[K, V]) LoadAndDelete(key K) (V, bool) {
	sd, _, k := m.getShard(key)
	sd.mu.Lock()
	defer sd.mu.Unlock()
	val, _, ok, existed := m.load(sd, key, k)
	if existed {
		m.delete(sd, key, k)
	}
	if !ok {
		return val, false
	}
	return val, true
}

// LoadAndDeleteIf 获取键的值并删除键值对，条件性删除
// 返回值第二个参数 如果 == true,表示删除成功
func (m *Map[K, V]) LoadAndDeleteIf(key K, condition func(value V) bool) (V, bool) {
	sd, _, k := m.getShard(key)
	sd.mu.Lock()
	defer sd.mu.Unlock()
	val, _, ok, existed := m.load(sd, key, k)
	if !ok {
		// 值无效，但是值存在，就删除。
		if existed {
			m.delete(sd, key, k)
		}
		return val, false
	}
	if condition == nil || condition(val) {
		m.delete(sd, key, k)
		return val, true
	}
	return val, false
}

func (m *Map[K, V]) LoadAndRefresh(key K, duration ...time.Duration) (V, bool) {
	sd, _, k := m.getShard(key)
	sd.mu.Lock()
	defer sd.mu.Unlock()
	val, item, ok, existed := m.load(sd, key, k)
	if !ok {
		// 值无效，但是值存在，就删除。
		if existed {
			m.delete(sd, key, k)
		}
		return val, false
	}
	if m.enableCleanup {
		if len(duration) > 0 {
			item.ttl = duration[0]
		}
		item.heartbeat = time.Now()
		// 如果ttl<=0，则表示不过期
		if item.ttl <= 0 {
			item.expire = 0
		} else {
			item.expire = time.Now().Add(item.ttl).UnixNano()
		}
	}
	return val, true
}

// LoadWithExpiry 获取键的值和过期时间
func (m *Map[K, V]) LoadWithExpiry(key K) (V, time.Time, bool) {
	sd, _, k := m.getShard(key)
	sd.mu.RLock()
	defer sd.mu.RUnlock()
	val, item, ok, _ := m.load(sd, key, k)
	if !ok {
		return val, time.Time{}, false
	}
	return val, time.Unix(0, item.expire), true
}

// Refresh 更新键的过期时间
func (m *Map[K, V]) Refresh(key K, duration ...time.Duration) bool {
	// 如果没有启用自动清理功能，则返回false
	if !m.enableCleanup {
		return false
	}
	sd, _, k := m.getShard(key)
	sd.mu.Lock()
	defer sd.mu.Unlock()
	_, item, ok, existed := m.load(sd, key, k)
	if !ok {
		// 值无效，但是值存在，就删除。
		if existed {
			m.delete(sd, key, k)
		}
		return false
	}
	if len(duration) > 0 {
		item.ttl = duration[0]
	}
	item.heartbeat = time.Now()
	// 如果ttl<=0，则表示不过期
	if item.ttl <= 0 {
		item.expire = 0
	} else {
		item.expire = time.Now().Add(item.ttl).UnixNano()
	}
	return true
}

// GetTTL 获取键的剩余存活时间
func (m *Map[K, V]) GetTTL(key K) (time.Duration, bool) {
	sd, _, k := m.getShard(key)
	sd.mu.RLock()
	defer sd.mu.RUnlock()
	_, item, ok, _ := m.load(sd, key, k)
	if !ok {
		return 0, false
	}
	if item.expire == 0 {
		return 0, true
	}
	return time.Duration(item.expire - time.Now().UnixNano()), true
}

// IsExpired 检测键是否已过期
func (m *Map[K, V]) IsExpired(key K) bool {
	sd, _, k := m.getShard(key)
	sd.mu.RLock()
	defer sd.mu.RUnlock()
	_, item, ok, _ := m.load(sd, key, k)
	if !ok {
		return false
	}
	if item.expire == 0 {
		return false
	}
	return item.expire < time.Now().UnixNano()
}

func (m *Map[K, V]) GetHeartbeat(key K) (time.Time, bool) {
	sd, _, k := m.getShard(key)
	sd.mu.RLock()
	defer sd.mu.RUnlock()
	_, item, ok, _ := m.load(sd, key, k)
	if !ok {
		return time.Time{}, false
	}
	return item.heartbeat, true
}

// Store 存储键的值。
func (m *Map[K, V]) Store(key K, val V, duration ...time.Duration) {
	sd, _, k := m.getShard(key)
	sd.mu.Lock()
	defer sd.mu.Unlock()
	m.store(sd, key, k, val, duration...)
}

// store 存储键的值。
func (m *Map[K, V]) store(sd *shard[K, V], key K, k uint64, val V, duration ...time.Duration) {
	v := &value[K, V]{
		key: key,
		val: val,
	}
	// 启用自动清理功能后，如果有私有TTL，则用私有TTL，否则用默认TTL
	// 最终还是需要判断 ttl是否大于0。如果ttl <= 0 就表示当前key不过期。
	if m.enableCleanup {
		v.ttl = m.defaultTTL
		if len(duration) > 0 && duration[0] > 0 {
			v.ttl = duration[0]
		}
		if v.ttl > 0 {
			v.expire = time.Now().Add(v.ttl).UnixNano()
			v.heartbeat = time.Now()
		}
	}

	if m.useKey {
		sd.keyItems[key] = v
	} else {
		sd.hashItems[k] = v
	}
}

// Delete 移除键的值。
func (m *Map[K, V]) Delete(key K) {
	sd, _, k := m.getShard(key)
	sd.mu.Lock()
	defer sd.mu.Unlock()
	m.delete(sd, key, k)
}

func (m *Map[K, V]) delete(sd *shard[K, V], key K, k uint64) {
	if m.useKey {
		delete(sd.keyItems, key)
	} else {
		delete(sd.hashItems, k)
	}
}

// DeleteAndGetCount 删除多个键并返回删除的数量
func (m *Map[K, V]) DeleteAndGetCount(keys ...K) int {
	count := 0
	type group struct {
		key     K
		hashKey uint64
	}
	shardMap := make(map[uint64][]group)
	for _, key := range keys {
		_, idx, k := m.getShard(key)
		shardMap[idx] = append(shardMap[idx], group{key: key, hashKey: k})
	}

	for idx, groups := range shardMap {
		sd := m.shards[idx]
		sd.mu.Lock()
		for _, g := range groups {
			if m.useKey {
				if _, ok := sd.keyItems[g.key]; ok {
					delete(sd.keyItems, g.key)
					count++
				}
			} else {
				if _, ok := sd.hashItems[g.hashKey]; ok {
					delete(sd.hashItems, g.hashKey)
					count++
				}
			}
		}
		sd.mu.Unlock()
	}
	return count
}

// DeleteAll 删除所有键值对
func (m *Map[K, V]) DeleteAll() {
	for _, sd := range m.shards {
		sd.mu.Lock()
		if m.useKey {
			sd.keyItems = make(map[K]*value[K, V])
		} else {
			sd.hashItems = make(map[uint64]*value[K, V])
		}
		sd.mu.Unlock()
	}
}

func (m *Map[K, V]) Range(f func(key K, value V) bool) {

	var rangeFunc = func(v *value[K, V]) bool {
		now := time.Now().UnixNano()
		if v.expire == 0 || v.expire > now {
			if !f(v.key, v.val) {
				return false
			}
		}
		return true
	}

	for _, sd := range m.shards {
		sd.mu.RLock()
		if m.useKey {
			for _, v := range sd.keyItems {
				if !rangeFunc(v) {
					sd.mu.RUnlock()
					return
				}
			}
		} else {
			for _, v := range sd.hashItems {
				if !rangeFunc(v) {
					sd.mu.RUnlock()
					return
				}
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
			if m.useKey {
				for k, v := range sd.keyItems {
					if strings.HasPrefix(tools.Any2string(v.key), prefix) {
						delete(sd.keyItems, k)
					}
				}
			} else {
				for k, v := range sd.hashItems {
					if strings.HasPrefix(tools.Any2string(v.key), prefix) {
						delete(sd.hashItems, k)
					}
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
			if m.useKey {
				for k, v := range sd.keyItems {
					if strings.HasSuffix(tools.Any2string(v.key), suffix) {
						delete(sd.keyItems, k)
					}
				}
			} else {
				for k, v := range sd.hashItems {
					if strings.HasSuffix(tools.Any2string(v.key), suffix) {
						delete(sd.hashItems, k)
					}
				}
			}
			sd.mu.Unlock()
		}
	}
}

// 获取分片
func (m *Map[K, V]) getShard(key K) (*shard[K, V], uint64, uint64) {
	k := m.hasher.Hash(key)
	idx := k & (m.shardSize - 1)
	return m.shards[idx], idx, k
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

	now := time.Now().UnixNano()
	var expiredKeys []K

	if m.useKey {
		if len(sd.keyItems) == 0 {
			return
		}
		for k, item := range sd.keyItems {
			// 判断过期时间是否 > 0，并且需要时到期。
			if item.expire > 0 && now > item.expire {
				delete(sd.keyItems, k)
				if m.onExpired != nil {
					if expiredKeys == nil {
						expiredKeys = make([]K, 0, 16) // 预先分配一个小缓冲队列
					}
					expiredKeys = append(expiredKeys, item.key)
				}
			}
		}
	} else {
		if len(sd.hashItems) == 0 {
			return
		}
		for k, item := range sd.hashItems {
			// 判断过期时间是否 > 0，并且需要时到期。
			if item.expire > 0 && now > item.expire {
				delete(sd.hashItems, k)
				if m.onExpired != nil {
					if expiredKeys == nil {
						expiredKeys = make([]K, 0, 16) // 预先分配一个小缓冲队列
					}
					expiredKeys = append(expiredKeys, item.key)
				}
			}
		}
	}

	if m.onExpired != nil && len(expiredKeys) > 0 {
		go m.onExpired(expiredKeys)
	}

}
