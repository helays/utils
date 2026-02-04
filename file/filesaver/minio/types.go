package minio

import (
	"fmt"

	"helay.net/go/utils/v3/config"
)

type Config struct {
	Endpoint        string  `json:"endpoint" yaml:"endpoint" ini:"endpoint"`                            // MinIO 节点地址（单点或集群）
	AccessKeyID     string  `json:"access_key_id" yaml:"access_key_id" ini:"access_key_id"`             // 访问密钥
	SecretAccessKey string  `json:"secret_access_key" yaml:"secret_access_key" ini:"secret_access_key"` // 秘密密钥
	UseSSL          bool    `json:"use_ssl" yaml:"use_ssl" ini:"use_ssl"`                               // 是否使用 HTTPS
	Options         Options `json:"options" yaml:"options" ini:"options"`                               // 配置项
}

type Options struct {
	Bucket        string `json:"bucket" yaml:"bucket" ini:"bucket"`                         // 存储桶名称
	Region        string `json:"region" yaml:"region" ini:"region"`                         //指定 Bucket 所在的区域（Region）。MinIO 默认使用 us-east-1 作为区域
	ObjectLocking bool   `json:"object_locking" yaml:"object_locking" ini:"object_locking"` //是否启用对象锁定（Object Locking）功能
}

func (c *Config) Valid() error {
	if c.Endpoint == "" {
		return fmt.Errorf("缺失地址")
	}
	return nil
}

func (c *Config) RemovePasswd() {
	c.SecretAccessKey = ""
}

func (c *Config) SetInfo(args ...any) {
	if len(args) != 2 {
		return
	}
	switch args[0].(string) {
	case config.ClientInfoHost:
		c.Endpoint = args[1].(string)
	case config.ClientInfoUser:
		c.AccessKeyID = args[1].(string)
	case config.ClientInfoPasswd:
		c.SecretAccessKey = args[1].(string)
	}
}
