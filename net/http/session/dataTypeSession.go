package session

import (
	"bytes"
	"database/sql/driver"
	"encoding/gob"
	"errors"
	"fmt"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

//
// ━━━━━━神兽出没━━━━━━
// 　　 ┏┓     ┏┓
// 　　┏┛┻━━━━━┛┻┓
// 　　┃　　　　　 ┃
// 　　┃　　━　　　┃
// 　　┃　┳┛　┗┳  ┃
// 　　┃　　　　　 ┃
// 　　┃　　┻　　　┃
// 　　┃　　　　　 ┃
// 　　┗━┓　　　┏━┛　Code is far away from bug with the animal protecting
// 　　　 ┃　　　┃    神兽保佑,代码无bug
// 　　　　┃　　　┃
// 　　　　┃　　　┗━━━┓
// 　　　　┃　　　　　　┣┓
// 　　　　┃　　　　　　┏┛
// 　　　　┗┓┓┏━┳┓┏┛
// 　　　　 ┃┫┫ ┃┫┫
// 　　　　 ┗┻┛ ┗┻┛
//
// ━━━━━━感觉萌萌哒━━━━━━
//
//
// User helay
// Date: 2024/12/8 17:53
//

type SessionValue struct {
	Val any
}

// Value return blob value, implement driver.Valuer interface
func (this SessionValue) Value() (driver.Value, error) {
	var buf bytes.Buffer
	err := gob.NewEncoder(&buf).Encode(this)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (this *SessionValue) Scan(val any) error {
	if val == nil {
		*this = SessionValue{}
		return nil
	}
	var ba []byte
	switch v := val.(type) {
	case []byte:
		ba = v
	default:
		return errors.New(fmt.Sprint("Failed to unmarshal JSONB value:", val))
	}
	return gob.NewDecoder(bytes.NewReader(ba)).Decode(this)
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
