package dataType

import (
	"database/sql/driver"
	"errors"

	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

type Binary []byte

// Scan 实现了 sql.Scanner 接口，用于从数据库读取二进制数据到 BinaryFile 类型
func (bf *Binary) Scan(value interface{}) error {
	if value == nil {
		*bf = nil
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}

	*bf = bytes
	return nil
}

// Value 实现了 driver.Valuer 接口，用于将 BinaryFile 类型的数据写入数据库
func (bf Binary) Value() (driver.Value, error) {
	return []byte(bf), nil
}

// GormDBDataType 返回数据库特定的数据类型名称
func (Binary) GormDBDataType(db *gorm.DB, field *schema.Field) string {
	return BlobDbDataType(db, field)
}
