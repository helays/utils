package apoolTools

import "reflect"

type RemovePasswd interface {
	RemovePasswd()
}

func RemovePasswdFunc(opt any) {
	if opt == nil {
		return
	}
	// 如果opt 实现了RemovePasswd接口，则调用RemovePasswd方法
	if rp, ok := opt.(RemovePasswd); ok {
		val := reflect.ValueOf(rp)
		if !val.IsNil() {
			rp.RemovePasswd()
		}
	}
}
