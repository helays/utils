package tools

import (
	"fmt"
	"reflect"
	"strings"
)

func StrSlice2AnySlice(inp []string) []any {
	var out []any
	for _, v := range inp {
		out = append(out, v)
	}
	return out
}

func AnySlice2StrSlice(slice []any) []string {
	var out []string
	for _, v := range slice {
		out = append(out, Any2string(v))
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

func AnySlice2StrWithEmpty(slice []any, _sep ...string) string {
	var builder strings.Builder
	l := len(slice)
	sep := ","
	if len(_sep) > 0 {
		sep = _sep[0]
	}
	for index, elem := range slice {
		// 使用 fmt.Sprint 将任何类型转换为字符串形式
		var strElem string
		if elem == nil {
			strElem = ""
		} else {
			strElem = Any2string(elem)
		}

		builder.WriteString(strElem)
		// 可以选择在此处添加分隔符，如空格、逗号等
		if index < (l - 1) {
			builder.WriteString(sep)
		}
	}
	return builder.String()
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
