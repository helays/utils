package response

import (
	"net/http"
	"strings"
)

// MethodNotAllow 405
func MethodNotAllow(w http.ResponseWriter) {
	http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
}
func InternalServerError(w http.ResponseWriter) {
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

// NotFound 设置返回 404
func NotFound(w http.ResponseWriter, msg ...string) {
	errMsg := http.StatusText(http.StatusNotFound)
	if len(msg) > 0 {
		errMsg = strings.Join(msg, "<br>")
	}
	http.Error(w, errMsg, http.StatusNotFound)
}

// Forbidden 设置系统返回403
func Forbidden(w http.ResponseWriter, msg ...string) {
	errMsg := http.StatusText(http.StatusForbidden)
	if len(msg) > 0 {
		errMsg = strings.Join(msg, "<br>")
	}
	http.Error(w, errMsg, http.StatusForbidden)
}

func Status(w http.ResponseWriter, code int, msg ...string) {
	errMsg := http.StatusText(code)
	if len(msg) > 0 {
		errMsg = strings.Join(msg, "<br>")
	}
	http.Error(w, errMsg, code)
}
