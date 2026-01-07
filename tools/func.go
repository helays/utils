package tools

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"math/rand/v2"
	"net"
	url2 "net/url"
	"strings"
	"time"

	"github.com/helays/utils/v2/config"
	"golang.org/x/exp/constraints"
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
			isPrevLowerOrDigit := s[i-1] >= 'a' && s[i-1] <= 'z' || s[i-1] >= '0' && s[i-1] <= '9'

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

// UrlEncode 将 query部分进行 url encode
func UrlEncode(url string) string {
	u, err := url2.Parse(url)
	if err != nil {
		return "-"
	}
	u.RawQuery = url2.PathEscape(u.RawQuery)
	return u.String()
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

// IsPowerOfTwo 位运算 - 最简洁高效
// 检测输入数字是否是2的幂
func IsPowerOfTwo[T constraints.Integer](n T) bool {
	if n <= 0 {
		return false
	}
	return (n & (n - 1)) == 0
}
