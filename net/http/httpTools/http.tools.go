package httpTools

import (
	"fmt"
	"github.com/helays/utils/tools"
	"net/http"
	"net/url"
	"regexp"
	"strings"
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

// SetDisposition 文件下载时候，设置中文文件名
func SetDisposition(w http.ResponseWriter, filename string) {
	encodedFileName := url.QueryEscape(filename)
	// 设置Content-Disposition头部，使用RFC5987兼容的方式指定文件名
	contentDisposition := fmt.Sprintf("attachment; filename=\"%s\"; filename*=UTF-8''%s", encodedFileName, encodedFileName)
	w.Header().Set("Content-Disposition", contentDisposition)
}

// RespCloneHeader 响应复制header
func RespCloneHeader(w http.ResponseWriter, header http.Header) {
	for k, v := range header {
		for _, vv := range v {
			w.Header().Add(k, vv)
		}
	}
}
