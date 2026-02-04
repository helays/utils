package safe

import (
	"time"
	"unsafe"

	"github.com/cespare/xxhash/v2"
	"golang.org/x/exp/constraints"
)

const (
	defaultShardSize = 1 << 8 // 默认分片数量
	defaultCapacity  = 1 << 6 // 默认缓存大小
)

type CacheConfig struct {
	EnableCleanup bool          `json:"enable_cleanup" yaml:"enable_cleanup" ini:"enable_cleanup"` // 是否启用自动清理功能
	ClearInterval time.Duration `json:"clear_interval" yaml:"clear_interval" ini:"clear_interval"` // 清理间隔,推荐值是设置成 ttl/2 或者 ttl/3
	TTL           time.Duration `json:"ttl" yaml:"ttl" ini:"ttl"`                                  // 默认TTL为0，表示不过期
	ShardSize     uint64        `json:"shard_size" yaml:"shard_size" ini:"shard_size"`             // 分片数量，默认为2的8次方
	UseKey        bool          `json:"use_key" yaml:"use_key" ini:"use_key"`                      // 是否使用 key 作为 hash 值
}

// OnExpired 过期回调
type OnExpired[K comparable] func(key []K) // 过期回调

// Hasher 编译时确定的哈希函数
type Hasher[K comparable] interface {
	Hash(K) uint64
}

// IntegerHasher 整数哈希器
type IntegerHasher[T constraints.Integer] struct{}

func (h IntegerHasher[T]) Hash(key T) uint64 {
	return uint64(key)
}

// FloatHasher 浮点数哈希器
type FloatHasher[T constraints.Float] struct{}

func (h FloatHasher[T]) Hash(key T) uint64 {
	f64 := float64(key)
	return xxhash.Sum64(unsafe.Slice((*byte)(unsafe.Pointer(&f64)), 8))
}

// StringHasher 字符串哈希器
type StringHasher struct{}

func (h StringHasher) Hash(key string) uint64 {
	return xxhash.Sum64String(key)
}

// BytesHasher []byte哈希器
type BytesHasher struct{}

func (h BytesHasher) Hash(key []byte) uint64 {
	return xxhash.Sum64(key)
}

// Array16Hasher 为常用长度创建专门的哈希器
type Array16Hasher struct{}

func (h Array16Hasher) Hash(key [16]byte) uint64 {
	return xxhash.Sum64(key[:])
}
