package tree

import (
	"fmt"
	"strings"
)

// shift bytes in array by n bytes left
func shiftNRuneBytes(rb [4]byte, n int) [4]byte {
	switch n {
	case 0:
		return rb
	case 1:
		return [4]byte{rb[1], rb[2], rb[3], 0}
	case 2:
		return [4]byte{rb[2], rb[3]}
	case 3:
		return [4]byte{rb[3]}
	default:
		return [4]byte{}
	}
}

func countParams(path string) uint8 {
	var n uint
	for i := 0; i < len(path); i++ {
		if path[i] != ':' && path[i] != '*' {
			continue
		}
		n++
	}
	if n >= uint(maxParamCount) {
		return maxParamCount
	}

	return uint8(n)
}

// ValidateRadixPath 通过Radix实现的路由搜索树 path验证函数
func ValidateRadixPath(path string) error {
	if path == "" {
		return fmt.Errorf("路径不能为空")
	}

	if path[0] != '/' {
		return fmt.Errorf("路径必须以 '/' 开头")
	}

	// 检查连续的斜杠
	if strings.Contains(path, "//") {
		return fmt.Errorf("路径不能包含连续的斜杠 '//'")
	}

	// 检查通配符使用
	asteriskPos := strings.Index(path, "*")
	if asteriskPos != -1 {
		// 通配符必须在路径末尾
		if asteriskPos != len(path)-1 {
			// 检查是否是命名通配符
			if asteriskPos+1 < len(path) && path[asteriskPos+1] != '/' {
				// 可能是命名通配符，继续检查
				remaining := path[asteriskPos+1:]
				if strings.Contains(remaining, "*") || strings.Contains(remaining, ":") {
					return fmt.Errorf("通配符后不能包含其他通配符")
				}
			} else {
				return fmt.Errorf("通配符 '*' 只能在路径末尾使用")
			}
		}
		// 通配符前必须是斜杠
		if asteriskPos > 0 && path[asteriskPos-1] != '/' {
			return fmt.Errorf("通配符前必须有斜杠 '/'")
		}
	}

	// 检查参数使用
	colonPos := strings.Index(path, ":")
	if colonPos != -1 {
		// 参数前必须是斜杠
		if colonPos > 0 && path[colonPos-1] != '/' {
			return fmt.Errorf("参数前必须有斜杠 '/'")
		}
		// 参数必须有名称
		if colonPos == len(path)-1 {
			return fmt.Errorf("参数必须有名称")
		}
		// 参数名称不能包含斜杠
		nextSlash := strings.Index(path[colonPos:], "/")
		if nextSlash != -1 {
			paramName := path[colonPos+1 : colonPos+nextSlash]
			if paramName == "" {
				return fmt.Errorf("参数必须有非空名称")
			}
			if strings.Contains(paramName, ":") || strings.Contains(paramName, "*") {
				return fmt.Errorf("参数名称不能包含 ':' 或 '*'")
			}
		} else {
			paramName := path[colonPos+1:]
			if paramName == "" {
				return fmt.Errorf("参数必须有非空名称")
			}
			if strings.Contains(paramName, ":") || strings.Contains(paramName, "*") {
				return fmt.Errorf("参数名称不能包含 ':' 或 '*'")
			}
		}
	}

	return nil
}
