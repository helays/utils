package route

import (
	"net/http"
)

func Favicon(w http.ResponseWriter) {
	// 预设置头部
	w.Header().Set("Content-Type", "image/x-icon")
	w.Header().Set("Cache-Control", "public, max-age=864000")
	_, _ = w.Write(favicon)
}
