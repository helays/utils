package httpTools

import (
	"regexp"
	"strings"
	"unicode/utf8"
)

// SubstrTag @params s string
// 表示被截取的字符串
// @params n int
// 表示截取的字数
func SubstrTag(s string, n int) string {

	if utf8.RuneCountInString(s) < n {
		return s
	}
	reg, err := regexp.Compile(`((?:<[a-z1-6]+>)?[\S\s]+?(?:</[a-z1-6]+>))`)
	if err != nil {
		return ""
	}
	re := reg.FindAllString(s, -1)
	if len(re) == 0 {
		// 截取字符串，利用rune 防止中文乱码
		return string([]rune(s)[0:n])
	}
	// 当第一段字数大于 指定数量时，直接返回第一段
	if utf8.RuneCountInString(re[0]) > n {
		return re[0]
	}
	// 当前第一段字数少于指定数量时，再加上第二段字数 如果大于指定字数的最上线就返回重0到当前。
	var (
		num = 0
		i   = 0
		j   string
	)
	for i, j = range re {
		onnum := utf8.RuneCountInString(j)
		num += onnum
		if num < n {
			continue
		}
		if 20 > (num-n) || i == 2 {
			i++
		}
		break
	}

	return strings.Join(re[0:i], "")

}

func Substr(s string, n int) string {
	reg := regexp.MustCompile(`\s?<\s?/?\s?[\S\s]+?>\s?`)
	s = reg.ReplaceAllString(s, "")

	return SubstrTag(s, n)
}
