package sftp

import (
	"database/sql/driver"
	"fmt"

	"github.com/helays/utils/v2/config"
	"github.com/helays/utils/v2/dataType"
	"github.com/helays/utils/v2/net/checkIp"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

type Config struct {
	Host           string         `json:"host" yaml:"host" ini:"host"`
	User           string         `json:"user" yaml:"user" ini:"user"`
	Pwd            string         `json:"pwd" yaml:"pwd" ini:"pwd"`
	Authentication Authentication `json:"authentication" yaml:"authentication" ini:"authentication"`
}

type Authentication string

const (
	Password  Authentication = "password"
	PublicKey Authentication = "public_key"
)

// noinspection all
func (c *Config) RemovePasswd() {
	c.Pwd = ""
}

// noinspection all
func (c *Config) Valid() error {
	if _, port, err := checkIp.ParseIPAndPort(c.Host); err != nil {
		return err
	} else if port < 1 {
		return fmt.Errorf("缺失端口号")
	}
	if c.User == "" {
		return fmt.Errorf("缺失账号")
	}
	if c.Pwd == "" {
		return fmt.Errorf("缺失密码")
	}
	if c.Authentication == "" {
		c.Authentication = Password
	} else if c.Authentication != Password && c.Authentication != PublicKey {
		return fmt.Errorf("无效的认证方式")
	}
	return nil
}

// noinspection all
func (c *Config) SetInfo(args ...any) {
	if len(args) != 2 {
		return
	}
	switch args[0].(string) {
	case config.ClientInfoHost:
		c.Host = args[1].(string)
	case config.ClientInfoUser:
		c.User = args[1].(string)
	case config.ClientInfoPasswd:
		c.Pwd = args[1].(string)
	}
}

// noinspection all
func (c Config) Value() (driver.Value, error) {
	return dataType.DriverValueWithJson(c)
}

// noinspection all
func (c *Config) Scan(val interface{}) error {
	return dataType.DriverScanWithJson(val, c)
}

func (c Config) GormDataType() string {
	return "json"
}

func (Config) GormDBDataType(db *gorm.DB, field *schema.Field) string {
	return dataType.JsonDbDataType(db, field)
}
