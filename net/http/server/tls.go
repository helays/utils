package server

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"os"
	"strings"

	"github.com/helays/utils/v2/tools"
)

// CipherSuiteMapping 密码套件映射
var CipherSuiteMapping = map[string]uint16{
	// TLS 1.3 密码套件
	"TLS_AES_128_GCM_SHA256":       tls.TLS_AES_128_GCM_SHA256,
	"TLS_AES_256_GCM_SHA384":       tls.TLS_AES_256_GCM_SHA384,
	"TLS_CHACHA20_POLY1305_SHA256": tls.TLS_CHACHA20_POLY1305_SHA256,

	// TLS 1.2 安全密码套件
	"TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256": tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
	"TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256":   tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
	"TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384": tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
	"TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384":   tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
	"TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305":  tls.TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305_SHA256,
	"TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305":    tls.TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305_SHA256,

	// TLS 1.2 兼容性密码套件
	"TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA256": tls.TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA256,
	"TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA":    tls.TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA,
	"TLS_ECDHE_ECDSA_WITH_AES_128_CBC_SHA":  tls.TLS_ECDHE_ECDSA_WITH_AES_128_CBC_SHA,
	"TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA":    tls.TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA,
	"TLS_ECDHE_ECDSA_WITH_AES_256_CBC_SHA":  tls.TLS_ECDHE_ECDSA_WITH_AES_256_CBC_SHA,

	// 不安全的密码套件 (通常不建议使用)
	"TLS_RSA_WITH_AES_128_GCM_SHA256": tls.TLS_RSA_WITH_AES_128_GCM_SHA256,
	"TLS_RSA_WITH_AES_256_GCM_SHA384": tls.TLS_RSA_WITH_AES_256_GCM_SHA384,
	"TLS_RSA_WITH_AES_128_CBC_SHA":    tls.TLS_RSA_WITH_AES_128_CBC_SHA,
	"TLS_RSA_WITH_AES_256_CBC_SHA":    tls.TLS_RSA_WITH_AES_256_CBC_SHA,
}

// 预定义的密码套件组合
var (
	ModernCipherSuites = []string{
		// TLS 1.3 密码套件
		"TLS_AES_128_GCM_SHA256",
		"TLS_AES_256_GCM_SHA384",
		"TLS_CHACHA20_POLY1305_SHA256",

		// 安全的 TLS 1.2 密码套件
		"TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256",
		"TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256",
		"TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384",
		"TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384",
		"TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305",
		"TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305",
	} // 现代密码套件 (TLS 1.3 + 安全的 TLS 1.2)

	CompatibleCipherSuites = []string{
		// TLS 1.3
		"TLS_AES_128_GCM_SHA256",
		"TLS_AES_256_GCM_SHA384",
		"TLS_CHACHA20_POLY1305_SHA256",

		// TLS 1.2
		"TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256",
		"TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256",
		"TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384",
		"TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384",
		"TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305",
		"TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305",

		// 更多兼容性套件
		"TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA256",
		"TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA",
		"TLS_ECDHE_ECDSA_WITH_AES_128_CBC_SHA",
	} // 兼容性密码套件 (包含更多旧版本支持)
)

// ParseCipherSuites 将字符串密码套件名称转换为对应的 ID
func ParseCipherSuites(cipherNames []string) ([]uint16, error) {
	if len(cipherNames) == 0 {
		return nil, nil
	}

	var cipherSuites []uint16
	var unknownCiphers []string

	for _, name := range cipherNames {
		if id, exists := CipherSuiteMapping[name]; exists {
			cipherSuites = append(cipherSuites, id)
		} else {
			unknownCiphers = append(unknownCiphers, name)
		}
	}

	if len(unknownCiphers) > 0 {
		return nil, fmt.Errorf("未知的密码套件: %s", strings.Join(unknownCiphers, ", "))
	}

	return cipherSuites, nil
}

// GetCipherSuiteName 根据 ID 获取密码套件名称
func GetCipherSuiteName(id uint16) string {
	for name, cipherID := range CipherSuiteMapping {
		if cipherID == id {
			return name
		}
	}
	return fmt.Sprintf("未知密码套件 (0x%04X)", id)
}

