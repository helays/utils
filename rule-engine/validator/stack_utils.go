package validator

import (
	"helay.net/go/utils/v3/rule-engine/validator/types"
	"sync"
)

// stackItem 栈元素定义
type stackItem struct {
	rule   *Rule
	result bool
	index  int                      // 当前处理到的conditions索引
	msgs   []*types.ValidationError // 收集的错误信息
}

// stack 栈结构定义
type stack struct {
	items []*stackItem
}

var stackItemPool = sync.Pool{
	New: func() interface{} {
		return &stackItem{
			msgs: make([]*types.ValidationError, 0, 4), // 预分配错误消息容量
		}
	},
}

// newStack 创建新栈（初始化容量）
func newStack() *stack {
	return &stack{
		items: make([]*stackItem, 0, 8), // 初始容量8
	}
}

// getStackItem 从对象池获取stackItem
func getStackItem() *stackItem {
	return stackItemPool.Get().(*stackItem)
}

// putStackItem 将stackItem放回对象池
func putStackItem(item *stackItem) {
	item.rule = nil
	item.index = 0
	item.result = false
	item.msgs = item.msgs[:0]
	stackItemPool.Put(item)
}

func (s *stack) push(item *stackItem) {
	s.items = append(s.items, item)
}

func (s *stack) pop() *stackItem {
	if len(s.items) == 0 {
		return nil
	}
	item := s.items[len(s.items)-1]
	s.items = s.items[:len(s.items)-1]
	return item
}

func (s *stack) peek() *stackItem {
	if len(s.items) == 0 {
		return nil
	}
	return s.items[len(s.items)-1]
}

func (s *stack) len() int {
	return len(s.items)
}
