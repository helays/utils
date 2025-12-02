package response

import (
	"net/http"
	"strings"
)

// 常用状态码的文本缓存
var (
	statusText400 = http.StatusText(http.StatusBadRequest)
	statusText401 = http.StatusText(http.StatusUnauthorized)
	statusText403 = http.StatusText(http.StatusForbidden)
	statusText404 = http.StatusText(http.StatusNotFound)
	statusText405 = http.StatusText(http.StatusMethodNotAllowed)
	statusText500 = http.StatusText(http.StatusInternalServerError)
)

// MethodNotAllow 405
func MethodNotAllow(w http.ResponseWriter) {
	http.Error(w, statusText405, http.StatusMethodNotAllowed)
}
func InternalServerError(w http.ResponseWriter) {
	http.Error(w, statusText500, http.StatusInternalServerError)
}

// NotFound 设置返回 404
func NotFound(w http.ResponseWriter, msg ...string) {
	errMsg := statusText404
	if len(msg) > 0 {
		errMsg = strings.Join(msg, "<br>")
	}
	http.Error(w, errMsg, http.StatusNotFound)
}

// Forbidden 设置系统返回403
func Forbidden(w http.ResponseWriter, msg ...string) {
	errMsg := statusText403
	if len(msg) > 0 {
		errMsg = strings.Join(msg, "<br>")
	}
	http.Error(w, errMsg, http.StatusForbidden)
}

// BadRequest 400 Bad Request （建议新增）
func BadRequest(w http.ResponseWriter, msg ...string) {
	errMsg := statusText400
	if len(msg) > 0 {
		errMsg = strings.Join(msg, "<br>")
	}
	http.Error(w, errMsg, http.StatusBadRequest)
}

// Unauthorized 401 Unauthorized （建议新增）
func Unauthorized(w http.ResponseWriter, msg ...string) {
	errMsg := statusText401
	if len(msg) > 0 {
		errMsg = strings.Join(msg, "<br>")
	}
	http.Error(w, errMsg, http.StatusUnauthorized)
}

func Status(w http.ResponseWriter, code int, msg ...string) {
	errMsg := http.StatusText(code)
	if len(msg) > 0 {
		errMsg = strings.Join(msg, "<br>")
	}
	http.Error(w, errMsg, code)
}
