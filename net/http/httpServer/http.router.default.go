package httpServer

import (
	"fmt"
	"github.com/helays/utils/close/httpClose"
	"github.com/helays/utils/close/vclose"
	"github.com/helays/utils/logger/ulogs"
	"github.com/helays/utils/net/http/mime"
	"github.com/helays/utils/tools"
	"io"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"strconv"
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
	// 注意，当发现响应状态非正常时，浏览器显示乱码，是标准库[http/fs.go]里面会删除Content-Encoding，需要处理下
	if !ro.dev && len(ro.staticEmbedFS) > 0 {
		for k, embedFS := range ro.staticEmbedFS {
			if strings.HasPrefix(r.URL.Path, k) {
				http.ServeFileFS(w, r, embedFS, _path)
				return
			}
		}
	}
	http.ServeFile(w, r, path.Join(ro.Root, _path))
}

// Index 默认页面
func (ro Router) Indexs(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		MethodNotAllow(w)
		return
	}
	files, path, status := SetRequestDefaultPage(ro.Default, ro.Root+r.URL.String())
	defer func() {
		if files != nil {
			for _, item := range files {
				if item != nil {
					_ = item.Close()
				}
			}
		}
		httpClose.CloseReq(r)
	}()
	var rmime string
	if len(path) < 1 {
		rmime = "text/html; charset=utf-8"
	} else {
		if len(filepath.Ext(path[0])) > 0 {
			rmime = mime.MimeMap[strings.ToLower(filepath.Ext(path[0])[1:])]
		}
		if rmime == "" {
			rmime = "text/html; charset=utf-8"
		}
	}
	w.Header().Set("Content-Type", rmime)
	if !status {
		if r.URL.Path == "/favicon.ico" {
			ro.favicon(w)
			return
		}
		NotFound(w, "404 not found")
		return
	}

	if len(files) == 1 {
		fileInfo, _ := files[0].Stat()
		fileSize := int(fileInfo.Size())
		total := strconv.Itoa(fileSize)
		w.Header().Set("last-modified", fileInfo.ModTime().Format(time.RFC822))
		w.Header().Set("Accept-Ranges", "bytes")

		ranges := int64(0)
		rangeSwap := strings.TrimSpace(r.Header.Get("Range"))
		if rangeSwap != "" {
			rangeSwap = rangeSwap[6:]
			rangeListSwap := strings.Split(rangeSwap, "-")
			if len(rangeListSwap) == 2 {
				if num, err := strconv.Atoi(rangeListSwap[0]); err == nil {
					ranges = int64(num)
				}
			}
		}
		w.Header().Set("Content-Length", strconv.Itoa(fileSize-int(ranges)))
		_, _ = files[0].Seek(ranges, 0)
		w.Header().Set("Etag", `W/"`+strconv.FormatInt(fileInfo.ModTime().Unix(), 16)+`-`+strconv.FormatInt(fileInfo.Size(), 16)+`"`)

		if ranges > 0 {
			w.Header().Set("Content-Range", "bytes "+strconv.Itoa(int(ranges))+"-"+strconv.Itoa(fileSize-1)+"/"+total) // 允许 range
			w.WriteHeader(206)
		} else {
			w.WriteHeader(200)
		}
	} else {
		w.WriteHeader(200)
	}

	if ro.HttpCache {
		w.Header().Set("cache-control", "max-age="+ro.HttpCacheMaxAge)
		if len(files) == 1 {
			fileInfo, _ := files[0].Stat()
			w.Header().Set("last-modified", fileInfo.ModTime().Format(time.RFC822))
		}
	}
	for _, file := range files {
		_, _ = io.Copy(w, file)
		_, _ = fmt.Fprintln(w)
	}
}

// SetRequestDefaultPage 设置 打开的默认页面
// defaultPage string 默认打开页面
// root 网站更目录
// path string
func SetRequestDefaultPage(defaultPage, path string) ([]*os.File, []string, bool) {
	sarr := strings.Split(path, "??")
	if len(sarr) == 1 {
		swapUrl, err := url.Parse(path)
		if err != nil {
			ulogs.Error("url 异常", err)
			return nil, nil, false
		}
		path = swapUrl.Path
		if filepath.Base(path) == "lib.js" {

		}
		f, err := os.OpenFile(path, os.O_RDONLY, 0644)
		if err != nil {
			return []*os.File{f}, []string{path}, false
		}
		fInfo, _ := f.Stat()
		if !fInfo.IsDir() {
			return []*os.File{f}, []string{path}, true
		}
		defaultPage = strings.TrimSpace(defaultPage)
		if defaultPage == "" {
			defaultPage = "index.html"
		}
		fp := path + "/" + defaultPage

		if strings.HasSuffix(path, "/") {
			fp = path + defaultPage
		}
		f, err = os.OpenFile(fp, os.O_RDONLY, 0644)
		if err != nil {
			return []*os.File{f}, []string{fp}, false
		}
		return []*os.File{f}, []string{fp}, true
	}

	var (
		swapList  []*os.File
		swapPaths []string
		status    bool
	)
	for _, item := range strings.Split(sarr[1], ",") {
		swapFile, swapPath, swapStatus := SetRequestDefaultPage(defaultPage, sarr[0]+item)
		if !swapStatus {
			continue
		}
		swapList = append(swapList, swapFile...)
		swapPaths = append(swapPaths, swapPath...)
	}
	if len(swapList) > 0 {
		status = true
	}
	return swapList, swapPaths, status
}
