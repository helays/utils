package local

import (
	"fmt"
	"github.com/helays/utils/close/vclose"
	"github.com/helays/utils/tools"
	"io"
	"os"
	"path"
)

type Local struct{}

// Write 写入文件
func (this Local) Write(p string, src io.Reader, existIgnores ...bool) (int64, error) {
	filePath := tools.Fileabs(p)
	if len(existIgnores) > 0 && existIgnores[0] {
		// 如果启用 文件存在就忽略，首先判断文件是否存在，
		// 如果文件存在，就中断处理
		// 如果err有问题，判断是否因为文件不存在导致的。
		if _, err := os.Stat(filePath); err == nil {
			return 0, nil
		} else if !os.IsNotExist(err) {
			return 0, err
		}
	}
	dir := path.Dir(filePath)
	if err := tools.Mkdir(dir); err != nil {
		return 0, fmt.Errorf("创建目录%s失败: %s", dir, err.Error())
	}
	file, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return 0, fmt.Errorf("打开文件%s失败: %s", filePath, err.Error())
	}
	var written int64
	written, err = io.Copy(file, src)
	if err != nil {
		return written, fmt.Errorf("写入文件%s失败: %s", filePath, err.Error())
	}
	return written, nil
}

func (this Local) Read(p string) (io.ReadCloser, error) {
	filePath := tools.Fileabs(p)
	file, err := os.Open(filePath)
	defer vclose.Close(file)
	if err != nil {
		return nil, fmt.Errorf("打开文件%s失败: %s", filePath, err.Error())
	}
	return file, nil
}

func (this Local) ListFiles(dirPath string) ([]string, error) {
	dirPath = tools.Fileabs(dirPath)
	entries, err := os.ReadDir(dirPath)
	if err != nil {
		return nil, fmt.Errorf("读取目录%s失败: %s", dirPath, err.Error())
	}
	var filePaths []string
	for _, entry := range entries {
		if !entry.IsDir() {
			filePaths = append(filePaths, entry.Name())
		}
	}
	return filePaths, nil
}

// Delete 删除文件
func (this Local) Delete(p string) error {
	filePath := tools.Fileabs(p)
	return os.Remove(filePath)
}

// DeleteAll 删除文件
func (this Local) DeleteAll(p string) error {
	filePath := tools.Fileabs(p)
	return os.RemoveAll(filePath)
}
