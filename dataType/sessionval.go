package dataType

import (
	"bytes"
	"database/sql/driver"
	"encoding/gob"

	"github.com/helays/utils/v2/tools"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

func NewSessionValue(val any) SessionValue {
	return SessionValue{Val: val}
}

// noinspection all
type SessionValue struct {
	Val any
}

// Value return blob value, implement driver.Valuer interface
// noinspection all
func (s SessionValue) Value() (driver.Value, error) {
	var buf bytes.Buffer
	err := gob.NewEncoder(&buf).Encode(s)
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

	err = gob.NewDecoder(bytes.NewReader(b)).Decode(s)

	return err
}

// GormDBDataType gorm db data type
// noinspection all
func (SessionValue) GormDBDataType(db *gorm.DB, field *schema.Field) string {
	return BlobDbDataType(db, field)
}

func (s SessionValue) GormDataType() string {
	return "blob"
}
