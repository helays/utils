package db

import (
	"database/sql/driver"
	"fmt"
	"github.com/helays/utils/config"
	"github.com/helays/utils/dataType"
	"github.com/helays/utils/logger/zaploger"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
	"net/url"
	"strings"
)

var (
	// SupportedDbType 支持的数据库类型
	SupportedDbType = []map[string]string{
		{"type": "mysql", "value": "mysql"},
		{"type": "pg", "value": "postgres"},
		{"type": "kafka", "value": "kafka"},
		{"type": "ftp", "value": "ftp"},
		{"type": "sftp", "value": "sftp"},
	}
	// FTPEpsv ftp模式
	FTPEpsv = []map[string]any{
		{"type": 0, "value": "被动模式"},
		{"type": 1, "value": "主动模式"},
	}
	// Authentication 认证方式
	Authentication = []map[string]string{
		{"type": "password", "value": "密码"},
		{"type": "public_key", "value": "密钥"},
	}
)

type Dbbase struct {
	DbIdentifier string `ini:"db_identifier" yaml:"db_identifier,omitempty" json:"db_identifier" gorm:"type:varchar(256);not null;uniqueIndex;comment:配置唯一标识"`
	DbType       string `ini:"db_type" yaml:"db_type" json:"db_type,omitempty" gorm:"type:varchar(32);not null;index;comment:数据库类型，mysql|pg"` // 数据库类型 mysql/pg

	// 这部分是公用的
	Host dataType.StringArray `ini:"host" yaml:"host" json:"host,omitempty" gorm:"not null;comment:连接信息"`
	User string               `ini:"user" yaml:"user" json:"user,omitempty" gorm:"type:varchar(256);not null;default:'';comment:数据库用户"`
	Pwd  string               `ini:"pwd" yaml:"pwd" json:"pwd,omitempty" gorm:"type:text;comment:数据库密码"`

	// 这部分是数据库独有
	Dbname        string          `ini:"dbname" yaml:"dbname" json:"dbname,omitempty" gorm:"type:varchar(128);not null;index;default:'';comment:默认连接的库"`
	Schema        string          `ini:"schema" yaml:"schema" json:"schema,omitempty" gorm:"type:varchar(128);not null;default:'';comment:数据库模式"`
	MaxIdleConns  int             `ini:"max_idle_conns" yaml:"max_idle_conns" json:"max_idle_conns,omitempty" gorm:"type:int;not null;default:2;comment:最大空闲连接数"`
	MaxOpenConns  int             `ini:"max_open_conns" yaml:"max_open_conns" json:"max_open_conns,omitempty" gorm:"type:int;not null;default:10;comment:最大连接数"`
	TablePrefix   string          `ini:"table_prefix" yaml:"table_prefix" json:"table_prefix,omitempty" gorm:"type:varchar(64);not null;default:'';comment:表前缀"`
	SingularTable int             `ini:"singular_table" yaml:"singular_table" json:"singular_table,omitempty" gorm:"type:int;not null;default:0;comment:是否启用单数表"` // 1 启用 0 不启用
	Logger        zaploger.Config `json:"logger" yaml:"logger" ini:"logger" gorm:"comment:日志配置"`
}

func (this Dbbase) Value() (driver.Value, error) {
	return dataType.DriverValueWithJson(this)
}

func (this *Dbbase) Scan(val interface{}) error {
	return dataType.DriverScanWithJson(val, this)
}

func (this Dbbase) GormDataType() string {
	return "json"
}

func (Dbbase) GormDBDataType(db *gorm.DB, field *schema.Field) string {
	return dataType.JsonDbDataType(db, field)
}

func (this Dbbase) Dsn() string {
	dsn := url.URL{
		User: url.UserPassword(this.User, this.Pwd),
		Host: strings.Join(this.Host, ","),
		Path: this.Dbname,
	}
	query := dsn.Query()
	switch this.DbType {
	case config.DbTypePostgres, config.DbTypePg:
		dsn.Scheme = "postgres"
		// 如果下面这里 设置成TimeZone ，有几率会出现时间异常
		query.Set("timezone", "Asia/Shanghai")
		if this.Schema != "" {
			query.Set("search_path", this.Schema)
		}
		dsn.RawQuery = query.Encode()
		return dsn.String()
	case config.DbTypeMysql:
		//dsn.Scheme = "mysql" // mysql 不需要这个
		// mysql 密码里面的特殊字符 不用序列化
		query.Set("charset", "utf8mb4")
		query.Set("parseTime", "True")
		query.Set("loc", "Local")
		return fmt.Sprintf("%s:%s@tcp(%s)/%s?%s", this.User, this.Pwd, dsn.Host, this.Dbname, query.Encode())
	}
	return ""
}

func (this *Dbbase) RemovePasswd() {
	this.Pwd = ""
}

func (this *Dbbase) SetInfo(args ...any) {
	if len(args) != 2 {
		return
	}
	switch args[0].(string) {
	case config.ClientInfoHost:
		this.Host = args[1].([]string)
	case config.ClientInfoUser:
		this.User = args[1].(string)
	case config.ClientInfoPasswd:
		this.Pwd = args[1].(string)
	}
}

// TableDefaultField 用于快速定义默认的表结构字段，包含id 创建时间 更新时间
type TableDefaultField struct {
	Id         int                 `json:"id,omitempty" gorm:"primaryKey;not null;autoIncrement;comment:行ID" form:"id"`
	CreateTime dataType.CustomTime `json:"create_time,omitempty" gorm:"autoCreateTime:true;index;not null;default:current_timestamp;comment:记录创建时间" form:"-"`
	UpdateTime dataType.CustomTime `json:"update_time,omitempty" gorm:"autoUpdateTime:true;index;comment:记录更新时间" form:"-"`
}

// TableDefaultTimeField 用于快速定义默认的表结构时间字段，这里不需要定义字段类型，因为会自动根据字段类型进行转换
type TableDefaultTimeField struct {
	CreateTime dataType.CustomTime `json:"create_time,omitempty" gorm:"autoCreateTime:true;index;not null;default:current_timestamp;comment:记录创建时间" form:"-"`
	UpdateTime dataType.CustomTime `json:"update_time,omitempty" gorm:"autoUpdateTime:true;index;comment:记录更新时间" form:"-"`
}

// TableDefaultUserField 用于快速定义默认的表结构用户字段，包含id 用户信息字段 创建时间 更新时间
type TableDefaultUserField struct {
	Id             int                 `json:"id,omitempty" gorm:"primaryKey;not null;autoIncrement;comment:行ID" form:"id"`
	CreateUserId   int                 `json:"create_user_id,omitempty" gorm:"not null;default:0;index;comment:创建人ID" form:"create_user_id"`
	CreateUserName string              `json:"create_user_name,omitempty" gorm:"not null;type:varchar(128);default:'';comment:创建人名称" form:"create_user_name"`
	CreateTime     dataType.CustomTime `json:"create_time,omitempty" gorm:"autoCreateTime:true;not null;index;default:current_timestamp;comment:记录创建时间" form:"-"`
	UpdateTime     dataType.CustomTime `json:"update_time,omitempty" gorm:"autoUpdateTime:true;index;comment:记录更新时间" form:"-"`
}
