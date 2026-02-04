package ftp

import (
	"database/sql/driver"
	"fmt"

	"gorm.io/gorm"
	"gorm.io/gorm/schema"
	"helay.net/go/utils/v3/config"
	"helay.net/go/utils/v3/dataType"
	"helay.net/go/utils/v3/net/checkIp"
)

// Config ftp 配置
// noinspection all
type Config struct {
	Host string `json:"host" yaml:"host" ini:"host"`
	User string `json:"user" yaml:"user" ini:"user"`
	Pwd  string `json:"pwd" yaml:"pwd" ini:"pwd"`
	Epsv Epsv   `ini:"epsv" yaml:"epsv" json:"epsv,omitempty"` // ftp连接模式
}

// Epsv ftp连接模式
// noinspection all
type Epsv int

// 0 被动模式 1 主动模式
// noinspection all
const (
	EpsvPassive Epsv = 0
	EpsvActive  Epsv = 1
)

// noinspection all
func (c Config) Value() (driver.Value, error) {
	return dataType.DriverValueWithJson(c)
}

// noinspection all
func (c *Config) Scan(val any) error {
	return dataType.DriverScanWithJson(val, c)
}

// noinspection all
func (c Config) GormDataType() string {
	return "json"
}

// noinspection all
func (Config) GormDBDataType(db *gorm.DB, field *schema.Field) string {
	return dataType.JsonDbDataType(db, field)
}

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
	if c.Epsv != EpsvPassive && c.Epsv != EpsvActive {
		return fmt.Errorf("无效的连接模式")
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
