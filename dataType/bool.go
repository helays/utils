package dataType

import (
	"database/sql/driver"

	"github.com/helays/utils/v2/config"
	"github.com/helays/utils/v2/tools"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

// Bool 注意当使用这个类型时，在定义模型时，默认值需要带上括号。不然pg数据库会报错。
type Bool struct {
	bool
}

// noinspection all
func (b Bool) Value() (driver.Value, error) {
	if b.bool {
		return int64(1), nil
	}
	return int64(0), nil
}

// noinspection all
func (b *Bool) Scan(value any) error {
	if value == nil {
		b.bool = false
		return nil
	}
	ok, e := tools.Any2bool(value)
	if e != nil {
		return e
	}
	b.bool = ok
	return nil
}

// noinspection all
func (b Bool) GormDataType() string {
	return "int"
}

// noinspection all
func (Bool) GormDBDataType(db *gorm.DB, _ *schema.Field) string {
	switch db.Dialector.Name() {
	case config.DbTypeSqlite:
		return "integer"
	case config.DbTypeMysql:
		return "tinyint(1)"
	case config.DbTypePostgres:
		return "int2"
	case config.DbTypeSqlserver:
		return "bit"
	}
	return "int"
}

// noinspection all
func (b Bool) Bool() bool {
	return b.bool
}

// noinspection all
func (b Bool) Int() int {
	if b.bool {
		return 1
	}
	return 0
}

// Resverse 反转
// noinspection all
func (b *Bool) Resverse() {
	b.bool = !b.bool
}

func (b *Bool) Set(b2 Bool) {
	b.bool = b2.bool
}

func (b *Bool) SetBool(b2 bool) {
	b.bool = b2
}

func (b *Bool) SetInt(i int) {
	b.bool = i != 0
}

func (b *Bool) SetString(s string) {
	b.bool = s != ""
}

func (b *Bool) Equals(b2 Bool) bool {
	return b.bool == b2.bool
}

func NewBool(b bool) Bool {
	return Bool{bool: b}
}
