package operators

import (
	"fmt"
	"github.com/araddon/dateparse"
	"github.com/helays/utils/rule-engine/validator/types"
	"github.com/helays/utils/tools"
	"strconv"
	"strings"
	"time"
)

func ValidateType(operator types.Operator, value any, rule []any) (any, string, bool) {
	var (
		ret any
		ok  bool
	)
	switch operator {
	case types.TypeString:
		ret, ok = value.(string)
	case types.TypeNumber:
		ret, ok = IsNumber(value)
	case types.TypeInteger:
		ret, ok = IsInteger(value)
	case types.TypeBoolean:
		ret, ok = value.(bool)
	case types.TypeTime:
		ret, ok = ValidateDate(value, rule)
	case types.TypeArray:
		ret, ok = value.([]any)
	case types.TypeObject:
		ret, ok = value.(map[string]any)
	case types.TypeNull:
		if value != nil {
			return value, types.DataTypeChineseMap[operator], false
		}
		return nil, "", true
	}
	if !ok {
		return ret, types.DataTypeChineseMap[operator], false
	}
	return ret, "", true
}

func IsNumber(val any) (any, bool) {
	switch _val := val.(type) {
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, float32, float64:
		return _val, true
	default:
		v, err := strconv.ParseFloat(fmt.Sprintf("%v", val), 64)
		return v, err == nil
	}
}

func IsInteger(value interface{}) (any, bool) {
	switch _val := value.(type) {
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
		return _val, true
	default:
		v, err := strconv.ParseInt(fmt.Sprintf("%v", value), 10, 64)
		return v, err == nil
	}
}

func ValidateDate(value any, rule []any) (time.Time, bool) {
	var realRule [][2]string
	for _, r := range rule {
		_r := strings.TrimSpace(tools.Any2string(r))
		if _r == "" {
			continue
		}
		arr := strings.Split(_r, ":")
		if len(arr) == 2 {
			realRule = append(realRule, [2]string{arr[0], arr[1]})
		} else if len(arr) == 1 {
			realRule = append(realRule, [2]string{"golang", arr[0]})
		} else {
			return time.Time{}, false
		}

	}
	dateStr := tools.Any2string(value)
	if len(realRule) > 0 {
		for _, r := range realRule {
			tpl := tools.ConvertTimeFormat(r[1], r[0])
			if tpl == "" {
				return time.Time{}, false
			}
			if _, err := time.Parse(tpl, dateStr); err == nil {
				return time.Time{}, false
			}
		}
		return time.Time{}, false
	}
	t, err := dateparse.ParseLocal(dateStr)
	if err != nil {
		return time.Time{}, false
	}
	return t, true
}
