package httpClose

import (
	"net/http"

	"github.com/quic-go/quic-go/http3"
)

func CloseResp(resp *http.Response) {
	if resp == nil || resp.Body == nil {
		return
	}
	_ = resp.Body.Close()
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

// Server 关闭http server
func Server(s *http.Server) {
	if s != nil {
		_ = s.Close()
	}
}

func ServerQuick(s *http3.Server) {
	if s != nil {
		_ = s.Close()
	}
}
