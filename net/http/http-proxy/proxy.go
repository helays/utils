package http_proxy

import (
	"bufio"
	"encoding/base64"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"strings"
	"sync"

	"github.com/helays/utils/v2/close/vclose"
	"golang.org/x/net/proxy"
)

var once sync.Once

// 初始化时注册 HTTP 拨号器
func init() {
	// 注册 CONNECT 方式的 HTTP 代理
	once.Do(func() {
		proxy.RegisterDialerType("http", func(uri *url.URL, forward proxy.Dialer) (proxy.Dialer, error) {
			return &httpProxyDialer{proxyURL: uri}, nil
		})
		proxy.RegisterDialerType("https", func(uri *url.URL, forward proxy.Dialer) (proxy.Dialer, error) {
			return &httpProxyDialer{proxyURL: uri}, nil
		})
	})
}

type httpProxyDialer struct {
	proxyURL *url.URL
}

func (d *httpProxyDialer) Dial(network, addr string) (net.Conn, error) {
	// 建立到 HTTP 代理服务器的连接
	conn, err := net.Dial("tcp", d.proxyURL.Host)
	if err != nil {
		return nil, err
	}
	// 发送 CONNECT 请求
	req := &http.Request{Method: "CONNECT", URL: &url.URL{Opaque: addr}, Host: addr, Header: make(http.Header)}
	if d.proxyURL.User != nil {
		if password, ok := d.proxyURL.User.Password(); ok {
			auth := d.proxyURL.User.Username() + ":" + password
			req.Header.Set("Proxy-Authorization", "Basic "+basicAuth(auth))
		}
	}
	if err = req.Write(conn); err != nil {
		vclose.Close(conn)
		return nil, err
	}
	reader := bufio.NewReader(conn)
	if err = readHeader(reader); err != nil {
		vclose.Close(conn)
		return nil, err
	}

	return conn, nil
}

// 基础认证辅助函数
func basicAuth(auth string) string {
	return base64.StdEncoding.EncodeToString([]byte(auth))
}

func readHeader(r *bufio.Reader) error {
	line, err := r.ReadString('\n')
	if err != nil {
		return err
	}
	if !strings.Contains(line, "200") {
		return fmt.Errorf("proxy connection failed: %s", line)
	}
	for {
		line, err = r.ReadString('\n')
		if err != nil {
			return err
		}
		if line == "\r\n" || line == "\n" {
			break
		}
	}
	return nil
}
