package tools

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
