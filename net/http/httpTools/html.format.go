package httpTools

import (
	"regexp"
	"strings"
)

// 默认需要删除的html 标签内容
var defaultcuttag []string = []string{"style", "script", "audio", "video", "iframe", "button"}

// Striptag 定义出去 html 标签的结构体
type Striptag struct {
	// 需要删除的标签
	Cuttag []string
	// 需要删除类似于 br input hr 这种单的标签。
	Deletetag []string
	// 匹配允许的标签时，使用 正则表达式的 负前瞻 <(?!abc|bff)
	Allowtag string
	Holdimg  bool
}

// StripTags 将非
func (h *Striptag) StripTags(s string) string {
	var err error
	s = strings.TrimSpace(s)
	// 第一步转换所有标签为小写
	if err = h.tolower(&s); err != nil {
		return s
	}
	// 过滤style script audio video iframe input
	if err = h.deleteThrow(&s); err != nil {
		return s
	}
	// 清除连续空白
	if err = h.clearSpace(&s); err != nil {
		return s
	}

	if err = h.tagsdiy(&s); err != nil {
		return s
	}

	return s
}

/*
*
自定义处理逻辑
*/
func (h *Striptag) tagsdiy(s *string) error {
	reg, err := regexp.Compile(`\s*?<[\S\s]+?>`)
	if err != nil {
		return err
	}
	if h.Allowtag == "" {
		h.Allowtag = "p"
	}
	*s = strings.TrimSpace(reg.ReplaceAllStringFunc(*s, h.diy))
	reg = regexp.MustCompile(`\n+`)
	if strings.TrimSpace(h.Allowtag) == "" {
		*s = reg.ReplaceAllString(*s, h.Allowtag)
		return nil
	}
	// 然后将"\n" 替换为</p><p>
	*s = "<" + h.Allowtag + ">" + reg.ReplaceAllString(*s, "</"+h.Allowtag+"><"+h.Allowtag+">") + "</" + h.Allowtag + ">"
	return nil
}

// 自定义处理逻辑，将非img 标签都替换为"\n"
func (h *Striptag) diy(s string) string {
	s = strings.TrimSpace(s)

	// 如果系统保留图片，
	if !h.Holdimg {
		return "\n"
	}

	reg := regexp.MustCompile(`[a-z1-6]+`)
	result := reg.FindStringSubmatch(s)
	if len(result) < 1 {
		return "\n"
	}
	tag := reg.FindStringSubmatch(s)[0]
	if tag != "img" {
		return "\n"
	}
	img := regexp.MustCompile(`src=['"](\S+?)['"]`)
	result = img.FindStringSubmatch(s)
	if len(result) < 2 {
		return "\n"
	}
	return `<img src="` + result[1] + `" />`
}

// 删除不需要的标签
func (h *Striptag) deleteThrow(s *string) error {
	if h.Cuttag == nil {
		h.Cuttag = defaultcuttag
	}
	r := strings.Join(h.Cuttag, "|")
	reg, err := regexp.Compile(`\s*?<\s*?(` + r + `)\s*?>[\S\s]*?<\s*?/(` + r + `)\s*?>\s*?`)
	if err != nil {
		return nil
	}
	*s = reg.ReplaceAllString(*s, "")
	return nil
}

// 清除连续的空白
func (h *Striptag) clearSpace(s *string) error {
	reg, err := regexp.Compile(`\s{2,}|[\n\r]`)
	if err != nil {
		return err
	}
	*s = reg.ReplaceAllString(*s, "")

	return nil
}

// 将所有标签转为小写
func (h *Striptag) tolower(s *string) error {
	reg, err := regexp.Compile(`<[\S\s]+?>`)
	if err != nil {
		return err
	}
	*s = reg.ReplaceAllStringFunc(*s, strings.ToLower)
	return nil
}
