package httpfile

import (
	"github.com/helays/utils/v2/close/vclose"
	"net/http"
	"os"
)

// NoSymlinkFileSystem 禁用符号链接的文件系统包装器
type NoSymlinkFileSystem struct {
	Fs http.FileSystem
}

func (nfs NoSymlinkFileSystem) Open(name string) (http.File, error) {
	f, err := nfs.Fs.Open(name)
	if err != nil {
		return nil, err
	}

	// 检查是否为符号链接
	if fi, statErr := f.Stat(); statErr == nil && fi.Mode()&os.ModeSymlink != 0 {
		vclose.Close(f)
		return nil, os.ErrNotExist
	}

	return f, nil
}
