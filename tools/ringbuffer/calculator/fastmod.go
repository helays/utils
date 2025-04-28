// Package calculator 实现快速条件判断索引计算
package calculator

// fastModCalculator 快速条件判断计算器
type fastModCalculator struct {
	size int
}

// NewFastModCalculator 创建快速条件判断计算器
func NewFastModCalculator(size int) IndexCalculator {
	return &fastModCalculator{size: size}
}

func (f *fastModCalculator) Next(i int) int {
	i++
	if i >= f.size {
		return 0
	}
	return i
}

func (f *fastModCalculator) Prev(i int) int {
	if i > 0 {
		return i - 1
	}
	return f.size - 1
}

func (f *fastModCalculator) Sub(i, n int) int {
	i -= n
	if i >= 0 {
		return i
	}
	return i + f.size
}
