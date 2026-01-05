package tools_test

import (
	"testing"

	"github.com/helays/utils/v2/tools"
)

type User struct {
	ID   int
	Name string
	Age  int
}

type Users []User

func (u *Users) Len() int {
	return len(*u)
}

func (u *Users) Range(f func(Idx int, v User)) {
	for i, v := range *u {
		f(i, v)
	}
}

func (u *Users) ExtractKey(v User) int {
	return v.ID
}

func TestDiffRangeable(t *testing.T) {
	// 准备测试数据
	srcUsers := Users{
		{ID: 1, Name: "Alice", Age: 25},
		{ID: 2, Name: "Bob", Age: 30},
		{ID: 3, Name: "Charlie", Age: 35},
		{ID: 4, Name: "David", Age: 40},
	}
	dstUsers := Users{
		{ID: 2, Name: "Bob Updated", Age: 31}, // 共同的，但有更新
		{ID: 3, Name: "Charlie", Age: 35},     // 共同的，完全相同
		{ID: 5, Name: "Eve", Age: 28},         // 只在 dst 中存在
		{ID: 6, Name: "Frank", Age: 45},       // 只在 dst 中存在
	}
	inSrc, common, inDst := tools.DiffRangeable(&srcUsers, &dstUsers)
	t.Logf("inSrc: %v\n", inSrc)
	t.Logf("common: %v\n", common)
	t.Logf("inDst: %v\n", inDst)
}
