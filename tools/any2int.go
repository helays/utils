package tools

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"unsafe"

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
		return T(v), nil
	case float64:
		return T(v), nil
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
		f, err := strconv.ParseFloat(v, 64)
		if err != nil {
			return zero, err
		}
		return T(f), nil
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

// parseWithBase 根据进制解析字符串到指定整数类型
func parseWithBase[T constraints.Integer](v string, base int) (T, error) {
	var zero T
	l := int(unsafe.Sizeof(zero) * 8)
	// 这里判断没问题，编译器会自动优化
	if T(0)-1 < 0 {
		val, err := strconv.ParseInt(v, base, l)
		if err != nil {
			return zero, err
		}
		return T(val), nil
	}
	val, err := strconv.ParseUint(v, base, l)
	if err != nil {
		return zero, err
	}
	return T(val), nil
}
