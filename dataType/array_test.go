package dataType

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"reflect"
	"strings"
	"testing"
)

func TestAnyArray_JSONMarshal(t *testing.T) {
	cases := []struct {
		name string
		arr  AnyArray[int]
		want string
	}{
		{"非空数组", NewAnyArray(1, 2, 3), `[1,2,3]`},
		{"显式空数组", NewAnyArrayFromSlice([]int{}), `[]`},
		{"未设置", AnyArray[int]{}, `null`},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			b, err := json.Marshal(c.arr)
			if err != nil {
				t.Fatalf("json.Marshal 出错: %v", err)
			}
			if string(b) != c.want {
				t.Fatalf("json.Marshal = %s, 期望 %s", b, c.want)
			}
		})
	}
}

func TestAnyArray_JSONUnmarshal(t *testing.T) {
	var nilArr AnyArray[int]
	if err := json.Unmarshal([]byte("null"), &nilArr); err != nil {
		t.Fatalf("解码 null 出错: %v", err)
	}
	if !nilArr.IsNil() {
		t.Fatal("null 应解码为未设置状态")
	}

	var emptyArr AnyArray[int]
	if err := json.Unmarshal([]byte("[]"), &emptyArr); err != nil {
		t.Fatalf("解码 [] 出错: %v", err)
	}
	if emptyArr.IsNil() {
		t.Fatal("显式空数组不应为 nil")
	}
	if emptyArr.Len() != 0 {
		t.Fatalf("空数组长度期望 0，实际 %d", emptyArr.Len())
	}

	var arr AnyArray[int]
	if err := json.Unmarshal([]byte("[1,2,3]"), &arr); err != nil {
		t.Fatalf("解码数组出错: %v", err)
	}
	if got := arr.ToSlice(); !reflect.DeepEqual(got, []int{1, 2, 3}) {
		t.Fatalf("解码结果 = %v, 期望 [1 2 3]", got)
	}

	var bad AnyArray[int]
	if err := json.Unmarshal([]byte(`{"a":1}`), &bad); err == nil {
		t.Fatal("解码 JSON 对象应返回错误")
	}
}

func TestAnyArray_JSONNestedElements(t *testing.T) {
	// []byte 元素按 base64 编码
	arr := NewAnyArray([]byte("a"), []byte("b"))
	b, err := json.Marshal(arr)
	if err != nil {
		t.Fatalf("json.Marshal 出错: %v", err)
	}
	if string(b) != `["YQ==","Yg=="]` {
		t.Fatalf("json.Marshal = %s, 期望 [\"YQ==\",\"Yg==\"]", b)
	}

	var got AnyArray[[]byte]
	if err := json.Unmarshal(b, &got); err != nil {
		t.Fatalf("json.Unmarshal 出错: %v", err)
	}
	if !reflect.DeepEqual(got.ToSlice(), arr.ToSlice()) {
		t.Fatalf("往返结果 = %v, 期望 %v", got.ToSlice(), arr.ToSlice())
	}

	// map 元素同样可往返
	m := NewAnyArray(map[string]int{"a": 1})
	mb, err := json.Marshal(m)
	if err != nil {
		t.Fatalf("json.Marshal map 元素出错: %v", err)
	}
	if string(mb) != `[{"a":1}]` {
		t.Fatalf("map 元素序列化结果 = %s", mb)
	}
}

func TestAnyArray_DBValue(t *testing.T) {
	v, err := NewAnyArray(1, 2).Value()
	if err != nil {
		t.Fatalf("Value 出错: %v", err)
	}
	if s, ok := v.(string); !ok || s != "[1,2]" {
		t.Fatalf("Value = %#v, 期望字符串 [1,2]", v)
	}

	// 未设置与空数组都写 "[]"，与改造前一致
	for name, arr := range map[string]AnyArray[int]{
		"未设置": AnyArray[int]{},
		"空数组": NewAnyArrayFromSlice([]int{}),
	} {
		v, err := arr.Value()
		if err != nil {
			t.Fatalf("%s Value 出错: %v", name, err)
		}
		if v.(string) != "[]" {
			t.Fatalf("%s Value = %v, 期望 []", name, v)
		}
	}
}

func TestAnyArray_Scan(t *testing.T) {
	var arr AnyArray[int]
	if err := arr.Scan([]byte("[1,2,3]")); err != nil {
		t.Fatalf("Scan []byte 出错: %v", err)
	}
	if got := arr.ToSlice(); !reflect.DeepEqual(got, []int{1, 2, 3}) {
		t.Fatalf("Scan 结果 = %v, 期望 [1 2 3]", got)
	}

	var fromString AnyArray[int]
	if err := fromString.Scan("[4,5]"); err != nil {
		t.Fatalf("Scan string 出错: %v", err)
	}
	if got := fromString.ToSlice(); !reflect.DeepEqual(got, []int{4, 5}) {
		t.Fatalf("Scan string 结果 = %v, 期望 [4 5]", got)
	}

	// Scan(nil) 不修改接收者，保持改造前行为
	if err := arr.Scan(nil); err != nil {
		t.Fatalf("Scan(nil) 出错: %v", err)
	}
	if got := arr.ToSlice(); !reflect.DeepEqual(got, []int{1, 2, 3}) {
		t.Fatalf("Scan(nil) 后结果 = %v, 期望保持不变", got)
	}

	var bad AnyArray[int]
	if err := bad.Scan("not-json"); err == nil {
		t.Fatal("Scan 非法 JSON 应返回错误")
	}
}

