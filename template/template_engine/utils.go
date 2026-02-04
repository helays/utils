package template_engine

import (
	"fmt"
	"github.com/helays/utils/v2/tools"
	"html"
	"net/url"
	"time"
)

//
// ━━━━━━神兽出没━━━━━━
// 　　 ┏┓     ┏┓
// 　　┏┛┻━━━━━┛┻┓
// 　　┃　　　　　 ┃
// 　　┃　　━　　　┃
// 　　┃　┳┛　┗┳  ┃
// 　　┃　　　　　 ┃
// 　　┃　　┻　　　┃
// 　　┃　　　　　 ┃
// 　　┗━┓　　　┏━┛　Code is far away from bug with the animal protecting
// 　　　 ┃　　　┃    神兽保佑,代码无bug
// 　　　　┃　　　┃
// 　　　　┃　　　┗━━━┓
// 　　　　┃　　　　　　┣┓
// 　　　　┃　　　　　　┏┛
// 　　　　┗┓┓┏━┳┓┏┛
// 　　　　 ┃┫┫ ┃┫┫
// 　　　　 ┗┻┛ ┗┻┛
//
// ━━━━━━感觉萌萌哒━━━━━━
//
//
// User helay
// Date: 2025/6/15 15:25
//

// 日期格式化函数
func formatDate(layout string, t time.Time) string {
	return t.Format(layout)
}

func timeSince(t time.Time) string {
	return time.Since(t).String()
}

func truncateString(s string, length int) string {
	if len(s) > length {
		return s[:length]
	}
	return s
}

// A 生成HTML超链接标签
// text: 链接显示的文本
// urlStr: 链接地址
// options: 可选参数，可以包含HTML属性如class, id等
func A(text string, urlStr string, options ...map[string]any) string {
	// 对文本和URL进行HTML转义
	escapedUrl := html.EscapeString(urlStr)

	// 构建属性字符串
	attrs := ""

	for _, option := range options {
		for key, value := range option {
			// 跳过href属性，因为我们会单独处理
			if key == "href" {
				continue
			}
			escapedValue := html.EscapeString(tools.Any2string(value))
			attrs += fmt.Sprintf(` %s="%s"`, key, escapedValue)
		}
	}

	// 返回完整的a标签
	return fmt.Sprintf(`<a href="%s"%s>%s</a>`, escapedUrl, attrs, text)
}

func ASafe(text, urlStr string, options ...map[string]any) string {
	// 对文本和URL进行HTML转义
	escapedText := html.EscapeString(text)
	escapedUrl := html.EscapeString(urlStr)

	// 构建属性字符串
	attrs := ""

	for _, option := range options {
		for key, value := range option {
			// 跳过href属性，因为我们会单独处理
			if key == "href" {
				continue
			}
			escapedValue := html.EscapeString(tools.Any2string(value))
			attrs += fmt.Sprintf(` %s="%s"`, key, escapedValue)
		}
	}

	// 返回完整的a标签
	return fmt.Sprintf(`<a href="%s"%s>%s</a>`, escapedUrl, attrs, escapedText)
}

// AWithQuery
// text: 链接显示的文本
// baseUrl: 基础URL
// queryParams: 查询参数
// options: 可选HTML属性
func AWithQuery(text string, baseUrl string, queryParams map[string]any, options ...map[string]any) string {
	// 解析基础URL
	u, err := url.Parse(baseUrl)
	if err != nil {
		return A(text, baseUrl, options...)
	}

	// 添加查询参数
	q := u.Query()
	for key, value := range queryParams {
		q.Add(key, tools.Any2string(value))
	}
	u.RawQuery = q.Encode()

	// 使用A函数生成链接
	return A(text, u.String(), options...)
}

// Dict 创建一个map，用于模板中传递多个键值对
// 只返回一个值，如果参数错误返回空map
func Dict(values ...any) map[string]any {
	if len(values)%2 != 0 {
		return map[string]any{} // 返回空map而不是错误
	}
	dict := make(map[string]any, len(values)/2)
	for i := 0; i < len(values); i += 2 {
		key, ok := values[i].(string)
		if !ok {
			continue
		}
		dict[key] = values[i+1]
	}

	return dict
}

func LoopInt(start, end int) []int {
	var s []int
	for i := start; i <= end; i++ {
		s = append(s, i)
	}
	return s
}
