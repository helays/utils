package httpTools

import "strings"

func HtmlExtractLink(body string, nofollow bool, fun func(u string)) {
	for _, item := range ALinkRegex.FindAllStringSubmatch(body, -1) {
		if len(item) != 2 {
			continue
		}
		if nofollow && Nofollow.MatchString(item[0]) {
			continue
		}
		if fun != nil {
			fun(strings.TrimSpace(item[1]))
		}
	}
	body = ""
}
