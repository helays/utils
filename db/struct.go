package db

import (
	"database/sql/driver"
	"time"

	"github.com/helays/utils/v2/config"
	"github.com/helays/utils/v2/dataType"
	"github.com/helays/utils/v2/logger/zaploger"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

var (
	// SupportedDbType 支持的数据库类型
	SupportedDbType = []map[string]string{
		{"type": config.DbTypeMysql, "value": config.DbTypeMysql},
		{"type": config.DbTypePg, "value": config.DbTypePostgres},
		{"type": config.DbClientEs, "value": config.DbClientEs},
		{"type": config.QueueTypeKafka, "value": config.QueueTypeKafka},
		{"type": config.FileTypeFtp, "value": config.FileTypeFtp},
		{"type": config.FileTypeSftp, "value": config.FileTypeSftp},
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
	Dbname          string          `ini:"dbname" yaml:"dbname" json:"dbname,omitempty" gorm:"type:varchar(128);not null;index;default:'';comment:默认连接的库"`
	Schema          string          `ini:"schema" yaml:"schema" json:"schema,omitempty" gorm:"type:varchar(128);not null;default:'';comment:数据库模式"`
	MaxIdleConns    int             `ini:"max_idle_conns" yaml:"max_idle_conns" json:"max_idle_conns,omitempty" gorm:"type:int;not null;default:2;comment:最大空闲连接数"`
	MaxOpenConns    int             `ini:"max_open_conns" yaml:"max_open_conns" json:"max_open_conns,omitempty" gorm:"type:int;not null;default:10;comment:最大连接数"`
	MaxConnLifetime time.Duration   `ini:"max_conn_lifetime" yaml:"max_conn_lifetime" json:"max_conn_lifetime,omitempty" gorm:"type:int;not null;default:0;comment:连接最大存活时间"`    // DB服务器wait_timeout的80-90%
	MaxConnIdleTime time.Duration   `ini:"max_conn_idle_time" yaml:"max_conn_idle_time" json:"max_conn_idle_time,omitempty" gorm:"type:int;not null;default:0;comment:连接最大空闲时间"` // 平均请求间隔 × 3
	TablePrefix     string          `ini:"table_prefix" yaml:"table_prefix" json:"table_prefix,omitempty" gorm:"type:varchar(64);not null;default:'';comment:表前缀"`
	SingularTable   int             `ini:"singular_table" yaml:"singular_table" json:"singular_table,omitempty" gorm:"type:int;not null;default:0;comment:是否启用单数表"` // 1 启用 0 不启用
	PostgresOpt     *PostgresOpt    `json:"postgres_opt" yaml:"postgres_opt" json:"postgres_opt,omitempty" gorm:"comment:Postgres专属配置"`
	Timeout         int             `ini:"timeout" yaml:"timeout" json:"timeout,omitempty" gorm:"type:int;not null;default:0;comment:sqlite 超时，单位毫秒"` // sqlite可用
	Logger          zaploger.Config `json:"logger" yaml:"logger" ini:"logger" gorm:"comment:日志配置"`
}

func (this Dbbase) Value() (driver.Value, error) {
	return dataType.DriverValueWithJson(this)
}

func (this *Dbbase) Scan(val interface{}) error {
	return dataType.DriverScanWithJson(val, this)
}

func (this Dbbase) GormDataType() string {
	return "dbbase"
}

func (Dbbase) GormDBDataType(db *gorm.DB, field *schema.Field) string {
	return dataType.JsonDbDataType(db, field)
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

type TableDefaultCreateField struct {
	Id         int                 `json:"id,omitempty" gorm:"primaryKey;not null;autoIncrement;comment:行ID" form:"id"`
	CreateTime dataType.CustomTime `json:"create_time,omitempty" gorm:"autoCreateTime:true;index;not null;default:current_timestamp;comment:记录创建时间" form:"-"`
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

type TableBaseModelAutoIncrement struct {
	Id int64 `json:"id,omitempty" gorm:"primaryKey;autoIncrement;comment:行ID" form:"id"`

	CreateTime dataType.CustomTime `json:"create_time,omitempty" gorm:"autoCreateTime:true;index;not null;default:current_timestamp;comment:记录创建时间" form:"-"`
	CreateBy   int64               `json:"create_by,omitempty" gorm:"comment:创建人ID" form:"create_by"`
	UpdateTime dataType.CustomTime `json:"update_time,omitempty" gorm:"autoUpdateTime:true;comment:记录更新时间" form:"-"`
	UpdateBy   int64               `json:"update_by,omitempty" gorm:"comment:更新人ID" form:"update_by"`
}

type TableBaseModelFull struct {
	Id int64 `json:"id,omitempty" gorm:"primaryKey;autoIncrement:false;comment:行ID" form:"id"`

	CreateTime dataType.CustomTime `json:"create_time,omitempty" gorm:"autoCreateTime:true;index;not null;default:current_timestamp;comment:记录创建时间" form:"-"`
	CreateBy   int64               `json:"create_by,omitempty" gorm:"comment:创建人ID" form:"create_by"`
	UpdateTime dataType.CustomTime `json:"update_time,omitempty" gorm:"autoUpdateTime:true;comment:记录更新时间" form:"-"`
	UpdateBy   int64               `json:"update_by,omitempty" gorm:"comment:更新人ID" form:"update_by"`
}

type SoftDeleteModel struct {
	IsDeleted   dataType.Bool       `json:"is_deleted,omitempty" gorm:"not null;index;default:0;comment:软删除标记 0 正常 1 删除" form:"is_deleted"`
	DeletedBy   int64               `json:"deleted_by,omitempty" gorm:"comment:删除人ID" form:"deleted_by"`
	DeletedTime dataType.CustomTime `json:"deleted_time,omitempty" gorm:"index;comment:删除时间" form:"-"`
}
