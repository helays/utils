package hdfs

import (
	"fmt"

	"helay.net/go/utils/v3/config"
)

type Config struct {
	Addresses []string `json:"addresses" yaml:"addresses" ini:"addresses,omitempty"` // 路径
	User      string   `json:"user" yaml:"user" ini:"user"`
	// 指定客户端是否通过主机名（而不是 IP 地址）连接 DataNode。
	UseDatanodeHostname bool `json:"use_datanode_hostname" yaml:"use_datanode_hostname" ini:"use_datanode_hostname"`
	// 指定 NameNode 的 Kerberos 服务主体名称（SPN）。格式为 <SERVICE>/<FQDN>，例如 nn/_HOST。
	KerberosServicePrincipleName string `json:"kerberos_service_principle_name" yaml:"kerberos_service_principle_name" ini:"kerberos_service_principle_name"`
	// 指定与 DataNode 通信时的数据保护级别。
	// authentication：仅认证;
	// integrity： 认证 + 数据完整性校验
	// integrity+privacy： 认证 + 数据完整性校验 + 数据加密
	DataTransferProtection string `json:"data_transfer_protection" yaml:"data_transfer_protection" ini:"data_transfer_protection"`
}

func (c *Config) Valid() error {
	if len(c.Addresses) < 1 {
		return fmt.Errorf("缺失地址")
	}

	return nil
}

func (c *Config) SetInfo(args ...any) {
	if len(args) != 2 {
		return
	}
	switch args[0].(string) {
	case config.ClientInfoHost:
		c.Addresses = args[1].([]string)
	case config.ClientInfoUser:
		c.User = args[1].(string)
	}
}
