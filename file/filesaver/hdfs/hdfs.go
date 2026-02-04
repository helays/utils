package hdfs

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/colinmarc/hdfs/v2"
	"github.com/helays/utils/v2/close/vclose"
)

type Saver struct {
	opt    *Config
	client *hdfs.Client
}

func New(cfg *Config) (*Saver, error) {
	s := &Saver{opt: cfg}
	if err := cfg.Valid(); err != nil {
		return nil, err
	}
	if err := s.login(); err != nil {
		return nil, err
	}
	return s, nil
}

func (s *Saver) Close() error {
	if s.client == nil {
		return nil
	}
	vclose.Close(s.client)
	s.client = nil
	return nil
}

func (s *Saver) Write(p string, src io.Reader, existIgnores ...bool) (int64, error) {
	if !filepath.IsAbs(p) {
		p = filepath.Join("/", p)
	}
	if ok, err := s.exist(p); ok {
		if len(existIgnores) > 0 && existIgnores[0] {
			return 0, nil
		}
		// 删除文件，重写
		if err = s.client.Remove(p); err != nil {
			return 0, fmt.Errorf("删除文件%s失败: %s", p, err.Error())
		}
	} else if err != nil {
		return 0, err
	}
	dir := filepath.Dir(p)
	if err := s.client.MkdirAll(dir, 0755); err != nil {
		return 0, fmt.Errorf("创建目录%s失败: %s", dir, err.Error())
	}
	remoteFile, err := s.client.Create(p)
	defer vclose.Close(remoteFile)
	if err != nil {
		return 0, fmt.Errorf("创建文件%s失败: %s", p, err.Error())
	}
	return io.Copy(remoteFile, src)
}

func (s *Saver) Read(p string) (io.ReadCloser, error) {
	if !filepath.IsAbs(p) {
		p = filepath.Join("/", p)
	}
	remoteFile, err := s.client.Open(p)
	if err != nil {
		return nil, err
	}
	return remoteFile, nil
}

// ListFiles 列出目录下的文件
func (s *Saver) ListFiles(p string) ([]string, error) {
	if !filepath.IsAbs(p) {
		p = filepath.Join("/", p)
	}
	entries, err := s.client.ReadDir(p)
	if err != nil {
		return nil, err
	}
	var files []string
	for _, entry := range entries {
		if !entry.IsDir() {
			files = append(files, entry.Name())
		}
	}
	return files, nil
}

func (s *Saver) Delete(p string) error {
	if !filepath.IsAbs(p) {
		p = filepath.Join("/", p)
	}
	if ok, err := s.exist(p); !ok {
		if err != nil {
			return err
		}
		return nil
	}
	return s.client.Remove(p)
}

func (s *Saver) DeleteAll(p string) error {
	if !filepath.IsAbs(p) {
		p = filepath.Join("/", p)
	}
	return s.client.RemoveAll(p)
}

func (s *Saver) login() error {
	if s.client != nil {
		return nil
	}
	var err error
	s.client, err = hdfs.NewClient(hdfs.ClientOptions{
		Addresses:                    s.opt.Addresses,                    // 指定要连接的 NameNode 地址列表。
		User:                         s.opt.User,                         // 指定客户端以哪个 HDFS 用户身份进行操作
		UseDatanodeHostname:          s.opt.UseDatanodeHostname,          // 指定客户端是否通过主机名（而不是 IP 地址）连接 DataNode。
		NamenodeDialFunc:             nil,                                // 自定义连接 NameNode 的拨号函数。
		DatanodeDialFunc:             nil,                                // 自定义连接 DataNode 的拨号函数。
		KerberosClient:               nil,                                // 于连接启用了 Kerberos 认证的 HDFS 集群。
		KerberosServicePrincipleName: s.opt.KerberosServicePrincipleName, // 指定 NameNode 的 Kerberos 服务主体名称（SPN）。格式为 <SERVICE>/<FQDN>，例如 nn/_HOST。
		DataTransferProtection:       s.opt.DataTransferProtection,       // 指定与 DataNode 通信时的数据保护级别。
	})
	if err != nil {
		return fmt.Errorf("hdfs连接失败 %v", err)
	}
	return nil
}

func (s *Saver) exist(p string) (bool, error) {
	if _, err := s.client.Stat(p); err == nil {
		return true, nil
	} else if !os.IsNotExist(err) {
		return false, err
	}
	return false, nil
}
