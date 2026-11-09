package dataType

import (
	"bytes"
	"context"
	"database/sql/driver"
	"encoding/gob"
	"encoding/json"
	"errors"
	"fmt"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/schema"
)

// AnyArray 泛型版本的数组
// data 为 nil 表示未设置，为 []T{} 表示显式空数组
type AnyArray[T any] struct {
	data []T
}

// NewAnyArray 从值列表创建，始终是非 nil 切片；无参时得到「未设置」
func NewAnyArray[T any](values ...T) AnyArray[T] {
	return AnyArray[T]{data: values}
}

// NewAnyArrayFromSlice 从普通切片创建，共享底层数组，不复制
func NewAnyArrayFromSlice[T any](s []T) AnyArray[T] {
	return AnyArray[T]{data: s}
}

// ToSlice 返回底层 data 视图（共享底层数组，零拷贝）：
// 修改返回的切片会影响 AnyArray 本身，需要独立副本请使用 Copy。
func (a AnyArray[T]) ToSlice() []T {
	return a.data
}

// Len 返回元素个数
func (a AnyArray[T]) Len() int {
	return len(a.data)
}

// IsNil 判断是否为「未设置」状态（data 为 nil），显式空数组返回 false
func (a AnyArray[T]) IsNil() bool {
	return a.data == nil
}

// IsEmpty 判空，nil 与空数组都返回 true
func (a AnyArray[T]) IsEmpty() bool {
	return len(a.data) == 0
}

// Copy 复制底层数组（元素为浅拷贝）
func (a AnyArray[T]) Copy() AnyArray[T] {
	tmp := make([]T, len(a.data))
	copy(tmp, a.data)
	return AnyArray[T]{data: tmp}
}

// Append 追加元素
func (a *AnyArray[T]) Append(elems ...T) {
	if a == nil {
		return
	}
	a.data = append(a.data, elems...)
}

// Get 获取元素（带边界检查）
func (a AnyArray[T]) Get(index int) (T, error) {
	if index < 0 || index >= len(a.data) {
		var zero T
		return zero, fmt.Errorf("index out of range: %d", index)
	}
	return a.data[index], nil
}

// Set 设置元素
func (a *AnyArray[T]) Set(index int, value T) error {
	if a == nil {
		return errors.New("AnyArray is nil")
	}
	if index < 0 || index >= len(a.data) {
		return fmt.Errorf("index out of range: %d", index)
	}
	a.data[index] = value
	return nil
}

// noinspection all
func (a AnyArray[T]) Value() (driver.Value, error) {
	return arrayValue(a.data)
}

// noinspection all
func (a *AnyArray[T]) Scan(val any) error {
	if a == nil {
		return errors.New("AnyArray is nil")
	}
	return arrayScan(&a.data, val)
}

// noinspection all
func (AnyArray[T]) GormDBDataType(db *gorm.DB, field *schema.Field) string {
	return JsonDbDataType(db, field)
}

// noinspection all
// GormValue 注意传入底层切片：marshalSlice 内部使用 reflect.ValueOf(v).Len()，
// 直接传 struct 会 panic。
func (a AnyArray[T]) GormValue(ctx context.Context, db *gorm.DB) clause.Expr {
	return arrayGormValue(a.data, db)
}

// MarshalJSON 保持与切片一致的 JSON 表现：nil 输出 null，空数组输出 []
// noinspection all
func (a AnyArray[T]) MarshalJSON() ([]byte, error) {
	if a.data == nil {
		return []byte("null"), nil
	}
	return json.Marshal(a.data)
}

// UnmarshalJSON 把 JSON 数组解码到底层 data
// noinspection all
func (a *AnyArray[T]) UnmarshalJSON(b []byte) error {
	if a == nil {
		return errors.New("AnyArray is nil")
	}
	return json.Unmarshal(b, &a.data)
}

// gobAnyArray 用于 Gob 序列化的镜像结构，字段必须导出
type gobAnyArray[T any] struct {
	Data []T
}

// GobEncode 实现 gob.GobEncoder，避免未导出字段被 gob 忽略而静默丢数据
// noinspection all
func (a AnyArray[T]) GobEncode() ([]byte, error) {
	var buf bytes.Buffer
	if err := gob.NewEncoder(&buf).Encode(gobAnyArray[T]{Data: a.data}); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// GobDecode 实现 gob.GobDecoder
// noinspection all
func (a *AnyArray[T]) GobDecode(data []byte) error {
	if a == nil {
		return errors.New("AnyArray is nil")
	}
	var tmp gobAnyArray[T]
	if err := gob.NewDecoder(bytes.NewReader(data)).Decode(&tmp); err != nil {
		return err
	}
	a.data = tmp.Data
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
