package elastic

import (
	"context"
	"crypto/tls"
	"github.com/elastic/go-elasticsearch/v8"
	"golang.org/x/net/proxy"
	"helay.net/go/utils/v3/config"
	"helay.net/go/utils/v3/net/http/client/simpleHttpClient"
	"helay.net/go/utils/v3/tools"
	"net"
	"net/http"
	"net/url"
	"time"
)

// Config elastic 配置
type Config struct {
	DbIdentifier string `ini:"db_identifier" yaml:"db_identifier,omitempty" json:"db_identifier"`
	Schema       string `json:"schema" yaml:"schema" ini:"schema"`

	Addresses                   []string                                            `json:"addresses" yaml:"addresses" ini:"addresses"`                                           // 一个包含 Elasticsearch 节点地址的列表，用于客户端连接到 Elasticsearch 集群。
	Username                    string                                              `json:"username" yaml:"username" ini:"username"`                                              // HTTP 基本认证的用户名。
	Password                    string                                              `json:"password" yaml:"password" ini:"password"`                                              // HTTP 基本认证的密码。
	CloudID                     string                                              `json:"cloud_id" yaml:"cloud_id" ini:"cloud_id"`                                              // 当使用 Elastic 提供的服务时（例如通过 https://elastic.co/cloud），这是服务的终端地址。
	APIKey                      string                                              `json:"api_key" yaml:"api_key" ini:"api_key"`                                                 // 用于鉴权的 Base64 编码令牌。如果设置了该值，则会覆盖 Username/Password 和服务令牌进行鉴权
	ServiceToken                string                                              `json:"service_token" yaml:"service_token" ini:"service_token"`                               // 另一种形式的鉴权令牌。如果设置，它会覆盖 Username 和 Password 进行鉴权。
	CertificateFingerprint      string                                              `json:"certificate_fingerprint" yaml:"certificate_fingerprint" ini:"certificate_fingerprint"` // 在首次启动 Elasticsearch 时提供的 SHA256 散列指纹，用于验证服务器证书的真实性。
	Header                      http.Header                                         `json:"header" yaml:"header" ini:"header"`                                                    // 全局 HTTP 请求头，可以用来添加自定义头部信息到所有请求中。
	CACert                      []byte                                              `json:"ca_cert" yaml:"ca_cert" ini:"ca_cert"`                                                 // PEM 编码的证书颁发机构(CA)证书列表。当设置时，会创建一个新的空证书池，并将这些证书添加进去。仅在未指定传输或使用的是 http.Transport 时有效。
	CAPath                      string                                              `json:"ca_path" yaml:"ca_path" ini:"ca_path"`
	RetryOnStatus               []int                                               `json:"retry_on_status" yaml:"retry_on_status" ini:"retry_on_status"`                                     // 状态码列表，指示哪些状态码发生时应该重试请求。默认包括 502, 503, 504。
	DisableRetry                bool                                                `json:"disable_retry" yaml:"disable_retry" ini:"disable_retry"`                                           // 是否禁用自动重试功能，默认为 false。
	MaxRetries                  int                                                 `json:"max_retries" yaml:"max_retries" ini:"max_retries"`                                                 //  自动重试的最大次数，默认为 3 次。
	CompressRequestBody         bool                                                `json:"compress_request_body" yaml:"compress_request_body" ini:"compress_request_body"`                   // 是否压缩请求体，默认为 false。
	CompressRequestBodyLevel    int                                                 `json:"compress_request_body_level" yaml:"compress_request_body_level" ini:"compress_request_body_level"` // 如果启用请求体压缩，此选项设定压缩级别，默认为 gzip.DefaultCompression。
	PoolCompressor              bool                                                `json:"pool_compressor" yaml:"pool_compressor" ini:"pool_compressor"`                                     // 如果设置为 true，则使用基于 sync.Pool 的 gzip writer 来管理压缩器对象，默认为 false。
	DiscoverNodesOnStart        bool                                                `json:"discover_nodes_on_start" yaml:"discover_nodes_on_start" ini:"discover_nodes_on_start"`             // 初始化客户端时是否发现节点，默认为 false。
	DiscoverNodesInterval       time.Duration                                       `json:"discover_nodes_interval" yaml:"discover_nodes_interval" ini:"discover_nodes_interval"`             // 定期发现节点的时间间隔，默认为禁用。单位秒
	EnableMetrics               bool                                                `json:"enable_metrics" yaml:"enable_metrics" ini:"enable_metrics"`                                        // 启用指标收集，默认为 false。
	EnableDebugLogger           bool                                                `json:"enable_debug_logger" yaml:"enable_debug_logger" ini:"enable_debug_logger"`                         // 启用调试日志记录，默认为 false。
	EnableCompatibilityMode     bool                                                `json:"enable_compatibility_mode" yaml:"enable_compatibility_mode" ini:"enable_compatibility_mode"`       //  启用发送兼容性头部，默认为 false。
	DisableMetaHeader           bool                                                `json:"disable_meta_header" yaml:"disable_meta_header" ini:"disable_meta_header"`                         // 禁用额外的 "X-Elastic-Client-Meta" HTTP 头部，默认为 false。
	*simpleHttpClient.Transport `json:"transport" yaml:"transport" ini:"transport"` // 传输配置，
}

