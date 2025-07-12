package sftp

import (
	"database/sql/driver"
	"fmt"
	"github.com/helays/utils/v2/close/vclose"
	"github.com/helays/utils/v2/config"
	"github.com/helays/utils/v2/dataType"
	"github.com/helays/utils/v2/net/checkIp"
	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
	"io"
	"path"
)

// Config sftp 配置
type Config struct {
	Host           string `json:"host" yaml:"host" ini:"host"` // 路径
	User           string `json:"user" yaml:"user" ini:"user"`
	Pwd            string `json:"pwd" yaml:"pwd" ini:"pwd"`                                  // 密码|密钥
	Authentication string `json:"authentication" yaml:"authentication" ini:"authentication"` // 认证方式 ，默认passwd,可选public_key

	sshClient  *ssh.Client
	sftpClient *sftp.Client
}

func (this *Config) RemovePasswd() {
	this.Pwd = ""
}

func (this *Config) Valid() error {
	if _, port, err := checkIp.ParseIPAndPort(this.Host); err != nil {
		return err
	} else if port < 1 {
		return fmt.Errorf("缺失端口号")
	}
	if this.User == "" {
		return fmt.Errorf("缺失账号")
	}
	if this.Pwd == "" {
		return fmt.Errorf("缺失密码")
	}
	if this.Authentication == "" {
		this.Authentication = "password"
	} else if this.Authentication != "password" && this.Authentication != "public_key" {
		return fmt.Errorf("无效的认证方式")
	}
	return nil
}

func (this *Config) SetInfo(args ...any) {
	if len(args) != 2 {
		return
	}
	switch args[0].(string) {
	case config.ClientInfoHost:
		this.Host = args[1].(string)
	case config.ClientInfoUser:
		this.User = args[1].(string)
	case config.ClientInfoPasswd:
		this.Pwd = args[1].(string)
	}
}

func (this Config) Value() (driver.Value, error) {
	return dataType.DriverValueWithJson(this)
}

func (this *Config) Scan(val interface{}) error {
	return dataType.DriverScanWithJson(val, this)
}

func (this Config) GormDataType() string {
	return "json"
}

func (Config) GormDBDataType(db *gorm.DB, field *schema.Field) string {
	return dataType.JsonDbDataType(db, field)
}

// Read 读取文件
func (this *Config) Read(p string) (io.ReadCloser, error) {
	if err := this.LoginSftp(); err != nil {
		return nil, err
	}
	filePath, err := SetPath(this.sftpClient, p)
	if err != nil {
		return nil, err
	}
	file, err := this.sftpClient.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("打开文件%s失败：%s", p, err.Error())
	}
	return file, nil
}

// Write 写入文件
func (this *Config) Write(p string, src io.Reader, existIgnores ...bool) (int64, error) {
	if err := this.LoginSftp(); err != nil {
		return 0, err
	}
	filePath, err := SetPath(this.sftpClient, p)
	if err != nil {
		return 0, err
	}
	// 判断是否需要覆盖写入
	if len(existIgnores) > 0 && existIgnores[0] {
		if ok, _err := Exist(this.sftpClient, filePath); ok {
			return 0, nil
		} else if _err != nil {
			return 0, _err
		}
	}

	dir := path.Dir(filePath)
	// 首先判断这个路径是否存在，然后创建
	if err = Mkdir(this.sftpClient, dir); err != nil {
		return 0, err
	}
	// 文件夹存在后，就开始创建文件
	file, err := this.sftpClient.Create(filePath)
	if err != nil {
		return 0, fmt.Errorf("创建文件%s失败：%s", filePath, err.Error())
	}
	defer vclose.Close(file)
	var written int64
	if written, err = io.Copy(file, src); err != nil {
		return written, fmt.Errorf("写入文件%s失败：%s", filePath, err.Error())
	}
	return written, nil
}

func (this *Config) ListFiles(dirPath string) ([]string, error) {
	if err := this.LoginSftp(); err != nil {
		return nil, err
	}
	filePath, err := SetPath(this.sftpClient, dirPath)
	if err != nil {
		return nil, err
	}
	entries, err := this.sftpClient.ReadDir(filePath)
	if err != nil {
		return nil, fmt.Errorf("获取目录%s失败：%s", filePath, err.Error())
	}
	var files []string
	for _, entry := range entries {
		if !entry.IsDir() {
			files = append(files, entry.Name())
		}
	}
	return files, nil
}

func (this *Config) Delete(p string) error {
	if err := this.LoginSftp(); err != nil {
		return err
	}
	filePath, err := SetPath(this.sftpClient, p)
	if err != nil {
		return err
	}
	if err = this.sftpClient.Remove(filePath); err != nil {
		return fmt.Errorf("删除文件%s失败：%s", filePath, err.Error())
	}
	return nil
}

func (this *Config) DeleteAll(p string) error {
	if err := this.LoginSftp(); err != nil {
		return err
	}
	filePath, err := SetPath(this.sftpClient, p)
	if err != nil {
		return err
	}
	if err = this.sftpClient.RemoveAll(filePath); err != nil {
		return fmt.Errorf("删除文件%s失败：%s", filePath, err.Error())
	}
	return nil
}

// LoginSsh 登录 ssh
func (this *Config) LoginSsh() error {
	if this.sshClient != nil {
		return nil
	}
	// 首先连接 ssh client
	clientConfig := &ssh.ClientConfig{
		User:            this.User,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}
	var auth ssh.AuthMethod
	if this.Authentication == "password" {
		auth = ssh.Password(this.Pwd)
	} else {
		var signer ssh.Signer
		signer, err := ssh.ParsePrivateKey([]byte(this.Pwd))
		if err != nil {
			return fmt.Errorf("ssh密钥解析失败：%s", err.Error())
		}
		auth = ssh.PublicKeys(signer)
	}
	clientConfig.Auth = []ssh.AuthMethod{auth}
	var err error
	this.sshClient, err = ssh.Dial("tcp", this.Host, clientConfig)
	if err != nil {
		return fmt.Errorf("ssh连接失败：%s", err.Error())
	}
	return nil
}

// LoginSftp ssh登录
// @return , error
func (this *Config) LoginSftp() error {
	if err := this.LoginSsh(); err != nil {
		return err
	}
	if this.sftpClient == nil {
		var err error
		this.sftpClient, err = sftp.NewClient(this.sshClient)
		if err != nil {
			return fmt.Errorf("sftp连接失败：%s", err.Error())
		}
	}
	return nil
}

// CloseSsh 关闭 ssh
func (this *Config) CloseSsh() {
	vclose.Close(this.sshClient)
	this.sshClient = nil
}

// CloseSftp 关闭 sftp
func (this *Config) CloseSftp() {
	vclose.Close(this.sftpClient)
	this.sftpClient = nil
}

// Close 关闭 ssh 和 sftp
func (this *Config) Close() {
	this.CloseSftp()
	this.CloseSsh()
}
