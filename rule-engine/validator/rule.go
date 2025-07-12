package validator

import (
	"fmt"
	"github.com/helays/utils/v2/rule-engine/validator/operators"
	"github.com/helays/utils/v2/rule-engine/validator/types"
	"strings"
	"sync"
)

type Logic string // 逻辑类型

const (
	LogicAnd Logic = "and"
	LogicOr  Logic = "or"
)

// Rule 定义校验规则
type Rule struct {
	// 逻辑类型
	Logic      Logic   `json:"logic,omitempty"` // 有逻辑类型时，条件类型无效
	Conditions []*Rule `json:"conditions,omitempty"`

	// 条件类型
	Field           string                   `json:"field,omitempty"`           // 字段名称
	FieldDataType   string                   `json:"field_data_type,omitempty"` // 字段数据类型
	Category        types.ValidationCategory `json:"category,omitempty"`        // 校验类型分类
	Operator        types.Operator           `json:"operator,omitempty"`
	Value           []any                    `json:"value,omitempty"`
	dataFormatCache sync.Map                 // 在数据类型校验过程中，缓存数据校验结果数据
}

// Validate 主入口方法
func (r *Rule) Validate(data map[string]any) (*types.ValidationError, bool) {
	s := newStack()
	defer func() {
		// 确保栈清空时回收所有对象
		for s.len() > 0 {
			if item := s.pop(); item != nil {
				putStackItem(item)
			}
		}
	}()

	// 初始化根节点
	rootItem := getStackItem()
	rootItem.rule = r
	s.push(rootItem)

	var (
		finalMsg    *types.ValidationError
		finalResult bool
	)

	for s.len() > 0 {
		current := s.peek()
		// 1. 处理叶子节点（基础条件）
		if current.rule.Logic == "" {
			msg, ok := r.validateCondition(current.rule, data)
			s.pop()
			putStackItem(current) // 立即回收
			if parent := s.peek(); parent != nil {
				// 短路优化：根据子结果决定是否跳过后续条件
				if stopEarly := updateParentResult(parent, ok); stopEarly {
					parent.index = len(parent.rule.Conditions)
				}
				if !ok {
					parent.msgs = append(parent.msgs, msg)
				}
			} else {
				finalMsg, finalResult = msg, ok
			}
			continue
		}

		// 2. 处理当前规则的子条件
		if current.index < len(current.rule.Conditions) {
			child := current.rule.Conditions[current.index]
			current.index++

			childItem := getStackItem()
			childItem.rule = child
			s.push(childItem)
			continue
		}

		// 3. 所有子条件处理完毕
		msg, ok := computeLogicResult(current)
		s.pop()
		putStackItem(current) // 立即回收

		if parent := s.peek(); parent != nil {
			parent.result = updateParentResult(parent, ok)
			if !ok {
				parent.msgs = append(parent.msgs, msg)
			}
		} else {
			finalMsg, finalResult = msg, ok
		}
	}
	if finalMsg != nil {
		finalMsg.Field = r.Field
	}
	return finalMsg, finalResult
}

// updateParentResult 更新父节点结果并返回是否需要短路
func updateParentResult(parent *stackItem, childResult bool) (stopEarly bool) {
	switch parent.rule.Logic {
	case LogicAnd:
		parent.result = parent.result && childResult
		return !childResult // AND遇到false时短路
	case LogicOr:
		parent.result = parent.result || childResult
		return childResult // OR遇到true时短路
	default:
		return false
	}
}

// 计算逻辑组合结果
func computeLogicResult(item *stackItem) (*types.ValidationError, bool) {
	if len(item.msgs) == 0 {
		return nil, true
	}

	switch item.rule.Logic {
	case LogicAnd:
		return item.msgs[0], false
	case LogicOr:
		if len(item.msgs) < len(item.rule.Conditions) {
			return nil, true // 至少有一个成功
		}
		return item.msgs[0], false
	default:
		return item.rule.setErr(nil, fmt.Sprintf("逻辑组合类型%s错误", item.rule.Logic)), false
	}
}

