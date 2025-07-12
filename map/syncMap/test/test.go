package main

import (
	"fmt"
	"github.com/helays/utils/v2/map/syncMap"
)

// 使用示例
// 使用示例
func main() {

	// 创建一个 SafeSegmentedMap
	safeMap := syncMap.NewMap[string, string]()
	safeMap.Store("key1", "value1")
	safeMap.Store("key2", "value2")

	// 加载值
	value, ok := safeMap.Load("key1")
	if ok {
		fmt.Println("Loaded value:", value)
	}

	// 删除值
	safeMap.Delete("key1")
	_, ok = safeMap.Load("key1")
	if !ok {
		fmt.Println("Key1 not found")
	}

	// 遍历所有键值对
	safeMap.Range(func(key string, value string) bool {
		fmt.Printf("Key: %s, Value: %v\n", key, value)
		return true
	})

	// 创建一个带有不同类型的 SafeSegmentedMap
	intSafeMap := syncMap.NewMap[int, int]()
	intSafeMap.Store(1, 100)
	intSafeMap.Store(2, 200)

	// 加载值
	intValue, ok := intSafeMap.Load(1)
	if ok {
		fmt.Println("Loaded value:", intValue)
	}

	// 删除值
	intSafeMap.Delete(1)
	_, ok = intSafeMap.Load(1)
	if !ok {
		fmt.Println("Key1 not found")
	}

	// 遍历所有键值对
	intSafeMap.Range(func(key int, value int) bool {
		fmt.Printf("Key: %d, Value: %d\n", key, value)
		return true
	})
}
