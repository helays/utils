package simpleHttpClient

import (
	"context"
	"crypto/tls"
	"fmt"
	"golang.org/x/net/proxy"
	"net"
	"net/http"
	"net/url"
	"time"
)

var simpleClient *http.Client

// InitHttpClient 初始化http client
func InitHttpClient(timeout time.Duration, args ...string) (*http.Client, error) {
	if simpleClient != nil {
		return simpleClient, nil
	}

	trans := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true, // client 不对https 证书进行校验
		},
		TLSHandshakeTimeout: 5 * time.Second, // tls 握手超时,
		DisableKeepAlives:   false,           // 是否禁用持久连接。当设为false时，支持HTTP keep-alive，允许在一个TCP连接上发送多个HTTP请求。
		DisableCompression:  false,           // 是否禁用对响应内容的自动解压缩。设为false表示接受编码过的响应，并由客户端自动解码。
		// 主机连接池管理
		MaxIdleConns:        1000, // 控制整个客户端的最大空闲连接数。值为0表示没有限制。
		MaxIdleConnsPerHost: 5,    // 限制每个主机的最大空闲连接数。同样地，0表示没有限制。
		MaxConnsPerHost:     0,    // 每个主机的最大连接数（包括活跃和空闲）。0表示无限制。
		IdleConnTimeout:     0,    // 设置空闲连接在被关闭前等待新请求的时间长度。0表示不主动关闭空闲连接。
		//超时配置
		ResponseHeaderTimeout: 0, // 等待响应头的超时时间。0表示没有超时。
		ExpectContinueTimeout: 0, // 在发送请求体之前等待服务器预检响应(100 Continue)的时间。0表示不等待。

		MaxResponseHeaderBytes: 0,     // 从服务器读取响应头的最大字节数。0表示没有限制。
		WriteBufferSize:        0,     // 写缓冲区大小。0表示使用系统默认值。
		ReadBufferSize:         0,     // 读缓冲区大小。0表示使用系统默认值。
		ForceAttemptHTTP2:      false, // 强制尝试使用HTTP/2。如果设置为true，即使服务端不明确支持HTTP/2，也会尝试升级。
	}

	if len(args) >= 2 {
		proxyAddr := args[1]
		switch args[0] {
		case "socks5":
			dialer, err := proxy.SOCKS5("tcp", proxyAddr, nil, proxy.Direct)
			if err != nil {
				return simpleClient, fmt.Errorf("newHttpClient socks5 proxy error: %v", err)
			}
			trans.DialContext = func(ctx context.Context, network, addr string) (net.Conn, error) {
				return dialer.Dial(network, addr)
			}
		case "http":
			u, err := url.Parse(proxyAddr)
			if err != nil {
				return simpleClient, fmt.Errorf("newHttpClient parse proxy url error: %v", err)
			}
			trans.Proxy = http.ProxyURL(u)
		}
	}
	simpleClient = &http.Client{

		Transport: trans,
		Timeout:   timeout,
	}
	return simpleClient, nil
}
