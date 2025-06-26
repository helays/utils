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
