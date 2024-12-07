package httpServer

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/helays/utils/close/osClose"
	"github.com/helays/utils/http/httpTools"
	"github.com/helays/utils/http/mime"
	"github.com/helays/utils/tools"
	"github.com/helays/utils/ulogs"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"
	"time"
)

// Play 公共函数文件
func Play(path string, w http.ResponseWriter, r *http.Request, args ...any) {
	f, err := os.OpenFile(path, os.O_RDONLY, 0644)
	defer osClose.CloseFile(f)
	if err != nil {
		ulogs.Error("文件不存在", path)
		http.NotFound(w, r)
		return
	}
	ranges := int64(0)
	rangeEnd := int(0)
	rangeSwap := strings.TrimSpace(r.Header.Get("Range"))
	if rangeSwap != "" {
		rangeSwap = rangeSwap[6:]
		rangListSwap := strings.Split(rangeSwap, "-")
		if len(rangeSwap) >= 1 {
			if num, err := strconv.Atoi(rangListSwap[0]); err == nil {
				ranges = int64(num)
			}
			if len(rangeSwap) > 1 {
				if num, err := strconv.Atoi(rangListSwap[1]); err == nil {
					rangeEnd = int(num)
				}
			}
		}
	}

	var (
		fileSize int
		//tmpF     []byte
	)
	fType := strings.ToLower(filepath.Ext(path)[1:])
	fInfo, err := f.Stat()
	if err != nil {
		Forbidden(w, "403 Forbidden!")
		return
	}

	//if fType=="mp4" {
	//	tmpF = GetMP4Duration(f)
	//	fileSize= len(tmpF)
	//	if rangeSwap!="" && ranges>0 {
	//		tmpF=tmpF[ranges:]
	//	}
	//}else{
	//
	//	if rangeSwap!="" && ranges>0 {
	//		_, _ = f.Seek(ranges, 0)
	//	}
	//	fileSize=int(fInfo.Size())
	//}
	//GetMP4Duration(f)
	if rangeSwap != "" && ranges > 0 {
		_, _ = f.Seek(ranges, 0)
	}
	fileSize = int(fInfo.Size())
	totalSize := fileSize
	if rangeSwap != "" && rangeEnd > 0 {
		totalSize = rangeEnd
	}
	total := strconv.Itoa(fileSize)
	m := mime.MimeMap[fType]
	if m == "" {
		m = "text/html;charset=utf-8"
	}
	w.Header().Set("Content-Type", m)
	w.Header().Set("Content-Length", strconv.Itoa(totalSize-int(ranges)))
	w.Header().Set("Last-Modified", fInfo.ModTime().Format(time.RFC822))
	w.Header().Set("Accept-Ranges", "bytes")
	w.Header().Set("Connection", "close")
	w.Header().Set("Etag", `W/"`+strconv.FormatInt(fInfo.ModTime().Unix(), 16)+`-`+strconv.FormatInt(fInfo.Size(), 16)+`"`)
	if len(args) > 0 {
		if args[0] == "downloader" {
			w.Header().Del("Accept-Ranges")
			httpTools.SetDisposition(w, filepath.Base(path))
		}
	}
	if rangeSwap != "" {
		w.Header().Set("Content-Range", "bytes "+strconv.Itoa(int(ranges))+"-"+strconv.Itoa(totalSize-1)+"/"+total)
		w.WriteHeader(206)
	} else {
		w.WriteHeader(200)
	}

	//if fType == "mp4" {
	//	_byt, _ := io.ReadAll(f)
	//	_, _ = w.Write(_byt)
	//	return
	//}
	_, _ = io.Copy(w, f)
}

// Forbidden 设置系统返回403
func Forbidden(w http.ResponseWriter, msg ...string) {
	w.WriteHeader(http.StatusForbidden)
	_html := `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1, shrink-to-fit=no">
    <title>403 Forbidden!</title></head><body><h3 style="text-align:center">` + strings.Join(msg, " ") + `</h3></body></html>`
	_, _ = fmt.Fprintln(w, _html)
	return
}