// 执行单条件校验
func (r *Rule) validateCondition(rule *Rule, data map[string]any) (*types.ValidationError, bool) {
	if rule.Field == "" {
		return rule.setErr(nil, "字段路径为空"), false
	}

	if strings.Contains(rule.Field, "*") {
		return r.validateWildcard(rule, data)
	}
	return r.validateSimple(rule, data[rule.Field])
}

// 简单条件校验
func (r *Rule) validateSimple(rule *Rule, value any) (*types.ValidationError, bool) {
	var (
		msg string
		ok  bool
		tv  any
	)
	switch rule.Category {
	case types.CategoryDataType: // 数据类型校验
		if msg, ok = validParams(rule, types.DataTypeAttributes); ok {
			if tv, msg, ok = operators.ValidateType(rule.Operator, value, rule.Value); ok {
				if rule.FieldDataType == "date" {
					// 当校验成功后，并且字段是日期来类型，需要将转换后的结果缓存下来，后续比较可以继续用
					r.dataFormatCache.Store(rule.Field, tv)
				}
				return nil, ok
			}
		}
	case types.CategoryLength: // 长度校验
		if msg, ok = validParams(rule, types.LengthAttributes); ok {
			if msg, ok = operators.ValidateLength(rule.Operator, value, rule.Value); ok {
				return nil, ok
			}
		}
	case types.CategoryFormat: // 格式校验
		if msg, ok = validParams(rule, types.FormatAttributes); ok {
			if msg, ok = operators.ValidateFormat(rule.Operator, value, rule.Value); ok {
				return nil, ok
			}
		}

	case types.CategoryContent: // 内容校验
		if msg, ok = validParams(rule, types.ContentAttributes); ok {
			// 先从 cache中获取是否有数据
			if tv, ok = r.dataFormatCache.Load(rule.Field); !ok {
				tv = nil
			}
			if msg, ok = operators.ValidateContent(rule.Operator, rule.FieldDataType, value, tv, rule.Value); ok {
				return nil, ok
			}
		}

	case types.CategoryAdvanced: // 高级校验
		if msg, ok = validParams(rule, types.AdvancedAttributes); ok {
			if msg, ok = operators.ValidateAdvanced(rule.Operator, value, rule.Value); ok {
				return nil, ok
			}
		}
	default:
		return rule.setErr(value, fmt.Sprintf("规则字段%s未知数据类型校验操作符：%s", rule.Field, rule.Operator)), false
	}
	return rule.setErr(value, fmt.Sprintf("规则字段%s%s", rule.Field, msg)), false
}

func validParams(rule *Rule, attrs map[types.Operator]types.Attribute) (string, bool) {
	attr, ok := attrs[rule.Operator]
	if !ok {
		return fmt.Sprintf("配置了未知条件运算类型【%s】", rule.Operator), false
	}
	if attr.ValueNum == 0 {
		return "", true
	} else if (attr.ValueNum < 0 && len(rule.Value) < 1) || (attr.ValueNum != len(rule.Value)) {
		// ValueNum<0 表示不固定参数，比较基准至少的有一个
		// 否则必须和ValueNum相等
		return fmt.Sprintf("%s校验基准参数数量错误【%d】", attr.Title, len(rule.Value)), false
	}
	return "", true
}

// validateWildcard 通配符条件校验
func (r *Rule) validateWildcard(rule *Rule, data map[string]any) (*types.ValidationError, bool) {
	regex := GetGlobalCache().Get(rule.Field)
	var errorMsgs []*types.ValidationError
	matched := false
	for key, value := range data {
		if regex.MatchString(key) {
			matched = true
			msg, ok := r.validateSimple(rule, value)
			if !ok {
				errorMsgs = append(errorMsgs, msg)
			}
		}
	}
	if !matched {
		return rule.setErr(nil, fmt.Sprintf("规则字段%s无法匹配到被校验数据", rule.Field)), false
	}

	if len(errorMsgs) > 0 {
		return errorMsgs[0], false
	}
	return nil, true
}

func (r *Rule) setErr(v any, m string) *types.ValidationError {
	return &types.ValidationError{
		RealField:  r.Field,
		InputValue: v,
		RuleValue:  r.Value,
		Operator:   r.Operator,
		Category:   r.Category,
		Message:    m,
	}
}
