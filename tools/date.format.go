package tools

import (
	"strings"

	"helay.net/go/utils/v3/config"
)

// 判断是否是Go语言的时间格式模板
func isGoTimeFormat(format string) bool {
	// Go语言时间格式必须包含2006、01、02、15、04、05等关键数字
	goFormatElements := []string{"2006", "06", "01", "1", "02", "2", "15", "03", "3", "04", "4", "05", "5", "PM", "pm", "MST", "-07:00", "-0700"}

	for _, elem := range goFormatElements {
		if strings.Contains(format, elem) {
			return true
		}
	}
	return false
}

// 定义各语言到Go的转换映射
var (
	conversionMapPHP = map[string]string{
		"Y": "2006", "y": "06",
		"m": "01", "n": "1",
		"d": "02", "j": "2",
		"H": "15", "G": "15",
		"h": "03", "g": "3",
		"i": "04",
		"s": "05",
		"A": "PM", "a": "pm",
		"T": "MST",
		"P": "-07:00", "O": "-0700",
	}
	conversionMapJAVA = map[string]string{
		"yyyy": "2006", "yy": "06",
		"MM": "01", "M": "1",
		"dd": "02", "d": "2",
		"HH": "15", "H": "15",
		"hh": "03", "h": "3",
		"mm": "04", "m": "4",
		"ss": "05", "s": "5",
		"a": "PM",
		"z": "MST", "Z": "-0700",
	}
	conversionMapPYTHON = map[string]string{
		"%Y": "2006", "%y": "06",
		"%m": "01", "%-m": "1",
		"%d": "02", "%-d": "2",
		"%H": "15", "%-H": "15",
		"%I": "03", "%-I": "3",
		"%M": "04", "%-M": "4",
		"%S": "05", "%-S": "5",
		"%p": "PM",
		"%Z": "MST", "%z": "-0700",
	}
	conversionMapLUA = map[string]string{
		"%Y": "2006", "%y": "06",
		"%m": "01", "%-m": "1",
		"%d": "02", "%-d": "2",
		"%H": "15", "%-H": "15",
		"%I": "03", "%-I": "3",
		"%M": "04", "%-M": "4",
		"%S": "05", "%-S": "5",
		"%p": "PM",
		"%Z": "MST", "%z": "-0700",
	}
)

// ConvertTimeFormat 多语言时间格式转换器
func ConvertTimeFormat(format string, lang string) string {
	// 如果输入已经是Go格式，直接返回
	if isGoTimeFormat(format) {
		return format
	}
	var convMap map[string]string
	// 获取对应语言的转换映射
	lang = strings.ToLower(lang)
	switch strings.ToLower(lang) {
	case config.ProgramLangJAVA:
		// 特殊处理Java格式(需要按长度排序确保长格式优先匹配)
		return convertTimeFormatJAVA(format)
	case config.ProgramLangPHP:
		convMap = conversionMapPHP
	case config.ProgramLangPYTHON:
		convMap = conversionMapPYTHON
	case config.ProgramLangLUA:
		convMap = conversionMapLUA
	default:
		return ""

	}
	return convertTimeFormatCommon(format, convMap)
}

// 将Java格式转换为Go格式
func convertTimeFormatJAVA(format string) string {
	sortedKeys := make([]string, 0, len(conversionMapJAVA))
	for k := range conversionMapJAVA {
		sortedKeys = append(sortedKeys, k)
	}
	// 按长度降序排序
	for i := 0; i < len(sortedKeys)-1; i++ {
		for j := i + 1; j < len(sortedKeys); j++ {
			if len(sortedKeys[i]) < len(sortedKeys[j]) {
				sortedKeys[i], sortedKeys[j] = sortedKeys[j], sortedKeys[i]
			}
		}
	}
	// 按排序后的键替换
	for _, k := range sortedKeys {
		format = strings.Replace(format, k, conversionMapJAVA[k], -1)
	}
	return format
}

func convertTimeFormatCommon(format string, convMap map[string]string) string {
	for from, to := range convMap {
		format = strings.Replace(format, from, to, -1)
	}
	return format
}
