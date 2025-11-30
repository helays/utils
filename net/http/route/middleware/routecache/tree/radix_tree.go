// Copyright 2013 Julien Schmidt. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be found
// in the LICENSE file.

// 基于 Julien Schmidt 的工作修改
// 原始项目：https://github.com/julienschmidt/httprouter

package tree

import (
	"fmt"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/helays/utils/v2/tools"
)

const maxParamCount = ^uint8(0)

type nodeType uint8

const (
	static nodeType = iota // default
	root
	param
	catchAll
)

type RadixTreeNode[T comparable] struct {
	path      string
	wildChild bool
	nType     nodeType
	maxParams uint8
	priority  uint32
	indices   string
	children  []*RadixTreeNode[T]
	handle    T
}

func NewRadixTree[T comparable]() *RadixTreeNode[T] {
	return &RadixTreeNode[T]{}
}

// 增加指定子节点的优先级 并 再必要时重新排序
func (n *RadixTreeNode[T]) incrementChildPriority(pos int) int {
	n.children[pos].priority++
	priority := n.children[pos].priority

	// 调整位置 （移动到前面）
	newPos := pos
	for newPos > 0 && n.children[newPos-1].priority < priority {
		// swap node positions
		n.children[newPos-1], n.children[newPos] = n.children[newPos], n.children[newPos-1]

		newPos--
	}

	// 构建新的索引字符 字符串
	if newPos != pos {
		n.indices = n.indices[:newPos] + // 未改变的前缀，可能未空
			n.indices[pos:pos+1] + // 要移动的索引字符
			n.indices[newPos:pos] + n.indices[pos+1:] // 剩下的部分，不包含 pos 处的字符
	}

	return newPos
}

// AddRoute 向路径添加一个带有有任意类型的节点
// Not concurrency-safe!
func (n *RadixTreeNode[T]) AddRoute(path string, handle T) error {
	fullPath := path
	n.priority++
	numParams := countParams(path)

	// 非空树
	if len(n.path) > 0 || len(n.children) > 0 {
	walk:
		for {
			// 更新当前节点的 maxParams
			if numParams > n.maxParams {
				n.maxParams = numParams
			}

			// 查找最长公共前缀。
			// 这也意味着公共前缀不包含 ':' 或 '*'，
			// 因为现有的键不能包含这些字符。
			i := 0

			m := tools.Min(len(path), len(n.path))
			for i < m && path[i] == n.path[i] {
				i++
			}

			// 分割边
			if i < len(n.path) {
				child := RadixTreeNode[T]{
					path:      n.path[i:],
					wildChild: n.wildChild,
					nType:     static,
					indices:   n.indices,
					children:  n.children,
					handle:    n.handle,
					priority:  n.priority - 1,
				}

				// 更新 maxParams（所有子节点的最大值）
				for j := range child.children {
					if child.children[j].maxParams > child.maxParams {
						child.maxParams = child.children[j].maxParams
					}
				}

				n.children = []*RadixTreeNode[T]{&child}
				// 使用 []byte 进行正确的 unicode 字符转换，参见 #65
				n.indices = string([]byte{n.path[i]})
				n.path = path[:i]
				var zero T
				n.handle = zero
				n.wildChild = false
			}

			// 使新节点成为此节点的子节点
			if i < len(path) {
				path = path[i:]

				if n.wildChild {
					n = n.children[0]
					n.priority++

					// 更新子节点的 maxParams
					if numParams > n.maxParams {
						n.maxParams = numParams
					}
					numParams--

					// 检查通配符是否匹配
					if len(path) >= len(n.path) && n.path == path[:len(n.path)] &&
						n.nType != catchAll && // 向 catchAll 节点添加子节点是不可能的
						(len(n.path) >= len(path) || path[len(n.path)] == '/') { // 检查更长的通配符，例如 :name 和 :names
						continue walk
					} else {
						// 通配符冲突
						var pathSeg string
						if n.nType == catchAll {
							pathSeg = path
						} else {
							pathSeg = strings.SplitN(path, "/", 2)[0]
						}
						prefix := fullPath[:strings.Index(fullPath, pathSeg)] + n.path
						return fmt.Errorf("路径段 '%s' 在新路径 '%s' 中与现有通配符 '%s' 冲突，冲突前缀为 '%s'", pathSeg, fullPath, n.path, prefix)
					}
				}

				c := path[0]

				// 参数后的斜杠
				if n.nType == param && c == '/' && len(n.children) == 1 {
					n = n.children[0]
					n.priority++
					continue walk
				}

				// 检查是否存在具有下一个路径字节的子节点
				for j := 0; j < len(n.indices); j++ {
					if c == n.indices[j] {
						j = n.incrementChildPriority(j)
						n = n.children[j]
						continue walk
					}
				}

				// 否则插入它
				if c != ':' && c != '*' {
					// 使用 []byte 进行正确的 unicode 字符转换，参见 #65
					n.indices += string([]byte{c})
					child := &RadixTreeNode[T]{
						maxParams: numParams,
					}
					n.children = append(n.children, child)
					n.incrementChildPriority(len(n.indices) - 1)
					n = child
				}
				return n.insertChild(numParams, path, fullPath, handle)

			} else if i == len(path) {
				// 使节点成为（路径内）叶子节点
				if !tools.IsZero(n.handle) {
					return fmt.Errorf("路径 '%s' 已经注册了处理函数", fullPath)
				}
				n.handle = handle
			}
			return nil
		}
	} else { // 空树
		err := n.insertChild(numParams, path, fullPath, handle)
		if err != nil {
			return err
		}
		n.nType = root
	}
	return nil
}

