// Package calculator 实现预计算表索引计算
package calculator

// precomputedCalculator 预计算余数表计算器
type precomputedCalculator struct {
	size      int
	nextTable []int         // next操作预计算表
	prevTable []int         // prev操作预计算表
	subTables map[int][]int // 常用sub操作的预计算表
}

// NewPrecomputedCalculator 创建预计算表计算器
func NewPrecomputedCalculator(size int) IndexCalculator {
	// 常用sub操作的n值(可根据实际场景调整)
	commonSubValues := []int{1, 2, 3, 4, 5, 8, 10, 16, 20, 32}

	pc := &precomputedCalculator{
		size:      size,
		nextTable: make([]int, size),
		prevTable: make([]int, size),
		subTables: make(map[int][]int),
	}

	// 预计算next表
	for i := 0; i < size; i++ {
		pc.nextTable[i] = (i + 1) % size
	}

	// 预计算prev表
	for i := 0; i < size; i++ {
		pc.prevTable[i] = (i - 1 + size) % size
	}

	// 预计算常用sub表
	for _, n := range commonSubValues {
		table := make([]int, size)
		for i := 0; i < size; i++ {
			table[i] = (i - n + size) % size
		}
		pc.subTables[n] = table
	}

	return pc
}

func (p *precomputedCalculator) Next(i int) int {
	return p.nextTable[i]
}

func (p *precomputedCalculator) Prev(i int) int {
	return p.prevTable[i]
}

func (p *precomputedCalculator) Sub(i, n int) int {
	if table, ok := p.subTables[n]; ok {
		return table[i]
	}
	return (i - n + p.size) % p.size
}
