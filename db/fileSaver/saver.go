package fileSaver

import (
	"fmt"
	"github.com/helays/utils/config"
	"github.com/helays/utils/db/fileSaver/ftp"
	"github.com/helays/utils/db/fileSaver/hdfs"
	"github.com/helays/utils/db/fileSaver/local"
	"github.com/helays/utils/db/fileSaver/minio"
	"github.com/helays/utils/db/fileSaver/sftp"
	"io"
	"path"
	"strings"
)

type Saver struct {
	StorageType string `json:"storage_type" yaml:"storage_type" ini:"storage_type"` // 存储类型 local、sftp、ftp、hdfs、miniio等
	Root        string `json:"root" yaml:"root" ini:"root"`                         // 文件跟路径，如果/开头，最终路径为/root/path，如果没有/,最终路径是current_path/root/path

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
		return this.MinioConfig.Read(p)
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
