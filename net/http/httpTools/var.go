package httpTools

import "regexp"

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
// Date: 2024/6/9 15:08
//

var (
	PageCharacterSetPreg = regexp.MustCompile(`(?i)<meta[\s\S]+?charset=['"]?([\w-]+)['"][\s\S]*?>`)              // 匹配页面 字符集
	ALinkRegex           = regexp.MustCompile(`(?i)\<(?:a|mip\-link)\s[\s\S]*?href=['"]([\s\S]*?)['"][\s\S]*?\>`) // 匹配页面 所有a标签的正则表达式
	Nofollow             = regexp.MustCompile(`(?i)rel=['"]nofollow['"]`)                                         // 匹配 rel=nofollow
	TitleRegex           = regexp.MustCompile(`(?i)\<title[\s\S]*?\>([\s\S]*?)\<\/title\>`)                       // 匹配页面 title的正则表达式
)
