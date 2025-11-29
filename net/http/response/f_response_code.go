package response

import (
	"encoding/json"
	"net/http"

	"github.com/helays/utils/v2/logger/ulogs"
	"github.com/helays/utils/v2/net/http/request"
)

type RespCode interface {
	HttpCode() int   // 设置 http 响应状态码
	RespCode() int   // 设置业务响应码
	Message() string // 设置业务响应信息
	EnableLog() bool // 是否记录日志
}

// RespErrWithCode 根据业务上配置的错误码，进行数据响应
func RespErrWithCode(w http.ResponseWriter, r *http.Request, resp RespCode, err error, d ...any) {
	RespJson(w)
	w.WriteHeader(resp.HttpCode())
	msg := resp.Message()
	respData := map[string]any{
		"code": resp.RespCode(),
		"msg":  msg,
		"err":  err.Error(),
	}
	dl := len(d)
	if dl == 1 {
		respData["data"] = d[0]
	} else if dl > 1 {
		respData["data"] = d
	}
	if resp.EnableLog() {
		ulogs.Errorf("IP[%s]调用接口接口[%s]失败，描述[%s]，错误信息 %v", request.Getip(r), r.URL.String(), msg, err)
	}
	_ = json.NewEncoder(w).Encode(respData)

}

// RespWithCode 根据业务上配置的响应码，进行数据响应
func RespWithCode(w http.ResponseWriter, resp RespCode, d ...any) {
	RespJson(w)
	w.WriteHeader(resp.HttpCode())
	respData := map[string]any{
		"code": resp.RespCode(),
		"msg":  resp.Message(),
	}
	dl := len(d)
	if dl == 1 {
		respData["data"] = d[0]
	} else if dl > 1 {
		respData["data"] = d
	}
	_ = json.NewEncoder(w).Encode(respData)

}
