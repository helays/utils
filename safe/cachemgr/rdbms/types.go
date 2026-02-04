package rdbms

import (
	"github.com/helays/utils/v2/dataType"
	"github.com/helays/utils/v2/db"
)

// CacheFast 采用双hash模式，缓存数据结构
// 效率足够高，但是在大规模数据下 可能存在hash冲突。
type CacheFast struct {
	InstanceHash dataType.Uint64       `json:"instance_hash" gorm:"primaryKey;autoIncrement:false;index;comment:缓存实例hash"`
	KeyHash      dataType.Uint64       `json:"key_hash" gorm:"primaryKey;autoIncrement:false;comment:缓存key hash"`
	InstanceID   string                `json:"instance_id" gorm:"type:varchar(32);not null;comment:缓存实例标识"`
	CacheKey     string                `json:"key" gorm:"type:varchar(512);not null;comment:缓存key"`
	Value        dataType.SessionValue `json:"value" gorm:"comment:缓存数据"`
	ExpiresTime  *dataType.CustomTime  `json:"expires_time" gorm:"index;comment:过期时间"`
	db.TableDefaultTimeField
}

func (CacheFast) TableName() string {
	return "cache_fast"
}

// CacheSafe 采用单hash模式，缓存数据结构
// 缓存实例标识，数量很少，采用 xxhash 足够安全
// 缓存key，数量很多，原样保存
type CacheSafe struct {
	InstanceHash dataType.Uint64       `json:"instance_hash" gorm:"primaryKey;autoIncrement:false;index;comment:缓存实例hash"`
	CacheKey     string                `json:"key" gorm:"primaryKey;type:varchar(512);not null;comment:缓存key"`
	InstanceID   string                `json:"instance_id" gorm:"type:varchar(32);not null;comment:缓存实例标识"`
	Value        dataType.SessionValue `json:"value" gorm:"comment:缓存数据"`
	ExpiresTime  *dataType.CustomTime  `json:"expires_time" gorm:"index;comment:过期时间"`
	db.TableDefaultTimeField
}

func (CacheSafe) TableName() string {
	return "cache_safe"
}
