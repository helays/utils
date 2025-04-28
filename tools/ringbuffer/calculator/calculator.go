// Package calculator 定义环形缓冲区索引计算策略
package calculator

import "math/bits"

// IndexCalculator 定义环形缓冲区索引计算接口
type IndexCalculator interface {
	Next(int) int     // 计算下一个索引
	Prev(int) int     // 计算上一个索引
	Sub(int, int) int // 计算前移n个位置的索引
}

// NewSmartCalculator 智能创建最优计算器
func NewSmartCalculator(size int) IndexCalculator {
	// 1. 如果是2的幂，使用位掩码
	if bits.OnesCount(uint(size)) == 1 {
		return NewMaskCalculator(size)
	}

	// 2. 小缓冲区使用预计算表（适合频繁操作的小缓冲区）
	if size <= 256 {
		return NewPrecomputedCalculator(size)
	}

	// 3. 中等大小缓冲区使用条件判断
	if size <= 1024*1024 {
		return NewFastModCalculator(size)
	}

	// 4. 超大缓冲区使用标准取模
	return NewModuloCalculator(size)
}
