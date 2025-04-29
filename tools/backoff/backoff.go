package backoff

import (
	"math"
	"math/rand"
	"time"
)

// BackoffType 定义递增曲线的类型
type BackoffType int

const (
	Exponential BackoffType = iota // 指数递增
	Linear                         // 线性递增
	Logarithmic                    // 对数递增
	Randomized                     // 随机化递增
	Step                           // 分段递增
)

// Backoff 封装了递增睡眠时间的逻辑
type Backoff struct {
	Type          BackoffType   `json:"type" yaml:"type" ini:"type"`                               // 递增曲线类型
	InitialSleep  time.Duration `json:"initial_sleep" yaml:"initial_sleep" ini:"initial_sleep"`    // 初始等待时间
	MaxSleep      time.Duration `json:"max_sleep" yaml:"max_sleep" ini:"max_sleep"`                // 最大等待时间
	Step          time.Duration `json:"step" yaml:"step" ini:"step"`                               // 线性递增的步长（仅用于 Linear 和 Step 类型）
	Base          float64       `json:"base" yaml:"base" ini:"base"`                               // 指数递增的基数（仅用于 Exponential 类型）
	StepThreshold int           `json:"step_threshold" yaml:"step_threshold" ini:"step_threshold"` // 分段递增的阈值（仅用于 Step 类型）
	currentSleep  time.Duration // 当前等待时间
	failureCount  int           // 失败次数
}

// NewBackoff 创建一个新的 Backoff 实例
func NewBackoff(backoffType BackoffType, initialSleep, maxSleep time.Duration, options ...any) *Backoff {
	b := &Backoff{
		Type:         backoffType,
		InitialSleep: initialSleep,
		MaxSleep:     maxSleep,
		currentSleep: initialSleep,
	}

	// 设置可选参数
	for _, option := range options {
		switch v := option.(type) {
		case time.Duration:
			b.Step = v // 设置步长
		case float64:
			b.Base = v // 设置基数
		case int:
			b.StepThreshold = v // 设置分段阈值
		}
	}

	return b
}

// Next 返回下一次的等待时间，并更新状态
func (b *Backoff) Next() time.Duration {
	// 如果当前等待时间已经达到最大值，直接返回
	if b.currentSleep >= b.MaxSleep {
		return b.MaxSleep
	}
	switch b.Type {
	case Exponential:
		b.currentSleep = time.Duration(float64(b.InitialSleep) * math.Pow(b.Base, float64(b.failureCount)))
	case Linear:
		b.currentSleep = b.InitialSleep + time.Duration(b.failureCount)*b.Step
	case Logarithmic:
		b.currentSleep = time.Duration(float64(b.InitialSleep) * math.Log2(float64(b.failureCount+2)))
	case Randomized:
		baseSleep := time.Duration(float64(b.InitialSleep) * math.Pow(b.Base, float64(b.failureCount)))
		jitter := time.Duration(rand.Int63n(int64(baseSleep))) // 随机抖动
		b.currentSleep = baseSleep + jitter
	case Step:
		if b.failureCount < b.StepThreshold {
			b.currentSleep = time.Duration(float64(b.InitialSleep) * math.Pow(b.Base, float64(b.failureCount)))
		} else {
			b.currentSleep = b.InitialSleep + time.Duration(b.failureCount)*b.Step
		}
	}

	// 确保不超过最大等待时间
	if b.currentSleep > b.MaxSleep {
		b.currentSleep = b.MaxSleep
	}

	b.failureCount++
	return b.currentSleep
}

// Reset 重置 Backoff 状态
func (b *Backoff) Reset() {
	b.currentSleep = b.InitialSleep
	b.failureCount = 0
}
