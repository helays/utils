package treeconv

// ListItem 列表项接口
// T 用于接收数据源的唯一标识，需要可比较
// S 是需要转换的目标数据结构
// M 是ListItem数据源本身
type ListItem[T comparable, S TreeNode[M, S], M any] interface {
	GetID() T       // 获取 ID
	GetParentID() T // 获取父 ID
}

// TreeNode 树节点接口
// M 是数据源，将数据换转传承TreeNode
// S 是树节点本身
type TreeNode[M any, S any] interface {
	PrepareChildren(n int)                          // 预分配子节点
	AddChild(Node S)                                // 添加子节点
	GetLastChild() (S, bool)                        // 获取最后一个子节点
	ToTreeNode(src M, level int, parentNode ...S) S // 转换成树节点
}

// ListToTree 列表转树
func ListToTree[L ListItem[T, S, L], T comparable, S TreeNode[L, S]](src []L, startNode ...T) []S {
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
	var tempNode S
	//对每个根节点分别构建数
	var zeroID T
	if len(startNode) > 0 {
		zeroID = startNode[0]
	}
	for _, rootId := range children[zeroID] {
		level := 1
		rootData := idToSrc[rootId]
		rootNode := tempNode.ToTreeNode(rootData, level)
		//rootNode := rootData.ToTreeNode(level)
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
				childNode := tempNode.ToTreeNode(childData, top.level+1, top.node)
				//childNode := childData.ToTreeNode(top.level+1, top.node)
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
