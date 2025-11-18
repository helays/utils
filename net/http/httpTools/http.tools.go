package httpTools

import (
	"fmt"
	"net/http"
	"net/url"
	"regexp"
	"strings"

	"github.com/helays/utils/v2/tools"
)

// QueryGetSlice 获取query参数，并分割
func QueryGetSlice(r *http.Request, key string, step string) []string {
	query := r.URL.Query()
	v := query.Get(key)
	if v == "" {
		return nil
	}
	return strings.Split(v, step)
}

// QueryGet 获取query参数，如果值不存在就设置默认值
func QueryGet(query url.Values, k, dfValue string) string {
	v := query.Get(k)
	return tools.Ternary(v == "", dfValue, v)
}

// QueryValid 验证query参数，如果规则不匹配就返回错误
func QueryValid(query url.Values, k string, rule *regexp.Regexp) (string, error) {
	v := query.Get(k)
	if !rule.MatchString(v) {
		return "", fmt.Errorf("参数%s值格式错误", k)
	}
	return v, nil
}

// QueryValidAll 验证query参数，所有规则都需要满足
func QueryValidAll(query url.Values, k string, rules []*regexp.Regexp) (string, error) {
	v := query.Get(k)
	for _, r := range rules {
		if !r.MatchString(v) {
			return "", fmt.Errorf("参数%s值格式错误", k)
		}
	}
	return v, nil
}

// QueryValidAny 验证query参数，只要满足其中一个规则就返回
func QueryValidAny(query url.Values, k string, rules []*regexp.Regexp) (string, error) {
	v := query.Get(k)
	for _, r := range rules {
		if r.MatchString(v) {
			return v, nil
		}
	}
	return "", fmt.Errorf("参数%s值格式错误", k)
}

// RespCloneHeader 响应复制header
func RespCloneHeader(w http.ResponseWriter, header http.Header) {
	for k, v := range header {
		for _, vv := range v {
			w.Header().Add(k, vv)
		}
	}
}

// SetDisposition 文件下载时候，设置中文文件名
func SetDisposition(w http.ResponseWriter, filename string) {
	encodedFileName := url.QueryEscape(filename)
	// 设置Content-Disposition头部，使用RFC5987兼容的方式指定文件名
	contentDisposition := fmt.Sprintf("attachment; filename=\"%s\"; filename*=UTF-8''%s", encodedFileName, encodedFileName)
	w.Header().Set("Content-Disposition", contentDisposition)
}

var (
	// 定义正则表达式来匹配文件名
	// 匹配 filename="..." 格式
	reQuoted = regexp.MustCompile(`filename="([^"]+)"`)
	// 匹配 filename*=UTF-8''... 格式
	reEncoded = regexp.MustCompile(`filename\*\s*=\s*UTF-8''([^;]+)`)
	// 匹配没有引号的filename=格式
	reUnquoted = regexp.MustCompile(`filename=([^;]+)`)
)

// ParseDisposition 从Content-Disposition头中解析出原始文件名
// 兼容处理各种格式，包括inline和其他类型的disposition
func ParseDisposition(disposition string) string {
	if disposition == "" {
		return ""
	}
	var filename string

	// 优先级1: 使用RFC5987编码的文件名（filename*=UTF-8''）
	if matches := reEncoded.FindStringSubmatch(disposition); len(matches) > 1 {
		encodedName := matches[1]
		// URL解码
		if decodedName, err := url.QueryUnescape(encodedName); err == nil {
			filename = decodedName
		} else {
			filename = encodedName
		}
		return filename
	}

	// 优先级2: 使用引号包围的文件名
	if matches := reQuoted.FindStringSubmatch(disposition); len(matches) > 1 {
		quotedName := matches[1]
		// 如果文件名是URL编码的，尝试解码
		if decodedName, err := url.QueryUnescape(quotedName); err == nil {
			filename = decodedName
		} else {
			filename = quotedName
		}
		return filename
	}

	// 优先级3: 使用没有引号的文件名
	if matches := reUnquoted.FindStringSubmatch(disposition); len(matches) > 1 {
		unquotedName := strings.TrimSpace(matches[1])
		// 如果文件名是URL编码的，尝试解码
		if decodedName, err := url.QueryUnescape(unquotedName); err == nil {
			filename = decodedName
		} else {
			filename = unquotedName
		}
		return filename
	}

	return ""

}
