package tools

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"math"
	"reflect"
	"strconv"
	"strings"
	"time"
)

func TrimStringHelper(v string) string {
	v = strings.TrimSpace(v)
	v = strings.TrimFunc(v, func(r rune) bool {
		return r == '"' || r == '`' || r == '\''
	})
	return strings.TrimSpace(v)
}

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
	case map[string]any, []map[string]any, []string, []int, []float64, []float32, struct{}, *struct{}:
		_byt, _ := json.Marshal(v)
		return string(_byt)
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

func MustAny2int(_v any) int64 {
	v, _ := Any2int(_v)
	return v
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
		return String2Int64(v)
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

func String2Int64(v string) (int64, error) {
	v = TrimStringHelper(v)
	if v == "" || v == "null" || v == "nil" || v == "undefined" {
		return 0, nil
	}

	if v == "true" {
		return 1, nil
	} else if v == "false" {
		return 0, nil
	}

	// 去除千分位逗号
	v = strings.ReplaceAll(v, ",", "")
	// 处理科学计数法
	if strings.ContainsAny(v, "eE") {
		// 先尝试用浮点数解析，再转整型
		if f, err := strconv.ParseFloat(v, 64); err == nil {
			return int64(f), nil
		}
	}

	// 处理不同进制
	if len(v) > 2 {
		prefix := v[:2]
		switch prefix {
		case "0x", "0X": // 十六进制
			return strconv.ParseInt(v[2:], 16, 64)
		case "0b", "0B": // 二进制
			return strconv.ParseInt(v[2:], 2, 64)
		case "0o", "0O": // 八进制
			return strconv.ParseInt(v[2:], 8, 64)
		}
	}

	// 包含小数点
	if strings.Contains(v, ".") {
		parts := strings.SplitN(v, ".", 2) // 只分割一次
		if parts[0] == "" && parts[1] == "" {
			return 0, fmt.Errorf("invalid number: %s", v)
		}
		if parts[0] == "" {
			return 0, nil // ".123" 的情况返回 0
		}
		return strconv.ParseInt(parts[0], 10, 64)
	}
	return strconv.ParseInt(v, 10, 64)
}

func MustAny2float64(_v any) float64 {
	v, _ := Any2float64(_v)
	return v
}

// Any2float64 尝试将任意类型转换为 float64
func Any2float64(_v any) (float64, error) {
	switch v := _v.(type) {
	case float32:
		return float64(v), nil
	case float64:
		return v, nil
	case int:
		return float64(v), nil
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
		return String2Float64(v)
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

func String2Float64(v string) (float64, error) {
	v = TrimStringHelper(v)
	if v == "" || v == "null" || v == "nil" || v == "undefined" {
		return 0, nil
	}

	if v == "true" {
		return 1, nil
	} else if v == "false" {
		return 0, nil
	}
	// 去除千分位逗号
	v = strings.ReplaceAll(v, ",", "")

	// 处理不同进制（浮点数通常只支持十进制，但为了兼容性可以特殊处理十六进制）
	if len(v) > 2 {
		prefix := v[:2]
		switch prefix {
		case "0x", "0X": // 十六进制浮点数
			// Go 的 strconv.ParseFloat 不支持十六进制浮点数，需要特殊处理
			if strings.ContainsAny(v, "pP") {
				// 如果是十六进制科学计数法格式，如 "0x1.2p3"
				return strconv.ParseFloat(v, 64)
			}
			// 普通十六进制整数，先转整型再转浮点数
			if i, err := strconv.ParseInt(v[2:], 16, 64); err == nil {
				return float64(i), nil
			}
		case "0b", "0B": // 二进制
			if i, err := strconv.ParseInt(v[2:], 2, 64); err == nil {
				return float64(i), nil
			}
		case "0o", "0O": // 八进制
			if i, err := strconv.ParseInt(v[2:], 8, 64); err == nil {
				return float64(i), nil
			}
		}
	}

	// 对于浮点数，不需要特殊处理小数点情况，因为 ParseFloat 可以直接处理
	return strconv.ParseFloat(v, 64)
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
	switch _v := v.(type) {
	case []byte:
		return _v, nil
	case string:
		return []byte(_v), nil
	case bool, int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, float32, float64:
		b, _ := json.Marshal(_v)
		return b, nil
	default:
		// 尝试使用 JSON 序列化其他类型
		b, err := json.Marshal(_v)
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
func Any2Reader(src any) *bytes.Reader {
	return bytes.NewReader(Any2Byte(src))
}

func Any2Map(src any) (any, error) {
	if src == "" || src == nil {
		return nil, nil
	}
	var dst map[string]any
	var dstSlice []any
	switch _src := src.(type) {
	case string:
		// 尝试解析JSON字符串
		_byt := []byte(_src)
		if err := json.Unmarshal(_byt, &dst); err == nil {
			return dst, nil
		}
		err := json.Unmarshal(_byt, &dstSlice)
		return dstSlice, err
	case map[string]any:
		return _src, nil
	case []byte:
		// 尝试解析JSON字节数组
		if err := json.Unmarshal(_src, &dst); err == nil {
			return dst, nil
		}
		err := json.Unmarshal(_src, &dstSlice)
		return dstSlice, err
		// 扩展的切片类型检测
	case []any, []string, []int, []int8, []int16, []int32, []int64,
		[]uint, []uint16, []uint32, []uint64, []uintptr,
		[]float32, []float64, []bool,
		[]complex64, []complex128,
		[][]byte,
		[]map[string]any, []map[any]any:
		return _src, nil
	case int, int8, int16, int32, int64,
		uint, uint8, uint16, uint32, uint64,
		float32, float64, bool:
		return [1]any{src}, nil
	default:
		// 尝试通过反射处理其他类型
		val := reflect.ValueOf(src)
		switch val.Kind() {
		case reflect.Map, reflect.Struct, reflect.Slice, reflect.Array:
			return _src, nil
		case reflect.Ptr:
			// 解引用指针
			if val.IsNil() {
				return nil, nil
			}
			return Any2Map(val.Elem().Interface())
		default:
			return [1]any{src}, nil
		}

	}
}

func Object2MapStringAny(input any) (map[string]any, bool) {
	if input == nil {
		return nil, false
	}

	// 情况1: 直接是 map[string]any
	if m, ok := input.(map[string]any); ok {
		return m, true
	}
	// 情况2: 其他类型的 map
	if m, ok := tryConvertOtherMap(input); ok {
		return m, true
	}
	// 情况3: []any 切片 - 只取第一个元素
	if slice, ok := input.([]any); ok && len(slice) > 0 {
		// 只尝试转换第一个元素
		return Object2MapStringAny(slice[0])
	}
	// 情况4: 结构体
	if m, ok := tryConvertStruct(input); ok {
		return m, true
	}

	return nil, false
}

// 尝试转换其他类型的 map
func tryConvertOtherMap(input any) (map[string]any, bool) {
	switch m := input.(type) {
	case map[string]string:
		result := make(map[string]any, len(m))
		for k, v := range m {
			result[k] = v
		}
		return result, true

	case map[string]int:
		result := make(map[string]any, len(m))
		for k, v := range m {
			result[k] = v
		}
		return result, true

	case map[string]float64:
		result := make(map[string]any, len(m))
		for k, v := range m {
			result[k] = v
		}
		return result, true

	case map[string]bool:
		result := make(map[string]any, len(m))
		for k, v := range m {
			result[k] = v
		}
		return result, true

	case map[int]any:
		result := make(map[string]any, len(m))
		for k, v := range m {
			result[fmt.Sprint(k)] = v
		}
		return result, true

	case map[any]any:
		result := make(map[string]any, len(m))
		for k, v := range m {
			if strKey, ok := k.(string); ok {
				result[strKey] = v
			} else {
				result[fmt.Sprint(k)] = v
			}
		}
		return result, true

	default:
		return nil, false
	}
}

// 尝试将结构体转换为 map
func tryConvertStruct(input any) (map[string]any, bool) {
	val := reflect.ValueOf(input)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	if val.Kind() != reflect.Struct {
		return nil, false
	}

	result := make(map[string]any)
	typ := val.Type()

	for i := 0; i < val.NumField(); i++ {
		field := typ.Field(i)
		// 跳过非导出字段
		if field.PkgPath != "" {
			continue
		}

		fieldValue := val.Field(i)
		// 获取字段名，可以使用json tag
		fieldName := field.Name
		if jsonTag := field.Tag.Get("json"); jsonTag != "" {
			if jsonTag == "-" {
				continue // 跳过明确忽略的字段
			}
			// 取json tag中逗号前的部分
			if commaIdx := strings.Index(jsonTag, ","); commaIdx != -1 {
				jsonTag = jsonTag[:commaIdx]
			}
			if jsonTag != "" {
				fieldName = jsonTag
			}
		}

		// 递归处理嵌套结构体
		if fieldValue.Kind() == reflect.Struct {
			if nestedMap, ok := tryConvertStruct(fieldValue.Interface()); ok {
				result[fieldName] = nestedMap
				continue
			}
		}

		result[fieldName] = fieldValue.Interface()
	}

	return result, true
}

// 预定义常见类型，减少反射调用
var (
	mapStringInterfaceType = reflect.TypeOf(map[string]any(nil))
	sliceInterfaceType     = reflect.TypeOf([]any(nil))
	timeType               = reflect.TypeOf(time.Time{})
)

// CheckIsObject 检查给定的值是否是一个对象（map、slice、struct、指针）
func CheckIsObject(v any) bool {
	// 先尝试类型断言（最快路径）
	switch v.(type) {
	case map[string]any, []any, // 最常见
		map[any]any,                // 任意 key 的 map
		[]map[string]any,           // map 数组
		[]int, []float64, []string, // 基本类型 slice
		struct{}, *struct{}: // 结构体及其指针
		return true
	case time.Time, *time.Time: // time.Time 通常不需要再序列化
		return false
	case nil, string, int, float64, bool: // 基本类型直接跳过
		return false
	}
	// 获取反射类型
	typ := reflect.TypeOf(v)
	if typ == nil {
		return false // nil 不需要序列化
	}
	// 检查是否是指针
	if typ.Kind() == reflect.Ptr {
		if reflect.ValueOf(v).IsNil() {
			return false // nil 指针不需要序列化
		}
		typ = typ.Elem() // 解引用指针
	}

	// 用预定义类型快速匹配（比 Kind() 更快）
	switch typ {
	case mapStringInterfaceType, sliceInterfaceType, timeType:
		return typ != timeType // time.Time 不序列化
	}

	// 检查 Kind（处理其他情况）
	switch typ.Kind() {
	case reflect.Map, reflect.Slice, reflect.Array, reflect.Struct:
		return true
	default:
		return false
	}
}

// Int32tobooltoint 将 int32 转换为 bool 并返回 int
// noinspection all
func Int32tobooltoint(i int32) int {
	if i > 0 {
		return 1
	}
	return 0
}

// Float32tostring 将 float32 转换为字符串
// noinspection all
func Float32tostring(f float32) string {
	return Float64tostring(float64(f))
}

// Float64tostring 将 float64 转换为字符串
// noinspection all
func Float64tostring(f float64) string {
	if math.IsNaN(f) || math.IsInf(f, 0) {
		return "0"
	}
	if f == math.Trunc(f) {
		return strconv.FormatInt(int64(f), 10)
	}
	return strconv.FormatFloat(f, 'f', 6, 64)
}

// StrToFloat64 字符串转 float 64
func StrToFloat64(s string) (float64, error) {
	return strconv.ParseFloat(s, 64)
}

func Bool1time(b bool) int {
	if b {
		return 1
	}
	return 0
}

// Booltostring 布尔转 1 0
func Booltostring(b bool) string {
	if b {
		return "1"
	}
	return "0"
}

// Uint64tostring uint64 转 string
// noinspection SpellCheckingInspection
func Uint64tostring(i uint64) string {
	return strconv.FormatUint(i, 10)
}

func Uint16ToBytes(n int) ([]byte, error) {
	var err error
	tmp := uint16(n)
	bytesBuffer := bytes.NewBuffer([]byte{})
	err = binary.Write(bytesBuffer, binary.BigEndian, tmp)
	return bytesBuffer.Bytes(), err
}

func Uint32ToBytes(n int) ([]byte, error) {
	var err error
	tmp := uint32(n)
	bytesBuffer := bytes.NewBuffer([]byte{})
	err = binary.Write(bytesBuffer, binary.BigEndian, tmp)
	return bytesBuffer.Bytes(), err
}

// BytesToInt 字节转换成整形
func BytesToInt(b []byte) int {
	bytesBuffer := bytes.NewBuffer(b)

	var x int32
	_ = binary.Read(bytesBuffer, binary.BigEndian, &x)

	return int(x)
}

func BytesToUint16(b []byte) uint16 {
	bytesBuffer := bytes.NewBuffer(b)
	var tmp uint16
	_ = binary.Read(bytesBuffer, binary.BigEndian, &tmp)
	return tmp
}

// EmptyString2 空字符串转为 -
func EmptyString2(s string) string {
	if s = strings.TrimSpace(s); s == "" {
		return "-"
	}
	return s
}

func NumberEmptyString(s string) string {
	if s = strings.TrimSpace(s); s == "" {
		return "0"
	}
	return s
}

// Struct2Map 将结构体转换为map
func Struct2Map(src any) map[string]any {
	var _map map[string]any
	byt, _ := json.Marshal(src)
	_ = json.Unmarshal(byt, &_map)
	return _map
}
