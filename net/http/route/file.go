package route

import (
	"fmt"
	"io"
	"io/fs"
	"net/http"
	"path/filepath"
	"strings"
	"time"

	"github.com/helays/utils/v2"
	"github.com/helays/utils/v2/close/vclose"
	"github.com/helays/utils/v2/config"
	"github.com/helays/utils/v2/net/http/httpkit"
	"github.com/helays/utils/v2/net/http/mime"
	"github.com/helays/utils/v2/tools"
)

// HttpFS 获取文件
func HttpFS(hfs http.FileSystem, path string) (http.File, fs.FileInfo, *ErrorResp) {
	if path == "" {
		return nil, nil, &ErrorResp{
			Code: http.StatusNotFound,
			Msg:  http.StatusText(http.StatusNotFound),
		}
	}
	f, err := hfs.Open(path)
	if err != nil {
		msg, code := toHTTPError(err)
		return nil, nil, &ErrorResp{
			Code: code,
			Msg:  msg,
		}
	}
	d, err := f.Stat()
	if err != nil {
		msg, code := toHTTPError(err)
		return nil, nil, &ErrorResp{
			Code: code,
			Msg:  msg,
		}
	}
	if d.IsDir() {
		return nil, nil, &ErrorResp{
			Code: http.StatusForbidden,
			Msg:  http.StatusText(http.StatusForbidden),
		}
	}
	return f, d, nil
}

// FilePlay 播放文件
func FilePlay(w http.ResponseWriter, r *http.Request, path string) {
	serveFile(w, r, path, "play")
}

// FileDownload 下载文件
func FileDownload(w http.ResponseWriter, r *http.Request, path string) {
	serveFile(w, r, path, "download")
}

func serveFile(w http.ResponseWriter, r *http.Request, path string, disposition string) {
	if tools.ContainsDotDot(path) {
		http.Error(w, "invalid URL path", http.StatusBadRequest)
		return
	}

	dir, file := filepath.Split(path)
	embedFs := http.Dir(dir)
	f, d, err := HttpFS(embedFs, file)
	defer vclose.Close(f)
	if err != nil {
		RenderErrorText(w, err)
		return
	}

	if disposition == "download" {
		httpkit.SetDisposition(w, file)
	}

	http.ServeContent(w, r, d.Name(), d.ModTime(), f)
}

func (ro *Route) Index(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	if ro.opt.URLPrefix != "" {
		path = strings.TrimPrefix(path, ro.opt.URLPrefix)
	}
	// noinspection SpellCheckingInspection
	w.Header().Set("X-Content-Type-Options", "nosniff") // 添加安全头，防止浏览器 mime嗅探
	if path == "/favicon.ico" {
		Favicon(w)
		return
	}

	if tools.ContainsDotDot(path) {
		RenderErrorText(w, &ErrorResp{Code: http.StatusBadRequest, Msg: "invalid URL path"})
		return
	}

	// 通过?? 拆分路径
	splitPath := strings.SplitN(path, "??", 2)
	if len(splitPath) == 1 {
		ro.singleFile(w, r, path)
		return
	}
	// 后面这是多文件请求
	files := strings.Split(splitPath[1], ",")
	ro.multipleFiles(w, r, path, files)
}

func (ro *Route) singleFile(w http.ResponseWriter, r *http.Request, path string) {
	// 检查最后一个字符 是不是 /
	if len(path) > 0 && path[len(path)-1] == '/' {
		path = path + ro.opt.Index
	}
	hfs, enableEmbed := ro.loadFile(path)
	// 获取文件
	f, d, err := HttpFS(hfs, path)
	defer vclose.Close(f)
	if err != nil {
		RenderErrorText(w, err)
		return
	}
	modTime := d.ModTime()
	if enableEmbed {
		var _err error
		if modTime, _err = time.ParseInLocation(time.DateTime, utils.BuildTime, config.CstSh); _err != nil {
			modTime = time.Date(1994, time.January, 1, 0, 0, 0, 0, time.UTC)
		}
	}
	http.ServeContent(w, r, d.Name(), modTime, f)
}

func (ro *Route) multipleFiles(w http.ResponseWriter, r *http.Request, path string, files []string) {
	if len(files) == 0 {
		RenderErrorText(w, &ErrorResp{Code: http.StatusBadRequest, Msg: "invalid URL path"})
		return
	}
	isFirst := true
	mimeStr := "text/html; charset=utf-8"
	for _, file := range files {
		tmp := filepath.Join(path, file)
		hfs, _ := ro.loadFile(tmp)
		f, _, errResp := HttpFS(hfs, tmp)
		if errResp != nil {
			vclose.Close(f)
			RenderErrorText(w, errResp)
			return
		}
		if isFirst {
			ext := filepath.Ext(file)
			if len(ext) > 1 {
				if m, ok := mime.MimeMap[strings.ToLower(ext[1:])]; ok {
					mimeStr = m
				}
			}
			w.Header().Set("Content-Type", mimeStr)
			isFirst = false
		}
		_, _ = io.Copy(w, f)
		vclose.Close(f)
		_, _ = fmt.Fprintln(w)
	}

}

func (ro *Route) loadFile(path string) (http.FileSystem, bool) {
	var (
		hfs         http.FileSystem
		enableEmbed = false
	)

	for _, cache := range ro.embed {
		if strings.HasPrefix(path, cache.Search) {
			if cache.Prefix != "" {
				path = filepath.Join(cache.Prefix, path)
			}
			hfs = http.FS(cache.FS)
			enableEmbed = true
			break
		}
	}

	// 判断 hfs 是否是空
	if hfs == nil {
		hfs = http.Dir(ro.opt.Root)
	}
	return hfs, enableEmbed
}
