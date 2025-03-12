package httpServer

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-playground/form/v4"
	"github.com/helays/utils/close/osClose"
	"github.com/helays/utils/close/vclose"
	"github.com/helays/utils/dataType/customWriter"
	"github.com/helays/utils/logger/ulogs"
	"github.com/helays/utils/net/http/httpTools"
	mime2 "github.com/helays/utils/net/http/mime"
	"github.com/helays/utils/tools"
	"io"
	"log"
	"mime/multipart"
	"net"
	"net/http"
	"net/textproto"
	"os"
	"path/filepath"
	"reflect"
	"regexp"
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
	m := mime2.MimeMap[fType]
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

func SetReturnCheckErr(w http.ResponseWriter, r *http.Request, err error, code int, msg any, data ...any) {
	if err == nil {
		SetReturnData(w, 0, "成功", data...)
		return
	}
	code = tools.Ternary(code != 0, code, http.StatusInternalServerError)
	SetReturnError(w, r, err, 500, msg)
}

// SetReturnCheckErrDisableLog 设置响应数据，根据err判断响应内容， 并不记录日志
func SetReturnCheckErrDisableLog(w http.ResponseWriter, r *http.Request, err error, code int, msg any, data ...any) {
	if err == nil {
		SetReturnData(w, 0, "成功", data...)
		return
	}
	code = tools.Ternary(code != 0, code, http.StatusInternalServerError)
	SetReturnErrorDisableLog(w, err, code, msg)
}

// SetReturnCheckErrWithoutError 设置响应数据，根据err判断响应内容， 不响应err信息
func SetReturnCheckErrWithoutError(w http.ResponseWriter, r *http.Request, err error, code int, msg any, data ...any) {
	if err == nil {
		SetReturnData(w, 0, "成功", data...)
		return
	}
	code = tools.Ternary(code != 0, code, http.StatusInternalServerError)
	SetReturnWithoutError(w, r, err, code, msg)
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
	if code != 0 && code != 200 && code != 404 {
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
	mimeType, _ := mime2.GetFilePathMimeType(file)
	w.Header().Set("Content-Type", mimeType)
	// 对文件名进行URL转义，以支持中文等非ASCII字符
	fileName := filepath.Base(file)
	httpTools.SetDisposition(w, fileName)
	_, _ = io.Copy(w, f)

}

// SetDownloadBytes 下载来源是字节数组
func SetDownloadBytes(w http.ResponseWriter, r *http.Request, b *[]byte, fileName string) {
	var rd io.Reader
	if len(*b) >= 512 {
		rd = bytes.NewReader((*b)[:512])
	} else {
		rd = bytes.NewReader(*b)
	}
	_m, err := mime2.GetFileMimeType(rd)
	if err != nil {
		SetReturnError(w, r, err, http.StatusInternalServerError, "下载失败")
		return
	}

	w.Header().Set("Content-Type", _m)
	httpTools.SetDisposition(w, fileName)
	w.Header().Set("Content-Length", strconv.Itoa(len(*b)))
	_, _ = w.Write(*b)

}

// SetReturnError 错误信息会记录下来，同时也会反馈给前端
func SetReturnError(w http.ResponseWriter, r *http.Request, err error, code int, msg ...any) {
	if code != 404 {
		ReqError(r, append([]any{err}, msg...)...)
	}
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
	if code != 404 {
		ReqError(r, append([]any{err}, msg...)...)
	}
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

// SetReturnErrorDisableLog 不记录日志,err 变量忽略不处理
func SetReturnErrorDisableLog(w http.ResponseWriter, err error, code int, msg ...any) {
	msg = append([]any{err.Error()}, msg...)
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
	} else if ip = r.Header.Get("HTTP_X_FORWARDED_FOR"); ip != "" {
		remoteAddr = ip
	} else if ip = r.Header.Get("X-Real-IP"); ip != "" {
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
		SetReturnErrorDisableLog(w, err, http.StatusInternalServerError, "参数解析失败", tools.MustStringReader(jd.Buffered()))
		return postData, false
	}
	return postData, true
}

type Files []File

type File struct {
	Filename string
	Size     int64
	Header   textproto.MIMEHeader
	Body     *bytes.Buffer
}

// FormDataDecode 解析表单数据并将其解码为指定类型T的实例。
// 该函数控制上传内容的大小，并处理表单数据的解析。
// 参数:
//
//	w: http.ResponseWriter，用于写入HTTP响应。
//	r: *http.Request，包含HTTP请求的详细信息。
//	size: int64，允许的最大上传大小，单位为MB。
//
// 返回值:
//
//	T: 解析后的表单数据实例。
//	bool: 表单数据是否成功解析。
func FormDataDecode[T any](w http.ResponseWriter, r *http.Request, sizes ...int64) (T, bool) {
	size := tools.Ternary(len(sizes) > 0 && sizes[0] > 0, sizes[0], 10) // 默认10M
	var formData T
	// 控制上传内容大小
	if err := r.ParseMultipartForm(size << 20); err != nil {
		SetReturnErrorDisableLog(w, err, http.StatusInternalServerError, "设置载荷大小失败")
		return formData, false
	}

	decoder := form.NewDecoder()
	if err := decoder.Decode(&formData, r.PostForm); err != nil {
		SetReturnErrorDisableLog(w, err, http.StatusInternalServerError, "参数解析失败")
		return formData, false
	}
	// 获取 T结构里面的 []字段
	t := reflect.TypeOf(formData)
	if t.Kind() != reflect.Struct {
		return formData, true
	}
	valsOf := reflect.ValueOf(&formData).Elem()
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)

		formTag := field.Tag.Get("form")
		if formTag == "" {
			continue
		}
		formTag = strings.Split(formTag, ";")[0]
		// 判断当前字段类型是否是 Files 结构体
		if field.Type != reflect.TypeOf(Files{}) {
			continue
		}
		rfs := r.MultipartForm.File[formTag]
		if len(rfs) < 1 {
			continue
		}
		var fs Files
		for _, fileHeader := range rfs {
			f, err := multipartUploader(fileHeader)
			if err != nil {
				SetReturnErrorDisableLog(w, err, http.StatusInternalServerError, "参数解析失败")
				return formData, false
			}
			f.Header = fileHeader.Header
			f.Size = fileHeader.Size
			f.Filename = fileHeader.Filename
			fs = append(fs, f)
		}
		// 这里批量上传文件
		if len(fs) > 0 {
			valsOf.FieldByName(field.Name).Set(reflect.ValueOf(fs))
		}
	}
	return formData, true
}

