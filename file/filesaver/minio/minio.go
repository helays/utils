package minio

import (
	"context"
	"fmt"
	"io"
	"path/filepath"
	"strings"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type Saver struct {
	opt    *Config
	client *minio.Client

	ctx    context.Context
	cancel context.CancelFunc
}

func New(cfg *Config) (*Saver, error) {
	s := &Saver{opt: cfg}
	if err := s.opt.Valid(); err != nil {
		return nil, err
	}
	if err := s.login(); err != nil {
		return nil, err
	}
	return s, nil
}

func (s *Saver) Write(p string, src io.Reader, existIgnores ...bool) (int64, error) {
	if len(existIgnores) > 0 && existIgnores[0] {
		if _, err := s.client.StatObject(s.ctx, s.opt.Options.Bucket, p, minio.StatObjectOptions{}); err == nil {
			return 0, err
		} else if _err := err.Error(); !strings.Contains(_err, "key does not exist") {
			return 0, fmt.Errorf("文件已存在: %s", _err)
		}
	}
	info, err := s.client.PutObject(s.ctx, s.opt.Options.Bucket, p, src, -1, minio.PutObjectOptions{})
	if err != nil {
		return 0, err
	}
	return info.Size, nil
}

func (s *Saver) Read(p string) (io.ReadCloser, error) {
	return s.client.GetObject(s.ctx, s.opt.Options.Bucket, p, minio.GetObjectOptions{})
}

func (s *Saver) ListFiles(p string) ([]string, error) {
	var files []string
	opts := minio.ListObjectsOptions{
		Prefix:    p,
		Recursive: true,
	}
	for obj := range s.client.ListObjects(s.ctx, s.opt.Options.Bucket, opts) {
		if obj.Err != nil {
			return nil, obj.Err
		}
		files = append(files, filepath.Base(obj.Key))
	}
	return files, nil
}

func (s *Saver) Delete(p string) error {
	return s.client.RemoveObject(s.ctx, s.opt.Options.Bucket, p, minio.RemoveObjectOptions{})
}

func (s *Saver) DeleteAll(p string) error {
	return s.client.RemoveObject(s.ctx, s.opt.Options.Bucket, p, minio.RemoveObjectOptions{})
}

func (s *Saver) Close() error {
	if s.client == nil {
		return nil
	}
	s.cancel()
	s.client = nil
	return nil
}

func (s *Saver) login() error {
	if s.client != nil {
		return nil
	}
	options := &minio.Options{
		Creds:  credentials.NewStaticV4(s.opt.AccessKeyID, s.opt.SecretAccessKey, ""),
		Secure: s.opt.UseSSL,
	}
	var err error
	if s.client, err = minio.New(s.opt.Endpoint, options); err != nil {
		return fmt.Errorf("连接MinIO节点失败: %s", err.Error())
	}

	s.ctx, s.cancel = context.WithCancel(context.Background())
	return s.createBucket()
}

// 创建 bucket
func (s *Saver) createBucket() error {
	if ok, err := s.client.BucketExists(s.ctx, s.opt.Options.Bucket); ok {
		return nil
	} else if err != nil {
		return fmt.Errorf("查询bucket %s失败: %s", s.opt.Options.Bucket, err.Error())
	}
	err := s.client.MakeBucket(s.ctx, s.opt.Options.Bucket, minio.MakeBucketOptions{
		Region:        s.opt.Options.Region,
		ObjectLocking: s.opt.Options.ObjectLocking,
	})
	if err != nil {
		return fmt.Errorf("创建bucket %s失败: %s", s.opt.Options.Bucket, err.Error())
	}
	return nil
}