// NotFound 设置返回 404
func NotFound(w http.ResponseWriter, msg ...string) {
	w.WriteHeader(http.StatusNotFound)
	_html := `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1, shrink-to-fit=no">
    <title>404 Not Found!</title></head><body><h3 style="text-align:center">` + strings.Join(msg, " ") + `</h3></body></html>`
	_, _ = fmt.Fprintln(w, _html)
	return
}

// MethodNotAllow 405
func MethodNotAllow(w http.ResponseWriter) {
	w.WriteHeader(http.StatusNotFound)
	_, _ = w.Write([]byte(`<h1 style="text-align:center;">405 Error!</h1>
<p style="text-align:center;">` + http.StatusText(http.StatusMethodNotAllowed) + `</p>`))
}

func InternalServerError(w http.ResponseWriter) {
	w.WriteHeader(http.StatusInternalServerError)
	_, _ = w.Write([]byte(`<h1 style="text-align:center;">500 Error!</h1>
<p style="text-align:center;">` + http.StatusText(http.StatusInternalServerError) + `</p>`))
}

func ReqError(r *http.Request, i ...any) {
	log.SetPrefix("")
	log.SetOutput(os.Stderr)
	var msg = []any{Getip(r), r.URL.String()}
	log.Println(append(msg, i...)...)
}

func RespJson(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
}

func SetReturnCheckErr(w http.ResponseWriter, r *http.Request, err error, msg any, data ...any) {
	if err == nil {
		SetReturnData(w, 0, "成功", data...)
		return
	}
	SetReturnError(w, r, err, 500, msg)
}

// SetReturn 设置 返回函数Play
func SetReturn(w http.ResponseWriter, code int, msg ...any) {
	RespJson(w)
	if len(msg) < 1 {
		if code == 0 {
			msg = []any{"成功"}
		} else {
			msg = []any{"失败"}
		}
	}
	ulogs.Checkerr(json.NewEncoder(w).Encode(map[string]any{
		"code": code,
		"msg":  msg[0],
	}), "SetReturn")
}

// SetReturnCode 设置返回函数
// code值异常，会记录日志
func SetReturnCode(w http.ResponseWriter, r *http.Request, code int, msg any, data ...any) {
	if code != 0 && code != 200 {
		ReqError(r, code, msg)
	}
	if _, ok := msg.(error); ok {
		if len(data) > 0 && reflect.TypeOf(data[0]).String() == "bool" && !data[0].(bool) {
			msg = "系统处理失败"
		} else {
			msg = msg.(error).Error()
		}
	}
	SetReturnData(w, code, msg, data...)
}

type resp struct {
	Code int `json:"code"`
	Msg  any `json:"msg"`
	Data any `json:"data,omitempty"`
}

// SetReturnData 设置返回函数
// 如果 code 异常，不想记录日志，就可以直接使用这个
func SetReturnData(w http.ResponseWriter, code int, msg any, data ...any) {
	RespJson(w)
	if code == 0 || code == 200 {
		w.WriteHeader(200)
	} else {
		w.WriteHeader(code)
	}
	r := resp{
		Code: code,
		Msg:  msg,
	}
	if len(data) == 1 {
		r.Data = data[0]
	} else if len(data) > 1 {
		r.Data = data
	}
	ulogs.Checkerr(json.NewEncoder(w).Encode(r), "SetReturnData")
}

// SetReturnFile 直接讲文件反馈给前端
func SetReturnFile(w http.ResponseWriter, r *http.Request, file string) {
	f, err := os.Open(file)
	defer osClose.CloseFile(f)
	if err != nil {
		SetReturnError(w, r, err, http.StatusForbidden, "模板下载失败")
	}
	// 设置响应头
	mimeType, _ := mime.GetFilePathMimeType(file)
	w.Header().Set("Content-Type", mimeType)
	// 对文件名进行URL转义，以支持中文等非ASCII字符
	fileName := filepath.Base(file)
	httpTools.SetDisposition(w, fileName)
	_, _ = io.Copy(w, f)
}

