package treeconv

import (
	"testing"
)

// 1. 定义数据源结构体
type User struct {
	ID       int
	ParentID int
	Name     string
	Level    int        // 用于存储层级
	Children []UserNode // 用于存储子节点
}

// 2. 实现 ListItem 接口
func (u *User) GetID() int {
	return u.ID
}

func (u *User) GetParentID() int {
	return u.ParentID
}

// 3. 定义树节点类型
type UserNode struct {
	ID       int
	Name     string
	Level    int
	Children []UserNode
}

// 4. 实现 TreeNode 接口
func (un *UserNode) PrepareChildren(n int) {
	if un.Children == nil {
		un.Children = make([]UserNode, 0, n)
	}
}

func (un *UserNode) AddChild(node *UserNode) {
	un.Children = append(un.Children, *node)
}

func (un *UserNode) GetLastChild() (*UserNode, bool) {
	if len(un.Children) == 0 {
		return nil, false
	}
	return &un.Children[len(un.Children)-1], true
}

func (un *UserNode) ToTreeNode(src *User, level int, parentNode ...*UserNode) *UserNode {
	// 清空现有数据
	var node = &UserNode{
		ID:    src.ID,
		Name:  src.Name,
		Level: level,
	}
	return node
}

// 5. 测试用例
func TestListToTree(t *testing.T) {
	// 测试数据：平铺的用户列表
	users := []*User{
		{ID: 1, ParentID: 0, Name: "Root 1"},
		{ID: 2, ParentID: 0, Name: "Root 2"},
		{ID: 11, ParentID: 1, Name: "Child 1.1"},
		{ID: 12, ParentID: 1, Name: "Child 1.2"},
		{ID: 111, ParentID: 11, Name: "Grandchild 1.1.1"},
		{ID: 21, ParentID: 2, Name: "Child 2.1"},
		{ID: 22, ParentID: 2, Name: "Child 2.2"},
		{ID: 221, ParentID: 22, Name: "Grandchild 2.2.1"},
	}

	// 执行转换
	trees := ListToTree[*User, int, *UserNode](users)

	// 验证根节点数量
	if len(trees) != 2 {
		t.Errorf("Expected 2 root nodes, got %d", len(trees))
	}

	// 验证第一个根节点的子节点数量
	if len(trees[0].Children) != 2 {
		t.Errorf("Expected 2 children for root 1, got %d", len(trees[0].Children))
	}

	// 验证层级关系
	if trees[0].Level != 1 {
		t.Errorf("Expected level 1 for root, got %d", trees[0].Level)
	}
	if trees[0].Children[0].Level != 2 {
		t.Errorf("Expected level 2 for child, got %d", trees[0].Children[0].Level)
	}
	if trees[0].Children[0].Children[0].Level != 3 {
		t.Errorf("Expected level 3 for grandchild, got %d", trees[0].Children[0].Children[0].Level)
	}
}

// 测试空列表
func TestListToTree_EmptyList(t *testing.T) {
	var users []*User
	trees := ListToTree[*User, int, *UserNode](users)
	if len(trees) != 0 {
		t.Errorf("Expected empty tree for empty list, got %v", trees)
	}
}

// 测试只有一个节点
func TestListToTree_SingleNode(t *testing.T) {
	users := []*User{
		{ID: 1, ParentID: 0, Name: "Single Root"},
	}

	trees := ListToTree[*User, int, *UserNode](users)

	if len(trees) != 1 {
		t.Errorf("Expected 1 root node, got %d", len(trees))
	}
	if trees[0].ID != 1 {
		t.Errorf("Expected ID 1, got %d", trees[0].ID)
	}
	if len(trees[0].Children) != 0 {
		t.Errorf("Expected no children, got %d", len(trees[0].Children))
	}
}

// 测试孤儿节点（没有父节点）
func TestListToTree_OrphanNodes(t *testing.T) {
	users := []*User{
		{ID: 1, ParentID: 999, Name: "Orphan"}, // 父节点不存在
	}

	trees := ListToTree[*User, int, *UserNode](users)

	// 孤儿节点应该不会出现在结果中（因为 ParentID 999 不存在）
	if len(trees) != 0 {
		t.Errorf("Expected no root nodes for orphan nodes, got %d", len(trees))
	}
}