// GetAvailableCipherSuites 获取所有可用的密码套件名称
func GetAvailableCipherSuites() []string {
	var names []string
	for name := range CipherSuiteMapping {
		names = append(names, name)
	}
	return names
}

// TLSVersionMapping TLS版本映射
var TLSVersionMapping = map[string]uint16{
	"TLS10": tls.VersionTLS10,
	"TLS11": tls.VersionTLS11,
	"TLS12": tls.VersionTLS12,
	"TLS13": tls.VersionTLS13,
}

// ParseTLSVersion 解析TLS版本字符串
func ParseTLSVersion(version string) (uint16, error) {
	if v, exists := TLSVersionMapping[version]; exists {
		return v, nil
	}
	return 0, fmt.Errorf("未知的TLS版本: %s", version)
}

// 支持的协议常量
const (
	ProtocolHTTP11 = "http/1.1"
	ProtocolHTTP2  = "h2"
	ProtocolHTTP3  = "h3"
)

// 预定义的协议组合
var (
	ModernProtocols     = []string{ProtocolHTTP11, ProtocolHTTP2}                // 现代协议 (HTTP/1.1 + HTTP/2)
	CompatibleProtocols = []string{ProtocolHTTP11}                               // 兼容性协议 (只支持 HTTP/1.1)
	AllProtocols        = []string{ProtocolHTTP11, ProtocolHTTP2, ProtocolHTTP3} // 全协议支持
)

type TLSConfig struct {
	Enable       bool          `json:"enable" yaml:"enable"`
	Certificates []Certificate `json:"certificates" yaml:"certificates"` // 证书
	RootCAFile   string        `json:"root_ca_file" yaml:"root_ca_file"` // 根CA文件
	NextProtos   []string      `json:"next_protos" yaml:"next_protos"`   // 支持的协议
	ServerName   string        `json:"server_name" yaml:"server_name"`   // 服务器名称

	// 0 不验证客户端证书
	// 1 请求客户端证书但不强制验证
	// 2 要求客户端证书但不验证CA
	// 3 如果提供客户端证书则验证
	// 4 要求并验证客户端证书
	ClientAuth                  tls.ClientAuthType       `json:"client_auth" yaml:"client_auth"`                                       // 客户端验证
	ClientCAFile                string                   `json:"client_ca_file" yaml:"client_ca_file"`                                 // 客户端CA文件
	InsecureSkipVerify          bool                     `json:"insecure_skip_verify" yaml:"insecure_skip_verify"`                     // 跳过验证
	CipherSuites                []string                 `json:"cipher_suites" yaml:"cipher_suites"`                                   // 密钥套件
	SessionTicketsDisabled      bool                     `json:"session_tickets_disabled" yaml:"session_tickets_disabled"`             // 禁用会话密钥
	MinVersion                  string                   `json:"min_version" yaml:"min_version"`                                       // 最低TLS版本
	MaxVersion                  string                   `json:"max_version" yaml:"max_version"`                                       // 最高TLS版本
	DynamicRecordSizingDisabled bool                     `json:"dynamic_record_sizing_disabled" yaml:"dynamic_record_sizing_disabled"` // 动态记录大小禁用
	Renegotiation               tls.RenegotiationSupport `json:"renegotiation" yaml:"renegotiation"`                                   // 重新协商
}

type Certificate struct {
	CertFile string `json:"cert_file" yaml:"cert_file"`
	KeyFile  string `json:"key_file" yaml:"key_file"`
}

