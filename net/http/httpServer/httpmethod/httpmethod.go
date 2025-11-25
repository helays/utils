// Package httpmethod package httpmethod
package httpmethod

import (
	"net/http"

	"github.com/helays/utils/v2/net/http/httpServer/response"
)

// Method 验证HTTP方法是否匹配
func Method(method string, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != method {
			response.MethodNotAllow(w)
			return
		}
		next.ServeHTTP(w, r)
	})
}

// Get 只允许GET请求
func Get(next http.HandlerFunc) http.Handler {
	return Method(http.MethodGet, next)
}

// Post 只允许POST请求
func Post(next http.HandlerFunc) http.Handler {
	return Method(http.MethodPost, next)
}

// Put 只允许PUT请求
func Put(next http.HandlerFunc) http.Handler {
	return Method(http.MethodPut, next)
}

// Delete 只允许DELETE请求
func Delete(next http.HandlerFunc) http.Handler {
	return Method(http.MethodDelete, next)
}

// Patch 只允许PATCH请求
func Patch(next http.HandlerFunc) http.Handler {
	return Method(http.MethodPatch, next)
}

// Head 只允许HEAD请求
func Head(next http.HandlerFunc) http.Handler {
	return Method(http.MethodHead, next)
}

// Options 只允许OPTIONS请求
func Options(next http.HandlerFunc) http.Handler {
	return Method(http.MethodOptions, next)
}

// Connect 只允许CONNECT请求
func Connect(next http.HandlerFunc) http.Handler {
	return Method(http.MethodConnect, next)
}

// Trace 只允许TRACE请求
func Trace(next http.HandlerFunc) http.Handler {
	return Method(http.MethodTrace, next)
}

// Any 允许任何HTTP方法
func Any(next http.HandlerFunc) http.Handler {
	return next
}

// Methods 允许指定的多个HTTP方法
func Methods(methods []string, next http.HandlerFunc) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		for _, m := range methods {
			if r.Method == m {
				next.ServeHTTP(w, r)
				return
			}
		}
		response.MethodNotAllow(w)
	})
}
