package router

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/dchest/captcha"
	"github.com/helays/utils/v2/close/vclose"
	"github.com/helays/utils/v2/logger/ulogs"
	"github.com/helays/utils/v2/net/http/httpServer/http_types"
	"github.com/helays/utils/v2/net/http/httpServer/response"
	"github.com/helays/utils/v2/net/http/mime"
	"github.com/helays/utils/v2/tools"
	"io"
	"io/fs"
	"net/http"
	"path"
	"path/filepath"
	"strings"
	"time"
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
	if r.Method == http.MethodPost || r.Method == http.MethodPut {
		response.MethodNotAllow(w)
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
	w.Header().Set("X-Content-Type-Options", "nosniff")
	embedFs := http.Dir(ro.Root)
	for _, v := range fl {
		f, _, errResp := ro.openEmbedFsFile(embedFs, path.Join(pathCache[0], v))
		if errResp != nil {
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
	// 添加安全头
	w.Header().Set("X-Content-Type-Options", "nosniff")
	// 判断path是否以/结尾
	if strings.HasSuffix(_path, "/") {
		_path = _path + defaultFile
	}
	if tools.ContainsDotDot(_path) {
		ro.error(w, http_types.ErrorResp{
			Code: http.StatusBadRequest,
			Msg:  "invalid URL path",
		})
		return
	}
	var embedFs http.FileSystem
	// 注意，当发现响应状态非正常时，浏览器显示乱码，是标准库[http/fs.go]里面会删除Content-Encoding，
	// 所以这里不用 http.ServeFile,http.ServeFileFS
	if len(ro.staticEmbedFS) > 0 {
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
	f, d, errResp := ro.openEmbedFsFile(embedFs, _path)
	defer vclose.Close(f)
	if errResp != nil {
		ro.error(w, *errResp)
		return
	}

	http.ServeContent(w, r, d.Name(), d.ModTime(), f)
}

func (ro *Router) openEmbedFsFile(embedFs http.FileSystem, _path string) (http.File, fs.FileInfo, *http_types.ErrorResp) {
	f, err := embedFs.Open(_path)
	if err != nil {
		msg, code := toHTTPError(err)
		return nil, nil, &http_types.ErrorResp{
			Code: code,
			Msg:  msg,
		}
	}
	d, _err := f.Stat()
	if _err != nil {
		msg, code := toHTTPError(err)
		return nil, nil, &http_types.ErrorResp{
			Code: code,
			Msg:  msg,
		}
	}
	if d.IsDir() {
		return nil, nil, &http_types.ErrorResp{
			Code: http.StatusForbidden,
			Msg:  http.StatusText(http.StatusForbidden),
		}
	}
	return f, d, nil
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

// 显示 favicon
func (ro Router) favicon(w http.ResponseWriter) {
	w.WriteHeader(200)
	rd := bytes.NewReader(favicon[:])
	_, _ = io.Copy(w, rd)
}

func (ro Router) Captcha(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	w.Header().Set("Pragma", "no-cache")
	w.Header().Set("Expires", "0")
	var content bytes.Buffer

	// 验证码存储在session中
	captchaId := captcha.NewLen(4)
	if err := ro.Store.Set(w, r, "captcha", captchaId, 4*time.Minute); err != nil {
		response.InternalServerError(w)
		return
	}

	if err := captcha.WriteImage(&content, captchaId, 106, 40); err != nil {
		response.InternalServerError(w)
		ulogs.Error(err, "captcha writeImage")
		return
	}
	w.Header().Set("Content-Type", "image/png")
	http.ServeContent(w, r, "", time.Time{}, bytes.NewReader(content.Bytes()))
}
