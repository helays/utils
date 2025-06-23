package http_types

type ErrorResp struct {
	Code  int
	Msg   string
	Error error
}
