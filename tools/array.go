package tools

import (
	"strings"
)

// ArrayChunk 高性能泛型切片分块函数
func ArrayChunk[T any](slice []T, size int) [][]T {
	if size <= 0 {
		return [][]T{slice}
	}

	length := len(slice)
	if length == 0 {
		return nil
	}

	chunkCount := (length + size - 1) / size
	chunks := make([][]T, chunkCount)

	for i := 0; i < chunkCount; i++ {
		start := i * size
		end := start + size
		if end > length {
			end = length
		}
		// 直接引用原切片的底层数组，避免额外内存分配
		chunks[i] = slice[start:end:end] // 使用完整切片表达式限制容量
	}

	return chunks
}

// ArrayChunkCopy 高性能且安全的版本，会复制数据而不是引用原切片
func ArrayChunkCopy[T any](slice []T, size int) [][]T {
	if size <= 0 {
		copied := make([]T, len(slice))
		copy(copied, slice)
		return [][]T{copied}
	}

	length := len(slice)
	if length == 0 {
		return nil
	}

	chunkCount := (length + size - 1) / size
	chunks := make([][]T, chunkCount)

	for i := 0; i < chunkCount; i++ {
		start := i * size
		end := start + size
		if end > length {
			end = length
		}
		chunk := make([]T, end-start)
		copy(chunk, slice[start:end])
		chunks[i] = chunk
	}

	return chunks
}

// Slice2Map 更清晰的参数命名
func Slice2Map[Key comparable, Elem any](slice []Elem, keyFunc func(Elem) Key) map[Key]Elem {
	result := make(map[Key]Elem)
	for _, item := range slice {
		result[keyFunc(item)] = item
	}
	return result
}

// SliceToMultiMap 将切片转换为映射，允许键重复，相同键的值会保存在切片中
// slice: 要转换的切片
// keyFunc: 从元素中提取键的函数
func SliceToMultiMap[Key comparable, Elem any](slice []Elem, keyFunc func(Elem) Key) map[Key][]Elem {
	result := make(map[Key][]Elem)
	for _, item := range slice {
		key := keyFunc(item)
		result[key] = append(result[key], item)
	}
	return result
}

// RemoveDuplicates 对slice去重
func RemoveDuplicates[T comparable](slice []T) []T {
	encountered := map[T]struct{}{}
	result := make([]T, 0, len(slice))

	for _, v := range slice {
		if _, ok := encountered[v]; !ok {
			encountered[v] = struct{}{}
			result = append(result, v)
		}
	}

	return result
}

// RemoveDuplicatesWithKeyFunc 对slice去重，使用keyFunc生成比较键
func RemoveDuplicatesWithKeyFunc[T any, K comparable](slice []T, keyFunc func(T) K) []T {
	encountered := map[K]struct{}{}
	result := make([]T, 0, len(slice))

	for _, v := range slice {
		key := keyFunc(v)
		if _, exists := encountered[key]; !exists {
			encountered[key] = struct{}{}
			result = append(result, v)
		}
	}

	return result
}

// Searchslice 在切片中判断某个值是否存在
func Searchslice(s string, o []string) bool {
	if o == nil {
		return false
	}
	s = strings.TrimSpace(s)
	for _, i := range o {
		i = strings.TrimSpace(i)
		if i == s {
			return true
		}
	}
	return false
}

// Contains 检查某个值是否在切片中（泛型实现）
func Contains[T comparable](slice []T, target T) bool {
	for _, v := range slice {
		if v == target {
			return true
		}
	}
	return false
}

// ContainsFunc 支持自定义比较逻辑（适用于不可比较类型，如 `struct` 带非 `comparable` 字段）
func ContainsFunc[T any](slice []T, target T, equal func(a, b T) bool) bool {
	for _, v := range slice {
		if equal(v, target) {
			return true
		}
	}
	return false
}

// ContainsAny 检查 `elems` 中是否有任意一个元素在 `targets` 里
func ContainsAny[T comparable](elems []T, targets []T) bool {
	if len(targets) == 0 || len(elems) == 0 {
		return false
	}
	// 选择较小的集合缓存
	var staticSet map[T]struct{}
	var dynamic []T
	if len(targets) < len(elems) {
		staticSet = make(map[T]struct{}, len(targets))
		for _, t := range targets {
			staticSet[t] = struct{}{}
		}
		dynamic = elems
	} else {
		staticSet = make(map[T]struct{}, len(elems))
		for _, e := range elems {
			staticSet[e] = struct{}{}
		}
		dynamic = targets
	}
	// 遍历动态集合
	for _, v := range dynamic {
		if _, ok := staticSet[v]; ok {
			return true
		}
	}
	return false
}

// ContainsAnyHashBest 检查 `elems` 中是否有任意一个元素在 `targets` 里，自定义比较函数
// elems 检查列表
// targets 被搜索列表
// hashFunc  计算hash
// equal 精确比较
func ContainsAnyHashBest[T any, H comparable](elems []T, targets []T, hashFunc func(T) H, equal func(a, b T) bool) bool {
	// 选择较小的集合作为缓存
	var staticSet map[H][]T
	var dynamic []T
	if len(targets) < len(elems) {
		staticSet = make(map[H][]T, len(targets))
		for _, t := range targets {
			h := hashFunc(t)
			staticSet[h] = append(staticSet[h], t)
		}
		dynamic = elems
	} else {
		staticSet = make(map[H][]T, len(elems))
		for _, e := range elems {
			h := hashFunc(e)
			staticSet[h] = append(staticSet[h], e)
		}
		dynamic = targets
	}

	// 遍历动态集合
	for _, v := range dynamic {
		h := hashFunc(v)
		for _, s := range staticSet[h] {
			if equal(v, s) {
				return true
			}
		}
	}
	return false
}