func (n *RadixTreeNode[T]) insertChild(numParams uint8, path, fullPath string, handle T) error {
	var offset int // 路径已处理的字节数

	// 查找前缀直到第一个通配符（以 ':' 或 '*' 开头）
	for i, m := 0, len(path); numParams > 0; i++ {
		c := path[i]
		if c != ':' && c != '*' {
			continue
		}

		// 查找通配符结束位置（'/' 或路径结束）
		end := i + 1
		for end < m && path[end] != '/' {
			switch path[end] {
			// 通配符名称不能包含 ':' 和 '*'
			case ':', '*':
				return fmt.Errorf("每个路径段只能有一个通配符，但在路径 '%s' 中发现多个: '%s'", fullPath, path[i:])
			default:
				end++
			}
		}

		// 检查如果在此处插入通配符，该节点现有的子节点是否会变得不可达
		if len(n.children) > 0 {
			return fmt.Errorf("通配符路由 '%s' 与路径 '%s' 中已存在的子节点冲突", path[i:end], fullPath)
		}

		if c == '*' {
			if end != m || numParams > 1 {
				return fmt.Errorf("在路径 '%s' 中，通配符路由只能在路径的末尾使用", fullPath)
			}

			if len(n.path) > 0 && n.path[len(n.path)-1] == '/' {
				return fmt.Errorf("在路径 '%s' 中，通配符与路径段根节点的现有处理函数冲突", fullPath)
			}

			// 当前固定宽度为 1，用于 '/'
			i--
			if path[i] != '/' {
				return fmt.Errorf("在路径 '%s' 中，通配符前缺少斜杠 '/'", fullPath)
			}

			n.path = path[offset:i]

			// 第一个节点：路径为空的全匹配节点
			child := &RadixTreeNode[T]{
				wildChild: true,
				nType:     catchAll,
				maxParams: 1,
			}
			// 更新父节点的 maxParams
			if n.maxParams < 1 {
				n.maxParams = 1
			}
			n.children = []*RadixTreeNode[T]{child}
			n.indices = string(path[i])
			n = child
			n.priority++

			// 第二个节点：保存变量的节点
			child = &RadixTreeNode[T]{
				path:      path[i:],
				nType:     catchAll,
				maxParams: 1,
				handle:    handle,
				priority:  1,
			}
			n.children = []*RadixTreeNode[T]{child}

			return nil
		}

		// 检查通配符是否有名称
		if end-i < 2 {
			return fmt.Errorf("在路径 '%s' 中，通配符必须使用非空名称", fullPath)
		}
		// 在通配符开始处分割路径
		if i > 0 {
			n.path = path[offset:i]
			offset = i
		}

		child := &RadixTreeNode[T]{
			nType:     param,
			maxParams: numParams,
		}
		n.children = []*RadixTreeNode[T]{child}
		n.wildChild = true
		n = child
		n.priority++
		numParams--
		// 如果路径不是以通配符结束，那么将会有另一个以 '/' 开头的非通配符子路径
		if end < m {
			n.path = path[offset:end]
			offset = end

			_child := &RadixTreeNode[T]{
				maxParams: numParams,
				priority:  1,
			}
			n.children = []*RadixTreeNode[T]{_child}
			n = _child
		}
	}

	// 将剩余的路径部分和处理函数插入到叶子节点
	n.path = path[offset:]
	n.handle = handle
	return nil
}

