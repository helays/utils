// Package impl 定义环形缓冲区的实现接口
package impl

// Buffer 定义环形缓冲区公开接口
type Buffer[T any] interface {
	Push(item T)        // 添加单个元素
	PushAll(items []T)  // 批量添加元素
	GetAll() []T        // 获取所有元素
	GetLast(n int) []T  // 获取最后n个元素
	Clear()             // 清空缓冲区
	Len() int           // 当前元素数量
	Cap() int           // 缓冲区容量
	IsFull() bool       // 是否已满
	IsEmpty() bool      // 是否为空
	Iterator() <-chan T // 创建元素迭代器
}
