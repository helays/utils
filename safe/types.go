package safe

import (
	"unsafe"

	"github.com/cespare/xxhash/v2"
	"golang.org/x/exp/constraints"
)

const (
	defaultShardSize = 1 << 8 // 默认分片数量
	defaultCapacity  = 1 << 6 // 默认缓存大小
)

type onExpired[K comparable] func(key []K) // 过期回调

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
