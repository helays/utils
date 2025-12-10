package tools

import "encoding/json"

func Map2SliceWithHeader(m map[string]any, header []string) []any {
	var slice = make([]any, len(header), len(header))
	for i, k := range header {
		slice[i] = m[k]
	}
	return slice
}

// NestedMapSet 自动初始化嵌套 map 并设置值
// 如果内层map不存在，会自动创建
func NestedMapSet[K comparable, V comparable, T any](m map[K]map[V]T, outerKey K, innerKey V, value T) {
	if _, ok := m[outerKey]; !ok {
		m[outerKey] = make(map[V]T)
	}
	m[outerKey][innerKey] = value
}

// MapDeepCopy map 深拷贝
func MapDeepCopy[T any](src T, dst *T) {
	byt, _ := json.Marshal(src)
	_ = json.Unmarshal(byt, dst)
}

// ReverseMapUnique 反转值唯一的 map
func ReverseMapUnique[K comparable, V comparable](m map[K]V) map[V]K {
	reversed := make(map[V]K)
	for k, v := range m {
		reversed[v] = k
	}
	return reversed
}

// GetLevel2MapValue 获取二级map的值
func GetLevel2MapValue[K any](inp map[string]map[string]K, key1, key2 string) (K, bool) {
	if v, ok := inp[key1]; ok {
		if vv, ok := v[key2]; ok {
			return vv, true
		}
	}
	var zeroValue K
	return zeroValue, false
}
