package session

import (
	"bytes"
	"database/sql/driver"
	"encoding/gob"
	"errors"
	"time"

	"github.com/helays/utils/v2/dataType"
	"github.com/helays/utils/v2/tools"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

func NewSessionValue(val any) SessionValue {
	return SessionValue{val: val}
}

// noinspection all
type SessionValue struct {
	val any
}

// Value return blob value, implement driver.Valuer interface
// noinspection all
func (s SessionValue) Value() (driver.Value, error) {
	var buf bytes.Buffer
	err := gob.NewEncoder(&buf).Encode(s.val)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// noinspection all
func (s *SessionValue) Scan(val any) error {
	if val == nil {
		*s = SessionValue{}
		return nil
	}

	b, err := tools.Any2bytes(val)
	if err != nil {
		return err
	}
	return gob.NewDecoder(bytes.NewReader(b)).Decode(&s.val)
}

// GormDBDataType gorm db data type
// noinspection all
func (SessionValue) GormDBDataType(db *gorm.DB, field *schema.Field) string {
	return dataType.BlobDbDataType(db, field)
}

func (s SessionValue) GormDataType() string {
	return "blob"
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

// noinspection all
const SessionID = "session_id"

type Callback struct {
	BeforeRenew func(expire dataType.CustomTime, data any) error
	AfterRenew  func(expire dataType.CustomTime, data any) error
}

type Value struct {
	SessionID string        // 可自定义session id
	Field     string        // session 值字段
	Value     any           // session 值
	TTL       time.Duration // 有效期
}
