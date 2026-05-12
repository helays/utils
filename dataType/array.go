package dataType

import (
	"context"
	"database/sql/driver"
	"fmt"

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

// ToSlice 转换为普通切片
func (a AnyArray[T]) ToSlice() []T {
	return []T(a)
}

// NewAnyArrayFromSlice 从普通切片创建
func NewAnyArrayFromSlice[T any](s []T) AnyArray[T] {
	return AnyArray[T](s)
}

// IsEmpty 判空
func (a AnyArray[T]) IsEmpty() bool {
	return len(a) == 0
}

// Copy 深拷贝
func (a AnyArray[T]) Copy() AnyArray[T] {
	tmp := make([]T, len(a))
	copy(tmp, a)
	return tmp
}

// Append 追加元素
func (a *AnyArray[T]) Append(elems ...T) {
	*a = append(*a, elems...)
}

// Get 获取元素（带边界检查）
func (a AnyArray[T]) Get(index int) (T, error) {
	if index < 0 || index >= len(a) {
		var zero T
		return zero, fmt.Errorf("index out of range: %d", index)
	}
	return a[index], nil
}

// Set 设置元素
func (a *AnyArray[T]) Set(index int, value T) error {
	if index < 0 || index >= len(*a) {
		return fmt.Errorf("index out of range: %d", index)
	}
	(*a)[index] = value
	return nil
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
