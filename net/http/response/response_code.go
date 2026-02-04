package response

import (
	"encoding/json"
	"net/http"

	"helay.net/go/utils/v3/logger/ulogs"
	"helay.net/go/utils/v3/net/http/request"
)

type RespCode interface {
	HttpCode() int   // 设置 http 响应状态码
	RespCode() int   // 设置业务响应码
	Message() string // 设置业务响应信息
	EnableLog() bool // 是否记录日志

}

// RespErrWithCode 根据业务上配置的错误码，进行数据响应
func RespErrWithCode(w http.ResponseWriter, r *http.Request, respInfo RespCode, err error, d ...any) {
	RespJson(w)
	w.WriteHeader(respInfo.HttpCode())
	msg := respInfo.Message()
	respData := resp{
		Code: respInfo.RespCode(),
		Msg:  msg,
	}
	if err != nil {
		respData.Err = err.Error()
	}
	dl := len(d)
	if dl == 1 {
		respData.Data = d[0]
	} else if dl > 1 {
		respData.Data = d
	}
	if respInfo.EnableLog() {
		ulogs.Errorf("IP[%s]调用接口接口[%s]失败，描述[%s]，错误信息 %v", request.Getip(r), r.URL.String(), msg, err)
	}
	_ = json.NewEncoder(w).Encode(respData)

}

// RespWithCode 根据业务上配置的响应码，进行数据响应
func RespWithCode(w http.ResponseWriter, respInfo RespCode, d ...any) {
	RespJson(w)
	w.WriteHeader(respInfo.HttpCode())

	respData := resp{
		Code: respInfo.RespCode(),
		Msg:  respInfo.Message(),
	}
	dl := len(d)
	if dl == 1 {
		respData.Data = d[0]
	} else if dl > 1 {
		respData.Data = d
	}
	_ = json.NewEncoder(w).Encode(respData)

}
