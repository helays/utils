package validator

import (
	"fmt"
	"helay.net/go/utils/v3/rule-engine/validator/types"
)

func (r *Rule) CheckRule(dataTypeMap map[string]string) error {
	if r == nil {
		return nil
	}

	// 使用栈来实现DFS遍历
	stackRule := []*Rule{r}
	for len(stackRule) > 0 {
		// 弹出栈顶元素
		current := stackRule[len(stackRule)-1]
		stackRule = stackRule[:len(stackRule)-1]

		// 检查当前节点
		if err := current.checkRule(); err != nil {
			return err
		}
		if current.Logic == "" {
			// 条件类型节点 - 补充FieldDataType并校验
			if current.Field == "" {
				return fmt.Errorf("条件节点缺少 'field' 属性")
			}
			// 从提供的map中获取数据类型
			dataType, exists := dataTypeMap[current.Field]
			if !exists {
				return fmt.Errorf("字段'%s'未在数据类型映射中找到", current.Field)
			}
			// 补充FieldDataType
			current.FieldDataType = dataType
		} else {
			// 将子节点压入栈中（逆序以保证原顺序）
			for i := len(current.Conditions) - 1; i >= 0; i-- {
				stackRule = append(stackRule, current.Conditions[i])
			}
		}
	}

	return nil
}

func (r *Rule) checkRule() error {
	if r.Logic != "" {
		return r.isLogic()
	}
	for _, category := range types.CategoryChineseNames {
		if category.Category == r.Category {
			_, ok := category.List[r.Operator]
			if ok {
				return nil
			}
			return fmt.Errorf("无效的操作符: %s", r.Operator)
		}
	}
	return fmt.Errorf("无效的分类: %s", r.Category)
}

func (l *Rule) isLogic() error {
	if l.Logic == LogicAnd || l.Logic == LogicOr {
		return nil
	}
	return fmt.Errorf("无效的逻辑类型: %s", l.Logic)
}
