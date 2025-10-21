package dataType

import (
	"database/sql/driver"

	"github.com/helays/utils/v2/config"
	"github.com/helays/utils/v2/tools"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

type Byte byte

func (b Byte) Value() (driver.Value, error) {
	return int(b), nil
}

func (b *Byte) Scan(value any) error {
	if value == nil {
		*b = 0
		return nil
	}
	switch t := value.(type) {
	case byte:
		*b = Byte(t)
	case int8:
		*b = Byte(t)
	case int:
		*b = Byte(t)
	default:
		v, err := tools.Any2int(value)
		if err != nil {
			return err
		}
		*b = Byte(v)
	}
	return nil
}

func (b Byte) GormDataType() string {
	return "byte"
}

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
