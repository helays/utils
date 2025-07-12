package fileSaver

import (
	"fmt"
	"github.com/helays/utils/v2/config"
	"github.com/helays/utils/v2/db/fileSaver/ftp"
	"github.com/helays/utils/v2/db/fileSaver/hdfs"
	"github.com/helays/utils/v2/db/fileSaver/local"
	"github.com/helays/utils/v2/db/fileSaver/minio"
	"github.com/helays/utils/v2/db/fileSaver/sftp"
	"io"
	"path"
	"strings"
)

type Saver struct {
	StorageType string `json:"storage_type" yaml:"storage_type" ini:"storage_type"` // 存储类型 local、sftp、ftp、hdfs、miniio等
	// 本地文件系统：/开头，最终路径为/root/path，如果没有/,最终路径是current_path/root/path
	// ftp、sftp：最终生成路径是/userHome/root/path
	Root string `json:"root" yaml:"root" ini:"root"`

	local.Local  `json:"local" yaml:"local" ini:"local"` // 本地文件系统
	SftpConfig   sftp.Config                             `json:"sftp_config" yaml:"sftp_config" ini:"sftp_config"`       // sftp客户端配置
	FtpConfig    ftp.Config                              `json:"ftp_config" yaml:"ftp_config" ini:"ftp_config"`          // ftp客户端配置
	HdfsConfig   hdfs.Config                             `json:"hdfs_config" yaml:"hdfs_config" ini:"hdfs_config"`       // hdfs客户端配置
	MinioConfig  minio.Config                            `json:"minio_config" yaml:"minio_config" ini:"minio_config"`    // minio客户端配置
	MinioOptions minio.Options                           `json:"minio_options" yaml:"minio_options" ini:"minio_options"` // minio客户端配置
}

// Write 写入文件
func (this *Saver) Write(p string, src io.Reader, existIgnores ...bool) (int64, error) {
	p = path.Join(this.Root, p)
	switch strings.ToLower(this.StorageType) {
	case config.FileTypeLocal: // 本地文件系统
		return this.Local.Write(p, src, existIgnores...)
	case config.FileTypeFtp: // ftp
		return this.FtpConfig.Write(p, src, existIgnores...)
	case config.FileTypeSftp: // sftp
		return this.SftpConfig.Write(p, src, existIgnores...)
	case config.FileTypeHdfs: // hdfs
		return this.HdfsConfig.Write(p, src, existIgnores...)
	case config.FileTypeMinio:
		return this.MinioConfig.Write(p, src, this.MinioOptions)
	default:
		return 0, fmt.Errorf("不支持的存储类型: %s", this.StorageType)
	}
}

// Read 读取文件
func (this *Saver) Read(p string) (io.ReadCloser, error) {
	p = path.Join(this.Root, p)
	switch strings.ToLower(this.StorageType) {
	case config.FileTypeLocal: // 本地文件系统
		return this.Local.Read(p)
	case config.FileTypeFtp: // ftp
		return this.FtpConfig.Read(p)
	case config.FileTypeSftp: // sftp
		return this.SftpConfig.Read(p)
	case config.FileTypeHdfs: // hdfs
		return this.HdfsConfig.Read(p)
	case config.FileTypeMinio:
		return this.MinioConfig.Read(p, this.MinioOptions)
	default:
		return nil, fmt.Errorf("不支持的存储类型: %s", this.StorageType)
	}
}

func (this *Saver) Delete(p string) error {
	p = path.Join(this.Root, p)
	switch strings.ToLower(this.StorageType) {
	case config.FileTypeLocal:
		return this.Local.Delete(p)
	case config.FileTypeFtp:
		return this.FtpConfig.Delete(p)
	case config.FileTypeSftp:
		return this.SftpConfig.Delete(p)
	case config.FileTypeHdfs:
		return this.HdfsConfig.Delete(p)
	case config.FileTypeMinio:
		return this.MinioConfig.Delete(p, this.MinioOptions)
	default:
		return fmt.Errorf("不支持的存储类型: %s", this.StorageType)
	}
}

func (this *Saver) DeleteAll(p string) error {
	p = path.Join(this.Root, p)
	switch strings.ToLower(this.StorageType) {
	case config.FileTypeLocal:
		return this.Local.DeleteAll(p)
	case config.FileTypeFtp:
		return this.FtpConfig.DeleteAll(p)
	case config.FileTypeSftp:
		return this.SftpConfig.DeleteAll(p)
	case config.FileTypeHdfs:
		return this.HdfsConfig.DeleteAll(p)
	case config.FileTypeMinio:
		return this.MinioConfig.Delete(p, this.MinioOptions)
	default:
		return fmt.Errorf("不支持的存储类型: %s", this.StorageType)
	}
}

func (this *Saver) ListFiles(dirPath string) ([]string, error) {
	dirPath = path.Join(this.Root, dirPath)
	switch strings.ToLower(this.StorageType) {
	case config.FileTypeLocal:
		return this.Local.ListFiles(dirPath)
	case config.FileTypeFtp:
		return this.FtpConfig.ListFiles(dirPath)
	case config.FileTypeSftp:
		return this.SftpConfig.ListFiles(dirPath)
	case config.FileTypeHdfs:
		return this.HdfsConfig.ListFiles(dirPath)
	case config.FileTypeMinio:
		return this.MinioConfig.ListFiles(dirPath, this.MinioOptions)
	default:
		return nil, fmt.Errorf("不支持的存储类型: %s", this.StorageType)
	}
}

// Close 关闭资源
func (this *Saver) Close() {
	switch strings.ToLower(this.StorageType) {
	case config.FileTypeLocal:
	case config.FileTypeFtp: // ftp
		this.FtpConfig.Close()
	case config.FileTypeSftp:
		this.SftpConfig.Close()
	case config.FileTypeHdfs: // hdfs
		this.HdfsConfig.Close()
	case config.FileTypeMinio:
		this.MinioConfig.Close()
	default:

	}
}
