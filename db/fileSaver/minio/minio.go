package minio

import (
	"context"
	"fmt"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"io"
	"strings"
)

type Options struct {
	Bucket        string `json:"bucket" yaml:"bucket" ini:"bucket"`                         // 存储桶名称
	Region        string `json:"region" yaml:"region" ini:"region"`                         //指定 Bucket 所在的区域（Region）。MinIO 默认使用 us-east-1 作为区域
	ObjectLocking bool   `json:"object_locking" yaml:"object_locking" ini:"object_locking"` //是否启用对象锁定（Object Locking）功能
	ExistIgnore   bool   `json:"exist_ignore" yaml:"exist_ignore" ini:"exist_ignore"`       // 如果文件已存在，是否忽略写入
}

type Config struct {
	Endpoint        string  `json:"endpoint" yaml:"endpoint" ini:"endpoint"`                            // MinIO 节点地址（单点或集群）
	AccessKeyID     string  `json:"access_key_id" yaml:"access_key_id" ini:"access_key_id"`             // 访问密钥
	SecretAccessKey string  `json:"secret_access_key" yaml:"secret_access_key" ini:"secret_access_key"` // 秘密密钥
	UseSSL          bool    `json:"use_ssl" yaml:"use_ssl" ini:"use_ssl"`                               // 是否使用 HTTPS
	options         Options // 配置项
	// 客户端
	ctx    context.Context
	cancel context.CancelFunc
	client *minio.Client
}

func (this *Config) Close() {
	if this.client == nil {
		return
	}
	this.cancel()
	this.client = nil
}

// 写入文件
func (this *Config) Write(filePath string, src io.Reader, options Options) (int64, error) {
	if err := this.login(); err != nil {
		return 0, err
	}
	this.options = options
	if err := this.createBucket(); err != nil {
		return 0, err
	}
	if this.options.ExistIgnore {
		_, err := this.client.StatObject(this.ctx, this.options.Bucket, filePath, minio.StatObjectOptions{})
		if err == nil {
			return 0, nil
		} else if _err := err.Error(); !strings.Contains(_err, "key does not exist") {
			return 0, fmt.Errorf("判断文件%s是否存在失败：%s", filePath, _err)
		}
		// 下面就是文件不存在，支持继续处理
	}
	if info, err := this.client.PutObject(this.ctx, this.options.Bucket, filePath, src, -1, minio.PutObjectOptions{}); err != nil {
		return 0, fmt.Errorf("上传文件%s失败：%s", filePath, err.Error())
	} else {
		return info.Size, nil
	}
}

func (this *Config) Read(filePath string, options Options) (io.ReadCloser, error) {
	if err := this.login(); err != nil {
		return nil, err
	}
	this.options = options
	if obj, err := this.client.GetObject(this.ctx, this.options.Bucket, filePath, minio.GetObjectOptions{}); err != nil {
		return nil, fmt.Errorf("打开文件%s失败：%s", filePath, err.Error())
	} else {
		return obj, nil
	}
}

// 登录
func (this *Config) login() error {
	if this.client != nil {
		return nil
	}

	options := &minio.Options{
		Creds:  credentials.NewStaticV4(this.AccessKeyID, this.SecretAccessKey, ""),
		Secure: this.UseSSL,
	}
	var err error
	if this.client, err = minio.New(this.Endpoint, options); err != nil {
		return fmt.Errorf("连接MinIO节点失败: %s", err.Error())
	}
	this.ctx, this.cancel = context.WithCancel(context.Background())
	return nil
}

// 创建bucket 存储桶
func (this *Config) createBucket() error {
	if ok, err := this.client.BucketExists(this.ctx, this.options.Bucket); ok {
		return nil
	} else if err != nil {
		return fmt.Errorf("查询bucket %s失败: %s", this.options.Bucket, err.Error())
	}
	// 这里创建bucket
	err := this.client.MakeBucket(this.ctx, this.options.Bucket, minio.MakeBucketOptions{
		Region:        this.options.Region,
		ObjectLocking: this.options.ObjectLocking,
	})
	if err != nil {
		return fmt.Errorf("创建bucket %s失败: %s", this.options.Bucket, err.Error())
	}
	return nil
}
