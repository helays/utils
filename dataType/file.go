package dataType

import (
	"database/sql/driver"
	"errors"

	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

type Binary []byte

// Scan 实现了 sql.Scanner 接口，用于从数据库读取二进制数据到 BinaryFile 类型
// noinspection all
func (bf *Binary) Scan(value interface{}) error {
	if value == nil {
		*bf = nil
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	// 必须复制：database/sql 会把驱动返回的 []byte 原样交给 Scanner，
	// 而 go-sql-driver/mysql 的结果切片指向连接读缓冲，后续任何查询都会覆盖它。
	// 直接持有（*bf = src）会导致读出的 BLOB 在下一次查询后变成另一行的内容。
	*bf = append([]byte(nil), bytes...)
	return nil
}

// Value 实现了 driver.Valuer 接口，用于将 BinaryFile 类型的数据写入数据库
// noinspection all
func (bf Binary) Value() (driver.Value, error) {
	return []byte(bf), nil
}

// GormDBDataType 返回数据库特定的数据类型名称
// noinspection all
func (Binary) GormDBDataType(db *gorm.DB, field *schema.Field) string {
	return BlobDbDataType(db, field)
}