// multipartUploader 用于上传multipart表单中的文件。
// 参数 fileHeader 是一个指向 multipart.FileHeader 的指针，包含了上传文件的信息。
// 返回值是一个 File 类型的结构体和一个错误值。
// 如果在文件上传过程中没有错误，错误值将为 nil。
func multipartUploader(fileHeader *multipart.FileHeader) (File, error) {
	var dst File
	f, err := fileHeader.Open()
	defer vclose.Close(f)
	if err != nil {
		return dst, fmt.Errorf("打开文件%s失败:%s", fileHeader.Filename, err.Error())
	}
	dst.Body = new(bytes.Buffer)
	_, err = io.Copy(dst.Body, f)
	if err != nil {
		return dst, fmt.Errorf("复制文件%s失败:%s", fileHeader.Filename, err.Error())
	}
	return dst, nil
}

// PostQueryFieldWithValidRegexp 检查POST请求中的查询参数是否符合指定的正则表达式规则，并返回匹配结果。
func PostQueryFieldWithValidRegexp(w http.ResponseWriter, r *http.Request, field string, rule *regexp.Regexp) (string, bool) {
	if !CheckReqPost(w, r) {
		return "", false
	}
	return QueryFieldWithValidRegexp(w, r, field, rule)
}

// QueryFieldWithValidRegexp 检查查询参数是否符合指定的正则表达式规则，并返回匹配结果。
func QueryFieldWithValidRegexp(w http.ResponseWriter, r *http.Request, field string, rule *regexp.Regexp) (string, bool) {
	id, err := httpTools.QueryValid(r.URL.Query(), field, rule)
	if err != nil {
		SetReturnErrorDisableLog(w, err, http.StatusBadRequest, err.Error())
		return "", false
	}
	return id, true
}

// ParseRequestBodyAsAnySliceAndLength 解析请求体为any切片,同时获取请求体长度
func ParseRequestBodyAsAnySliceAndLength(w http.ResponseWriter, r *http.Request) ([]any, int, error) {
	var (
		_postData any
		counter   = &customWriter.SizeCounter{}
		postData  []any
		ok        bool
	)
	teeReader := io.TeeReader(r.Body, counter)
	dec := json.NewDecoder(teeReader)
	dec.UseNumber()
	if err := dec.Decode(&_postData); err != nil {
		if err == io.EOF {
			SetReturnCode(w, r, http.StatusInternalServerError, fmt.Errorf("请求体为空"))
		} else {
			SetReturnCode(w, r, http.StatusInternalServerError, err)
		}

		return nil, int(counter.TotalSize), err
	}

	if postData, ok = _postData.([]any); !ok {
		postData = []any{_postData}
	}
	return postData, int(counter.TotalSize), nil
}
