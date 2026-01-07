package tools

import (
	"fmt"
	"strings"
)

// ToHttpRouterSyntax 将 {field} 和 {*field} 格式转换为 :field 和 *field 格式
// 输入示例： "/users/{id}/profile/{*filepath}" -> "/users/:id/profile/*filepath"
func ToHttpRouterSyntax(path string) (string, error) {
	var result strings.Builder
	result.Grow(len(path))

	inWildcard := false
	isCatchAll := false

	for i := 0; i < len(path); i++ {
		c := path[i]

		switch {
		case c == '{':
			if inWildcard {
				return "", fmt.Errorf("非法嵌套的花括号，位置: %d", i)
			}
			inWildcard = true
			isCatchAll = false
			// 检查下一个字符是否是 '*'（表示全匹配参数）
			if i+1 < len(path) && path[i+1] == '*' {
				isCatchAll = true
				i++ // 跳过 '*'
			}
			result.WriteByte(':')
			if isCatchAll {
				result.WriteByte('*')
			}

		case c == '}':
			if !inWildcard {
				return "", fmt.Errorf("未匹配的右花括号，位置: %d", i)
			}
			inWildcard = false
			// 不需要写入 '}'

		default:
			result.WriteByte(c)
		}
	}

	if inWildcard {
		return "", fmt.Errorf("未闭合的花括号")
	}

	return result.String(), nil
}

// ToBraceSyntax 将 :field 和 *field 格式转换为 {field} 和 {*field} 格式
// 输入示例： "/users/:id/profile/*filepath" -> "/users/{id}/profile/{*filepath}"
func ToBraceSyntax(path string) (string, error) {
	var result strings.Builder
	result.Grow(len(path) + 10) // 预分配额外空间给花括号

	for i := 0; i < len(path); i++ {
		c := path[i]

		switch {
		case c == ':':
			// 检查是否是通配符
			if i+1 < len(path) && path[i+1] == '*' {
				// 这是 :*field 格式（全匹配参数）
				result.WriteString("{*")
				i++ // 跳过 ':'
				// 跳过 '*' 但记录参数名
				paramStart := i + 1
				for i+1 < len(path) && path[i+1] != '/' {
					i++
				}
				if paramStart > i {
					return "", fmt.Errorf("通配符参数名缺失，位置: %d", i)
				}
				result.WriteString(path[paramStart : i+1])
				result.WriteByte('}')
			} else {
				// 这是 :field 格式（单段参数）
				result.WriteByte('{')
				paramStart := i + 1
				for i+1 < len(path) && path[i+1] != '/' {
					i++
				}
				if paramStart > i {
					return "", fmt.Errorf("参数名缺失，位置: %d", i)
				}
				result.WriteString(path[paramStart : i+1])
				result.WriteByte('}')
			}

		default:
			result.WriteByte(c)
		}
	}

	return result.String(), nil
}