func (this *Config) SetInfo(args ...any) {
	if len(args) != 2 {
		return
	}
	switch args[0].(string) {
	case config.ClientInfoHost:
		this.Addresses = args[1].([]string)
	case config.ClientInfoUser:
		this.Username = args[1].(string)
	case config.ClientInfoPasswd:
		this.Password = args[1].(string)
	}
}

func (this Config) NewClient() (*elasticsearch.Client, error) {
	var err error
	if this.CAPath != "" {
		this.CAPath = tools.Fileabs(this.CAPath)
		this.CACert, err = tools.FileGetContents(this.CAPath)
		if err != nil {
			return nil, err
		}
	}
	cfg := elasticsearch.Config{
		Addresses:                this.Addresses,
		Username:                 this.Username,
		Password:                 this.Password,
		CloudID:                  this.CloudID,
		APIKey:                   this.APIKey,
		ServiceToken:             this.ServiceToken,
		CertificateFingerprint:   this.CertificateFingerprint,
		Header:                   this.Header,
		CACert:                   this.CACert,
		RetryOnStatus:            this.RetryOnStatus,
		DisableRetry:             this.DisableRetry,
		MaxRetries:               this.MaxRetries,
		CompressRequestBody:      this.CompressRequestBody,
		CompressRequestBodyLevel: this.CompressRequestBodyLevel,
		PoolCompressor:           this.PoolCompressor,
		DiscoverNodesOnStart:     this.DiscoverNodesOnStart,
		DiscoverNodesInterval:    tools.AutoTimeDuration(this.DiscoverNodesInterval, time.Second),
		EnableMetrics:            this.EnableMetrics,
		EnableDebugLogger:        this.EnableDebugLogger,
		EnableCompatibilityMode:  this.EnableCompatibilityMode,
		DisableMetaHeader:        this.DisableMetaHeader,
	}
	if this.Transport != nil {
		// 配置自定义传输结构
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
		// 配置代理
		switch this.Transport.ProxyType {
		case simpleHttpClient.ProxyHttp, simpleHttpClient.ProxyHttps:
			u, _err := url.Parse(this.Transport.ProxyAddr)
			if _err != nil {
				return nil, _err
			}
			trans.Proxy = http.ProxyURL(u)
		case simpleHttpClient.ProxySocks5:
			dialer, _err := proxy.SOCKS5("tcp", this.Transport.ProxyAddr, nil, proxy.Direct)
			if _err != nil {
				return nil, _err
			}
			trans.DialContext = func(ctx context.Context, network, addr string) (net.Conn, error) {
				return dialer.Dial(network, addr)
			}
		}
		cfg.Transport = trans
	}
	return elasticsearch.NewClient(cfg)
}
