package dataType

import (
	"database/sql/driver"
	"fmt"

	"github.com/helays/utils/v2/config"
	"github.com/helays/utils/v2/tools"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

type Byte byte

// noinspection all
func (b Byte) Value() (driver.Value, error) {
	// 注意，这里只能接受int64类型
	return int64(b), nil
}

// noinspection all
func (b *Byte) Scan(value any) error {
	if value == nil {
		*b = 0
		return nil
	}
	v, err := tools.Any2Int[byte](value)
	if err != nil {
		return fmt.Errorf("Byte.Scan: unknown type %T", value)
	}
	*b = Byte(v)
	return nil
}

// noinspection all
func (Byte) GormDBDataType(db *gorm.DB, field *schema.Field) string {
	switch db.Dialector.Name() {
	case config.DbTypeSqlite:
		return "integer"
	case config.DbTypeMysql:
		return "tinyint(1)"
	case config.DbTypePostgres:
		return "int2"
	case config.DbTypeSqlserver:
		return "tinyint"
	}
	return "int"
}

type Uint64 struct {
	uint64
}

func NewUint64(v uint64) Uint64 {
	return Uint64{uint64: v}
}

func (u *Uint64) GetValue() uint64 {
	return u.uint64
}

func (u *Uint64) SetValue(v uint64) {
	u.uint64 = v
}

func (u *Uint64) Equals(other Uint64) bool {
	return u.uint64 == other.uint64
}

func (u *Uint64) EqualsInt(other int) bool {
	return u.uint64 == uint64(other)
}

func (u *Uint64) EqualsUint64(other uint64) bool {
	return u.uint64 == other
}

// noinspection all
func (u Uint64) Value() (driver.Value, error) {
	return u.uint64, nil
}

// noinspection all
func (u *Uint64) Scan(value any) error {
	if value == nil {
		return nil
	}
	v, err := tools.Any2Int[uint64](value)
	if err != nil {
		return err
	}
	u.uint64 = v
	return nil
}

// noinspection all
func (u Uint64) GormDBDataType(db *gorm.DB, field *schema.Field) string {
	switch db.Dialector.Name() {
	case config.DbTypeSqlite:
		return "integer"
	case config.DbTypeMysql:
		return "BIGINT UNSIGNED"
	case config.DbTypePostgres:
		return "BIGINT"
	case config.DbTypeSqlserver:
		return "BIGINT"
	}
	return "bigint"
}
