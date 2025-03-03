package tools

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/helays/utils/close/osClose"
	"github.com/helays/utils/config"
	"github.com/helays/utils/logger/ulogs"
	"io"
	"math"
	"math/rand"
	"net"
	url2 "net/url"
	"os"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"
	"time"
)

// PadRight 在字符串后面补齐固定字符，并达到n个长度
func PadRight(str string, padStr string, lenght int) string {
	if len(str) >= lenght {
		return str
	}
	for i := len(str); i < lenght; i++ {
		str += padStr
	}
	return str
}

// SnakeString 将驼峰命名法的字符串转换为蛇形命名法（小写字母加下划线）
func SnakeString(s string) string {
	var data []byte // 用于存储转换后的字符
	num := len(s)   // 获取字符串长度

	for i := 0; i < num; i++ {
		d := s[i] // 当前字符

		// 检查当前字符是否为大写字母且不是第一个字符
		if d >= 'A' && d <= 'Z' && i > 0 {
			// 向前看是否跟着一个小写字母
			isNextLower := i+1 < num && s[i+1] >= 'a' && s[i+1] <= 'z'
			// 向后看是否前面是小写字母或数字
			isPrevLowerOrDigit := i > 0 && (s[i-1] >= 'a' && s[i-1] <= 'z' || s[i-1] >= '0' && s[i-1] <= '9')

			// 如果当前字符是大写，并且前面是小写字母或数字，或者后面是小写字母，则添加下划线
			if isPrevLowerOrDigit || isNextLower {
				data = append(data, '_') // 添加下划线
			}
		}

		data = append(data, d) // 添加当前字符到结果中
	}

	return strings.ToLower(string(data)) // 返回转换为小写后的结果字符串
}

// CamelString 蛇形转驼峰
func CamelString(s string) string {
	data := make([]byte, 0, len(s))
	j := false
	k := false
	num := len(s) - 1
	for i := 0; i <= num; i++ {
		d := s[i]
		if k == false && d >= 'A' && d <= 'Z' {
			k = true
		}
		if d >= 'a' && d <= 'z' && (j || k == false) {
			d = d - 32
			j = false
			k = true
		}
		if k && d == '_' && num > i && s[i+1] >= 'a' && s[i+1] <= 'z' {
			j = true
			continue
		}
		data = append(data, d)
	}
	return string(data[:])
}

func JsonEncode(j any) ([]byte, error) {
	bf := bytes.NewBuffer([]byte{})
	jsonEncoder := json.NewEncoder(bf)
	jsonEncoder.SetEscapeHTML(false)
	err := jsonEncoder.Encode(j)
	if err != nil {
		return nil, err
	}

	return bf.Bytes(), nil
}

// DeleteStrarr 删除字符串切片的某一个元素
func DeleteStrarr(arr []string, val string) []string {
	for index, _id := range arr {
		if _id == val {
			arr = append(arr[:index], arr[index+1:]...)
			break
		}
	}
	return arr
}

// FilePutContents 快速简易写文件
func FilePutContents(path, content string) error {
	if err := Mkdir(filepath.Dir(path)); err != nil {
		return err
	}
	file, err := os.OpenFile(path, os.O_TRUNC|os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		return err
	}
	_, err = file.WriteString(content)
	osClose.CloseFile(file)
	return err
}

func FilePutContentsbytes(path string, content []byte) error {
	_path := filepath.Dir(path)
	if _, err := os.Stat(_path); err != nil {
		if err := Mkdir(_path); err != nil {
			return err
		}
	}

	file, err := os.OpenFile(path, os.O_TRUNC|os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	_, err = file.Write(content)
	osClose.CloseFile(file)
	return err
}

// FileAppendContents 快速简易写文件（追加）
func FileAppendContents(path, content string) error {
	_path := filepath.Dir(path)
	if _, err := os.Stat(_path); err != nil {
		if err := Mkdir(_path); err != nil {
			return err
		}
	}
	file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		return err
	}
	_, err = file.WriteString(content)
	osClose.CloseFile(file)
	return err
}

// FileGetContents 快速简易读取文件
func FileGetContents(path string) ([]byte, error) {
	file, err := os.OpenFile(path, os.O_RDONLY, 0755)
	if err != nil {
		return nil, err
	}
	defer osClose.CloseFile(file)
	return io.ReadAll(file)
}

