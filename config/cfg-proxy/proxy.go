package cfg_proxy

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"

	_ "github.com/helays/utils/v2/net/http/http-proxy"
	"github.com/pkg/errors"
	"golang.org/x/net/proxy"
)

type ProxyTypes string

const (
	ProxyNone    ProxyTypes = "none"    // 无代理
	ProxyHttp    ProxyTypes = "http"    // http 代理
	ProxyHttps   ProxyTypes = "https"   // https 代理
	ProxySocks5  ProxyTypes = "socks5"  // socks 代理
	ProxyUnknown ProxyTypes = "unknown" // 无效的代理类型
)

type Proxy struct {
	Addr      string `json:"addr" yaml:"addr" ini:"addr"`
	Account   string `json:"account" yaml:"account" ini:"account"`
	Password  string `json:"password" yaml:"password" ini:"password"`
	proxyType ProxyTypes
	dialer    proxy.Dialer
	httpProxy func(r *http.Request) (*url.URL, error)
	u         *url.URL
}

func (p *Proxy) Valid() error {
	p.Addr = strings.TrimSpace(p.Addr)
	if p.Addr == "" {
		p.proxyType = ProxyNone
		return nil
	}
	u, err := url.Parse(p.Addr)
	if err != nil {
		p.proxyType = ProxyUnknown
		return fmt.Errorf("无效的代理地址 %v", errors.WithStack(err))
	}
	p.u = u

	switch u.Scheme {
	case "http":
		p.proxyType = ProxyHttp
		p.httpProxy = http.ProxyURL(u)
		return nil
	case "https":
		p.httpProxy = http.ProxyURL(u)
		p.proxyType = ProxyHttps
		return nil
	case "socks5":
		p.proxyType = ProxySocks5
		var auth *proxy.Auth
		if p.Account == "" && p.Password == "" && u.User != nil {
			p.Account = u.User.Username()
			p.Password, _ = u.User.Password()
		}
		if p.Account != "" && p.Password != "" {
			auth = &proxy.Auth{User: p.Account, Password: p.Password}
		}
		p.dialer, err = proxy.SOCKS5("tcp", u.Host, auth, proxy.Direct)
		if err != nil {
			p.proxyType = ProxyUnknown
			return fmt.Errorf("socks5拨号器初始化失败 %v", errors.WithStack(err))
		}
		return nil
	default:
		p.proxyType = ProxyUnknown
		return nil
	}
}

func (p *Proxy) ProxyType() ProxyTypes {
	return p.proxyType
}

func (p *Proxy) HttpProxy() func(r *http.Request) (*url.URL, error) {
	return p.httpProxy
}

func (p *Proxy) Dialer() proxy.Dialer {
	return p.dialer
}

func (p *Proxy) GetProxyUrl() *url.URL {
	return p.u
}

func (p *Proxy) AutoDialer() (proxy.Dialer, error) {
	if p.proxyType == ProxySocks5 {
		return p.dialer, nil
	}
	if p.proxyType == ProxyUnknown || p.proxyType == ProxyNone {
		return nil, nil
	}
	return proxy.FromURL(p.u, proxy.Direct)
}
