package tools

import "strings"

// Searchslice 在切片中判断某个值是否存在
// Deprecated: 使用 Contains
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

// SearchIntSlice 在整数切片中搜索指定的元素，并返回是否找到。
// 参数:
//
//	s - 待搜索的整数。
//	arr - 整数切片，将被搜索。
//
// 返回值:
//
//	如果找到 s 在 arr 中，则返回 true；否则返回 false。
//
// Deprecated: 使用 Contains
func SearchIntSlice(s int, arr []int) bool {
	if arr == nil {
		return false
	}
	for _, i := range arr {
		if i == s {
			return true
		}
	}
	return false
}

// Deprecated: 使用 Contains
func SearchInt64Slice(s int64, arr []int64) bool {
	if arr == nil {
		return false
	}
	for _, i := range arr {
		if i == s {
			return true
		}
	}
	return false
}

// Deprecated: 使用 Contains
func SearchStringSlice(key string, arr []string) bool {
	if arr == nil {
		return false
	}
	for _, v := range arr {
		if v == key {
			return true
		}
	}
	return false
}