// Mkdir 判断目录是否存在，否则创建目录
func Mkdir(path string) error {
	if _, err := os.Stat(path); err == nil {
		return nil
	}
	return os.MkdirAll(path, 0755)
}

// UrlEncode 将 query部分进行 url encode
func UrlEncode(url string) string {
	u, err := url2.Parse(url)
	if err != nil {
		return "-"
	}
	u.RawQuery = url2.PathEscape(u.RawQuery)
	return u.String()
}

func Int32tostring(i int32) string {
	return strconv.Itoa(int(i))
}

func Int32tobooltoint(i int32) int {
	if i > 0 {
		return 1
	}
	return 0
}

func Int64tostring(i int64) string {
	return strconv.FormatInt(i, 10)
}

func Float32tostring(f float32) string {
	f64 := float64(f)
	if math.IsNaN(f64) || math.IsInf(f64, 0) {
		return "0"
	}

	return strconv.FormatFloat(f64, 'f', 6, 64)
}

func Float64tostring(f float64) string {
	if math.IsNaN(f) || math.IsInf(f, 0) {
		return "0"
	}
	return strconv.FormatFloat(f, 'f', 6, 64)
}

func MaxInt32(d1, d2 int32) int32 {
	if d1 > d2 {
		return d1
	}
	return d2
}

// AvgInt32 计算平均数
func AvgInt32(d1, d2 int32, isf bool) int32 {
	if isf {
		if d1 > d2 {
			return d1
		}
		return d2
	}
	return (d1 + d2) / 2
}

func MinInt32(d1, d2 int32) int32 {
	if d1 > d2 {
		return d2
	}
	return d1
}

func MaxInt64(d1, d2 int64) int64 {
	if d1 > d2 {
		return d1
	}
	return d2
}

func MinInt64(d1, d2 int64) int64 {
	if d1 > d2 {
		return d2
	}
	return d1
}

func AvgInt64(d1, d2 int64, isf bool) int64 {
	if isf {
		if d1 > d2 {
			return d1
		}
		return d2
	}
	return (d1 + d2) / 2
}

func MaxUint64(d1, d2 uint64) uint64 {
	if d1 > d2 {
		return d1
	}
	return d2
}

func MinUint64(d1, d2 uint64) uint64 {
	if d1 > d2 {
		return d2
	}
	return d1
}

func AvgUint64(d1, d2 uint64, isf bool) uint64 {
	if isf {
		if d1 > d2 {
			return d1
		}
		return d2
	}
	return (d1 + d2) / 2
}

func MaxFloat32(d1, d2 float32) float32 {
	if d1 > d2 {
		return d1
	}
	return d2
}

func MinFloat32(d1, d2 float32) float32 {
	if d1 > d2 {
		return d2
	}
	return d1
}

