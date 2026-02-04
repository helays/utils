package response

import (
	"fmt"
	"net/http"

	"helay.net/go/utils/v3/net/http/request"
	"helay.net/go/utils/v3/tools/decode/json_decode_tee"
)

func JsonDecodeResp[T any](w http.ResponseWriter, r *http.Request) (T, bool) {
	postData, err := request.JsonDecode[T](r)
	if err != nil {
		SetReturnErrorDisableLog(w, fmt.Errorf(err.Error()), http.StatusInternalServerError, err.(*json_decode_tee.JsonDecodeError).Raw.String())
		return postData, false
	}
	return postData, true
}

func JsonDecodePtrResp[T interface{ *E }, E any](w http.ResponseWriter, r *http.Request) (T, bool) {
	postData, err := request.JsonDecodePtr[T](r)
	if err != nil {
		SetReturnErrorDisableLog(w, fmt.Errorf(err.Error()), http.StatusInternalServerError, err.(*json_decode_tee.JsonDecodeError).Raw.String())
		return postData, false
	}
	return postData, true
}

// FormDataDecodeResp 解析表单数据并将其解码为指定类型T的实例。
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
func FormDataDecodeResp[T any](w http.ResponseWriter, r *http.Request, sizes ...int64) (*T, bool) {
	data, err := request.FormDataDecode[T](r, sizes...)
	if err != nil {
		SetReturnErrorDisableLog(w, err, http.StatusInternalServerError)
		return data, false
	}
	return data, true
}
