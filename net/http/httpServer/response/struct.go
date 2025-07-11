package response

type resp struct {
	Code  int `json:"code"`
	Msg   any `json:"msg"`
	Data  any `json:"data,omitempty"`
	AddOn any `json:"add_on,omitempty"`
}