func AvgFloat32(d1, d2 float32, isf bool) float32 {
	if isf {
		if d1 > d2 {
			return d1
		}
		return d2
	}
	return (d1 + d2) / 2
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

// // 字节转换成整形
// func BytesToInt(b []byte) (int, error) {
// 	bytesBuffer := bytes.NewBuffer(b)
//
// 	var x int32
// 	err = binary.Read(bytesBuffer, binary.BigEndian, &x)
//
// 	return int(x), err
// }

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
	return uint16(tmp)
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

// Searchslice 在切片中判断某个值是否存在
func Searchslice(s string, o []string) bool {
	if o == nil {
		return false
	}
	s = strings.TrimSpace(s)
	for _, i := range o {
		i = strings.TrimSpace(i)
		if i == s {
			return true
		}
	}
	return false
}

// StringUniq 对字符串切片进行去重
func StringUniq(tmp []string) []string {
	var tmpMap = make(map[string]bool)
	var result []string
	for _, item := range tmp {
		if tmpMap[item] {
			continue
		}
		result = append(result, item)
		tmpMap[item] = true
	}
	return result
}

// CreateSignature 带有 密钥的 sha1 hash
func CreateSignature(s, key string) string {
	h := hmac.New(sha1.New, []byte(key))
	h.Write([]byte(s))

	return base64.StdEncoding.EncodeToString(h.Sum(nil))
}

// ByteFormat 字节格式化
func ByteFormat(s int) string {
	var (
		p      float64
		format = ` Bytes`
		swap   = float64(s)
	)
	if s <= 1 {
		return "-"
	} else if s > 0 && s < 1024 {
		return strconv.Itoa(s) + format
	} else if s >= 1024 && swap < math.Pow(1024, 2) {
		p = 1
		format = ` KB`
	} else if swap >= math.Pow(1024, 2) && float64(s) < math.Pow(1024, 3) {
		p = 2
		format = ` MB`
	} else if swap >= math.Pow(1024, 3) && float64(s) < math.Pow(1024, 4) {
		p = 3
		format = ` GB`
	} else if swap >= math.Pow(1024, 4) && float64(s) < math.Pow(1024, 5) {
		p = 4
		format = ` TB`
	}
	return fmt.Sprintf("%.2f%s", swap/math.Pow(1024, p), format)
}

// MapDeepCopy map 深拷贝
func MapDeepCopy(src map[string]interface{}) map[string]interface{} {
	byt, _ := json.Marshal(src)
	var _tmp = new(map[string]interface{})
	_ = json.Unmarshal(byt, _tmp)
	return *_tmp
}

// SearchIntSlice 在整数切片中搜索指定的元素，并返回是否找到。
// 参数:
//
//	s - 待搜索的整数。
//	arr - 整数切片，将被搜索。
//
// 返回值:
//
//	如果找到 s 在 arr 中，则返回 true；否则返回 false。
func SearchIntSlice(s int, arr []int) bool {
	if arr == nil {
		return false
	}
	for _, i := range arr {
		if i == s {
			return true
		}
	}
	return false
}

func SearchInt64Slice(s int64, arr []int64) bool {
	if arr == nil {
		return false
	}
	for _, i := range arr {
		if i == s {
			return true
		}
	}
	return false
}

func SearchStringSlice(key string, arr []string) bool {
	if arr == nil {
		return false
	}
	for _, v := range arr {
		if v == key {
			return true
		}
	}
	return false
}

// CutStrSlice2Slice 获取切片的子切片
func CutStrSlice2Slice(s []string, key string, direct int) []string {
	for idx, v := range s {
		if v == key {
			if idx+direct < len(s) {
				return s[idx+direct:]
			} else {
				return []string{} // 索引越界时返回空切片
			}
		}
	}
	return []string{}
}

// Fileabs 生成文件的绝对路径
func Fileabs(cpath string) string {
	if filepath.IsAbs(cpath) {
		return cpath
	}
	return filepath.Join(config.Appath, cpath)
}

// FileAbsWithCurrent 生成文件的绝对路径,根目录手动指定
func FileAbsWithCurrent(current, cpath string) string {
	if filepath.IsAbs(cpath) {
		return cpath
	}
	return filepath.Join(current, cpath)
}

func RemoveAll(path string) {
	ulogs.Checkerr(os.RemoveAll(path), "删除文件失败")
}

var defaultLetters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

// RandomString 伪随机字符串
func RandomString(n int, allowedChars ...[]rune) string {
	var letters []rune
	if len(allowedChars) == 0 {
		letters = defaultLetters
	} else {
		letters = allowedChars[0]
	}
	rd := rand.New(rand.NewSource(time.Now().UnixNano()))
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rd.Intn(len(letters))]
	}
	return string(b)
}

// MinMaxAvgSum 获取数组中最大值最小值平均值和求和
func MinMaxAvgSum(nums []int) (min int, max int, avg float64, sum int) {
	if len(nums) == 0 {
		return 0, 0, 0, 0
	}
	min, max, sum = nums[0], nums[0], nums[0]
	for _, num := range nums[1:] {
		if num < min {
			min = num
		}
		if num > max {
			max = num
		}
		sum += num
	}
	avg = float64(sum) / float64(len(nums))
	return
}

// AnySlice2Str 将任意切片转成字符串
func AnySlice2Str(slice []any, _sep ...string) string {
	var builder strings.Builder
	l := len(slice)
	sep := ","
	if len(_sep) > 0 {
		sep = _sep[0]
	}
	for index, elem := range slice {
		// 使用 fmt.Sprint 将任何类型转换为字符串形式
		strElem := fmt.Sprint(elem)
		if strElem == "" {
			continue
		}
		builder.WriteString(strElem)
		// 可以选择在此处添加分隔符，如空格、逗号等
		if index < (l - 1) {
			builder.WriteString(sep)
		}
	}

	return builder.String()
}

