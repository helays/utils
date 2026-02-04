package localfile

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"helay.net/go/utils/v3/close/vclose"
	"helay.net/go/utils/v3/tools"
)

type Config struct {
	Root string `json:"root" yaml:"root" ini:"root"`
}

type Saver struct {
	opt *Config
}

func New(cfg *Config) (*Saver, error) {
	return &Saver{opt: cfg}, nil
}

func (s *Saver) Write(p string, src io.Reader, existIgnores ...bool) (int64, error) {
	path, err := s.realPath(p)
	if err != nil {
		return 0, err
	}
	if len(existIgnores) > 0 && existIgnores[0] {
		// 如果启用 文件存在就忽略，首先判断文件是否存在，
		// 如果文件存在，就中断处理
		// 如果err有问题，判断是否因为文件不存在导致的。
		if _, err = os.Stat(path); err == nil {
			return 0, nil
		} else if !os.IsNotExist(err) {
			return 0, err
		}
	}
	dir := filepath.Dir(path)
	if err = tools.Mkdir(dir); err != nil {
		return 0, fmt.Errorf("创建目录[%s]失败:%s", dir, err)
	}
	file, err := os.OpenFile(path, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
	defer vclose.Close(file)
	if err != nil {
		return 0, fmt.Errorf("创建文件[%s]失败:%s", path, err)
	}
	var n int64
	n, err = io.Copy(file, src)
	if err != nil {
		return n, fmt.Errorf("写入文件[%s]失败:%s", path, err)
	}
	return n, nil
}

func (s *Saver) Read(p string) (io.ReadCloser, error) {
	path, err := s.realPath(p)
	if err != nil {
		return nil, err
	}
	file, err := os.Open(path)
	defer vclose.Close(file)
	if err != nil {
		return nil, fmt.Errorf("打开文件[%s]失败:%s", path, err)
	}
	return file, nil
}

func (s *Saver) ListFiles(p string) ([]string, error) {
	path, err := s.realPath(p)
	if err != nil {
		return nil, err
	}
	entries, err := os.ReadDir(path)
	if err != nil {
		return nil, fmt.Errorf("读取目录[%s]失败:%s", path, err)
	}
	var filePaths []string
	for _, entry := range entries {
		if !entry.IsDir() {
			filePaths = append(filePaths, entry.Name())
		}
	}
	return filePaths, nil
}

func (s *Saver) Delete(p string) error {
	path, err := s.realPath(p)
	if err != nil {
		return err
	}
	return os.Remove(path)
}

func (s *Saver) DeleteAll(p string) error {
	path, err := s.realPath(p)
	if err != nil {
		return err
	}
	return os.RemoveAll(path)
}

func (s *Saver) Close() error {
	return nil
}

func (s *Saver) realPath(p string) (string, error) {
	if tools.ContainsDotDot(p) {
		return "", fmt.Errorf("路径[%s]包含 '..'", p)
	}
	return tools.Fileabs(filepath.Join(s.opt.Root, p)), nil
}
