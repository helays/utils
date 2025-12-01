package response

type resp struct {
	Code  int    `json:"code"`
	Msg   any    `json:"msg,omitempty"`
	Data  any    `json:"data,omitempty"`
	Err   string `json:"err,omitempty"`
	AddOn any    `json:"add_on,omitempty"`
}
