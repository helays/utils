package tools

import (
	"context"
	"math/rand"
	"time"
)

// RunSyncFunc 同步运行
func RunSyncFunc(enable bool, f func()) {
	if enable {
		f()
	}
}

// RunAsyncFunc 异步运行
func RunAsyncFunc(enable bool, f func()) {
	if enable && f != nil {
		go f()
	}
}

// RunAsyncTickerFunc 异步运行，并定时执行
// ctx 用于控制循环的退出
// enable 是否启用
// d 执行间隔
// f 要执行的函数
// runFirst 是否先执行一次
func RunAsyncTickerFunc(ctx context.Context, enable bool, d time.Duration, f func(), runFirst ...bool) {
	if !enable {
		return
	}
	if f == nil {
		return
	}
	if len(runFirst) < 1 || runFirst[0] {
		f()
	}
	go func() {
		tck := time.NewTicker(d)
		defer tck.Stop()
		for {
			select {
			case <-ctx.Done(): // 退出循环
				return
			case <-tck.C:
				f()
			}
		}
	}()
}

// RunAsyncTickerProbabilityFunc 异步运行，并定时执行，概率触发
func RunAsyncTickerProbabilityFunc(ctx context.Context, enable bool, d time.Duration, probability float64, f func()) {
	if !enable {
		return
	}
	if f == nil {
		return
	}
	go func() {
		tck := time.NewTicker(d)
		defer tck.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-tck.C:
				if ProbabilityTrigger(probability) {
					f()
				}
			}
		}
	}()
}

// threadSafeRand 是一个全局变量，用于提供线程安全的随机数。
var threadSafeRand = rand.New(rand.NewSource(time.Now().UnixNano()))

// ProbabilityTrigger 使用线程安全的随机数生成器根据给定的概率触发事件
func ProbabilityTrigger(probability float64) bool {
	if probability <= 0 {
		return false
	} else if probability >= 1 {
		return true
	}
	// 生成一个0到1之间的随机浮点数
	randomNumber := threadSafeRand.Float64()
	// 比较随机数和概率
	return randomNumber < probability
}
