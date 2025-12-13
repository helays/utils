package request

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"reflect"
	"strings"

	"github.com/go-playground/form/v4"
	"github.com/gorilla/schema"
	"github.com/helays/utils/v2/close/vclose"
	"github.com/helays/utils/v2/tools"
	"github.com/helays/utils/v2/tools/decode/json_decode_tee"
)

// QueryDecode 解析query数据
func QueryDecode(u url.Values, dst any) error {
	decoder := schema.NewDecoder()
	decoder.SetAliasTag("query")
	decoder.IgnoreUnknownKeys(true)
	return decoder.Decode(dst, u)
}

// FormDecode 解析 application/x-www-form-urlencoded 编码的 form 数据
func FormDecode(r *http.Request, dst any) error {
	decoder := schema.NewDecoder()
	decoder.SetAliasTag("form")
	decoder.IgnoreUnknownKeys(true)
	return decoder.Decode(dst, r.PostForm)
}

// JsonDecode 解析json数据
// 值类型，适合小结构体，当字段少于10的时候，缺点是返回时会复制整个结构体
func JsonDecode[T any](r *http.Request) (T, error) {
	var postData T
	err := json_decode_tee.JsonDecode(r.Body, &postData)
	return postData, err
}

// JsonDecodePtr 解析json数据
// 处理指针类型（调用方需确保T是指针类型，如 *YourStruct）
func JsonDecodePtr[T interface{ *E }, E any](r *http.Request, target ...T) (T, error) {
	// 处理目标指针
	var postData T
	if len(target) > 0 {
		postData = target[0]
	} else {
		postData = new(E)
	}
	err := json_decode_tee.JsonDecode(r.Body, &postData)

	return postData, err
}

func FormDataDecode[T any](r *http.Request, sizes ...int64) (*T, error) {
	size := tools.Ternary(len(sizes) > 0 && sizes[0] > 0, sizes[0], 10) // 默认10M
	var formData T
	// 控制上传内容大小
	if err := r.ParseMultipartForm(size << 20); err != nil {
		return nil, fmt.Errorf("设置载荷大小失败 %v", err)
	}

	decoder := form.NewDecoder()
	if err := decoder.Decode(&formData, r.PostForm); err != nil {
		return nil, fmt.Errorf("参数解析失败 %v", err)
	}
	// 获取 T结构里面的 []字段
	t := reflect.TypeOf(formData)
	if t.Kind() != reflect.Struct {
		return &formData, nil
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
				return nil, fmt.Errorf("参数解析失败 %v", err)
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
	return &formData, nil
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
