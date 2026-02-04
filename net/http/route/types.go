package route

import (
	"database/sql/driver"
	"embed"

	"gorm.io/gorm"
	"gorm.io/gorm/schema"
	"helay.net/go/utils/v3/dataType"
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
	Root  string `json:"root" yaml:"root"`   // 文件根目录
	Index string `json:"index" yaml:"index"` // 默认首页

	// 静态文件的请求地址加了一个前缀，比如/html的时候，
	// 就可以通过这个参数来去掉请求 path部分中的前面这部分内容。
	// 用最后部分去匹配文件系统。
	URLPrefix string `json:"url_prefix" yaml:"url_prefix"` // 路由前缀
}

type ErrorResp struct {
	Code  int
	Msg   string
	Error error
}

type EmbedInfo struct {
	Search string // 用search 去匹配请求path的前面部分，看是否包含 string.HasPrefix(path,search)

	// 这是缓存系统的前缀
	// 当缓存系统设置的路径是 html/static/xxx的时候
	// 实际请求路径是 static/xxx
	// 这个情况下实际路径不包含html,如果直接用search匹配，就会失败，可以通过这个prefix前缀补充，再去匹配。
	Prefix string
	FS     *embed.FS // 内置embed fs 系统
}
