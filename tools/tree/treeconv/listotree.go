package treeconv

type ListItem[T comparable, S TreeNode[S]] interface {
	GetID() T                                // 获取 ID
	GetParentID() T                          // 获取父 ID
	ToTreeNode(level int, parentNode ...S) S // 转换成树节点
}

type TreeNode[S any] interface {
	PrepareChildren(n int)   // 预分配子节点
	AddChild(Node S)         // 添加子节点
	GetLastChild() (S, bool) // 获取最后一个子节点
}

// ListToTree 列表转树
func ListToTree[L ListItem[T, S], T comparable, S TreeNode[S]](src []L) []S {
	var (
		result []S
		n      = len(src)
	)
	if n == 0 {
		return result
	}
	// 构建映射关系
	children := make(map[T][]T, n)
	idToSrc := make(map[T]L, n)
	for _, item := range src {
		parentID := item.GetParentID()
		id := item.GetID()
		children[parentID] = append(children[parentID], id)
		idToSrc[id] = item
	}
	//对每个根节点分别构建数
	var zeroID T
	for _, rootId := range children[zeroID] {
		level := 1
		rootData := idToSrc[rootId]
		rootNode := rootData.ToTreeNode(level)
		// 使用栈进行深度优先遍历
		type stackItem struct {
			node  S
			id    T
			level int
		}
		stack := []stackItem{{node: rootNode, id: rootId, level: level}}
		for len(stack) > 0 {
			// 弹出栈顶
			top := stack[len(stack)-1]
			stack = stack[:len(stack)-1]
			// 获取子节点 ID
			childIds := children[top.id]
			if len(childIds) == 0 {
				continue
			}
			// 初始化子节点切片
			top.node.PrepareChildren(len(childIds))
			// 遍历子节点
			for _, childId := range childIds {
				childData := idToSrc[childId]
				// 创建子节点
				childNode := childData.ToTreeNode(top.level+1, top.node)
				// 添加父节点
				top.node.AddChild(childNode)
				// 由于上面的SetChildren 里面会解引用，所以这里不能直接用childNode
				lastChild, ok := top.node.GetLastChild()
				// 确保 GetLastChild() 不为 nil
				if !ok {
					continue
				}
				stack = append(stack, stackItem{
					node:  lastChild,
					id:    childId,
					level: top.level + 1,
				})
			}
		}

		result = append(result, rootNode)
	}

	return result
}
