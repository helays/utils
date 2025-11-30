package session

import (
	"database/sql/driver"
	"errors"

	"github.com/helays/utils/v2/dataType"
	"github.com/helays/utils/v2/tools"
	"github.com/vmihailenco/msgpack/v5"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

func NewSessionValue(val any) SessionValue {
	return SessionValue{val: val}
}

type SessionValue struct {
	val any
}

// Value return blob value, implement driver.Valuer interface
func (s SessionValue) Value() (driver.Value, error) {
	return msgpack.Marshal(s.val)
}

func (s *SessionValue) Scan(val any) error {
	if val == nil {
		*s = SessionValue{}
		return nil
	}
	b, err := tools.Any2bytes(val)
	if err != nil {
		return err
	}
	// 这里应该使用 msgpack.Unmarshal 而不是 gob
	return msgpack.Unmarshal(b, &s.val)
}

// GormDBDataType gorm db data type
func (SessionValue) GormDBDataType(db *gorm.DB, field *schema.Field) string {
	return dataType.BlobDbDataType(db, field)
}

//  这个需要移除 上级Session 已经实现了二进制序列化
//func (s SessionValue) GobEncode() ([]byte, error) {
//	return msgpack.Marshal(s.val)
//}
//
//func (s *SessionValue) GobDecode(data []byte) error {
//	return msgpack.Unmarshal(data, &s.val)
//}

type CookieCarrier string

func (c CookieCarrier) String() string {
	return string(c)
}

const (
	CookieCarrierCookie = "cookie"
	CookieCarrierHeader = "header"
)

var (
	ErrUnSupport  = errors.New("不支持的session载体")
	ErrNotFound   = errors.New("session不存在")
	ErrNotPointer = errors.New("session变量目标必须是指针")
)

type Engine string

func (e Engine) String() string {
	return string(e)
}

const (
	EngineRedis  Engine = "redis"
	EngineRdbms  Engine = "rdbms"
	EngineMemory Engine = "memory"
	EngineFile   Engine = "file"
)
