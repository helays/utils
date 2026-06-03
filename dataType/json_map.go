package dataType

import (
	"bytes"
	"context"
	"database/sql/driver"
	"encoding/json"
	"errors"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/schema"
)

// JSONMap defined JSON data type, need to implements driver.Valuer, sql.Scanner interface
// noinspection all
type JSONMap struct {
	data map[string]any
}

// Value return json value, implement driver.Valuer interface
// noinspection all
func (m JSONMap) Value() (driver.Value, error) {
	return DriverValueWithJson(m.data)
}

// Scan scan value into Jsonb, implements sql.Scanner interface
// noinspection all
func (m *JSONMap) Scan(val any) error {
	if m == nil {
		return errors.New("JSONMap is nil")
	}
	return DriverScanWithJson(val, &m.data) // 这里暂时先用这个版本
}

// MarshalJSON to output non base64 encoded []byte
// noinspection all
func (m JSONMap) MarshalJSON() ([]byte, error) {
	if m.data == nil {
		return []byte("null"), nil
	}
	return json.Marshal(m.data)
}

// UnmarshalJSON to deserialize []byte
// 这个函数很重要，Scan的时候  序列化通用函数会用到这个
// noinspection all
func (m *JSONMap) UnmarshalJSON(b []byte) error {
	rd := bytes.NewReader(b)
	decoder := json.NewDecoder(rd)
	decoder.UseNumber()
	return decoder.Decode(&m.data)
}

// GormDataType gorm common data type
// noinspection all
func (m JSONMap) GormDataType() string {
	return "custom_json_map"
}

// GormDBDataType gorm db data type
// noinspection all
func (JSONMap) GormDBDataType(db *gorm.DB, field *schema.Field) string {
	return JsonDbDataType(db, field)
}

// noinspection all
func (m JSONMap) GormValue(_ context.Context, db *gorm.DB) clause.Expr {
	data, _ := m.MarshalJSON()
	return MapGormValue(string(data), db)
}

func (m *JSONMap) Add(k string, v any) {
	if m == nil {
		return
	}
	if m.data == nil {
		m.data = make(map[string]any)
	}
	m.data[k] = v
}

func (m *JSONMap) Get(k string) (any, bool) {
	if m == nil || m.data == nil {
		return nil, false
	}
	v, ok := m.data[k]
	return v, ok
}

func (m *JSONMap) Remove(k string) {
	if m == nil || m.data == nil {
		return
	}

	delete(m.data, k)
}
