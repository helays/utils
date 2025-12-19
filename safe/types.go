package safe

import (
	"github.com/cespare/xxhash/v2"
	"golang.org/x/exp/constraints"
)

const (
	defaultShardSize = 1 << 8 // 默认分片数量

)

// Hasher 编译时确定的哈希函数
type Hasher[K comparable] interface {
	Hash(K) uint64
}

// IntegerHasher 整数哈希器
type IntegerHasher[T constraints.Integer] struct{}

func (h IntegerHasher[T]) Hash(key T) uint64 {
	return uint64(key)
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
