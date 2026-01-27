package filesaver

import (
	"io"

	"github.com/helays/utils/v2/file/filesaver/ftp"
	"github.com/helays/utils/v2/file/filesaver/hdfs"
	"github.com/helays/utils/v2/file/filesaver/localfile"
	"github.com/helays/utils/v2/file/filesaver/minio"
	"github.com/helays/utils/v2/file/filesaver/sftp"
)

type FileSaver interface {
	Write(p string, src io.Reader, existIgnores ...bool) (int64, error) // 写入文件
	Read(p string) (io.ReadCloser, error)                               // 读取文件
	Delete(p string) error                                              // 删除指定文件
	DeleteAll(p string) error                                           // 删除文件夹
	ListFiles(p string) ([]string, error)                               // 列出指定目录下的所有文件
	Close() error                                                       // 关闭资源
}

type Driver string

// noinspection all
const (
	DriverLocal Driver = "local"
	DriverSftp  Driver = "sftp"
	DriverFtp   Driver = "ftp"
	DriverHdfs  Driver = "hdfs"
	DriverMinio Driver = "minio"
	DriverCeph  Driver = "ceph"
)

// noinspection all
type Config struct {
	Driver Driver `json:"driver" yaml:"driver" ini:"driver"`

	Local localfile.Config `json:"local" yaml:"local" ini:"local"` // 本地文件系统
	FTP   ftp.Config       `json:"ftp" yaml:"ftp" ini:"ftp"`       // ftp
	SFTP  sftp.Config      `json:"sftp" yaml:"sftp" ini:"sftp"`    // sftp
	HDFS  hdfs.Config      `json:"hdfs" yaml:"hdfs" ini:"hdfs"`    // hdfs
	Minio minio.Config     `json:"minio" yaml:"minio" ini:"minio"` // minio
}
