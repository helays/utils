package dataType

import (
	"database/sql"
	"database/sql/driver"
	"fmt"

	"github.com/helays/utils/v2/config"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

type String string

func (s *String) Scan(val interface{}) (err error) {
	return HelperStringScan(val, s)
}

func (s String) Value() (driver.Value, error) {
	return string(s), nil
}

func (s String) GormDataType() string {
	return "string"
}

// GormDBDataType gorm db data type
func (String) GormDBDataType(db *gorm.DB, field *schema.Field) string {
	switch db.Dialector.Name() {
	case config.DbTypeSqlite:
		return "text"
	case config.DbTypeMysql:
		return "LONGTEXT"
	case config.DbTypePostgres:
		return "text"
	case config.DbTypeSqlserver:
		return "VARCHAR(MAX)"
	}
	return ""
}

// HelperStringScan  通用的 string 自定义类型 SCAN函数
func HelperStringScan[T ~string](val any, dst *T) error {
	nullStr := sql.NullString{}
	if err := nullStr.Scan(val); err != nil {
		return err
	}
	if !nullStr.Valid {
		*dst = *new(T)
		return nil
	}
	*dst = T(nullStr.String)
	return nil
}

func HelperStringGormDBDataType(db *gorm.DB, field *schema.Field, length int) string {
	switch db.Dialector.Name() {
	case config.DbTypeSqlite, config.DbTypeMysql, config.DbTypePostgres:
		return fmt.Sprintf("varchar(%d)", length)
	case config.DbTypeSqlserver:
		return fmt.Sprintf("VARCHAR(%d)", length)
	}
	return fmt.Sprintf("varchar(%d)", length)
}
