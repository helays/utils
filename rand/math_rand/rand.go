package math_rand

import (
	"math/rand"
	"time"
)

//
// ━━━━━━神兽出没━━━━━━
// 　　 ┏┓     ┏┓
// 　　┏┛┻━━━━━┛┻┓
// 　　┃　　　　　 ┃
// 　　┃　　━　　　┃
// 　　┃　┳┛　┗┳  ┃
// 　　┃　　　　　 ┃
// 　　┃　　┻　　　┃
// 　　┃　　　　　 ┃
// 　　┗━┓　　　┏━┛　Code is far away from bug with the animal protecting
// 　　　 ┃　　　┃    神兽保佑,代码无bug
// 　　　　┃　　　┃
// 　　　　┃　　　┗━━━┓
// 　　　　┃　　　　　　┣┓
// 　　　　┃　　　　　　┏┛
// 　　　　┗┓┓┏━┳┓┏┛
// 　　　　 ┃┫┫ ┃┫┫
// 　　　　 ┗┻┛ ┗┻┛
//
// ━━━━━━感觉萌萌哒━━━━━━
//
//
// User helay
// Date: 2025/6/15 18:42
//

var rng = rand.New(rand.NewSource(time.Now().UnixNano()))

// RandomElement 从任意切片或数组中随机返回一个元素
func RandomElement[T any](collection []T) (T, bool) {
	var zero T // 用于返回类型的零值

	if len(collection) == 0 {
		return zero, false
	}

	// 生成随机索引
	randomIndex := rng.Intn(len(collection))

	return collection[randomIndex], true
}

func RandomInt(min, max int) int {
	return rng.Intn(max-min+1) + min
}

func RandomFloat(min, max float64) float64 {
	return min + rng.Float64()*(max-min)
}