// SetReturnError 错误信息会记录下来，同时也会反馈给前端
func SetReturnError(w http.ResponseWriter, r *http.Request, err error, code int, msg ...any) {
	ReqError(r, append([]any{err}, msg...)...)
	if len(msg) < 1 {
		msg = []any{err.Error()}
	} else {
		msg = append(msg, err.Error())
	}
	RespJson(w)
	if code == 0 || code == 200 {
		w.WriteHeader(200)
	} else {
		w.WriteHeader(code)
	}
	ulogs.Checkerr(json.NewEncoder(w).Encode(map[string]interface{}{
		"code": code,
		"msg":  tools.AnySlice2Str(msg),
	}), "SetReturnError")
}

// SetReturnWithoutError ，错误信息会记录下来，但是只会反馈msg
func SetReturnWithoutError(w http.ResponseWriter, r *http.Request, err error, code int, msg ...any) {
	ReqError(r, append([]any{err}, msg...)...)
	if len(msg) < 1 {
		msg = []any{"数据处理失败"}
	}
	RespJson(w)
	if code == 0 || code == 200 {
		w.WriteHeader(200)
	} else {
		w.WriteHeader(code)
	}
	ulogs.Checkerr(json.NewEncoder(w).Encode(map[string]interface{}{
		"code": code,
		"msg":  tools.AnySlice2Str(msg),
	}), "SetReturnError")
}

// CheckReqPost 检查请求是否post
func CheckReqPost(w http.ResponseWriter, r *http.Request) bool {
	if r.Method != http.MethodPost {
		Forbidden(w, "Forbidden")
		return false
	}
	return true
}

// Getip 获取客户端IP
func Getip(r *http.Request) string {
	remoteAddr := r.RemoteAddr
	if ip := r.Header.Get("HTTP_CLIENT_IP"); ip != "" {
		remoteAddr = ip
	} else if ip := r.Header.Get("HTTP_X_FORWARDED_FOR"); ip != "" {
		remoteAddr = ip
	} else if ip := r.Header.Get("X-Real-IP"); ip != "" {
		remoteAddr = ip
	} else {
		remoteAddr, _, _ = net.SplitHostPort(remoteAddr)
	}
	if remoteAddr == "::1" {
		remoteAddr = "127.0.0.1"
	}
	return remoteAddr
}

// flushingWriter 是一个带有自动 Flush 的 io.Writer
type flushingWriter struct {
	w       io.Writer
	flusher http.Flusher
}

func (fw *flushingWriter) Write(p []byte) (int, error) {
	n, err := fw.w.Write(p)
	if err != nil {
		return n, err
	}
	fw.flusher.Flush()
	return n, nil
}

// Copy 复制数据，并自动刷新
func Copy(dst io.Writer, src io.Reader) (written int64, err error) {
	// 使用 Flusher 确保数据及时发送
	flusher, ok := dst.(http.Flusher)
	if !ok {
		return 0, errors.New("Streaming unsupported!")
	}
	return io.Copy(&flushingWriter{w: dst, flusher: flusher}, src)
}

// CopyBuffer 复制数据，并自动刷新
func CopyBuffer(dst io.Writer, src io.Reader, buf []byte) (written int64, err error) {
	// 使用 Flusher 确保数据及时发送
	flusher, ok := dst.(http.Flusher)
	if !ok {
		return 0, errors.New("Streaming unsupported!")
	}
	return io.CopyBuffer(&flushingWriter{w: dst, flusher: flusher}, src, buf)
}

// JsonDecode 解析json数据
func JsonDecode[T any](w http.ResponseWriter, r *http.Request) (T, bool) {
	var postData T
	jd := json.NewDecoder(r.Body)
	err := jd.Decode(&postData)
	if err != nil && !errors.Is(err, io.EOF) {
		SetReturnError(w, r, err, http.StatusInternalServerError, "参数解析失败", tools.MustStringReader(jd.Buffered()))
		return postData, false
	}
	return postData, true
}
