// Package calculator 实现标准取模索引计算
package calculator

// moduloCalculator 标准取模计算器
type moduloCalculator struct {
	size int
}

// NewModuloCalculator 创建标准取模计算器
func NewModuloCalculator(size int) IndexCalculator {
	return &moduloCalculator{
		size: size,
	}
}

func (m *moduloCalculator) Next(i int) int {
	return (i + 1) % m.size
}

func (m *moduloCalculator) Prev(i int) int {
	return (i - 1 + m.size) % m.size
}

func (m *moduloCalculator) Sub(i, n int) int {
	return (i - n + m.size) % m.size
}
