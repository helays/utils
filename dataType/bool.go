package dataType

import (
	"database/sql/driver"

	"github.com/helays/utils/v2/config"
	"github.com/helays/utils/v2/tools"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

type Bool bool

func (b Bool) Value() (driver.Value, error) {
	if b {
		return 1, nil
	}
	return 0, nil
}

func (b *Bool) Scan(value any) error {
	if value == nil {
		*b = false
		return nil
	}
	ok, e := tools.Any2bool(value)
	if e != nil {
		return e
	}
	*b = Bool(ok)
	return nil
}

func (b Bool) GormDataType() string {
	return "custom_bool"
}

func (Bool) GormDBDataType(db *gorm.DB, field *schema.Field) string {
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
