package tools

import (
	"bytes"
	"cmp"
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"math/rand/v2"
	"net"
	url2 "net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/helays/utils/v2/config"
	"github.com/helays/utils/v2/logger/ulogs"
)

// NestedMapSet 自动初始化嵌套 map 并设置值
// 如果内层map不存在，会自动创建
func NestedMapSet[K comparable, V comparable, T any](m map[K]map[V]T, outerKey K, innerKey V, value T) {
	if _, ok := m[outerKey]; !ok {
		m[outerKey] = make(map[V]T)
	}
	m[outerKey][innerKey] = value
}

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
	return Float64tostring(float64(f))
}

func Float64tostring(f float64) string {
	if math.IsNaN(f) || math.IsInf(f, 0) {
		return "0"
	}
	if f == math.Trunc(f) {
		return strconv.FormatInt(int64(f), 10)
	}
	return strconv.FormatFloat(f, 'f', 6, 64)
}

func Max[T cmp.Ordered](d1, d2 T) T {
	if d1 > d2 {
		return d1
	}
	return d2
}

func Min[T cmp.Ordered](d1, d2 T) T {
	if d1 < d2 {
		return d1
	}
	return d2
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
	if s < 0 {
		return "-" // 对于负值，返回"-"
	}
	if s == 0 {
		return "0 Bytes" // 特殊处理0的情况
	}

	units := []string{"Bytes", "KB", "MB", "GB", "TB"}
	exponent := math.Floor(math.Log(float64(s)) / math.Log(1024))

	// 限制 exponent 的范围，防止数组越界
	exponent = math.Min(exponent, float64(len(units)-1))
	converted := float64(s) / math.Pow(1024, exponent)

	// 使用 %.2f 来格式化输出，确保小数点后有两位数字
	return fmt.Sprintf("%.2f %s", converted, units[int(exponent)])
}

// MapDeepCopy map 深拷贝
func MapDeepCopy[T any](src T, dst *T) {
	byt, _ := json.Marshal(src)
	_ = json.Unmarshal(byt, dst)
}

// Fileabs 生成文件的绝对路径
// noinspection SpellCheckingInspection
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

	rng := config.RandPool.Get().(*rand.Rand)
	defer config.RandPool.Put(rng)

	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rng.IntN(len(letters))]
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

// Struct2Map 将结构体转换为map
func Struct2Map(src any) map[string]any {
	var _map map[string]any
	byt, _ := json.Marshal(src)
	_ = json.Unmarshal(byt, &_map)
	return _map
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

// IsZero isZero 检查值是否为类型的零值
func IsZero[T comparable](v T) bool {
	var zero T
	return v == zero
}
