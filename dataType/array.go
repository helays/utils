package dataType

import (
	"context"
	"database/sql/driver"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/schema"
)

// AnyArray 泛型版本的数组
type AnyArray[T any] []T

// noinspection all
func (a AnyArray[T]) Value() (driver.Value, error) {
	return arrayValue(a)
}

// noinspection all
func (a *AnyArray[T]) Scan(val interface{}) error {
	return arrayScan(a, val)
}

// noinspection all
func (AnyArray[T]) GormDBDataType(db *gorm.DB, field *schema.Field) string {
	return JsonDbDataType(db, field)
}

// noinspection all
func (a AnyArray[T]) GormValue(ctx context.Context, db *gorm.DB) clause.Expr {
	return arrayGormValue(a, db)
}

type Array []any

// Value return json value, implement driver.Valuer interface
// noinspection all
func (m Array) Value() (driver.Value, error) {
	return arrayValue(m)
}

// Scan scan value into Jsonb, implements sql.Scanner interface
// noinspection all
func (m *Array) Scan(val interface{}) error {
	return arrayScan(m, val)
}

// GormDBDataType gorm db data type
// noinspection all
func (Array) GormDBDataType(db *gorm.DB, field *schema.Field) string {
	return JsonDbDataType(db, field)
}

// noinspection all
func (jm Array) GormValue(ctx context.Context, db *gorm.DB) clause.Expr {
	return arrayGormValue(jm, db)
}

type StringArray []string

// noinspection all
func (m StringArray) Value() (driver.Value, error) {
	b, err := arrayValue(m)
	return b, err
}

// noinspection all
func (m *StringArray) Scan(val interface{}) error {
	return arrayScan(m, val)
}

// noinspection all
func (StringArray) GormDBDataType(db *gorm.DB, field *schema.Field) string {
	return JsonDbDataType(db, field)
}

// noinspection all
func (jm StringArray) GormValue(ctx context.Context, db *gorm.DB) clause.Expr {
	return arrayGormValue(jm, db)
}
