package response

import (
	"net/http"
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
	http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
}

// Forbidden 设置系统返回403
func Forbidden(w http.ResponseWriter, msg ...string) {
	http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
}
