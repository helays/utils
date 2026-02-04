package operators

import (
	"fmt"
	"github.com/araddon/dateparse"
	"helay.net/go/utils/v3/rule-engine/validator/types"
	"helay.net/go/utils/v3/tools"
	"regexp"
	"strings"
	"time"
)

func ValidateContent(operator types.Operator, dataType string, value, valued any, rule []any) (string, bool) {
	switch operator {
	case types.Required:
		if value == nil || value == "" {
			return types.ContentChineseMap[operator], false
		}
	case types.NotBlank:
		if strings.TrimSpace(tools.Any2string(value)) == "" {
			return types.ContentChineseMap[operator], false
		}
	case types.InEnum:
		if !searchSlice(valued, rule) {
			return fmt.Sprintf(types.ContentChineseMap[operator], rule), false
		}
	case types.NotInEnum:
		if searchSlice(valued, rule) {
			return fmt.Sprintf(types.ContentChineseMap[operator], rule), false
		}
	case types.GreaterThan: // 大于
		ret, ok := compare(dataType, value, valued, rule[0])
		if !ok {
			return fmt.Sprintf("数据无法比较，入参%v 基准%v", value, rule[0]), false
		} else if ret <= 0 {
			return fmt.Sprintf(types.ContentChineseMap[operator], rule[0]), false
		}
	case types.GreaterEqual: // 大于等于
		ret, ok := compare(dataType, value, valued, rule[0])
		if !ok {
			return fmt.Sprintf("数据无法比较，入参%v 基准%v", value, rule[0]), false
		} else if ret < 0 {
			return fmt.Sprintf(types.ContentChineseMap[operator], rule[0]), false
		}
	case types.LessThan: // 小于
		ret, ok := compare(dataType, value, valued, rule[0])
		if !ok {
			return fmt.Sprintf("数据无法比较，入参%v 基准%v", value, rule[0]), false
		} else if ret >= 0 {
			return fmt.Sprintf(types.ContentChineseMap[operator], rule[0]), false
		}
	case types.LessEqual: // 小于等于
		ret, ok := compare(dataType, value, valued, rule[0])
		if !ok {
			return fmt.Sprintf("数据无法比较，入参%v 基准%v", value, rule[0]), false
		} else if ret > 0 {
			return fmt.Sprintf(types.ContentChineseMap[operator], rule[0]), false
		}
	case types.Equal: // 等于
		ret, ok := compare(dataType, value, valued, rule[0])
		if !ok {
			return fmt.Sprintf("数据无法比较，入参%v 基准%v", value, rule[0]), false
		} else if ret != 0 {
			return fmt.Sprintf(types.ContentChineseMap[operator], rule[0]), false
		}
	case types.NotEqual: // 不等于
		ret, ok := compare(dataType, value, valued, rule[0])
		if !ok {
			return fmt.Sprintf("数据无法比较，入参%v 基准%v", value, rule[0]), false
		} else if ret == 0 {
			return fmt.Sprintf(types.ContentChineseMap[operator], rule[0]), false
		}
	case types.RegexMatch: // 正则匹配
		matched, _ := regexp.MatchString(tools.Any2string(rule[0]), tools.Any2string(value))
		if !matched {
			return fmt.Sprintf(types.ContentChineseMap[operator], rule[0]), false
		}
	}
	return "", true
}

func searchSlice(inp any, o []any) bool {
	_s := tools.Any2string(inp)
	for _, i := range o {
		if tools.Any2string(i) == _s {
			return true
		}
	}
	return false
}

// compare 比较数字值，支持int, float等类型
func compare(dataType string, a, aed, b any) (int, bool) {
	switch dataType {
	case "string":
		return compareString(a, b), true
	case "int", "float", "byte":
		return compareNumbers(a, b)
	case "bool":
		return compareBool(a, b)
	case "date":
		baseB, err := dateparse.ParseLocal(tools.Any2string(b))
		if err != nil {
			return 0, false
		}
		if aed == nil {
			aed, err = dateparse.ParseLocal(tools.Any2string(a))
			if err != nil {
				return 0, false
			}
		}
		return compareTime(aed.(time.Time), baseB), true
	default:
		return 0, false
	}
}

func compareNumbers(a, b any) (int, bool) {
	// 尝试将a和b转换为float64进行比较
	af, aErr := tools.Any2float64(a)
	bf, bErr := tools.Any2float64(b)
	if aErr != nil || bErr != nil {
		return 0, false // 转换失败，视为相等
	}

	if af > bf {
		return 1, true
	} else if af < bf {
		return -1, true
	}
	return 0, true
}

func compareString(a, b any) int {
	return strings.Compare(tools.Any2string(a), tools.Any2string(b))
}

func compareTime(a, b time.Time) int {
	if a.After(b) {
		return 1
	} else if a.Before(b) {
		return -1
	}
	return 0
}

func compareBool(a, b any) (int, bool) {
	_a, err := tools.Any2bool(a)
	if err != nil {
		return 0, false
	}
	_b, err := tools.Any2bool(b)
	if err != nil {
		return 0, false
	}
	if _a == _b {
		return 0, true
	}
	if _a {
		return 1, true
	}
	return -1, true
}