// Struct2Map 将结构体转换为map
func Struct2Map(src any) map[string]any {
	var _map map[string]any
	byt, _ := json.Marshal(src)
	_ = json.Unmarshal(byt, &_map)
	return _map
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
func Any2float64(v any) (float64, error) {
	switch v := v.(type) {
	case float32:
		return float64(v), nil
	case float64:
		return v, nil
	case string:
		if v == "" {
			return 0, nil
		}
		return strconv.ParseFloat(v, 64)
	case nil:
		return 0, nil
	case json.Number:
		return v.Float64()
	default:
		return 0.0, fmt.Errorf("无法将类型 %T 转换为 float64", v)
	}
}

// Any2bool 尝试将任意类型转换为 bool
func Any2bool(v any) (bool, error) {
	switch v := v.(type) {
	case bool:
		return v, nil
	case string:
		return strconv.ParseBool(v)
	default:
		return false, fmt.Errorf("无法将类型 %T 转换为 bool", v)
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

// GetIpVersion 解析ip地址，确认ip版本
func GetIpVersion(ip string) (string, error) {
	_ip := net.ParseIP(ip)
	if _ip == nil {
		return "", errors.New("ip地址不合法")
	}
	if _ip.To4() != nil {
		return "ipv4", nil
	}
	return "ipv6", nil
}

// Str2StrSlice 字符串转切片
func Str2StrSlice(values string) ([]string, error) {
	values = strings.TrimSpace(values)
	if values == "" {
		return nil, nil
	}
	var slice []string
	if err := json.Unmarshal([]byte(values), &slice); err != nil {
		return nil, fmt.Errorf("解析数据 [%s] 失败：%v", values, err)
	}
	return slice, nil
}

// Slice2MapWithHeader 主要是将excel 或者 csv的每一行转为map，键为header，值为cell
func Slice2MapWithHeader(rows any, header []string) map[string]any {
	// 获取 rows 的反射值
	rowsValue := reflect.ValueOf(rows)
	// 检查 rows 是否为切片类型
	if rowsValue.Kind() != reflect.Slice && rowsValue.Kind() != reflect.Ptr {
		return nil
	}
	// 如果 rows 是切片的指针，则获取指向的切片
	if rowsValue.Kind() == reflect.Ptr {
		if rowsValue.IsNil() {
			return nil
		}
		rowsValue = rowsValue.Elem()
	}
	fieldLen := len(header)
	var tmp = make(map[string]any)
	//判断rows是切片，或者是切片的指针，如果是就遍历，不是就返回nil
	// 遍历 rows 切片
	for i := 0; i < rowsValue.Len(); i++ {
		if i >= fieldLen {
			continue
		}
		tmp[header[i]] = rowsValue.Index(i).Interface()
	}
	return tmp
}

// Ternary 是一个通用的三元运算函数。
// 它接受一个布尔条件和两个参数 a 和 b。
// 如果条件为 true，则返回 a；否则返回 b。
func Ternary[Type any](condition bool, a, b Type) Type {
	if condition {
		return a
	}
	return b
}

// AutoTimeDuration 自动转换时间单位，主要是用于 ini json yaml 几种配置文件 解析出来的时间单位不一致。
func AutoTimeDuration(input time.Duration, unit time.Duration, dValue ...time.Duration) time.Duration {
	if input < 1 {
		if len(dValue) < 1 {
			return 0
		}
		return dValue[0]
	}
	// 然后这里就要开始自适应ini json yaml 几种配置文件解析出来的时间勒
	if input < time.Microsecond {
		// 这表示输入时间就是默认单位,要更新单位
		return input * unit
	}
	return input
}

// ReverseMapUnique 反转值唯一的 map
func ReverseMapUnique[K comparable, V comparable](m map[K]V) map[V]K {
	reversed := make(map[V]K)
	for k, v := range m {
		reversed[v] = k
	}
	return reversed
}

// GetLevel2MapValue 获取二级map的值
func GetLevel2MapValue[K any](inp map[string]map[string]K, key1, key2 string) (K, bool) {
	if v, ok := inp[key1]; ok {
		if vv, ok := v[key2]; ok {
			return vv, true
		}
	}
	var zeroValue K
	return zeroValue, false
}
