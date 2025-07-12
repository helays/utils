package operators

import (
	"fmt"
	"github.com/helays/utils/v2/rule-engine/validator/types"
	"github.com/helays/utils/v2/tools"
)

func ValidateLength(operator types.Operator, value any, rule []any) (string, bool) {
	str := tools.Any2string(value)
	length := len(str)
	switch operator {
	case types.LenMin:
		_min, err := tools.Any2int(rule[0])
		if err != nil {
			return fmt.Sprintf("长度校验操作%s失败，参数%v不是数字", operator, rule[0]), false
		}
		if length < int(_min) {
			return fmt.Sprintf(types.LengthChineseMap[operator], _min), false
		}
	case types.LenMax:
		_max, err := tools.Any2int(rule[0])
		if err != nil {
			return fmt.Sprintf("长度校验操作%s失败，参数%v不是数字", operator, rule[0]), false
		}
		if length > int(_max) {
			return fmt.Sprintf(types.LengthChineseMap[operator], _max), false
		}
	case types.LenEquals:
		_equals, err := tools.Any2int(rule[0])
		if err != nil {
			return fmt.Sprintf("长度校验操作%s失败，参数%v不是数字", operator, rule[0]), false
		}
		if length != int(_equals) {
			return fmt.Sprintf(types.LengthChineseMap[operator], _equals), false
		}
	case types.LenRange:
		_min, err := tools.Any2int(rule[0])
		if err != nil {
			return fmt.Sprintf("长度校验操作%s失败，参数%v不是数字", operator, rule[0]), false
		}
		_max, err := tools.Any2int(rule[1])
		if err != nil {
			return fmt.Sprintf("长度校验操作%s失败，参数%v不是数字", operator, rule[1]), false
		}
		if length < int(_min) || length > int(_max) {
			return fmt.Sprintf(types.LengthChineseMap[operator], _min, _max), false
		}
	}
	return "", true
}
