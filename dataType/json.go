package dataType

import (
	"context"
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/schema"
)

type JSONRaw json.RawMessage

// Value return json value, implement driver.Valuer interface
func (j JSONRaw) Value() (driver.Value, error) {
	if len(j) == 0 {
		return nil, nil
	}
	return string(j), nil
}

// Scan value into Jsonb, implements sql.Scanner interface
func (j *JSONRaw) Scan(value interface{}) error {
	if value == nil {
		*j = JSONRaw("null")
		return nil
	}
	var bytes []byte
	if s, ok := value.(fmt.Stringer); ok {
		bytes = []byte(s.String())
	} else {
		switch v := value.(type) {
		case []byte:
			if len(v) > 0 {
				bytes = make([]byte, len(v))
				copy(bytes, v)
			}
		case string:
			bytes = []byte(v)
		default:
			return errors.New(fmt.Sprint("Failed to unmarshal JSONB value:", value))
		}
	}

	result := json.RawMessage(bytes)
	*j = JSONRaw(result)
	return nil
}

// GormDataType gorm common data type
func (JSONRaw) GormDataType() string {
	return "json"
}

// GormDBDataType gorm db data type
func (JSONRaw) GormDBDataType(db *gorm.DB, field *schema.Field) string {
	return JsonDbDataType(db, field)
}

func (js JSONRaw) GormValue(ctx context.Context, db *gorm.DB) clause.Expr {
	if len(js) == 0 {
		return gorm.Expr("NULL")
	}
	data, _ := js.MarshalJSON()
	return MapGormValue(string(data), db)
}

// MarshalJSON to output non base64 encoded []byte
func (js JSONRaw) MarshalJSON() ([]byte, error) {
	return json.RawMessage(js).MarshalJSON()
}

// UnmarshalJSON to deserialize []byte
func (js *JSONRaw) UnmarshalJSON(b []byte) error {
	result := json.RawMessage{}
	err := result.UnmarshalJSON(b)
	*js = JSONRaw(result)
	return err
}
