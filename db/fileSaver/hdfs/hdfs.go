package hdfs

import (
	"fmt"
	"github.com/colinmarc/hdfs/v2"
	"github.com/helays/utils/close/vclose"
	"io"
	"os"
	"path"
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

	client *hdfs.Client
}

// Close 关闭 hdfs client
func (this *Config) Close() {
	if this.client == nil {
		return
	}
	vclose.Close(this.client)
	this.client = nil
}

// Write 写入文件
func (this *Config) Write(filePath string, src io.Reader, existIgnores ...bool) error {
	if err := this.login(); err != nil {
		return err
	}
	if !path.IsAbs(filePath) {
		filePath = path.Join("/", filePath)
	}
	if ok, err := this.exist(filePath); ok {
		if len(existIgnores) > 0 && existIgnores[0] {
			// 忽略写入
			return nil
		}
		// 删除文件，重写
		if err = this.client.Remove(filePath); err != nil {
			return fmt.Errorf("删除文件%s失败: %s", filePath, err.Error())
		}
	} else if err != nil {
		return err
	}
	dir := path.Dir(filePath)
	if err := this.client.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("创建目录%s失败: %s", dir, err.Error())
	}

	remoteFile, err := this.client.Create(filePath)
	if err != nil {
		return fmt.Errorf("创建文件%s失败: %s", filePath, err.Error())
	}
	defer vclose.Close(remoteFile)
	if _, err = io.Copy(remoteFile, src); err != nil {
		return fmt.Errorf("写入文件%s失败: %s", filePath, err.Error())
	}
	return nil
}

// login 登录
func (this *Config) login() error {
	if this.client != nil {
		return nil
	}
	var err error
	this.client, err = hdfs.NewClient(hdfs.ClientOptions{
		Addresses:                    this.Addresses,                    // 指定要连接的 NameNode 地址列表。
		User:                         this.User,                         // 指定客户端以哪个 HDFS 用户身份进行操作
		UseDatanodeHostname:          this.UseDatanodeHostname,          // 指定客户端是否通过主机名（而不是 IP 地址）连接 DataNode。
		NamenodeDialFunc:             nil,                               // 自定义连接 NameNode 的拨号函数。
		DatanodeDialFunc:             nil,                               // 自定义连接 DataNode 的拨号函数。
		KerberosClient:               nil,                               // 于连接启用了 Kerberos 认证的 HDFS 集群。
		KerberosServicePrincipleName: this.KerberosServicePrincipleName, // 指定 NameNode 的 Kerberos 服务主体名称（SPN）。格式为 <SERVICE>/<FQDN>，例如 nn/_HOST。
		DataTransferProtection:       this.DataTransferProtection,       // 指定与 DataNode 通信时的数据保护级别。
	})
	if err != nil {
		return fmt.Errorf("hdfs连接失败：%s", err.Error())
	}
	return nil
}

// exist 判断文件是否存在
func (this Config) exist(sPath string) (bool, error) {
	if _, err := this.client.Stat(sPath); err == nil {
		return true, nil
	} else if !os.IsNotExist(err) {
		return false, err
	}
	return false, nil
}
