package httpServer

import (
	"errors"
	"fmt"
	"github.com/helays/utils/close/vclose"
	"github.com/helays/utils/net/http/mime"
	"github.com/helays/utils/tools"
	"io"
	"io/fs"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"
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
// Date: 2025/6/14 21:37
//

const defaultIndexPage = "index.html"

//	if ro.HttpCache {
//		w.Header().Set("cache-control", "max-age="+ro.HttpCacheMaxAge)
//		if len(files) == 1 {
//			fileInfo, _ := files[0].Stat()
//			w.Header().Set("last-modified", fileInfo.ModTime().Format(time.RFC822))
//		}
//	}
//
// 上面的后续待定
func (ro *Router) Index(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		MethodNotAllow(w)
		return
	}
	_path := r.URL.Path
	if _path == "/favicon.ico" {
		w.Header().Set("Content-Type", mime.MimeMap["ico"])
		ro.favicon(w)
		return
	}
	defaultFile := tools.Ternary(ro.Default == "", defaultIndexPage, ro.Default)
	// 根据 ?? 分割
	pathCache := strings.Split(r.URL.String(), "??")
	if len(pathCache) == 1 {
		ro.singleFile(w, r, _path, defaultFile)
		return
	}
	if len(pathCache) != 2 {
		return
	}
	fl := strings.Split(pathCache[1], ",")
	rmime := "text/html; charset=utf-8"
	isFirst := false
	for _, v := range fl {
		rp := path.Join(ro.Root, pathCache[0], v)
		f, err := os.Open(rp)
		if err != nil {
			vclose.Close(f)
			continue
		}
		if !isFirst {
			rmime = mime.MimeMap[strings.ToLower(filepath.Ext(v)[1:])]
			w.Header().Set("Content-Type", rmime)
			isFirst = true
		}
		_, _ = io.Copy(w, f)
		vclose.Close(f)
		_, _ = fmt.Fprintln(w)
	}
}

func (ro *Router) singleFile(w http.ResponseWriter, r *http.Request, _path, defaultFile string) {
	// 判断path是否以/结尾
	if strings.HasSuffix(_path, "/") && defaultFile != defaultIndexPage {
		_path = _path + defaultFile
	}
	if tools.ContainsDotDot(_path) {
		http.Error(w, "invalid URL path", http.StatusBadRequest)
		return
	}
	var embedFs http.FileSystem
	// 注意，当发现响应状态非正常时，浏览器显示乱码，是标准库[http/fs.go]里面会删除Content-Encoding，
	// 所以这里不用 http.ServeFile,http.ServeFileFS
	if !ro.dev && len(ro.staticEmbedFS) > 0 {
		for k, _embedFS := range ro.staticEmbedFS {
			if strings.HasPrefix(r.URL.Path, k) {
				embedFs = http.FS(_embedFS)
				break
			}
		}
	}
	if embedFs == nil {
		embedFs = http.Dir(ro.Root)
	}
	f, d, ok := openEmbedFsFile(w, embedFs, _path)
	if !ok {
		return
	}
	http.ServeContent(w, r, d.Name(), d.ModTime(), f)
}

func openEmbedFsFile(w http.ResponseWriter, embedFs http.FileSystem, _path string) (http.File, fs.FileInfo, bool) {
	f, err := embedFs.Open(_path)
	if err != nil {
		msg, code := toHTTPError(err)
		http.Error(w, msg, code)
		return nil, nil, false
	}
	d, _err := f.Stat()
	if _err != nil {
		msg, code := toHTTPError(err)
		http.Error(w, msg, code)
		return nil, nil, false
	}
	if d.IsDir() {
		Forbidden(w, "403 Forbidden")
		return nil, nil, false
	}
	return f, d, true
}

func toHTTPError(err error) (msg string, httpStatus int) {
	if errors.Is(err, fs.ErrNotExist) {
		return "404 page not found", http.StatusNotFound
	}
	if errors.Is(err, fs.ErrPermission) {
		return "403 Forbidden", http.StatusForbidden
	}
	// Default:
	return "500 Internal Server Error", http.StatusInternalServerError
}
