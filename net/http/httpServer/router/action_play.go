package router

import (
	"net/http"
	"path/filepath"

	"github.com/helays/utils/v2/close/vclose"
	"github.com/helays/utils/v2/net/http/httpkit"
	"github.com/helays/utils/v2/net/http/route"
	"github.com/helays/utils/v2/tools"
)

//
// ━━━━━━神兽出没━━━━━━
// 　　 ┏┓     ┏┓
// 　　┏┛┻━━━━━┛┻┓
// 　　┃　　　　　 ┃
// 　　┃　　━　　　┃
// 　　┃　┳┛　┗┳  ┃
// 　　┃　　　　　 ┃
// 　　┃　　┻　　　┃
// 　　┃　　　　　 ┃
// 　　┗━┓　　　┏━┛　Code is far away from bug with the animal protecting
// 　　　 ┃　　　┃    神兽保佑,代码无bug
// 　　　　┃　　　┃
// 　　　　┃　　　┗━━━┓
// 　　　　┃　　　　　　┣┓
// 　　　　┃　　　　　　┏┛
// 　　　　┗┓┓┏━┳┓┏┛
// 　　　　 ┃┫┫ ┃┫┫
// 　　　　 ┗┻┛ ┗┻┛
//
// ━━━━━━感觉萌萌哒━━━━━━
//
//
// User helay
// Date: 2025/6/15 11:57
//

// Deprecated: 请使用 route.FilePlay 或者 route.FileDownload
func (ro *Router) Play(w http.ResponseWriter, r *http.Request, fname string, args ...any) {
	if tools.ContainsDotDot(fname) {
		http.Error(w, "invalid URL path", http.StatusBadRequest)
		return
	}
	dir, file := filepath.Split(fname)
	embedFs := http.Dir(dir)

	f, d, respErr := route.HttpFS(embedFs, file)
	defer vclose.Close(f)
	if respErr != nil {
		route.RenderErrorText(w, respErr)
		return
	}
	if len(args) > 0 && args[0] == "downloader" {
		httpkit.SetDisposition(w, file)
	}
	http.ServeContent(w, r, d.Name(), d.ModTime(), f)
}