// GetValue 返回注册到给定路径（键）的处理函数。通配符的值保存到映射中。
// 如果找不到处理函数，且对于给定路径存在一个带（或不带）尾部斜杠的处理函数，
// 则会给出 TSR（尾部斜杠重定向）建议。
func (n *RadixTreeNode[T]) GetValue(path string) (handle T, p Params, tsr bool, err error) {
walk: // 用于遍历树的外部循环
	for {
		if len(path) > len(n.path) {
			if path[:len(n.path)] == n.path {
				path = path[len(n.path):]
				// 如果此节点没有通配符（参数或全匹配）子节点，
				// 我们可以直接查找下一个子节点并继续向下遍历树
				if !n.wildChild {
					c := path[0]
					for i := 0; i < len(n.indices); i++ {
						if c == n.indices[i] {
							n = n.children[i]
							continue walk
						}
					}

					// 未找到任何内容。
					// 如果存在该路径的叶子节点，我们可以建议重定向到不带尾部斜杠的相同 URL。
					tsr = path == "/" && !tools.IsZero(n.handle)
					return

				}

				// 处理通配符子节点
				n = n.children[0]
				switch n.nType {
				case param:
					// 查找参数结束位置（'/' 或路径结束）
					end := 0
					for end < len(path) && path[end] != '/' {
						end++
					}

					// 保存参数值
					if p == nil {
						// 延迟分配
						p = make(Params, 0, n.maxParams)
					}
					i := len(p)
					p = p[:i+1] // 在预分配容量内扩展切片
					p[i].Key = n.path[1:]
					p[i].Value = path[:end]

					// 我们需要继续深入！
					if end < len(path) {
						if len(n.children) > 0 {
							path = path[end:]
							n = n.children[0]
							continue walk
						}

						// ... 但是我们无法继续
						tsr = len(path) == end+1
						return
					}

					if handle = n.handle; !tools.IsZero(handle) {
						return
					} else if len(n.children) == 1 {
						// 未找到处理函数。检查是否存在此路径加上尾部斜杠的处理函数，用于 TSR 建议
						n = n.children[0]
						tsr = n.path == "/" && !tools.IsZero(n.handle)
					}

					return

				case catchAll:
					// 保存参数值
					if p == nil {
						// 延迟分配
						p = make(Params, 0, n.maxParams)
					}
					i := len(p)
					p = p[:i+1] // 在预分配容量内扩展切片

					if n.path == "/*" {
						p[i].Key = "*"
					} else {
						p[i].Key = n.path[2:]
					}
					p[i].Value = path

					handle = n.handle
					return

				default:
					err = fmt.Errorf("无效的节点类型: %v", n.nType)
					return
				}
			}
		} else if path == n.path {
			// 我们应该已经到达包含处理函数的节点。
			// 检查此节点是否注册了处理函数。
			if handle = n.handle; !tools.IsZero(handle) {
				return
			}

			if path == "/" && n.wildChild && n.nType != root {
				tsr = true
				return
			}

			// 未找到处理函数。检查是否存在此路径加上尾部斜杠的处理函数，用于尾部斜杠建议
			for i := 0; i < len(n.indices); i++ {
				if n.indices[i] == '/' {
					n = n.children[i]
					tsr = (len(n.path) == 1 && !tools.IsZero(n.handle)) || (n.nType == catchAll && !tools.IsZero(n.children[0].handle))
					return
				}
			}

			return
		}

		// 未找到任何内容。如果存在该路径的叶子节点，我们可以建议重定向到带有一个额外尾部斜杠的相同 URL
		tsr = (path == "/") ||
			(len(n.path) == len(path)+1 && n.path[len(path)] == '/' && path == n.path[:len(n.path)-1] && !tools.IsZero(n.handle))
		return
	}
}

// FindCaseInsensitivePath 对给定路径进行不区分大小写的查找，尝试找到处理函数。
// 它还可以选择性地修复尾部斜杠。
// 它返回大小写校正后的路径和一个布尔值，指示查找是否成功。
func (n *RadixTreeNode[T]) FindCaseInsensitivePath(path string, fixTrailingSlash bool) (ciPath []byte, found bool, err error) {
	return n.findCaseInsensitivePathRec(
		path,
		make([]byte, 0, len(path)+1), // 为新路径预分配足够的内存
		[4]byte{},                    // 空的 rune 缓冲区
		fixTrailingSlash,
	)
}

