package tools

import (
	"context"
	"math/rand/v2"
	"time"

	"github.com/helays/utils/v2/config"
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

func RunAsyncTickerWithContext(ctx context.Context, enable bool, d time.Duration, f func(ctx context.Context)) {
	if !enable {
		return
	}
	go func() {
		f(ctx)
		tck := time.NewTicker(d)
		defer tck.Stop()
		for {
			select {
			case <-ctx.Done(): // 退出循环
				return
			case <-tck.C:
				f(ctx)
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

func RunAsyncTickerProbabilityWithContext(ctx context.Context, enable bool, d time.Duration, probability float64, f func(ctx context.Context)) {
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
					f(ctx)
				}
			}
		}
	}()
}

// ProbabilityTrigger 使用线程安全的随机数生成器根据给定的概率触发事件
func ProbabilityTrigger(probability float64) bool {
	if probability <= 0 {
		return false
	} else if probability >= 1 {
		return true
	}

	rng := config.RandPool.Get().(*rand.Rand)
	defer config.RandPool.Put(rng)
	// 生成一个0到1之间的随机浮点数
	randomNumber := rng.Float64()
	// 比较随机数和概率
	return randomNumber < probability
}

func AutoRetry(retryCount int, retryInterval time.Duration, f func() bool) bool {
	for i := 0; i < retryCount; i++ {
		if f() {
			return true
		}
		// 最后一次错误不睡眠
		if i < retryCount-1 {
			time.Sleep(retryInterval)
		}
	}
	return false
}

func AutoRetryWithErr(retryCount int, retryInterval time.Duration, f RetryCallbackFunc) error {
	var err error
	for i := 0; i < retryCount; i++ {
		if err = f(); err == nil {
			return nil
		}
		// 最后一次错误不睡眠
		if i < retryCount-1 {
			time.Sleep(retryInterval)
		}
	}
	return err
}

type RetryCallbackFunc func() error

// RetryRunner 重试执行函数
func RetryRunner(retry int, sleep time.Duration, callback RetryCallbackFunc) {
	for i := 0; i < retry; i++ {
		err := callback()
		if err == nil {
			return
		}
		// 最后一次不睡眠
		if i == retry-1 {
			break
		}
		time.Sleep(sleep)
	}
}

func WaitForCondition(ctx context.Context, condition func() bool) bool {
	ticker := time.NewTicker(time.Millisecond * 100)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return false
		case <-ticker.C:
			if condition() {
				return true
			}
		}
	}
}
