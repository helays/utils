package userDb

import (
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// FilterSoftDelete 软删除过滤器
func FilterSoftDelete(tableName ...string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		column := clause.Column{Name: "is_deleted"}
		if len(tableName) > 0 {
			column.Table = tableName[0]
		}
		return db.Where(clause.Eq{
			Column: column,
			Value:  1,
		})
	}
}
