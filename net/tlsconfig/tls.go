package tlsconfig

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"os"
	"strings"

	"github.com/helays/utils/v2/tools"
)

var CurvePreferencesMap = map[string]tls.CurveID{
	"CurveP256":      tls.CurveP256,      // 标准椭圆曲线 （传统加密）
	"CurveP384":      tls.CurveP384,      // 标准椭圆曲线 （高安全性）
	"CurveP521":      tls.CurveP521,      // 标准椭圆曲线 （高安全性）
	"X25519":         tls.X25519,         // 密码学曲线 最快，比P256快3-4倍
	"X25519MLKEM768": tls.X25519MLKEM768, // 后量子混合曲线
}

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
		// TLS 1.2 密码套件（按性能排序，从最快到最慢）

		"TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305",  // 1. ECDSA证书 + ChaCha20-Poly1305（移动设备/无AES-NI时最快）
		"TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305",    // 2. RSA证书 + ChaCha20-Poly1305（兼容移动设备）
		"TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256", // 3. ECDSA证书 + AES-128-GCM（服务器有AES-NI时最快）
		"TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256",   // 4. RSA证书 + AES-128-GCM（广泛兼容且快）
		"TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384", // 5. ECDSA证书 + AES-256-GCM（更高安全性）
		"TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384",   // 6. RSA证书 + AES-256-GCM（高安全性兼容）

		// TLS 1.3 密码套件（Go会自动处理，放最后）
		"TLS_AES_128_GCM_SHA256",       // TLS 1.3最快
		"TLS_CHACHA20_POLY1305_SHA256", // TLS 1.3移动友好
		"TLS_AES_256_GCM_SHA384",       // TLS 1.3高安全
	} // 现代密码套件 (TLS 1.3 + 安全的 TLS 1.2)

	CompatibleCipherSuites = []string{
		// TLS 1.2 密码套件（按性能排序）
		// 第1梯队：现代快速套件
		"TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305",  // 最快（ECDSA+ChaCha20）
		"TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305",    // 快（RSA+ChaCha20）
		"TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256", // 快（ECDSA+AES-128）
		"TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256",   // 广泛兼容且快

		// 第2梯队：高安全性套件
		"TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384", // 高安全 ECDSA
		"TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384",   // 高安全 RSA

		// 第3梯队：兼容性套件（性能较差，必要时使用）
		"TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA256", // CBC模式，有AES-NI时还行
		"TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA",    // CBC+ECDSA
		"TLS_ECDHE_ECDSA_WITH_AES_128_CBC_SHA",  // CBC+RSA

		// TLS 1.3 密码套件（放最后，Go会忽略）
		"TLS_AES_128_GCM_SHA256",
		"TLS_AES_256_GCM_SHA384",
		"TLS_CHACHA20_POLY1305_SHA256",
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

// noinspection all
type TLSConfig struct {
	Enable       bool          `json:"enable" yaml:"enable"`
	Certificates []Certificate `json:"certificates" yaml:"certificates"` // 证书
	RootCAFile   []string      `json:"root_ca_file" yaml:"root_ca_file"` // 根CA文件
	NextProtos   []string      `json:"next_protos" yaml:"next_protos"`   // 支持的协议
	ServerName   string        `json:"server_name" yaml:"server_name"`   // 服务器名称

	// 0 不验证客户端证书
	// 1 请求客户端证书但不强制验证
	// 2 要求客户端证书但不验证CA
	// 3 如果提供客户端证书则验证
	// 4 要求并验证客户端证书
	ClientAuth                  tls.ClientAuthType       `json:"client_auth" yaml:"client_auth"`                                       // 客户端验证
	ClientCAFile                []string                 `json:"client_ca_file" yaml:"client_ca_file"`                                 // 客户端CA文件
	InsecureSkipVerify          bool                     `json:"insecure_skip_verify" yaml:"insecure_skip_verify"`                     // 跳过验证
	CipherSuites                []string                 `json:"cipher_suites" yaml:"cipher_suites"`                                   // 密钥套件
	CurvePreferences            []string                 `json:"curve_preferences" yaml:"curve_preferences"`                           // 曲线偏好
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
		ClientAuth:                  t.ClientAuth,
	}

	// 自定义曲线
	if len(t.CurvePreferences) > 0 {
		for _, curveName := range t.CurvePreferences {
			curveID, exists := CurvePreferencesMap[curveName]
			if !exists {
				return nil, fmt.Errorf("未知的曲线: %s", curveName)
			}
			config.CurvePreferences = append(config.CurvePreferences, curveID)
		}
	}

	// 解析 TLS 版本
	{
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
	}

	// 解析密码套件
	{
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
			cipherSuites, err = ParseCipherSuites(CompatibleCipherSuites)
			if err != nil {
				return nil, fmt.Errorf("解析默认密码套件失败: %v", err)
			}
			config.CipherSuites = cipherSuites
		}
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
	// noinspection all
	if pool, err := t.loadCaKit(t.RootCAFile...); err != nil {
		return nil, fmt.Errorf("服务端证书添加失败 %v", err)
	} else {
		config.RootCAs = pool
	}

	// 加载客户端证书
	if len(t.ClientCAFile) > 0 {
		pool, err := t.loadCaKit(t.ClientCAFile...)
		if err != nil {
			return nil, fmt.Errorf("客户端证书添加失败 %v", err)
		}
		config.ClientCAs = pool
	} else if t.ClientAuth != tls.NoClientCert {
		// 如果设置了客户端认证但没有提供 CA 文件，返回错误
		return nil, fmt.Errorf("客户端认证需要提供 client_ca_file")
	}
	return config, nil
}

// 载入 ca 证书工具方法
// 支持pem和der两种格式的证书。
func (t *TLSConfig) loadCaKit(cas ...string) (*x509.CertPool, error) {
	if len(cas) < 1 {
		return nil, nil
	}
	pool := x509.NewCertPool()
	for _, ca := range cas {
		caFile := tools.Fileabs(ca)
		caCert, err := os.ReadFile(caFile)
		if err != nil {
			return nil, fmt.Errorf("读取证书[%s]内容失败：%v", ca, err)
		}
		if pool.AppendCertsFromPEM(caCert) {
			continue
		}
		cert, err := x509.ParseCertificate(caCert)
		if err != nil {
			return nil, fmt.Errorf("解析证书[%s]失败：%v（既不是PEM也不是DER格式）", ca, err)
		}
		pool.AddCert(cert)
	}
	return pool, nil
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
