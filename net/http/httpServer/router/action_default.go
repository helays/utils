package router

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/dchest/captcha"
	"helay.net/go/utils/v3"
	"helay.net/go/utils/v3/close/vclose"
	"helay.net/go/utils/v3/config"
	"helay.net/go/utils/v3/logger/ulogs"
	"helay.net/go/utils/v3/net/http/mime"
	"helay.net/go/utils/v3/net/http/response"
	"helay.net/go/utils/v3/net/http/route"
	"helay.net/go/utils/v3/net/http/session"
	"helay.net/go/utils/v3/tools"
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
// Deprecated: 请使用route.Index
func (ro *Router) Index(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost || r.Method == http.MethodPut {
		response.MethodNotAllow(w)
		return
	}
	_path := r.URL.Path
	if _path == "/favicon.ico" {
		route.Favicon(w)
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
		f, _, errResp := route.HttpFS(embedFs, filepath.Join(pathCache[0], v))
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
		route.RenderErrorText(w, &route.ErrorResp{
			Code: http.StatusBadRequest,
			Msg:  "invalid URL path",
		})
		return
	}
	var embedFs http.FileSystem
	isEmbedFs := false
	// 注意，当发现响应状态非正常时，浏览器显示乱码，是标准库[http/fs.go]里面会删除Content-Encoding，
	// 所以这里不用 http.ServeFile,http.ServeFileFS
	if len(ro.staticEmbedFS) > 0 {
		for k, _embedFSInfo := range ro.staticEmbedFS {
			if strings.HasPrefix(r.URL.Path, k) {
				if _embedFSInfo.prefix != "" {
					_path = path.Join(_embedFSInfo.prefix, _path)
				}
				embedFs = http.FS(_embedFSInfo.embedFS)

				isEmbedFs = true
				break
			}
		}
	}
	if embedFs == nil {
		embedFs = http.Dir(ro.Root)
	}

	f, d, errResp := route.HttpFS(embedFs, _path)
	defer vclose.Close(f)
	if errResp != nil {
		route.RenderErrorText(w, errResp)
		return
	}
	modTime := d.ModTime()
	if isEmbedFs {
		var err error
		if modTime, err = time.ParseInLocation(time.DateTime, utils.BuildTime, config.CstSh); err != nil {
			modTime = time.Date(1994, time.January, 1, 0, 0, 0, 0, time.UTC)
		}
	}

	http.ServeContent(w, r, d.Name(), modTime, f)
}

const CaptchaID = "captcha"

// Deprecated: 请使用captcha.Text
func (ro *Router) Captcha(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	w.Header().Set("Pragma", "no-cache")
	w.Header().Set("Expires", "0")
	var content bytes.Buffer

	// 验证码存储在session中
	captchaId := captcha.NewLen(4)
	sv := session.Value{
		Field: CaptchaID,
		Value: captchaId,
		TTL:   4 * time.Minute,
	}
	if err := ro.session.Set(w, r, &sv); err != nil {
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