// 由 n.findCaseInsensitivePath 使用的递归不区分大小写查找函数
func (n *RadixTreeNode[T]) findCaseInsensitivePathRec(path string, ciPath []byte, rb [4]byte, fixTrailingSlash bool) ([]byte, bool, error) {
	npLen := len(n.path)

walk: // 用于遍历树的外部循环
	for len(path) >= npLen && (npLen == 0 || strings.EqualFold(path[1:npLen], n.path[1:])) {
		// 将公共前缀添加到结果中

		oldPath := path
		path = path[npLen:]
		ciPath = append(ciPath, n.path...)

		if len(path) > 0 {
			// 如果此节点没有通配符（参数或全匹配）子节点，
			// 我们可以直接查找下一个子节点并继续向下遍历树
			if !n.wildChild {
				// 跳过已处理的 rune 字节
				rb = shiftNRuneBytes(rb, npLen)

				if rb[0] != 0 {
					// 旧的 rune 未完成
					for i := 0; i < len(n.indices); i++ {
						if n.indices[i] == rb[0] {
							// 继续处理子节点
							n = n.children[i]
							npLen = len(n.path)
							continue walk
						}
					}
				} else {
					// 处理一个新的 rune
					var rv rune

					// 查找 rune 起始位置
					// runes 最多 4 字节长，
					// -4 肯定是另一个 rune
					var off int
					for m := tools.Min(npLen, 3); off < m; off++ {
						if i := npLen - off; utf8.RuneStart(oldPath[i]) {
							// 从缓存的路径读取 rune
							rv, _ = utf8.DecodeRuneInString(oldPath[i:])
							break
						}
					}

					// 计算当前 rune 的小写字节
					lo := unicode.ToLower(rv)
					utf8.EncodeRune(rb[:], lo)

					// 跳过已处理的字节
					rb = shiftNRuneBytes(rb, off)

					for i := 0; i < len(n.indices); i++ {
						// 小写匹配
						if n.indices[i] == rb[0] {
							// 必须使用递归方法，因为大写字节和小写字节都可能作为索引存在
							if out, found, err := n.children[i].findCaseInsensitivePathRec(
								path, ciPath, rb, fixTrailingSlash,
							); err != nil {
								return nil, false, err
							} else if found {
								return out, true, nil
							}
							break
						}
					}

					// 如果我们没有找到匹配项，同样检查大写 rune（如果它与小写不同）
					if up := unicode.ToUpper(rv); up != lo {
						utf8.EncodeRune(rb[:], up)
						rb = shiftNRuneBytes(rb, off)

						for i, c := 0, rb[0]; i < len(n.indices); i++ {
							// 大写匹配
							if n.indices[i] == c {
								// 继续处理子节点
								n = n.children[i]
								npLen = len(n.path)
								continue walk
							}
						}
					}
				}

				// 未找到任何内容。如果存在该路径的叶子节点，我们可以建议重定向到不带尾部斜杠的相同 URL
				return ciPath, fixTrailingSlash && path == "/" && !tools.IsZero(n.handle), nil
			}

			n = n.children[0]
			switch n.nType {
			case param:
				// 查找参数结束位置（'/' 或路径结束）
				k := 0
				for k < len(path) && path[k] != '/' {
					k++
				}

				// 将参数值添加到不区分大小写的路径
				ciPath = append(ciPath, path[:k]...)

				// 我们需要继续深入！
				if k < len(path) {
					if len(n.children) > 0 {
						// 继续处理子节点
						n = n.children[0]
						npLen = len(n.path)
						path = path[k:]
						continue
					}

					// ... 但是我们无法继续
					if fixTrailingSlash && len(path) == k+1 {
						return ciPath, true, nil
					}
					return ciPath, false, nil
				}

				if !tools.IsZero(n.handle) {
					return ciPath, true, nil
				} else if fixTrailingSlash && len(n.children) == 1 {
					// 未找到处理函数。检查是否存在此路径加上尾部斜杠的处理函数
					n = n.children[0]
					if n.path == "/" && !tools.IsZero(n.handle) {
						return append(ciPath, '/'), true, nil
					}
				}
				return ciPath, false, nil

			case catchAll:
				return append(ciPath, path...), true, nil

			default:
				return nil, false, fmt.Errorf("无效的节点类型: %v", n.nType)
			}
		} else {
			// 我们应该已经到达包含处理函数的节点。
			// 检查此节点是否注册了处理函数。
			if !tools.IsZero(n.handle) {
				return ciPath, true, nil
			}

			// 未找到处理函数。
			// 尝试通过添加尾部斜杠来修复路径
			if fixTrailingSlash {
				for i := 0; i < len(n.indices); i++ {
					if n.indices[i] == '/' {
						n = n.children[i]
						if (len(n.path) == 1 && !tools.IsZero(n.handle)) ||
							(n.nType == catchAll && !tools.IsZero(n.children[0].handle)) {
							return append(ciPath, '/'), true, nil
						}
						return ciPath, false, nil
					}
				}
			}
			return ciPath, false, nil
		}
	}

	// 未找到任何内容。
	// 尝试通过添加/删除尾部斜杠来修复路径
	if fixTrailingSlash {
		if path == "/" {
			return ciPath, true, nil
		}
		if len(path)+1 == npLen && n.path[len(path)] == '/' &&
			strings.EqualFold(path[1:], n.path[1:len(path)]) && !tools.IsZero(n.handle) {
			return append(ciPath, n.path...), true, nil
		}
	}
	return ciPath, false, nil
}
