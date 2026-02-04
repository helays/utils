package simpleHttpClient

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"time"

	"golang.org/x/net/proxy"
)

var simpleClient *http.Client

// InitHttpClient 初始化http client
// Deprecated: 请使用 NewWithCache 或者 New
func InitHttpClient(timeout time.Duration, args ...string) (*http.Client, error) {
	if simpleClient != nil {
		return simpleClient, nil
	}
	var err error
	simpleClient, err = newClient(timeout, args...)
	if err != nil {
		return nil, err
	}
	return simpleClient, nil
}

func newClient(timeout time.Duration, args ...string) (*http.Client, error) {
	trans := initTransportHelper()
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
	return &http.Client{
		Transport: trans,
		Timeout:   timeout,
	}, nil
}
