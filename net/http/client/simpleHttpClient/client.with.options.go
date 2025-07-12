package simpleHttpClient

import (
	"context"
	"crypto/tls"
	"github.com/helays/utils/v2/tools"
	"golang.org/x/net/proxy"
	"net"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"time"
)

const (
	ProxyHttp   = "http"
	ProxyHttps  = "https"
	ProxySocks5 = "socks5"
)

type HttpClient struct {
	NewClient    bool                                                `json:"new_client" yaml:"new_client" ini:"new_client"`          // 是否启用新的client
	EnableCookie bool                                                `json:"enable_cookie" yaml:"enable_cookie" ini:"enable_cookie"` // 是否启用cookie
	Timeout      time.Duration                                       `json:"timeout" yaml:"timeout" ini:"timeout"`                   // 超时时间
	Transport    `json:"transport" yaml:"transport" ini:"transport"` // 传输配置
}

type Transport struct {
	ProxyType              string        `json:"proxy_type" yaml:"proxy_type" ini:"proxy_type"`                                        // 代理类型 http,https,socks5
	ProxyAddr              string        `json:"proxy_addr" yaml:"proxy_addr" ini:"proxy_addr"`                                        // 代理地址
	InsecureSkipVerify     bool          `json:"insecure_skip_verify" yaml:"insecure_skip_verify" ini:"insecure_skip_verify"`          // 是否跳过证书验证
	TLSHandshakeTimeout    time.Duration `json:"tls_handshake_timeout" yaml:"tls_handshake_timeout" ini:"tls_handshake_timeout"`       // tls握手超时时间
	DisableKeepAlives      bool          `json:"disable_keep_alives" yaml:"disable_keep_alives" ini:"disable_keep_alives"`             // 是否禁用keepalives
	DisableCompression     bool          `json:"disable_compression" yaml:"disable_compression" ini:"disable_compression"`             // 是否禁用压缩
	MaxIdleConns           int           `json:"max_idle_conns" yaml:"max_idle_conns" ini:"max_idle_conns"`                            // 控制整个客户端的最大空闲连接数。值为0表示没有限制。
	MaxIdleConnsPerHost    int           `json:"max_idle_conns_per_host" yaml:"max_idle_conns_per_host" ini:"max_idle_conns_per_host"` // 限制每个主机的最大空闲连接数。同样地，0表示没有限制。
	MaxConnsPerHost        int           `json:"max_conns_per_host" yaml:"max_conns_per_host" ini:"max_conns_per_host"`                // 每个主机的最大连接数（包括活跃和空闲）。0表示无限制。
	IdleConnTimeout        time.Duration `json:"idle_conn_timeout" yaml:"idle_conn_timeout" ini:"idle_conn_timeout"`                   // 设置空闲连接在被关闭前等待新请求的时间长度。0表示不主动关闭空闲连接。
	ResponseHeaderTimeout  time.Duration `json:"response_header_timeout" yaml:"response_header_timeout" ini:"response_header_timeout"`
	ExpectContinueTimeout  time.Duration `json:"expect_continue_timeout" yaml:"expect_continue_timeout" ini:"expect_continue_timeout"`
	MaxResponseHeaderBytes int64         `json:"max_response_header_bytes" yaml:"max_response_header_bytes" ini:"max_response_header_bytes"`
	WriteBufferSize        int           `json:"write_buffer_size" yaml:"write_buffer_size" ini:"write_buffer_size"`
	ReadBufferSize         int           `json:"read_buffer_size" yaml:"read_buffer_size" ini:"read_buffer_size"`
	ForceAttemptHTTP2      bool          `json:"force_attempt_http2" yaml:"force_attempt_http2" ini:"force_attempt_http2"`
}

var client *http.Client

func (this HttpClient) InitClient() (*http.Client, error) {
	if !this.NewClient && client != nil {
		return client, nil
	}
	trans := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: this.Transport.InsecureSkipVerify, // client 不对https 证书进行校验
		},
		TLSHandshakeTimeout: tools.AutoTimeDuration(this.Transport.TLSHandshakeTimeout, time.Second), // tls 握手超时,
		DisableKeepAlives:   this.Transport.DisableKeepAlives,                                        // 是否禁用持久连接。当设为false时，支持HTTP keep-alive，允许在一个TCP连接上发送多个HTTP请求。
		DisableCompression:  this.Transport.DisableCompression,                                       // 是否禁用对响应内容的自动解压缩。设为false表示接受编码过的响应，并由客户端自动解码。
		// 主机连接池管理
		MaxIdleConns:        this.Transport.MaxIdleConns,                                         // 控制整个客户端的最大空闲连接数。值为0表示没有限制。
		MaxIdleConnsPerHost: this.Transport.MaxIdleConnsPerHost,                                  // 限制每个主机的最大空闲连接数。同样地，0表示没有限制。
		MaxConnsPerHost:     this.Transport.MaxConnsPerHost,                                      // 每个主机的最大连接数（包括活跃和空闲）。0表示无限制。
		IdleConnTimeout:     tools.AutoTimeDuration(this.Transport.IdleConnTimeout, time.Second), // 设置空闲连接在被关闭前等待新请求的时间长度。0表示不主动关闭空闲连接。
		//超时配置
		ResponseHeaderTimeout: tools.AutoTimeDuration(this.Transport.ResponseHeaderTimeout, time.Second), // 等待响应头的超时时间。0表示没有超时。
		ExpectContinueTimeout: tools.AutoTimeDuration(this.Transport.ExpectContinueTimeout, time.Second), // 在发送请求体之前等待服务器预检响应(100 Continue)的时间。0表示不等待。

		MaxResponseHeaderBytes: this.Transport.MaxResponseHeaderBytes, // 从服务器读取响应头的最大字节数。0表示没有限制。
		WriteBufferSize:        this.Transport.WriteBufferSize,        // 写缓冲区大小。0表示使用系统默认值。
		ReadBufferSize:         this.Transport.ReadBufferSize,         // 读缓冲区大小。0表示使用系统默认值。
		ForceAttemptHTTP2:      this.Transport.ForceAttemptHTTP2,      // 强制尝试使用HTTP/2。如果设置为true，即使服务端不明确支持HTTP/2，也会尝试升级。
	}
	switch this.Transport.ProxyType {
	case ProxyHttp, ProxyHttps:
		u, err := url.Parse(this.Transport.ProxyAddr)
		if err != nil {
			return nil, err
		}
		trans.Proxy = http.ProxyURL(u)
	case ProxySocks5:
		dialer, err := proxy.SOCKS5("tcp", this.Transport.ProxyAddr, nil, proxy.Direct)
		if err != nil {
			return nil, err
		}
		trans.DialContext = func(ctx context.Context, network, addr string) (net.Conn, error) {
			return dialer.Dial(network, addr)
		}
	}
	_client := &http.Client{
		Transport: trans,
		Timeout:   tools.AutoTimeDuration(this.Timeout, time.Second),
	}
	if this.EnableCookie {
		_client.Jar, _ = cookiejar.New(nil)
	}
	if this.NewClient {
		return _client, nil
	}
	client = _client
	return client, nil
}
