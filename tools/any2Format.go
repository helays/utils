package tools

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"reflect"
	"strconv"
	"strings"
	"time"
)

// Map2Struct 将map转换为结构体
// dst 需要传入一个变量的指针
func Map2Struct(dst any, src map[string]any, customConvert map[string]func(dst any, src map[string]any) error) error {
	var err error
	// 这里通过反射，将map转换为结构体
	val := reflect.ValueOf(dst).Elem()
	typ := val.Type()
	for i := 0; i < val.NumField(); i++ {
		field := typ.Field(i)
		// 如果这个字段是一个匿名结构体，还需要递归处理
		if field.Type.Kind() == reflect.Struct && field.Anonymous {
			// 递归处理嵌套结构体，注意这里不能使用field
			if err := Map2Struct(val.Field(i).Addr().Interface(), src, customConvert); err != nil {
				return err
			}
			continue
		}
		// 如果有自定义转换函数，则使用自定义转换函数
		if f, ok := customConvert[field.Name]; ok {
			if err = f(dst, src); err != nil {
				return fmt.Errorf("自定义转换函数%s执行失败：%v", typ.Field(i).Name, err)
			}
			continue
		}
		jsonTag := field.Tag.Get("json")
		if jsonTag == "" {
			continue
		}
		jsonTag = strings.Split(jsonTag, ",")[0]
		// 获取结构体json标签
		// 检查 map 中是否存在对应的键
		value, ok := src[jsonTag]
		if !ok {
			continue
		}
		// 设置字段的值
		fieldVal := val.Field(i)
		switch fieldVal.Kind() {
		case reflect.String:
			// 如果 value是nil
			if value == nil {
				continue
			}
			fieldVal.SetString(fmt.Sprintf("%v", value))
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			_tv, err := Any2int(value)
			if err != nil {
				return fmt.Errorf("字段%s转int失败：%v", jsonTag, err)
			}
			fieldVal.SetInt(_tv)
		case reflect.Float32, reflect.Float64:
			tv, err := Any2float64(value)
			if err != nil {
				return fmt.Errorf("字段 %s 转 float64 失败: %v", jsonTag, err)
			}
			fieldVal.SetFloat(tv)
		case reflect.Bool:
			tv, err := Any2bool(value)
			if err != nil {
				return fmt.Errorf("字段 %s 转 bool 失败: %v", jsonTag, err)
			}
			fieldVal.SetBool(tv)
		case reflect.Slice:
			if field.Type.Elem().Kind() == reflect.Uint8 { // []byte
				tv, err := Any2bytes(value)
				if err != nil {
					return fmt.Errorf("字段 %s 转 []byte 失败: %v", jsonTag, err)
				}
				fieldVal.SetBytes(tv)
			}
		case reflect.Struct, reflect.Map:
			err := json.Unmarshal(Any2Byte(value), fieldVal.Addr().Interface())
			if err != nil {
				return fmt.Errorf("字段 %s 转 struct 失败: %v", jsonTag, err)
			}
		default:
		}
	}
	return nil
}

// Any2string 将任意类型转换为字符串
func Any2string(v any) string {
	if v == nil {
		return ""
	}
	switch _v := v.(type) {
	case string:
		return _v
	case int:
		return strconv.Itoa(_v)
	case int64:
		return strconv.FormatInt(_v, 10)
	case int32:
		return strconv.FormatInt(int64(_v), 10)
	case float32:
		return Float32tostring(_v)
	case float64:
		return Float64tostring(_v)
	case bool:
		return strconv.FormatBool(_v)
	case uint:
		return strconv.FormatUint(uint64(_v), 10)
	case uint64:
		return strconv.FormatUint(_v, 10)
	case uint32:
		return strconv.FormatUint(uint64(_v), 10)
	case []byte:
		return string(_v)
	case time.Time:
		return _v.Format(time.DateTime)
	case time.Duration:
		return _v.String()
	case fmt.Stringer:
		return _v.String()
	case error:
		return _v.Error()
	case nil:
		return ""
	default:
		// 使用反射处理更多类型
		rv := reflect.ValueOf(v)
		// 处理指针类型
		if rv.Kind() == reflect.Ptr {
			if rv.IsNil() {
				return ""
			}
			rv = rv.Elem()
		}
		switch rv.Kind() {
		case reflect.Slice, reflect.Array, reflect.Map, reflect.Struct:
			_byt, _ := json.Marshal(v)
			return string(_byt)
		default:
			return fmt.Sprintf("%v", v)
		}
	}
}

