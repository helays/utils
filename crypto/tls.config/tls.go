package tls_config

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"os"

	"github.com/helays/utils/v2/tools"
)

type Config struct {
	CaFile             string `json:"ca_file" yaml:"ca_file" ini:"ca_file"`                                        // CA证书文件
	CertFile           string `json:"cert_file" yaml:"cert_file" ini:"cert_file"`                                  // 证书文件
	KeyFile            string `json:"key_file" yaml:"key_file" ini:"key_file"`                                     // 密钥文件
	InsecureSkipVerify bool   `json:"insecure_skip_verify" yaml:"insecure_skip_verify" ini:"insecure_skip_verify"` // 跳过证书验证
}

// NewTLSConfig 创建一个新的TLS配置
func (this Config) NewTLSConfig() (*tls.Config, error) {
	// 创建一个新的TLS配置
	tlsConfig := &tls.Config{
		InsecureSkipVerify: this.InsecureSkipVerify, // 跳过证书验证
	}

	// 如果CA文件不为空，则加载CA证书
	if this.CaFile != "" {
		caCert, err := os.ReadFile(tools.Fileabs(this.CaFile))
		if err != nil {
			return nil, fmt.Errorf("载入CA文件失败：%s", err.Error())
		}
		tlsConfig.RootCAs = x509.NewCertPool()
		tlsConfig.RootCAs.AppendCertsFromPEM(caCert)
	}
	cert, err := tls.LoadX509KeyPair(tools.Fileabs(this.CertFile), tools.Fileabs(this.KeyFile))
	if err != nil {
		return nil, fmt.Errorf("载入证书或者密钥文件失败：%s", err.Error())
	}
	tlsConfig.Certificates = []tls.Certificate{cert}
	return tlsConfig, nil
}
