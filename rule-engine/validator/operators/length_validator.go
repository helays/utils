package operators

import (
	"fmt"

	"helay.net/go/utils/v3/rule-engine/validator/types"
	"helay.net/go/utils/v3/tools"
)

func ValidateLength(operator types.Operator, value any, rule []any) (string, bool) {
	str := tools.Any2string(value)
	length := int64(len(str))
	switch operator {
	case types.LenMin:
		_min, err := tools.Any2Int[int64](rule[0])
		if err != nil {
			return fmt.Sprintf("长度校验操作%s失败，参数%v不是数字", operator, rule[0]), false
		}
		if length < _min {
			return fmt.Sprintf(types.LengthChineseMap[operator], _min), false
		}
	case types.LenMax:
		_max, err := tools.Any2Int[int64](rule[0])
		if err != nil {
			return fmt.Sprintf("长度校验操作%s失败，参数%v不是数字", operator, rule[0]), false
		}
		if length > _max {
			return fmt.Sprintf(types.LengthChineseMap[operator], _max), false
		}
	case types.LenEquals:
		_equals, err := tools.Any2Int[int64](rule[0])
		if err != nil {
			return fmt.Sprintf("长度校验操作%s失败，参数%v不是数字", operator, rule[0]), false
		}
		if length != _equals {
			return fmt.Sprintf(types.LengthChineseMap[operator], _equals), false
		}
	case types.LenRange:
		_min, err := tools.Any2Int[int64](rule[0])
		if err != nil {
			return fmt.Sprintf("长度校验操作%s失败，参数%v不是数字", operator, rule[0]), false
		}
		_max, err := tools.Any2Int[int64](rule[1])
		if err != nil {
			return fmt.Sprintf("长度校验操作%s失败，参数%v不是数字", operator, rule[1]), false
		}
		if length < _min || length > _max {
			return fmt.Sprintf(types.LengthChineseMap[operator], _min, _max), false
		}
	}
	return "", true
}
