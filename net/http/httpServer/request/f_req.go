package request

import (
	"fmt"
	"github.com/helays/utils/v2/tools"
	"net"
	"net/http"
)

// Getip 获取客户端IP
func Getip(r *http.Request) string {
	remoteAddr := r.RemoteAddr
	if ip := r.Header.Get("HTTP_CLIENT_IP"); ip != "" {
		remoteAddr = ip
	} else if ip = r.Header.Get("HTTP_X_FORWARDED_FOR"); ip != "" {
		remoteAddr = ip
	} else if ip = r.Header.Get("X-Real-IP"); ip != "" {
		remoteAddr = ip
	} else {
		remoteAddr, _, _ = net.SplitHostPort(remoteAddr)
	}
	if remoteAddr == "::1" {
		remoteAddr = "127.0.0.1"
	}
	return remoteAddr
}

// GetHeaderValueFormatInt 获取多个header字段，并转为int
func GetHeaderValueFormatInt(header http.Header, fields ...string) (map[string]int, error) {
	resp := make(map[string]int)
	for _, field := range fields {
		_v := header.Get(field)
		v, err := tools.Any2int(_v)
		if err != nil {
			return nil, fmt.Errorf("不合法的字段[%s]值 %s %v", field, _v, err)
		}
		resp[field] = int(v)
	}
	return resp, nil
}
