package response

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strconv"

	"helay.net/go/utils/v3/close/osClose"
	"helay.net/go/utils/v3/dataType/customWriter"
	"helay.net/go/utils/v3/net/http/httpkit"
	mime2 "helay.net/go/utils/v3/net/http/mime"
	"helay.net/go/utils/v3/net/http/request"
	"helay.net/go/utils/v3/tools"
)

func ReqError(r *http.Request, i ...any) {
	log.SetPrefix("")
	log.SetOutput(os.Stderr)
	var msg = []any{request.Getip(r), r.URL.String()}
	log.Println(append(msg, i...)...)
}

func RespHtml(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
}

func RespJson(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
}

// SetReturnData 设置返回函数
// 如果 code 异常，不想记录日志，就可以直接使用这个
func SetReturnData(w http.ResponseWriter, code int, msg any, data ...any) {
	RespJson(w)
	setHttpCode(w, code)
	r := resp{Code: code, Msg: msg}
	if len(data) > 0 {
		r.Data = data[0]
		if len(data) > 1 {
			r.AddOn = data[1:]
		}
	}
	_ = json.NewEncoder(w).Encode(r)
}

// SetReturnFile 直接讲文件反馈给前端
func SetReturnFile(w http.ResponseWriter, r *http.Request, file string) {
	f, err := os.Open(file)
	defer osClose.CloseFile(f)
	if err != nil {
		SetReturnError(w, r, err, http.StatusForbidden, "文件打开失败")
	}
	// 设置响应头
	mimeType, _ := mime2.GetFilePathMimeType(file)
	w.Header().Set("Content-Type", mimeType)
	// 对文件名进行URL转义，以支持中文等非ASCII字符
	fileName := filepath.Base(file)
	httpkit.SetDisposition(w, fileName)
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
	httpkit.SetDisposition(w, fileName)
	w.Header().Set("Content-Length", strconv.Itoa(len(*b)))
	_, _ = w.Write(*b)

}

// SetReturnError 错误信息会记录下来，同时也会反馈给前端
func SetReturnError(w http.ResponseWriter, r *http.Request, err error, code int, msg ...any) {
	if code != 404 {
		ReqError(r, append([]any{err}, msg...)...)
	}
	if err != nil {
		if len(msg) < 1 {
			msg = []any{err.Error()}
		} else {
			msg = append(msg, err.Error())
		}
	}

	RespJson(w)
	setHttpCode(w, code)
	_ = json.NewEncoder(w).Encode(resp{Code: code, Msg: tools.AnySlice2Str(msg)})
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

	setHttpCode(w, code)
	_ = json.NewEncoder(w).Encode(resp{Code: code, Msg: tools.AnySlice2Str(msg)})
}

// SetReturnErrorDisableLog 不记录日志,err 变量忽略不处理
func SetReturnErrorDisableLog(w http.ResponseWriter, err error, code int, msg ...any) {
	RespJson(w)
	setHttpCode(w, code)
	rsp := resp{Code: code, Err: err.Error()}
	ml := len(msg)
	if ml == 1 {
		rsp.Msg = msg[0]
	} else if ml > 1 {
		rsp.Msg = msg
	}
	_ = json.NewEncoder(w).Encode(rsp)
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

// CheckReqPost 检查请求是否post
func CheckReqPost(w http.ResponseWriter, r *http.Request) bool {
	if r.Method != http.MethodPost {
		Forbidden(w, "Forbidden")
		return false
	}
	return true
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
		return 0, errors.New("streaming unsupported")
	}
	return io.Copy(&flushingWriter{w: dst, flusher: flusher}, src)
}

// CopyBuffer 复制数据，并自动刷新
func CopyBuffer(dst io.Writer, src io.Reader, buf []byte) (written int64, err error) {
	// 使用 Flusher 确保数据及时发送
	flusher, ok := dst.(http.Flusher)
	if !ok {
		return 0, errors.New("streaming unsupported")
	}
	return io.CopyBuffer(&flushingWriter{w: dst, flusher: flusher}, src, buf)
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
	id, err := httpkit.QueryValid(r.URL.Query(), field, rule)
	if err != nil {
		SetReturnErrorDisableLog(w, err, http.StatusBadRequest)
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
			SetReturnError(w, r, fmt.Errorf("请求体为空"), http.StatusInternalServerError)
		} else {
			SetReturnError(w, r, err, http.StatusInternalServerError)
		}

		return nil, int(counter.TotalSize), err
	}

	if postData, ok = _postData.([]any); !ok {
		postData = []any{_postData}
	}
	return postData, int(counter.TotalSize), nil
}

func AddContentEncoding(w http.ResponseWriter, encoding string) {
	if encoding != "" {
		w.Header().Set("Content-Encoding", encoding)
	}
}

func setHttpCode(w http.ResponseWriter, code int) {
	if code != 0 && code != 200 {
		if code < 100 {
			code = http.StatusInternalServerError
		}
		w.WriteHeader(code)
	}
}
