package tools

import "strings"

// HasKeyWithPrefix 检查 map 中是否存在以给定前缀开头的 key
func HasKeyWithPrefix[T any](m map[string]T, prefix string) bool {
	for k := range m {
		if strings.HasPrefix(k, prefix) {
			return true
		}
	}
	return false
}

// GetKeysWithPrefix 返回 map 中所有以给定前缀开头的 keys
func GetKeysWithPrefix[T any](m map[string]T, prefix string) []string {
	var keys []string
	for k := range m {
		if strings.HasPrefix(k, prefix) {
			keys = append(keys, k)
		}
	}
	return keys
}
