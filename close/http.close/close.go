package http_close

import (
	"io"
	"net/http"
)

func CloseResp(resp *http.Response) {
	if resp == nil || resp.Body == nil {
		return
	}
	_, _ = io.Copy(io.Discard, resp.Body)
	_ = resp.Close
}

func CloseReq(resp *http.Request) {
	if resp == nil || resp.Body == nil {
		return
	}
	_ = resp.Body.Close()
}

func Closeresponse(resp *http.Response) {
	if resp == nil || resp.Body == nil {
		return
	}
	_ = resp.Body.Close()
}