// ToTLSConfig 将配置转换为标准库的 tls.Config
func (t *TLSConfig) ToTLSConfig() (*tls.Config, error) {
	if t == nil || !t.Enable {
		return nil, nil
	}

	config := &tls.Config{
		ServerName:                  t.ServerName,
		NextProtos:                  t.NextProtos,
		SessionTicketsDisabled:      t.SessionTicketsDisabled,
		InsecureSkipVerify:          t.InsecureSkipVerify,
		DynamicRecordSizingDisabled: t.DynamicRecordSizingDisabled,
		Renegotiation:               t.Renegotiation,
	}

	// 解析TLS版本
	if t.MinVersion != "" {
		minVersion, err := ParseTLSVersion(t.MinVersion)
		if err != nil {
			return nil, fmt.Errorf("解析最低TLS版本失败: %v", err)
		}
		config.MinVersion = minVersion
	} else {
		// 默认使用 TLS 1.2
		config.MinVersion = tls.VersionTLS12
	}
	if t.MaxVersion != "" {
		maxVersion, err := ParseTLSVersion(t.MaxVersion)
		if err != nil {
			return nil, fmt.Errorf("解析最高TLS版本失败: %v", err)
		}
		config.MaxVersion = maxVersion
	}

	// 解析密码套件
	var cipherSuites []uint16
	var err error
	if len(t.CipherSuites) > 0 {
		cipherSuites, err = ParseCipherSuites(t.CipherSuites)
		if err != nil {
			return nil, fmt.Errorf("解析密码套件失败: %v", err)
		}
		config.CipherSuites = cipherSuites
	} else {
		// 默认使用现代密码套件
		cipherSuites, err = ParseCipherSuites(ModernCipherSuites)
		if err != nil {
			return nil, fmt.Errorf("解析默认密码套件失败: %v", err)
		}
		config.CipherSuites = cipherSuites
	}

	// 加载服务器证书
	if len(t.Certificates) > 0 {
		certs := make([]tls.Certificate, 0, len(t.Certificates))
		for i, cert := range t.Certificates {
			if cert.CertFile == "" || cert.KeyFile == "" {
				return nil, fmt.Errorf("证书 %d: 必须同时提供 cert_file 和 key_file", i)
			}

			certFile := tools.Fileabs(cert.CertFile)
			keyFile := tools.Fileabs(cert.KeyFile)

			certificate, err := tls.LoadX509KeyPair(certFile, keyFile)
			if err != nil {
				return nil, fmt.Errorf("加载证书 %d 失败: %v", i, err)
			}
			certs = append(certs, certificate)
		}
		config.Certificates = certs
	}

	// 加载根 CA 证书
	if t.RootCAFile != "" {
		rootCAFile := tools.Fileabs(t.RootCAFile)
		rootCACert, err := os.ReadFile(rootCAFile)
		if err != nil {
			return nil, fmt.Errorf("读取根证书文件失败: %v", err)
		}
		rootCertPool := x509.NewCertPool()
		if !rootCertPool.AppendCertsFromPEM(rootCACert) {
			return nil, fmt.Errorf("解析根证书失败")
		}
		config.RootCAs = rootCertPool
	}
	// 加载客户端 CA 证书
	if t.ClientCAFile != "" {
		clientCACert, err := os.ReadFile(tools.Fileabs(t.ClientCAFile))
		if err != nil {
			return nil, fmt.Errorf("读取客户端CA文件失败: %v", err)
		}
		clientCertPool := x509.NewCertPool()
		if !clientCertPool.AppendCertsFromPEM(clientCACert) {
			return nil, fmt.Errorf("解析客户端CA证书失败")
		}
		config.ClientCAs = clientCertPool
		config.ClientAuth = t.ClientAuth
	} else if t.ClientAuth != tls.NoClientCert {
		// 如果设置了客户端认证但没有提供 CA 文件，返回错误
		return nil, fmt.Errorf("客户端认证需要提供 client_ca_file")
	}
	return config, nil
}

// DefaultTLSConfig 返回一个安全的默认 TLS 配置
func DefaultTLSConfig() *TLSConfig {
	return &TLSConfig{
		Enable:        true,
		NextProtos:    ModernProtocols,
		MinVersion:    "TLS12",
		CipherSuites:  ModernCipherSuites,
		ClientAuth:    tls.NoClientCert,
		Renegotiation: tls.RenegotiateNever,
	}
}

// CompatibleTLSConfig 返回一个兼容性更好的 TLS 配置
func CompatibleTLSConfig() *TLSConfig {
	return &TLSConfig{
		Enable:        true,
		NextProtos:    CompatibleProtocols,
		MinVersion:    "TLS12",
		CipherSuites:  CompatibleCipherSuites,
		ClientAuth:    tls.NoClientCert,
		Renegotiation: tls.RenegotiateNever,
	}
}
