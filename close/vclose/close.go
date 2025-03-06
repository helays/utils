package vclose

import (
	"io"
)

// 这是一个通用的牛逼的关闭资源的封装

// Close 通用关闭函数，不用操心有没有报错的情况
func Close(conn io.Closer) {
	if conn == nil {
		return
	}
	defer func() {
		recover()
	}()
	_ = conn.Close()
}
