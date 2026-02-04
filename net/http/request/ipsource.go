package request

import (
	"net"
	"net/http"
)

// DefaultIPSourcePriority 默认IP来源解析顺序（按可信度从高到低）
var DefaultIPSourcePriority = []IPSource{
	// 第一梯队：专有可信代理头部（通常由可信CDN/代理设置）
	HeaderTrueClientIP,   // Akamai, AWS等可信CDN
	HeaderCFConnectingIP, // CloudFlare可信代理

	// 第二梯队：反向代理标准头部
	HeaderXRealIP, // Nginx等反向代理

	// 第三梯队：RFC标准头部
	HeaderForwarded, // RFC 7239标准

	// 第四梯队：事实标准但可能被篡改
	HeaderXForwardedFor, // 标准代理链

	// 第五梯队：非标准变体（兼容旧系统）
	HeaderHTTPXForwardedFor,
	HeaderHTTPClientIP,

	// 最后：直接连接地址（最低信任度，但不会被伪造）
	RemoteAddr,

	// 注意：以下头部不用于IP提取，仅用于其他信息
	// HeaderXForwardedHost,  // 用于主机名
	// HeaderXForwardedProto, // 用于协议
	// HeaderVia,             // 用于代理链跟踪
}

type IPSource string

func (s IPSource) String() string {
	return string(s)
}

// IPHeader 定义标准的HTTP头部字段名
// noinspection no
const (
	RemoteAddr IPSource = "RemoteAddr"

	//noinspection no 标准或广泛使用的头部
	HeaderForwarded      IPSource = "Forwarded"        // RFC 7239
	HeaderXForwardedFor  IPSource = "X-Forwarded-For"  // 事实标准
	HeaderXRealIP        IPSource = "X-Real-IP"        // Nginx等
	HeaderCFConnectingIP IPSource = "CF-Connecting-IP" // CloudFlare
	HeaderTrueClientIP   IPSource = "True-Client-IP"   // Akamai等

	//noinspection no 常见但非标准的变体（保持向后兼容）
	HeaderHTTPClientIP      IPSource = "HTTP_CLIENT_IP"
	HeaderHTTPXForwardedFor IPSource = "HTTP_X_FORWARDED_FOR"
	HeaderXForwardedHost    IPSource = "X-Forwarded-Host"
	HeaderXForwardedProto   IPSource = "X-Forwarded-Proto"

	//noinspection no 特殊代理头部
	HeaderVia IPSource = "Via" // 代理服务器链
)

type IPSources struct {
	RemoteAddr string `json:"remote_addr"` // 直接连接信息

	// 标准头部
	Forwarded      string `json:"forwarded,omitempty"`        // RFC 7239
	XForwardedFor  string `json:"x_forwarded_for,omitempty"`  // 事实标准
	XRealIP        string `json:"x_real_ip,omitempty"`        // Nginx等
	CFConnectingIP string `json:"cf_connecting_ip,omitempty"` // CloudFlare
	TrueClientIP   string `json:"true_client_ip,omitempty"`   // Akamai等

	// 非标准变体（向后兼容）
	HTTPClientIP      string `json:"http_client_ip,omitempty"`
	HTTPXForwardedFor string `json:"http_x_forwarded_for,omitempty"`

	// 其他相关头部
	XForwardedHost  string `json:"x_forwarded_host,omitempty"`
	XForwardedProto string `json:"x_forwarded_proto,omitempty"`
	Via             string `json:"via,omitempty"`
}

// ParseRequestIPSources 解析请求中的所有IP相关信息
func ParseRequestIPSources(r *http.Request) IPSources {
	remoteAddr, _, _ := net.SplitHostPort(r.RemoteAddr)
	if remoteAddr == "::1" {
		remoteAddr = "127.0.0.1"
	}
	sources := IPSources{
		RemoteAddr: remoteAddr,
	}

	// 收集所有标准头部
	sources.Forwarded = r.Header.Get(HeaderForwarded.String())
	sources.XForwardedFor = r.Header.Get(HeaderXForwardedFor.String())
	sources.XRealIP = r.Header.Get(HeaderXRealIP.String())
	sources.CFConnectingIP = r.Header.Get(HeaderCFConnectingIP.String())
	sources.TrueClientIP = r.Header.Get(HeaderTrueClientIP.String())

	// 收集非标准变体
	sources.HTTPClientIP = r.Header.Get(HeaderHTTPClientIP.String())
	sources.HTTPXForwardedFor = r.Header.Get(HeaderHTTPXForwardedFor.String())

	// 收集其他相关头部
	sources.XForwardedHost = r.Header.Get(HeaderXForwardedHost.String())
	sources.XForwardedProto = r.Header.Get(HeaderXForwardedProto.String())
	sources.Via = r.Header.Get(HeaderVia.String())

	return sources
}

// Getip 获取客户端IP
// noinspection SpellCheckingInspection
func Getip(r *http.Request, policys ...IPSource) string {
	var remoteAddr string
	ps := DefaultIPSourcePriority
	if len(policys) > 0 {
		ps = policys
	}
outerLoop: // 定义标签
	for _, ipSource := range ps {
		switch ipSource {
		case RemoteAddr:
			continue
		default:
			if ip := r.Header.Get(ipSource.String()); ip != "" {
				remoteAddr = ip
				break outerLoop
			}
		}
	}
	if remoteAddr == "" {
		remoteAddr, _, _ = net.SplitHostPort(r.RemoteAddr)
	}
	if remoteAddr == "::1" {
		remoteAddr = "127.0.0.1"
	}
	return remoteAddr
}

type IPParsePolicy []IPSource

func (p IPParsePolicy) GetIP(r *http.Request) string {
	return Getip(r, p...)
}

func (p IPParsePolicy) GetFromIPSources(sources *IPSources) string {
	if sources == nil {
		return ""
	}
	ps := DefaultIPSourcePriority
	if len(p) > 0 {
		ps = p
	}
	for _, ipSource := range ps {
		var ip string
		switch ipSource {
		case RemoteAddr:
			ip = sources.RemoteAddr
		case HeaderForwarded:
			ip = sources.Forwarded
		case HeaderXForwardedFor:
			ip = sources.XForwardedFor
		case HeaderXRealIP:
			ip = sources.XRealIP
		case HeaderCFConnectingIP:
			ip = sources.CFConnectingIP
		case HeaderTrueClientIP:
			ip = sources.TrueClientIP
		case HeaderHTTPClientIP:
			ip = sources.HTTPClientIP
		case HeaderHTTPXForwardedFor:
			ip = sources.HTTPXForwardedFor
		case HeaderXForwardedHost:
			ip = sources.XForwardedHost
		case HeaderXForwardedProto:
			ip = sources.XForwardedProto
		case HeaderVia:
			ip = sources.Via
		}
		if ip != "" {
			return ip
		}
	}
	return sources.RemoteAddr
}
