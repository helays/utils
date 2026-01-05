package tools

import "cmp"

func Max[T cmp.Ordered](d1, d2 T) T {
	if d1 > d2 {
		return d1
	}
	return d2
}

func Min[T cmp.Ordered](d1, d2 T) T {
	if d1 < d2 {
		return d1
	}
	return d2
}

func MaxInt32(d1, d2 int32) int32 {
	if d1 > d2 {
		return d1
	}
	return d2
}

// AvgInt32 计算平均数
func AvgInt32(d1, d2 int32, isf bool) int32 {
	if isf {
		if d1 > d2 {
			return d1
		}
		return d2
	}
	return (d1 + d2) / 2
}

func MinInt32(d1, d2 int32) int32 {
	if d1 > d2 {
		return d2
	}
	return d1
}

func MaxInt64(d1, d2 int64) int64 {
	if d1 > d2 {
		return d1
	}
	return d2
}

func MinInt64(d1, d2 int64) int64 {
	if d1 > d2 {
		return d2
	}
	return d1
}
func MaxUint64(d1, d2 uint64) uint64 {
	if d1 > d2 {
		return d1
	}
	return d2
}

func MinUint64(d1, d2 uint64) uint64 {
	if d1 > d2 {
		return d2
	}
	return d1
}
func MaxFloat32(d1, d2 float32) float32 {
	if d1 > d2 {
		return d1
	}
	return d2
}

func MinFloat32(d1, d2 float32) float32 {
	if d1 > d2 {
		return d2
	}
	return d1
}

// MinMaxAvgSum 获取数组中最大值最小值平均值和求和
func MinMaxAvgSum(nums []int) (min int, max int, avg float64, sum int) {
	if len(nums) == 0 {
		return 0, 0, 0, 0
	}
	min, max, sum = nums[0], nums[0], nums[0]
	for _, num := range nums[1:] {
		if num < min {
			min = num
		}
		if num > max {
			max = num
		}
		sum += num
	}
	avg = float64(sum) / float64(len(nums))
	return
}

// IsZero isZero 检查值是否为类型的零值
func IsZero[T comparable](v T) bool {
	var zero T
	return v == zero
}

// Rangeable 定义可用于range的接口
// Idx: 集合的键类型（数组索引、map键等）
// V: 元素值类型
// C: 用于比较的键类型
// noinspection all
type Rangeable[Idx comparable, V any, C comparable] interface {
	Range(func(Idx, V))
	Len() int
	ExtractKey(V) C
}

// DiffRangeable 比较两个Rangeable
// 返回两个Rangeable的差异
// src: 源数据
// dst: 目标数据
// 返回值:
// inSrc: 目标数据中不存在的元素
// common: 两个Rangeable中都存在的元素
// inDst: 源数据中不存在的元素
// noinspection all
func DiffRangeable[Idx comparable, V any, C comparable](src, dst Rangeable[Idx, V, C]) (inSrc, common, inDst []V) {
	elementMap := make(map[C]int8, src.Len()+dst.Len())
	// 遍历源数据,标记成1
	src.Range(func(_ Idx, v V) {
		elementMap[src.ExtractKey(v)] = 1
	})
	// 遍历目标数据
	dst.Range(func(_ Idx, v V) {
		val := dst.ExtractKey(v)
		if status, exists := elementMap[val]; exists {
			// 原数组中存在的
			if status == 1 {
				elementMap[val] = 3 // 标记共同存在
				common = append(common, v)
			}
		} else {
			elementMap[val] = 2
			inDst = append(inDst, v)
		}
	})
	// 收集只在源数组中存在的元素
	src.Range(func(_ Idx, v V) {
		if status := elementMap[src.ExtractKey(v)]; status == 1 {
			inSrc = append(inSrc, v)
		}
	})
	return
}
