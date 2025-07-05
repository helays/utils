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
	"github.com/helays/utils/close/osClose"
	"github.com/helays/utils/close/vclose"
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
	file, err := os.Open(path)
	defer vclose.Close(file)
	if err != nil {
		return nil, err
	}
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

func SearchAnySlice(in any, lst []any) bool {
	// 第一轮：快速尝试直接比较（所有类型）
	if lst == nil || len(lst) < 1 {
		return false
	}
	for _, v := range lst {
		if v == in {
			return true
		}
	}

	// 如果是基本类型且==比较失败，直接返回false
	switch in.(type) {
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, float32, float64, string, bool:
		return false
	}

	// 其他类型使用深度比较
	for _, v := range lst {
		if reflect.DeepEqual(in, v) {
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

type RetryCallbackFunc func() error

// RetryRunner 重试执行函数
func RetryRunner(retry int, sleep time.Duration, callback RetryCallbackFunc) {
	for i := 0; i < retry; i++ {
		err := callback()
		if err == nil {
			return
		}
		// 最后一次不睡眠
		if i == retry-1 {
			break
		}
		time.Sleep(sleep)
	}
}
