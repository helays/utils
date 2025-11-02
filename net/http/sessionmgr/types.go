package sessionmgr

import (
	"bytes"
	"database/sql/driver"
	"encoding/gob"
	"errors"
	"fmt"

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
func (v SessionValue) Value() (driver.Value, error) {
	var buf bytes.Buffer
	err := gob.NewEncoder(&buf).Encode(v)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (v *SessionValue) Scan(val any) error {
	if val == nil {
		*v = SessionValue{}
		return nil
	}
	var ba []byte
	switch v := val.(type) {
	case []byte:
		ba = v
	default:
		return errors.New(fmt.Sprint("Failed to unmarshal JSONB value:", val))
	}
	return gob.NewDecoder(bytes.NewReader(ba)).Decode(v)
}

// GormDBDataType gorm db data type
func (SessionValue) GormDBDataType(db *gorm.DB, field *schema.Field) string {
	switch db.Dialector.Name() {
	case "sqlite":
		return "BLOB"
	case "mysql":
		return "BLOB"
	case "postgres":
		return "BYTEA"
	}
	return ""
}

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
