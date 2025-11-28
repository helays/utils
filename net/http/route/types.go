package route

import (
	"database/sql/driver"

	"github.com/helays/utils/v2/dataType"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

// noinspection all
type RouteType string

// noinspection all
func (r *RouteType) Scan(val any) (err error) {
	return dataType.HelperStringScan(val, r)
}

// noinspection all
func (r RouteType) Value() (driver.Value, error) {
	return string(r), nil
}

// GormDBDataType gorm db data type
// noinspection all
func (RouteType) GormDBDataType(db *gorm.DB, field *schema.Field) string {
	return dataType.HelperStringGormDBDataType(db, field, 32)
}

// noinspection all
const (
	Static   RouteType = "static"   // 静态路由
	Param    RouteType = "param"    // 路径参数路由（如 /users/:id）
	Regex    RouteType = "regex"    // 正则表达式路由
	Wildcard RouteType = "wildcard" // 通配符路由
)

type Config struct {
	Root  string   `json:"root" yaml:"root"`   // 文件根目录
	Index []string `json:"index" yaml:"index"` // 默认首页
}

type FileRoute struct {
	Name       string // 文件名
	Downloader bool   // 是否允许下载
}
