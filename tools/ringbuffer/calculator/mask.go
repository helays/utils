// Package calculator 实现位掩码索引计算
package calculator

// maskCalculator 位掩码计算器(适用于2的幂大小)
type maskCalculator struct {
	mask int
}

// NewMaskCalculator 创建位掩码计算器
func NewMaskCalculator(size int) IndexCalculator {
	return &maskCalculator{
		mask: size - 1, // 例如size=8, mask=7(0b0111)
	}
}

func (m *maskCalculator) Next(i int) int {
	return (i + 1) & m.mask
}

func (m *maskCalculator) Prev(i int) int {
	return (i - 1) & m.mask
}

func (m *maskCalculator) Sub(i, n int) int {
	return (i - n) & m.mask
}
