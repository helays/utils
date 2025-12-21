package response

// 通用响应结构体
type resp struct {
	Code  int    `json:"code"`
	Msg   any    `json:"msg,omitempty"`
	Data  any    `json:"data,omitempty"`
	Err   string `json:"err,omitempty"`
	AddOn any    `json:"add_on,omitempty"`
}

type DataList[T any] struct {
	List  []*T  `json:"list"`
	Total int64 `json:"total"`
}
