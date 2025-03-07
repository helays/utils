package sftp

import (
	"database/sql/driver"
	"fmt"
	"github.com/helays/utils/close/vclose"
	"github.com/helays/utils/dataType"
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
