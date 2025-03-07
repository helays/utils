package ftp

import (
	"database/sql/driver"
	"fmt"
	"github.com/helays/utils/close/ftpClose"
	"github.com/helays/utils/dataType"
	"github.com/helays/utils/dataType/customWriter"
	"github.com/jlaffaye/ftp"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
	"io"
	"path"
)

// Config ftp 配置
type Config struct {
	Host string `json:"host" yaml:"host" ini:"host"` // ftp地址:端口
	User string `json:"user" yaml:"user" ini:"user"`
	Pwd  string `json:"pwd" yaml:"pwd" ini:"pwd"` // 密码
	// 这部分是ftp的
	Epsv int `ini:"epsv" yaml:"epsv" json:"epsv,omitempty" gorm:"type:int;not null;default:0;comment:连接模式"` // ftp 连接模式，0 被动模式 1 主动模式

	client *ftp.ServerConn
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

func (this *Config) RemovePasswd() {
	this.Pwd = ""
}

// Write 写入文件
func (this *Config) Write(p string, src io.Reader, existIgnores ...bool) (int64, error) {
	if err := this.Login(); err != nil {
		return 0, err
	}
	filePath, err := SetPath(this.client, p)
	if err != nil {
		return 0, err
	}
	// 判断是否需要覆盖写入
	if len(existIgnores) > 0 && existIgnores[0] {
		if ok, _err := Exist(this.client, filePath); ok {
			return 0, nil
		} else if _err != nil {
			return 0, _err
		}
	}
	dir := path.Dir(filePath)
	// 首先判断这个路径是否存在，然后创建
	if err = Mkdir(this.client, dir); err != nil {
		return 0, err
	}
	counter := &customWriter.SizeCounter{}
	teeReader := io.TeeReader(src, counter)
	if err = this.client.Stor(filePath, teeReader); err != nil {
		return counter.TotalSize, fmt.Errorf("写入文件%s失败：%s", filePath, err.Error())
	}
	return counter.TotalSize, nil
}

func (this *Config) Read(p string) (io.ReadCloser, error) {
	if err := this.Login(); err != nil {
		return nil, err
	}
	filePath, err := SetPath(this.client, p)
	if err != nil {
		return nil, err
	}
	return this.client.Retr(filePath)
}

// Login ftp登录
func (this *Config) Login() error {
	if this.client != nil {
		return nil
	}
	var err error
	this.client, err = ftp.Dial(this.Host, ftp.DialWithDisabledEPSV(this.Epsv == 1))
	if err != nil {
		return fmt.Errorf("ftp连接失败：%s", err.Error())
	}
	if err = this.client.Login(this.User, this.Pwd); err != nil {
		return fmt.Errorf("ftp登录失败：%s", err.Error())
	}
	return nil
}

func (this *Config) Close() {
	if this.client == nil {
		return
	}
	ftpClose.CloseFtpClient(this.client)
	this.client = nil
}
