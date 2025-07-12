package syncMap

import (
	"fmt"
	"math/rand"
	"sync"
	"testing"
	"time"

	"github.com/helays/utils/v2/map/syncMap" // 替换为你的包路径
)

// go test -bench=BenchmarkAll -benchmem 测试指令

func init() {
	rand.Seed(time.Now().UnixNano())
}

// benchmarkMaps 是一个辅助函数，用于实际执行性能测试。
func benchmarkMaps(b *testing.B, mapType string, readRatio float64, concurrency int, dataCount int) {
	var sm interface{}
	if mapType == "syncMap" {
		sm = syncMap.NewMap[string, string]()
		defer sm.(*syncMap.Map[string, string]).Close()
	} else {
		sm = &sync.Map{}
	}

	keys := make([]string, dataCount)
	for i := 0; i < dataCount; i++ {
		keys[i] = fmt.Sprintf("key-%d", i)
	}

	// 初始化数据
	for _, key := range keys {
		value := fmt.Sprintf("value-%s", key)
		switch v := sm.(type) {
		case *syncMap.Map[string, string]:
			v.Store(key, value)
		case *sync.Map:
			v.Store(key, value)
		}
	}

	var wg sync.WaitGroup
	wg.Add(concurrency)

	b.ResetTimer()

	for i := 0; i < concurrency; i++ {
		go func() {
			defer wg.Done()
			for j := 0; j < b.N/concurrency; j++ {
				if rand.Float64() < readRatio {
					key := keys[rand.Intn(dataCount)]
					switch v := sm.(type) {
					case *syncMap.Map[string, string]:
						v.Load(key)
					case *sync.Map:
						v.Load(key)
					}
				} else {
					key := keys[rand.Intn(dataCount)]
					value := fmt.Sprintf("value-%s", key)
					switch v := sm.(type) {
					case *syncMap.Map[string, string]:
						v.Store(key, value)
					case *sync.Map:
						v.Store(key, value)
					}
				}
			}
		}()
	}

	wg.Wait()
}

// BenchmarkAll 测试 syncMap.Map 和 sync.Map 在多种条件下的性能，并输出结果
func BenchmarkAll(b *testing.B) {
	testCases := []struct {
		Name       string
		MapType    string
		ReadRatio  float64
		Concurrent int
		DataCount  int
	}{
		{"SyncMapReadHeavy", "syncMap", 0.9, 10, 1000},
		{"SyncMapWriteHeavy", "syncMap", 0.1, 10, 1000},
		{"SyncMapMixed", "syncMap", 0.5, 10, 1000},
		{"StdSyncMapReadHeavy", "sync.Map", 0.9, 10, 1000},
		{"StdSyncMapWriteHeavy", "sync.Map", 0.1, 10, 1000},
		{"StdSyncMapMixed", "sync.Map", 0.5, 10, 1000},
		// 您可以根据需要添加更多的测试用例
	}

	for _, tc := range testCases {
		b.Run(tc.Name, func(b *testing.B) {
			benchmarkMaps(b, tc.MapType, tc.ReadRatio, tc.Concurrent, tc.DataCount)
		})
	}
}
