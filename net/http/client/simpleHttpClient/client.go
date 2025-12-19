package simpleHttpClient

import (
	"context"
	"crypto/tls"
	"net"
	"net/http"
	"time"

	cfg_proxy "github.com/helays/utils/v2/config/cfg-proxy"
	"github.com/helays/utils/v2/safe"
)

var cache = safe.NewMap[string, *http.Client](safe.StringHasher{})

func NewWithProxy(timeout time.Duration, proxy *cfg_proxy.Proxy) (*http.Client, error) {
	return initClientHelper(timeout, proxy)
}

func New(timeout time.Duration, args ...string) (*http.Client, error) {
	var proxy *cfg_proxy.Proxy
	if len(args) >= 1 {
		proxy = &cfg_proxy.Proxy{Addr: args[0]}
	}
	return initClientHelper(timeout, proxy)
}

func NewWithCache(ck string, timeout time.Duration, args ...string) (*http.Client, error) {
	c, ok := cache.Load(ck)
	if ok {
		return c, nil
	}
	c, err := New(timeout, args...)
	if err != nil {
		return nil, err
	}
	cache.Store(ck, c)
	return c, nil
}

func initClientHelper(timeout time.Duration, proxy *cfg_proxy.Proxy) (*http.Client, error) {
	trans := initTransportHelper()
	if proxy != nil {
		if err := proxy.Valid(); err != nil {
			return nil, err
		}
		switch proxy.ProxyType() {
		case cfg_proxy.ProxySocks5:
			trans.DialContext = func(ctx context.Context, network, addr string) (net.Conn, error) {
				return proxy.Dialer().Dial(network, addr)
			}
		case cfg_proxy.ProxyHttp, cfg_proxy.ProxyHttps:
			trans.Proxy = proxy.HttpProxy()
		}
	}
	return &http.Client{Transport: trans, Timeout: timeout}, nil
}

func initTransportHelper() *http.Transport {
	return &http.Transport{
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
}
