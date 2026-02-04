package sftp

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
	"helay.net/go/utils/v3/close/vclose"
)

type Saver struct {
	opt        *Config
	sshClient  *ssh.Client
	sftpClient *sftp.Client
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

func (s *Saver) Write(p string, src io.Reader, existIgnores ...bool) (int64, error) {
	path, err := s.setPath(p)
	if err != nil {
		return 0, err
	}
	if len(existIgnores) > 0 && existIgnores[0] {
		if ok, err := s.exist(path); err != nil {
			return 0, err
		} else if ok {
			return 0, nil
		}
	}
	dir := filepath.Dir(path)
	if err = s.mkdir(dir); err != nil {
		return 0, err
	}
	file, err := s.sftpClient.Create(path)
	defer vclose.Close(file)
	if err != nil {
		return 0, fmt.Errorf("创建文件%s失败：%s", path, err.Error())
	}
	return io.Copy(file, src)
}

func (s *Saver) Read(p string) (io.ReadCloser, error) {
	path, err := s.setPath(p)
	if err != nil {
		return nil, err
	}
	file, err := s.sftpClient.Open(path)
	if err != nil {
		return nil, fmt.Errorf("打开文件%s失败：%s", path, err.Error())
	}
	return file, nil
}

func (s *Saver) ListFiles(dirPath string) ([]string, error) {
	path, err := s.setPath(dirPath)
	if err != nil {
		return nil, err
	}
	entries, err := s.sftpClient.ReadDir(path)
	if err != nil {
		return nil, fmt.Errorf("获取目录%s失败：%s", path, err.Error())
	}
	var files []string
	for _, entry := range entries {
		if !entry.IsDir() {
			files = append(files, entry.Name())
		}
	}
	return files, nil
}

func (s *Saver) Delete(p string) error {
	path, err := s.setPath(p)
	if err != nil {
		return err
	}
	return s.sftpClient.Remove(path)
}

func (s *Saver) DeleteAll(p string) error {
	path, err := s.setPath(p)
	if err != nil {
		return err
	}
	return s.sftpClient.RemoveAll(path)
}

func (s *Saver) Close() error {
	vclose.Close(s.sftpClient)
	s.sftpClient = nil
	vclose.Close(s.sshClient)
	s.sshClient = nil
	return nil
}

func (s *Saver) login() error {
	if err := s.loginSSH(); err != nil {
		return err
	}
	if s.sftpClient == nil {
		var err error
		s.sftpClient, err = sftp.NewClient(s.sshClient)
		if err != nil {
			return fmt.Errorf("sftp连接失败：%s", err.Error())
		}
	}
	return nil
}

func (s *Saver) loginSSH() error {
	if s.sshClient != nil {
		return nil
	}
	cfg := &ssh.ClientConfig{
		User:            s.opt.User,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}
	var auth ssh.AuthMethod
	if s.opt.Authentication == Password {
		auth = ssh.Password(s.opt.Pwd)
	} else {
		signer, err := ssh.ParsePrivateKey([]byte(s.opt.Pwd))
		if err != nil {
			return fmt.Errorf("ssh密钥解析失败：%s", err.Error())
		}
		auth = ssh.PublicKeys(signer)
	}
	cfg.Auth = []ssh.AuthMethod{auth}
	var err error
	s.sshClient, err = ssh.Dial("tcp", s.opt.Host, cfg)
	if err != nil {
		return fmt.Errorf("ssh连接失败：%s", err.Error())
	}
	return nil
}

// setPath 设置当前 文件全路径
// p 如果是绝对路径，那么直接返回p
// p 如果是相对路径，会跟上当前目录
func (s *Saver) setPath(p string) (string, error) {
	current, err := s.sftpClient.Getwd()
	if err != nil {
		return "", err
	}
	return s.sftpClient.Join(current, p), nil
}

func (s *Saver) exist(p string) (bool, error) {
	if _, err := s.sftpClient.Stat(p); err == nil {
		return true, nil
	} else if !os.IsNotExist(err) {
		return false, err
	}
	return false, nil
}

func (s *Saver) mkdir(p string) error {
	if ok, err := s.exist(p); err != nil {
		return err
	} else if ok {
		return nil
	}
	if err := s.sftpClient.MkdirAll(p); err != nil {
		return fmt.Errorf("创建文件夹%s失败：%s", p, err.Error())
	}
	return nil
}