func TestAnyArray_ToSliceSharesBackingArray(t *testing.T) {
	arr := NewAnyArray(1, 2, 3)
	s := arr.ToSlice()
	s[0] = 9
	if v, err := arr.Get(0); err != nil || v != 9 {
		t.Fatalf("修改 ToSlice 返回值后 Get(0) = %d, err = %v, 期望 9", v, err)
	}
}

func TestAnyArray_CopyIsIndependent(t *testing.T) {
	arr := NewAnyArray(1, 2, 3)
	cp := arr.Copy()
	if err := cp.Set(0, 100); err != nil {
		t.Fatalf("Set 出错: %v", err)
	}
	if v, _ := arr.Get(0); v != 1 {
		t.Fatalf("修改副本后原数组 Get(0) = %d, 期望 1", v)
	}

	// nil 输入拷贝后为空非 nil 切片，与改造前一致
	empty := AnyArray[int]{}.Copy()
	if empty.IsNil() {
		t.Fatal("nil 数组 Copy 后应为空非 nil")
	}
	if !empty.IsEmpty() {
		t.Fatal("nil 数组 Copy 后应为空")
	}
}

func TestAnyArray_AppendAndLen(t *testing.T) {
	var arr AnyArray[int]
	arr.Append(1, 2)
	if arr.Len() != 2 {
		t.Fatalf("Append 后 Len = %d, 期望 2", arr.Len())
	}
	if got := arr.ToSlice(); !reflect.DeepEqual(got, []int{1, 2}) {
		t.Fatalf("Append 后内容 = %v", got)
	}

	// nil 接收者调用不 panic
	var p *AnyArray[int]
	p.Append(1)
}

func TestAnyArray_GetSetBounds(t *testing.T) {
	arr := NewAnyArray(1, 2)

	if _, err := arr.Get(-1); err == nil || !strings.Contains(err.Error(), "index out of range") {
		t.Fatalf("Get(-1) 错误 = %v", err)
	}
	if _, err := arr.Get(2); err == nil || !strings.Contains(err.Error(), "index out of range") {
		t.Fatalf("Get(2) 错误 = %v", err)
	}
	if err := arr.Set(2, 3); err == nil || !strings.Contains(err.Error(), "index out of range") {
		t.Fatalf("Set(2) 错误 = %v", err)
	}
	if err := arr.Set(0, 8); err != nil {
		t.Fatalf("Set 出错: %v", err)
	}
	if v, _ := arr.Get(0); v != 8 {
		t.Fatalf("Set 后 Get(0) = %d, 期望 8", v)
	}

	var p *AnyArray[int]
	if err := p.Set(0, 1); err == nil {
		t.Fatal("nil 接收者 Set 应返回错误")
	}
}

func TestAnyArray_ZeroValue(t *testing.T) {
	var arr AnyArray[int]
	if !arr.IsNil() {
		t.Fatal("零值 IsNil 应为 true")
	}
	if !arr.IsEmpty() {
		t.Fatal("零值 IsEmpty 应为 true")
	}
	if arr.Len() != 0 {
		t.Fatalf("零值 Len = %d, 期望 0", arr.Len())
	}
	if arr.ToSlice() != nil {
		t.Fatal("零值 ToSlice 应为 nil")
	}
}

func TestAnyArray_GobRoundTrip(t *testing.T) {
	arr := NewAnyArray(1, 2, 3)
	var buf bytes.Buffer
	if err := gob.NewEncoder(&buf).Encode(arr); err != nil {
		t.Fatalf("gob 编码出错: %v", err)
	}

	var got AnyArray[int]
	if err := gob.NewDecoder(&buf).Decode(&got); err != nil {
		t.Fatalf("gob 解码出错: %v", err)
	}
	if !reflect.DeepEqual(got.ToSlice(), arr.ToSlice()) {
		t.Fatalf("gob 往返 = %v, 期望 %v", got.ToSlice(), arr.ToSlice())
	}

	// 作为结构体字段时同样可往返
	type holder struct {
		Items AnyArray[int]
	}
	var buf2 bytes.Buffer
	if err := gob.NewEncoder(&buf2).Encode(holder{Items: arr}); err != nil {
		t.Fatalf("gob 编码结构体出错: %v", err)
	}
	var h holder
	if err := gob.NewDecoder(&buf2).Decode(&h); err != nil {
		t.Fatalf("gob 解码结构体出错: %v", err)
	}
	if !reflect.DeepEqual(h.Items.ToSlice(), arr.ToSlice()) {
		t.Fatalf("gob 结构体往返 = %v, 期望 %v", h.Items.ToSlice(), arr.ToSlice())
	}
}

func TestAnyArray_JSONInStruct(t *testing.T) {
	type holder struct {
		Items AnyArray[int]
	}

	b, err := json.Marshal(holder{Items: NewAnyArray(1, 2)})
	if err != nil {
		t.Fatalf("json.Marshal 结构体出错: %v", err)
	}
	if string(b) != `{"Items":[1,2]}` {
		t.Fatalf("结构体序列化 = %s, 期望 {\"Items\":[1,2]}", b)
	}

	var h holder
	if err := json.Unmarshal([]byte(`{"Items":[3,4]}`), &h); err != nil {
		t.Fatalf("json.Unmarshal 结构体出错: %v", err)
	}
	if got := h.Items.ToSlice(); !reflect.DeepEqual(got, []int{3, 4}) {
		t.Fatalf("结构体解码 = %v, 期望 [3 4]", got)
	}
}
