package cors

import (
	"context"
	"database/sql/driver"
	"encoding/json"

	"github.com/vmihailenco/msgpack/v5"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/schema"
	"helay.net/go/utils/v3/dataType"
)

// Config CORS配置结构体
type Config struct {
	Enabled          bool     `json:"enabled,omitempty" yaml:"enabled" ini:"enabled"`                               // 是否启用CORS
	AllowOrigins     []string `json:"allow_origins,omitempty" yaml:"allow_origins" ini:"allow_origins"`             // 允许的源
	AllowMethods     []string `json:"allow_methods,omitempty" yaml:"allow_methods" ini:"allow_methods"`             // 允许的HTTP方法
	AllowHeaders     []string `json:"allow_headers,omitempty" yaml:"allow_headers" ini:"allow_headers"`             // 允许的请求头
	AllowCredentials bool     `json:"allow_credentials,omitempty" yaml:"allow_credentials" ini:"allow_credentials"` // 是否允许发送 Cookie
	ExposeHeaders    []string `json:"expose_headers,omitempty" yaml:"expose_headers" ini:"expose_headers"`          // 暴露给前端的响应头
	MaxAge           int      `json:"max_age,omitempty" yaml:"max_age" ini:"max_age"`                               // 预检请求缓存时间(秒)
	Strict           bool     `json:"strict,omitempty" yaml:"strict" ini:"strict"`                                  // 严格模式 严格模式下，不允许的Origin阻值后不继续后续流程
}

func (c Config) Value() (driver.Value, error) {
	return dataType.DriverValueWithJson(c)
}
func (c *Config) Scan(val any) (err error) {
	return dataType.DriverScanWithJson(val, c)
}

func (c Config) GormDataType() string {
	return "json"
}

func (Config) GormDBDataType(db *gorm.DB, field *schema.Field) string {
	return dataType.JsonDbDataType(db, field)
}

func (c Config) GormValue(_ context.Context, db *gorm.DB) clause.Expr {
	byt, _ := json.Marshal(c)
	return dataType.MapGormValue(string(byt), db)
}

func (c Config) GobEncode() ([]byte, error) {
	return msgpack.Marshal(c)
}

func (c *Config) GobDecode(data []byte) error {
	return msgpack.Unmarshal(data, c)
}
