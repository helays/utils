package tools

import (
	"fmt"
	"strings"
)

func StrSlice2AnySlice(inp []string) []any {
	var out []any
	for _, v := range inp {
		out = append(out, v)
	}
	return out
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