// Any2int 尝试将任意类型转换为 int
func Any2int(_v any) (int64, error) {
	switch v := _v.(type) {
	case nil:
		return 0, nil
	case int:
		return int64(v), nil
	case int8:
		return int64(v), nil
	case int16:
		return int64(v), nil
	case int32:
		return int64(v), nil
	case int64:
		return v, nil
	case string:
		if v == "" {
			return 0, nil
		}
		// 包含小数点
		if strings.Contains(v, ".") {
			return strconv.ParseInt(strings.Split(v, ".")[0], 10, 64)
		}
		return strconv.ParseInt(v, 10, 64)
	case json.Number:
		if v == "" {
			return 0, nil
		}
		cache := string(v)
		// 包含小数点
		if strings.Contains(cache, ".") {
			return strconv.ParseInt(strings.Split(cache, ".")[0], 10, 64)
		}
		return strconv.ParseInt(cache, 10, 64)
	case float32:
		return int64(v), nil
	case float64:
		return int64(v), nil
	case bool:
		if v {
			return 1, nil
		}
		return 0, nil
	default:
		// 使用反射尝试获取基础整数值
		val := reflect.ValueOf(v)
		switch val.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			return val.Int(), nil
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			return int64(val.Uint()), nil
		default:
			return 0, fmt.Errorf("无法将类型 %T 转换为 int", _v)
		}
	}
}

// Any2float64 尝试将任意类型转换为 float64
func Any2float64(_v any) (float64, error) {
	switch v := _v.(type) {
	case float32:
		return float64(v), nil
	case float64:
		return v, nil
	case int8:
		return float64(v), nil
	case int16:
		return float64(v), nil
	case int32:
		return float64(v), nil
	case int64:
		return float64(v), nil
	case uint:
		return float64(v), nil
	case uint8:
		return float64(v), nil
	case uint16:
		return float64(v), nil
	case uint32:
		return float64(v), nil
	case uint64:
		return float64(v), nil
	case string:
		if v == "" {
			return 0, nil
		}
		return strconv.ParseFloat(v, 64)
	case nil:
		return 0, nil
	case json.Number:
		return v.Float64()
	case bool:
		if v {
			return 1, nil
		}
		return 0, nil
	default:
		// 使用反射处理更多类型
		rv := reflect.ValueOf(v)
		// 处理指针类型
		if rv.Kind() == reflect.Ptr {
			if rv.IsNil() {
				return 0, nil
			}
			rv = rv.Elem()
		}
		switch rv.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			return float64(rv.Int()), nil
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			return float64(rv.Uint()), nil
		case reflect.Float32, reflect.Float64:
			return rv.Float(), nil
		case reflect.Bool:
			if rv.Bool() {
				return 1, nil
			}
			return 0, nil
		default:
			return 0, fmt.Errorf("无法将类型 %T 转换为 float64", v)
		}
	}
}

// Any2bool 尝试将任意类型转换为 bool
func Any2bool(_v any) (bool, error) {
	switch v := _v.(type) {
	case nil:
		return false, nil
	case bool:
		return v, nil
	case string:
		if v == "" {
			return false, nil
		}
		return strconv.ParseBool(v)
	case int, int8, int16, int32, int64:
		// 非零值为 true，零值为 false
		val := reflect.ValueOf(v).Int()
		return val != 0, nil
	case uint, uint8, uint16, uint32, uint64:
		// 非零值为 true，零值为 false
		val := reflect.ValueOf(v).Uint()
		return val != 0, nil
	case float32, float64:
		// 非零值为 true，零值为 false
		val := reflect.ValueOf(v).Float()
		return val != 0, nil
	case json.Number:
		if v == "" {
			return false, nil
		}
		// 尝试解析为 float64
		f, err := v.Float64()
		if err != nil {
			return false, fmt.Errorf("无法将 json.Number %v 转换为 bool: %w", v, err)
		}
		return f != 0, nil
	default:
		// 使用反射处理更多类型
		rv := reflect.ValueOf(v)
		// 处理指针类型
		if rv.Kind() == reflect.Ptr {
			if rv.IsNil() {
				return false, nil
			}
			rv = rv.Elem()
		}
		switch rv.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			return rv.Int() != 0, nil
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			return rv.Uint() != 0, nil
		case reflect.Float32, reflect.Float64:
			return rv.Float() != 0, nil
		case reflect.Bool:
			return rv.Bool(), nil
		case reflect.String:
			return strconv.ParseBool(rv.String())
		default:
			return false, fmt.Errorf("无法将类型 %T 转换为 bool", v)
		}
	}
}

// Any2bytes 尝试将任意类型转换为 []byte
func Any2bytes(v any) ([]byte, error) {
	switch v := v.(type) {
	case []byte:
		return v, nil
	case string:
		return []byte(v), nil
	case bool, int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, float32, float64:
		b, _ := json.Marshal(v)
		return b, nil
	default:
		// 尝试使用 JSON 序列化其他类型
		b, err := json.Marshal(v)
		if err != nil {
			return nil, fmt.Errorf("无法将类型 %T 转换为 []byte: %w", v, err)
		}
		return b, nil
	}
}

// Any2Byte 将任意类型转换为字节数组
func Any2Byte(src any) []byte {
	b, _ := json.Marshal(src)
	return b
}

// Any2Reader 将任意类型转换为 io.Reader
func Any2Reader(src any) io.Reader {
	return bytes.NewReader(Any2Byte(src))
}
