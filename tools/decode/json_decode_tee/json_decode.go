package json_decode_tee

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
)

// JsonDecodeError json解析错误
// 包含原始数据
type JsonDecodeError struct {
	Err error
	Raw bytes.Buffer
}

func (e *JsonDecodeError) Error() string {
	return fmt.Sprintf("%v", e.Err)
}

// Unwrap 返回被包装的原始错误，用于errors.Is和errors.As
func (e *JsonDecodeError) Unwrap() error {
	return e.Err
}

func ReadErrRaw(err error) []byte {
	var jsonErr *JsonDecodeError
	if errors.As(err, &jsonErr) {
		return jsonErr.Raw.Bytes()
	}
	return nil
}

// JsonDecode 解析JSON数据，同时在解析失败时保留原始数据
// 参数:
//   - rd: 数据源，实现io.Reader接口
//   - dst: 目标结构体指针，解析后的数据将存储在此
//
// 返回值:
//   - error: 解析成功返回nil，失败返回JsonDecodeError包含原始数据
//
// 示例:
//
//	var user User
//	err := JsonDecode(response.Body, &user)
//	if err != nil {
//	    var jsonErr *JsonDecodeError
//	    if errors.As(err, &jsonErr) {
//	        log.Printf("Failed to parse: %s, raw data: %s", err, jsonErr.Raw.String())
//	    }
//	}
func JsonDecode[T any](rd io.Reader, dst *T) error {
	var buf bytes.Buffer
	tee := io.TeeReader(rd, &buf)
	jd := json.NewDecoder(tee)
	if err := jd.Decode(dst); err != nil {
		return &JsonDecodeError{
			Err: err,
			Raw: buf,
		}
	}
	return nil
}
