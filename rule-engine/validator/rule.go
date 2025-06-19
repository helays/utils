package validator

import (
	"fmt"
	"github.com/helays/utils/rule-engine/validator/operators"
	"github.com/helays/utils/rule-engine/validator/types"
	"reflect"
	"strconv"
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
func (r *Rule) Validate(data map[string]any) (string, bool) {
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
		finalMsg    string
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

// computeLogicResult 计算逻辑组合结果
func computeLogicResult(item *stackItem) (string, bool) {
	if len(item.msgs) == 0 {
		return "", true
	}

	switch item.rule.Logic {
	case LogicAnd:
		return fmt.Sprintf("AND条件失败: %s", strings.Join(item.msgs, "; ")), false
	case LogicOr:
		if len(item.msgs) < len(item.rule.Conditions) {
			return "", true // 至少有一个成功
		}
		return fmt.Sprintf("OR条件失败: %s", strings.Join(item.msgs, "; ")), false
	default:
		return "未知逻辑类型", false
	}
}

// 执行单条件校验
func (r *Rule) validateCondition(rule *Rule, data map[string]any) (string, bool) {
	if rule.Field == "" {
		return "字段路径为空", false
	}

	if strings.Contains(rule.Field, "*") {
		return r.validateWildcard(rule, data)
	}
	return r.validateSimple(rule, data[rule.Field])
}

// 简单条件校验
func (r *Rule) validateSimple(rule *Rule, value any) (string, bool) {
	switch rule.Category {
	case types.CategoryDataType: // 数据类型校验
		if msg, ok := validParams(rule, types.DataTypeAttributes); !ok {
			return msg, ok
		}
		tv, msg, ok := operators.ValidateType(rule.Operator, value, rule.Value)
		if ok && rule.FieldDataType == "date" {
			// 当校验成功后，并且字段是日期来类型，需要将转换后的结果缓存下来，后续比较可以继续用
			r.dataFormatCache.Store(rule.Field, tv)
		}
		return msg, ok

	case types.CategoryLength: // 长度校验
		if msg, ok := validParams(rule, types.LengthAttributes); !ok {
			return msg, ok
		}
		return operators.ValidateLength(rule.Operator, value, rule.Value)

	case types.CategoryFormat: // 格式校验
		if msg, ok := validParams(rule, types.FormatAttributes); !ok {
			return msg, ok
		}
		return operators.ValidateFormat(rule.Operator, value, rule.Value)

	case types.CategoryContent: // 内容校验
		if msg, ok := validParams(rule, types.ContentAttributes); !ok {
			return msg, ok
		}
		// 先从 cache中获取是否有数据
		tv, ok := r.dataFormatCache.Load(rule.Field)
		if !ok {
			tv = nil
		}
		return operators.ValidateContent(rule.Operator, rule.FieldDataType, value, tv, rule.Value)

	case types.CategoryAdvanced: // 高级校验
		if msg, ok := validParams(rule, types.AdvancedAttributes); !ok {
			return msg, ok
		}
		return operators.ValidateAdvanced(rule.Operator, value, rule.Value)
	default:
		return fmt.Sprintf("未知数据类型校验操作符：%s", rule.Operator), false
	}
}

func validParams(rule *Rule, attrs map[types.Operator]types.Attribute) (string, bool) {
	attr, ok := attrs[rule.Operator]
	if !ok {
		return fmt.Sprintf("未知%s：%s", rule.Category, rule.Operator), false
	}
	if attr.ValueNum == 0 {
		return "", true
	} else if (attr.ValueNum < 0 && len(rule.Value) < 1) || (attr.ValueNum != len(rule.Value)) {
		// ValueNum<0 表示不固定参数，比较基准至少的有一个
		// 否则必须和ValueNum相等
		return fmt.Sprintf("%s校验基准参数数量错误：%d", attr.Title, len(rule.Value)), false
	}
	return "", true
}

// validateWildcard 通配符条件校验
func (r *Rule) validateWildcard(rule *Rule, data map[string]any) (string, bool) {
	regex := GetGlobalCache().Get(rule.Field)
	var errorMsgs []string

	for key, value := range data {
		if regex.MatchString(key) {
			if msg, ok := r.validateSimple(rule, value); !ok {
				errorMsgs = append(errorMsgs, msg)
			} else if rule.Logic == LogicOr {
				return "", true // OR逻辑下任意成功即返回
			}
		}
	}

	if len(errorMsgs) > 0 {
		return fmt.Sprintf("通配符[%s]校验失败: %s", rule.Field, strings.Join(errorMsgs, "; ")), false
	}
	return fmt.Sprintf("未匹配到通配符字段[%s]", rule.Field), false
}

// toString 高效类型转换
func toString(v any) string {
	switch val := v.(type) {
	case string:
		return val
	case int, int8, int16, int32, int64:
		return strconv.FormatInt(reflect.ValueOf(val).Int(), 10)
	case uint, uint8, uint16, uint32, uint64:
		return strconv.FormatUint(reflect.ValueOf(val).Uint(), 10)
	case float32, float64:
		return strconv.FormatFloat(reflect.ValueOf(val).Float(), 'f', -1, 64)
	case bool:
		return strconv.FormatBool(val)
	default:
		return fmt.Sprintf("%v", v)
	}
}

// compareGreater 比较大小
func compareGreater(a, b string) bool {
	// 尝试解析为数字
	floatA, errA := strconv.ParseFloat(a, 64)
	floatB, errB := strconv.ParseFloat(b, 64)
	if errA == nil && errB == nil {
		return floatA > floatB
	}
	// 回退到字符串比较
	return a > b
}
