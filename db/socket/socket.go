package socket

import (
	"database/sql/driver"
	"fmt"
	"github.com/helays/utils/config"
	"github.com/helays/utils/dataType"
	"github.com/helays/utils/net/checkIp"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

type Config struct {
	Protocol string `json:"protocol" yaml:"protocol" ini:"protocol"`
	Addr     string `json:"addr" yaml:"addr" ini:"addr"`
	Timeout  int    `json:"timeout" yaml:"timeout" ini:"timeout"` // 超时配置
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

var protocolMap = map[string]bool{
	config.ProtocolTCP: true,
	config.ProtocolUDP: true,
}

func (this *Config) Valid() error {
	if _, ok := protocolMap[this.Protocol]; !ok {
		return config.ErrProtocolInvalid
	}
	// 验证ip:port是否有效
	if _, port, err := checkIp.ParseIPAndPort(this.Addr); err != nil {
		return err
	} else if port < 1 {
		return fmt.Errorf("缺失端口号")
	}
	return nil
}
