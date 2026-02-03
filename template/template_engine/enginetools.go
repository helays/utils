package template_engine

import (
	"html/template"
	"time"

	"github.com/helays/utils/v2/tools"
)

func builtinFuncMap() template.FuncMap {
	return template.FuncMap{
		// 时间处理
		"now":        time.Now,
		"timestamp":  func(t time.Time) int64 { return t.Unix() },
		"formatDate": formatDate,
		"timeAgo":    timeSince, // 实现相对时间显示

		// 字符串处理
		"truncate": truncateString,
		"safeHTML": func(s string) template.HTML { return template.HTML(s) },

		// 数组
		"listString": func(s ...string) []string { return s },
		"listInt":    func(s ...int) []int { return s },
		"loopInt":    LoopInt,

		// 数学计算
		"add":    func(a, b int) int { return a + b },
		"mul":    func(a, b int) int { return a * b },
		"divide": func(a, b int) float64 { return float64(a) / float64(b) },

		// 链接处理
		"a":          A,
		"aSafe":      ASafe,
		"aWithQuery": AWithQuery,
		"dict":       Dict,

		// 数据转换函数
		"toBool": func(v any) bool {
			ok, _ := tools.Any2bool(v)
			return ok
		},
		"toString": func(v any) string {
			return tools.Any2string(v)
		},
		"toInt": func(v any) int64 {
			return tools.MustAny2Int[int64](v)
		},
		"toFloat": func(v any) float64 {
			n, _ := tools.Any2float64(v)
			return n
		},
	}
}
