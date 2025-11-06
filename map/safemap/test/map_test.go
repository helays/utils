package test

import (
	"fmt"
	"github.com/helays/utils/v2/map/safemap"
	"sync"
	"testing"
)

// 假设 SyncMapWrapper 已经按照之前提供的代码实现。

func BenchmarkSyncMaps(b *testing.B) {
	// Direct sync.Map benchmarks
	b.Run("DirectSyncMap_StoreLoad", func(b *testing.B) {
		m := &sync.Map{}

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			key := fmt.Sprintf("key-%d", i%100)
			value := i

			m.Store(key, value)
			if _, loaded := m.Load(key); !loaded {
				b.Fatal("Failed to load stored key")
			}
		}
	})

	b.Run("GenericSyncMapWrapper_StoreLoad", func(b *testing.B) {
		m := &safemap.SyncMap[string, int]{}

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			key := fmt.Sprintf("key-%d", i%100)
			value := i

			m.Store(key, value)
			if val, loaded := m.Load(key); !loaded || val != value {
				b.Fatal("Failed to load stored key or value mismatch")
			}
		}
	})

	// Pre-populate keys for delete tests
	keys := make([]string, 0, 100)
	for i := 0; i < 100; i++ {
		keys = append(keys, fmt.Sprintf("key-%d", i))
	}

	b.Run("DirectSyncMap_Delete", func(b *testing.B) {
		m := &sync.Map{}
		for _, key := range keys {
			m.Store(key, 0) // Pre-populate the map.
		}

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			for _, key := range keys {
				m.Delete(key)
				m.Store(key, i) // Re-store to ensure we have something to delete next iteration.
			}
		}
	})

	b.Run("GenericSyncMapWrapper_Delete", func(b *testing.B) {
		m := &safemap.SyncMap[string, int]{}
		for _, key := range keys {
			m.Store(key, 0) // Pre-populate the map.
		}

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			for _, key := range keys {
				m.Delete(key)
				m.Store(key, i) // Re-store to ensure we have something to delete next iteration.
			}
		}
	})
}

func main() {
	// 如果你想直接从main函数运行基准测试，可以这样做：
	// 注意：通常情况下，基准测试应该通过命令行工具（如 go test）来执行。
	test := testing.B{}
	BenchmarkSyncMaps(&test)
}
