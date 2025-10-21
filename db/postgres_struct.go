package db

import (
	"database/sql/driver"

	"github.com/helays/utils/v2/dataType"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

type PostgresOpt struct {
	Passfile       string `ini:"passfile" yaml:"passfile" json:"passfile,omitempty"`                      //  密码文件路径
	ConnectTimeout string `ini:"connect_timeout" yaml:"connect_timeout" json:"connect_timeout,omitempty"` // 连接超时(秒)
	ClientEncoding string `ini:"client_encoding" yaml:"client_encoding" json:"client_encoding,omitempty"` // 客户端编码
	//连接行为参数
	Sslmode            string `ini:"sslmode" yaml:"sslmode" json:"sslmode,omitempty"`                                        //  SSL模式(disable|allow|prefer|require|verify-ca|verify-full)
	Sslcompression     string `ini:"sslcompression" yaml:"sslcompression" json:"sslcompression,omitempty"`                   // SSL压缩(0/1)
	Sslcert            string `ini:"sslcert" yaml:"sslcert" json:"sslcert,omitempty"`                                        // SSL客户端证书路径
	Sslkey             string `ini:"sslkey" yaml:"sslkey" json:"sslkey,omitempty"`                                           // SSL客户端密钥路径
	Sslrootcert        string `ini:"sslrootcert" yaml:"sslrootcert" json:"sslrootcert,omitempty"`                            // SSL根证书路径
	Sslcrl             string `ini:"sslcrl" yaml:"sslcrl" json:"sslcrl,omitempty"`                                           //  SSL证书撤销列表路径
	Sslpassword        string `ini:"sslpassword" yaml:"sslpassword" json:"sslpassword,omitempty"`                            // SSL密钥密码
	Service            string `ini:"service" yaml:"service" json:"service,omitempty"`                                        //  服务名(用于pg_service.conf)
	TargetSessionAttrs string `ini:"target_session_attrs" yaml:"target_session_attrs" json:"target_session_attrs,omitempty"` // 目标会话属性(read-write|primary)
	//应用行为参数
	ApplicationName         string `ini:"application_name" yaml:"application_name" json:"application_name,omitempty"`                            //  应用名称
	FallbackApplicationName string `ini:"fallback_application_name" yaml:"fallback_application_name" json:"fallback_application_name,omitempty"` // 备用应用名称
	Keepalives              string `ini:"keepalives" yaml:"keepalives" json:"keepalives,omitempty"`                                              // 是否启用TCP保持连接(1/0)
	KeepalivesIdle          string `ini:"keepalives_idle" yaml:"keepalives_idle" json:"keepalives_idle,omitempty"`                               //  TCP保持连接空闲时间(秒)
	KeepalivesInterval      string `ini:"keepalives_interval" yaml:"keepalives_interval" json:"keepalives_interval,omitempty"`                   //  TCP保持连接间隔(秒)
	KeepalivesCount         string `ini:"keepalives_count" yaml:"keepalives_count" json:"keepalives_count,omitempty"`                            // TCP保持连接探测次数
	TcpUserTimeout          string `ini:"tcp_user_timeout" yaml:"tcp_user_timeout" json:"tcp_user_timeout,omitempty"`                            // TCP用户超时(毫秒)
	// 性能参数
	StatementCacheMode   string `ini:"statement_cache_mode" yaml:"statement_cache_mode" json:"statement_cache_mode,omitempty"`       // 语句缓存模式
	StatementCacheSize   string `ini:"statement_cache_size" yaml:"statement_cache_size" json:"statement_cache_size,omitempty"`       //语句缓存大小
	PreferSimpleProtocol string `ini:"prefer_simple_protocol" yaml:"prefer_simple_protocol" json:"prefer_simple_protocol,omitempty"` //偏好简单协议(布尔)
}

func (d PostgresOpt) Value() (driver.Value, error) {
	return dataType.DriverValueWithJson(d)
}

func (d *PostgresOpt) Scan(val interface{}) error {
	return dataType.DriverScanWithJson(val, d)
}

func (d PostgresOpt) GormDataType() string {
	return "json"
}

func (PostgresOpt) GormDBDataType(db *gorm.DB, field *schema.Field) string {
	return dataType.JsonDbDataType(db, field)
}

func (d *PostgresOpt) dsn() []string {
	var builds []string
	if d.Passfile != "" {
		builds = append(builds, "passfile="+d.Passfile)
	}
	if d.ConnectTimeout != "" {
		builds = append(builds, "connect_timeout="+d.ConnectTimeout)
	}
	if d.ClientEncoding != "" {
		builds = append(builds, "client_encoding="+d.ClientEncoding)
	}
	if d.Sslmode != "" {
		builds = append(builds, "sslmode="+d.Sslmode)
	}
	if d.Sslcompression != "" {
		builds = append(builds, "sslcompression="+d.Sslcompression)
	}
	if d.Sslcert != "" {
		builds = append(builds, "sslcert="+d.Sslcert)
	}
	if d.Sslkey != "" {
		builds = append(builds, "sslkey="+d.Sslkey)
	}
	if d.Sslrootcert != "" {
		builds = append(builds, "sslrootcert="+d.Sslrootcert)
	}
	if d.Sslcrl != "" {
		builds = append(builds, "sslcrl="+d.Sslcrl)
	}
	if d.Sslpassword != "" {
		builds = append(builds, "sslpassword="+d.Sslpassword)
	}
	if d.Service != "" {
		builds = append(builds, "service="+d.Service)
	}
	if d.TargetSessionAttrs != "" {
		builds = append(builds, "target_session_attrs="+d.TargetSessionAttrs)
	}
	if d.ApplicationName != "" {
		builds = append(builds, "application_name="+d.ApplicationName)
	}
	if d.FallbackApplicationName != "" {
		builds = append(builds, "fallback_application_name="+d.FallbackApplicationName)
	}
	if d.Keepalives != "" {
		builds = append(builds, "keepalives="+d.Keepalives)
	}
	if d.KeepalivesIdle != "" {
		builds = append(builds, "keepalives_idle="+d.KeepalivesIdle)
	}
	if d.KeepalivesInterval != "" {
		builds = append(builds, "keepalives_interval="+d.KeepalivesInterval)
	}
	if d.KeepalivesCount != "" {
		builds = append(builds, "keepalives_count="+d.KeepalivesCount)
	}
	if d.TcpUserTimeout != "" {
		builds = append(builds, "tcp_user_timeout="+d.TcpUserTimeout)
	}
	if d.StatementCacheMode != "" {
		builds = append(builds, "statement_cache_mode="+d.StatementCacheMode)
	}
	if d.StatementCacheSize != "" {
		builds = append(builds, "statement_cache_size="+d.StatementCacheSize)
	}
	if d.PreferSimpleProtocol != "" {
		builds = append(builds, "prefer_simple_protocol="+d.PreferSimpleProtocol)
	}
	return builds
}
