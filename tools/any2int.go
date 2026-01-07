package tools

import (
	"encoding/json"
	"fmt"
	"math"
	"reflect"
	"strconv"
	"strings"

	"golang.org/x/exp/constraints"
)

func MustAny2Int[T constraints.Integer](_v any) T {
	v, _ := Any2Int[T](_v)
	return v
}

func Any2Int[T constraints.Integer](_v any) (T, error) {
	var zero T
	switch v := _v.(type) {
	case nil:
		return zero, nil
	case int:
		return T(v), nil
	case int8:
		return T(v), nil
	case int16:
		return T(v), nil
	case int32:
		return T(v), nil
	case int64:
		return T(v), nil
	case uint:
		return T(v), nil
	case uint8:
		return T(v), nil
	case uint16:
		return T(v), nil
	case uint32:
		return T(v), nil
	case uint64:
		return T(v), nil
	case string:
		return String2Int[T](v)
	case json.Number:
		if v == "" {
			return 0, nil
		}
		return String2Int[T](string(v))
	case float32:
		return numericConvert[T](v)
	case float64:
		return numericConvert[T](v)
	case bool:
		if v {
			return 1, nil
		}
		return zero, nil
	default:
		// 使用反射尝试获取基础数值
		val := reflect.ValueOf(v)
		switch val.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			return Any2Int[T](val.Int())
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			return Any2Int[T](val.Uint())
		case reflect.Float32, reflect.Float64:
			return Any2Int[T](val.Float())
		case reflect.Bool:
			if val.Bool() {
				return 1, nil
			}
			return zero, nil
		default:
			return zero, fmt.Errorf("无法将类型 %T 转换为 int", _v)
		}
	}
}

func String2Int[T constraints.Integer](v string) (T, error) {
	var zero T
	v = TrimStringHelper(v)
	if v == "" || v == "null" || v == "nil" || v == "undefined" {
		return zero, nil
	}

	if v == "true" {
		return T(1), nil
	} else if v == "false" {
		return zero, nil
	}

	// 去除千分位逗号
	v = strings.ReplaceAll(v, ",", "")

	// 处理科学计数法
	if strings.ContainsAny(v, "eE") {
		// 根据目标类型决定解析方式
		switch any(zero).(type) {
		case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
			// 先尝试用浮点数解析，再转整型
			if f, err := strconv.ParseFloat(v, 64); err == nil {
				// 检查是否超出目标类型的范围
				return numericConvert[T](f)
			}
		}
	}

	if len(v) > 2 {
		prefix := v[:2]
		var base int
		switch prefix {
		case "0x", "0X": // 十六进制
			base = 16
			v = v[2:]
		case "0b", "0B": // 二进制
			base = 2
			v = v[2:]
		case "0o", "0O": // 八进制
			base = 8
			v = v[2:]
		default:
			base = 0
		}
		if base != 0 {
			// 根据目标类型进行解析
			return parseWithBase[T](v, base)
		}
	}
	// 包含小数点
	if strings.Contains(v, ".") {
		parts := strings.SplitN(v, ".", 2)
		if parts[0] == "" && parts[1] == "" {
			return zero, fmt.Errorf("invalid number: %s", v)
		}
		if parts[0] == "" {
			return zero, nil // ".123" 的情况返回 0
		}
		v = parts[0]
	}

	// 默认按十进制解析
	return parseWithBase[T](v, 10)
}

// numericConvert 将数值转换为目标类型，检查溢出
func numericConvert[T constraints.Integer, F constraints.Float](f F) (T, error) {
	var zero T
	switch any(zero).(type) {
	case int:
		if f > math.MaxInt || f < math.MinInt {
			return zero, strconv.ErrRange
		}
		return T(f), nil
	case int8:
		if f > math.MaxInt8 || f < math.MinInt8 {
			return zero, strconv.ErrRange
		}
		return T(f), nil
	case int16:
		if f > math.MaxInt16 || f < math.MinInt16 {
			return zero, strconv.ErrRange
		}
		return T(f), nil
	case int32:
		if f > math.MaxInt32 || f < math.MinInt32 {
			return zero, strconv.ErrRange
		}
		return T(f), nil
	case int64:

		if f > math.MaxInt64 || f < math.MinInt64 {
			return zero, strconv.ErrRange
		}
		return T(f), nil

	case uint:
		if f > math.MaxUint || f < 0 {
			return zero, strconv.ErrRange
		}
		return T(f), nil

	case uint8:
		if f > math.MaxUint8 || f < 0 {
			return zero, strconv.ErrRange
		}
		return T(f), nil

	case uint16:
		if f > math.MaxUint16 || f < 0 {
			return zero, strconv.ErrRange
		}
		return T(f), nil

	case uint32:
		if f > math.MaxUint32 || f < 0 {
			return zero, strconv.ErrRange
		}
		return T(f), nil

	case uint64:
		if f > math.MaxUint64 || f < 0 {
			return zero, strconv.ErrRange
		}
		return T(f), nil
	default:
		return T(f), nil
	}
}

// parseWithBase 根据进制解析字符串到指定整数类型
func parseWithBase[T constraints.Integer](v string, base int) (T, error) {
	var zero T

	switch any(zero).(type) {
	case int:
		val, err := strconv.ParseInt(v, base, strconv.IntSize)
		if err != nil {
			return zero, err
		}
		// 检查是否超出 int 范围
		if val > int64(^uint(0)>>1) || val < -int64(^uint(0)>>1)-1 {
			return zero, &strconv.NumError{
				Func: "ParseInt",
				Num:  v,
				Err:  strconv.ErrRange,
			}
		}
		return T(val), nil

	case int8:
		val, err := strconv.ParseInt(v, base, 8)
		return T(val), err

	case int16:
		val, err := strconv.ParseInt(v, base, 16)
		return T(val), err

	case int32:
		val, err := strconv.ParseInt(v, base, 32)
		return T(val), err

	case int64:
		val, err := strconv.ParseInt(v, base, 64)
		return T(val), err

	case uint:
		val, err := strconv.ParseUint(v, base, strconv.IntSize)
		return T(val), err

	case uint8:
		val, err := strconv.ParseUint(v, base, 8)
		return T(val), err

	case uint16:
		val, err := strconv.ParseUint(v, base, 16)
		return T(val), err

	case uint32:
		val, err := strconv.ParseUint(v, base, 32)
		return T(val), err

	case uint64:
		val, err := strconv.ParseUint(v, base, 64)
		return T(val), err

	default:
		// 对于自定义类型，尝试使用 int64
		val, err := strconv.ParseInt(v, base, 64)
		return T(val), err
	}
}
