package cachemgr

import (
	"context"
	"time"

	"helay.net/go/utils/v3/safe"
)

// Cache 缓存接口
type Cache[K comparable, V any] interface {
	SetOnExpired(onExpired safe.OnExpired[K])                                                          // 设置过期回调函数。
	Load(key K) (V, bool)                                                                              // 获取键的值。
	LoadOrStore(key K, val V, duration ...time.Duration) (actual V, loaded bool)                       // 获取键的值，如果键不存在则存储键值对。
	LoadOrStoreFunc(key K, valueFunc func(k K) (V, error), duration ...time.Duration) (V, bool, error) // 获取键的值，如果键不存在则存储键值对。
	LoadAndDelete(key K) (V, bool)                                                                     // 获取键的值并删除键。
	LoadAndDeleteIf(key K, condition func(value V) bool) (V, bool)                                     // 获取键的值并删除键，如果满足条件。
	LoadAndRefresh(key K, duration ...time.Duration) (V, bool)                                         // 获取键的值并刷新键的过期时间。
	LoadWithExpiry(key K) (V, time.Time, bool)                                                         // 获取键的值和过期时间。
	Refresh(key K, duration ...time.Duration) bool                                                     // 刷新键的过期时间。
	GetTTL(key K) (time.Duration, bool)                                                                // 获取键的剩余过期时间。
	IsExpired(key K) bool                                                                              // 判断键是否已过期。
	GetHeartbeat(key K) (time.Time, bool)                                                              // 获取键的更新时间。

	Store(key K, val V, duration ...time.Duration) // 存储键值对。
	Delete(key K)                                  // 删除键。
	DeleteAndGetCount(keys ...K) int               // 删除多个键并返回删除的键数量。
	DeleteAll()                                    // 删除所有键。
	Range(f func(key K, value V) bool)             // 遍历缓存中的键值对。
	DeletePrefix(prefix string)                    // 删除以指定前缀开头的键。
	DeleteSuffix(suffix string)                    // 删除以指定后缀结尾的键。
}

type Driver string

// noinspection all
const (
	// 内存缓存
	DriverMemory Driver = "memory"
	// 关系数据库缓存
	// 缓存数据都在一张表中，但是由于有多个缓存实例，所以需要定义一个缓存标识来区分。
	// 这是一个联合索引。identity+key
	DriverRdbms Driver = "rdbms"
	DriverFile  Driver = "file"  // 文件缓存
	DriverRedis Driver = "redis" // redis 缓存
)

type Config struct {
	Driver           Driver `json:"driver" yaml:"driver" ini:"driver"`       // 缓存驱动
	Identity         string `json:"identity" yaml:"identity" ini:"identity"` // 缓存标识，在非内存存储下，非常有用。用于隔离多个缓存实例里面的数据。
	safe.CacheConfig `json:"memory" yaml:"memory" ini:"memory"`
}

func New[K comparable, V any](ctx context.Context, hasher safe.Hasher[K], cfg Config) Cache[K, V] {
	var (
		driver Cache[K, V]
	)
	switch cfg.Driver {
	case DriverMemory:
		driver = safe.NewMap[K, V](ctx, hasher, cfg.CacheConfig)
	default:
		return nil
	}

	return NewWithDriver[K, V](driver)
}

func NewWithDriver[K comparable, V any](driver Cache[K, V]) Cache[K, V] {
	return driver
}
