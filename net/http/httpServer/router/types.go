package router

import (
	"database/sql/driver"

	"github.com/helays/utils/v2/dataType"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

type RouteType string

func (r *RouteType) Scan(val any) (err error) {
	return dataType.HelperStringScan(val, r)
}

func (r RouteType) Value() (driver.Value, error) {
	return string(r), nil
}

// GormDBDataType gorm db data type
func (RouteType) GormDBDataType(db *gorm.DB, field *schema.Field) string {
	return dataType.HelperStringGormDBDataType(db, field, 32)
}

const (
	RouteTypeStatic   RouteType = "static"   // 静态路由
	RouteTypeParam    RouteType = "param"    // 路径参数路由（如 /users/:id）
	RouteTypeRegex    RouteType = "regex"    // 正则表达式路由
	RouteTypeWildcard RouteType = "wildcard" // 通配符路由
)
