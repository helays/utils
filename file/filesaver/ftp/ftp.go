package ftp

import (
	"fmt"
	"io"

	"path/filepath"
	"strings"

	"github.com/jlaffaye/ftp"
	"helay.net/go/utils/v3/close/ftpClose"
	"helay.net/go/utils/v3/dataType/customWriter"
)

type Saver struct {
	opt    *Config
	client *ftp.ServerConn
}

func New(cfg *Config) (*Saver, error) {
	s := &Saver{opt: cfg}
	if err := cfg.Valid(); err != nil {
		return nil, err
	}
	if err := s.login(); err != nil {
		return nil, err
	}
	return s, nil
}

// Write 写入文件
func (s *Saver) Write(p string, src io.Reader, existIgnores ...bool) (int64, error) {
	path, err := s.setPath(p)
	if err != nil {
		return 0, err
	}
	if len(existIgnores) > 0 && existIgnores[0] {
		if exist, err := s.exist(path); err != nil {
			return 0, err
		} else if exist {
			return 0, nil
		}
	}
	dir := filepath.Dir(path)
	if err = s.mkdir(dir); err != nil {
		return 0, err
	}
	counter := &customWriter.SizeCounter{}
	teeReader := io.TeeReader(src, counter)
	if err = s.client.Stor(path, teeReader); err != nil {
		return 0, err
	}
	return counter.TotalSize, nil
}

// Read 读取文件
func (s *Saver) Read(p string) (io.ReadCloser, error) {
	path, err := s.setPath(p)
	if err != nil {
		return nil, err
	}
	return s.client.Retr(path)
}

func (s *Saver) ListFiles(dirPath string) ([]string, error) {
	path, err := s.setPath(dirPath)
	if err != nil {
		return nil, err
	}
	entries, err := s.client.List(path)
	if err != nil {
		return nil, err
	}
	var fileNames []string
	for _, entry := range entries {
		if entry.Type == ftp.EntryTypeFile {
			fileNames = append(fileNames, entry.Name)
		}
	}
	return fileNames, nil
}

func (s *Saver) Delete(p string) error {
	path, err := s.setPath(p)
	if err != nil {
		return err
	}
	return s.client.Delete(path)
}

func (s *Saver) DeleteAll(p string) error {
	path, err := s.setPath(p)
	if err != nil {
		return err
	}
	return s.client.RemoveDirRecur(path)
}

func (s *Saver) Close() error {
	if s.client == nil {
		return nil
	}
	ftpClose.CloseFtpClient(s.client)
	s.client = nil
	return nil
}

// 设置当前文件全路径
func (s *Saver) setPath(p string) (string, error) {
	//if path.IsAbs(p) {
	//	return p, nil
	//}
	// 获取当前目录
	current, err := s.client.CurrentDir()
	if err != nil {
		return "", fmt.Errorf("获取当前目录失败：%s", err.Error())
	}
	return filepath.Join(current, p), nil
}

// 判断文件是否存在
func (s *Saver) exist(p string) (bool, error) {
	remotePath := filepath.Dir(p)
	lst, err := s.client.List(remotePath)
	if err != nil {
		return false, fmt.Errorf("判断文件%s是否存在失败：%s", p, err.Error())
	}
	remoteFileName := filepath.Base(p)
	for _, v := range lst {
		if v.Name == remoteFileName {
			return true, nil
		}
	}
	return false, nil
}

func (s *Saver) mkdir(p string) error {
	if ok, err := s.exist(p); err != nil {
		return err
	} else if ok {
		return nil
	}
	return s.mkdirALL(p)
}

func (s *Saver) mkdirALL(p string) error {
	currentDir, err := s.client.CurrentDir()
	if err != nil {
		return fmt.Errorf("获取当前目录失败：%s", err.Error())
	}
	path := strings.TrimPrefix(p, currentDir)
	var currentPath string
	for _, part := range strings.Split(path, "/") {
		if part == "" {
			continue // Skip empty parts which can happen with leading/trailing slashes or double slashes.
		}
		currentPath = fmt.Sprintf("%s/%s", currentPath, part)
		err = s.client.ChangeDir(currentPath)
		if err != nil {
			// Directory does not exist, so create it.
			if err = s.client.MakeDir(part); err != nil {
				return err
			}
			// Change to the newly created directory.
			if err = s.client.ChangeDir(part); err != nil {
				return err
			}
		}
	}
	return s.client.ChangeDir(currentDir)
}

// 登录
func (s *Saver) login() error {
	if s.client != nil {
		return nil
	}
	var err error
	s.client, err = ftp.Dial(s.opt.Host, ftp.DialWithDisabledEPSV(s.opt.Epsv == EpsvActive))
	if err != nil {
		return fmt.Errorf("ftp连接失败：%s", err.Error())
	}
	if err = s.client.Login(s.opt.User, s.opt.Pwd); err != nil {
		return fmt.Errorf("ftp登录失败：%s", err.Error())
	}
	return nil
}
