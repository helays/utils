package tools

import (
	"reflect"
)

func IsStruct(ipt any) (name string, isStruct bool, isAnonymous bool, dst reflect.Type) {
	t := reflect.TypeOf(ipt)
	// 处理指针类型
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	// 检查是否为结构体
	if t.Kind() != reflect.Struct {
		return "", false, false, t
	}

	// 获取类型名称
	// 获取类型名称
	name = t.Name()
	isAnonymous = name == ""

	// 如果是匿名结构体，返回空名称
	if isAnonymous {
		return "", true, true, t
	}
	return name, true, false, t
}
