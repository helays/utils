package httpServer

import (
	"github.com/helays/utils/close/vclose"
	"github.com/helays/utils/net/http/httpTools"
	"github.com/helays/utils/tools"
	"net/http"
	"path/filepath"
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

func (ro *Router) Play(w http.ResponseWriter, r *http.Request, fname string, args ...any) {
	if r.Method == "POST" {
		MethodNotAllow(w)
		return
	}
	if tools.ContainsDotDot(fname) {
		http.Error(w, "invalid URL path", http.StatusBadRequest)
		return
	}
	dir, file := filepath.Split(fname)
	embedFs := http.Dir(dir)
	f, d, respErr := ro.openEmbedFsFile(embedFs, file)
	defer vclose.Close(f)
	if respErr != nil {
		ro.error(w, *respErr)
		return
	}
	if len(args) > 0 && args[0] == "downloader" {
		httpTools.SetDisposition(w, file)
	}
	http.ServeContent(w, r, d.Name(), d.ModTime(), f)
}
