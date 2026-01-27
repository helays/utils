package filesaver

import (
	"fmt"

	"github.com/helays/utils/v2/file/filesaver/ftp"
	"github.com/helays/utils/v2/file/filesaver/hdfs"
	"github.com/helays/utils/v2/file/filesaver/localfile"
	"github.com/helays/utils/v2/file/filesaver/minio"
	"github.com/helays/utils/v2/file/filesaver/sftp"
)

func New(cfg *Config) (FileSaver, error) {
	switch cfg.Driver {
	case DriverLocal:
		return localfile.New(&cfg.Local)
	case DriverSftp:
		return sftp.New(&cfg.SFTP)
	case DriverFtp:
		return ftp.New(&cfg.FTP)
	case DriverHdfs:
		return hdfs.New(&cfg.HDFS)
	case DriverMinio:
		return minio.New(&cfg.Minio)
	//case DriverCeph:
	default:
		panic(fmt.Errorf("不支持的文件系统驱动 %s", cfg.Driver))
	}
}
