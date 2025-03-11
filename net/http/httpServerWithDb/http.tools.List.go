package httpServerWithDb

import (
	"fmt"
	"github.com/helays/utils/config"
	"github.com/helays/utils/db/userDb"
	"github.com/helays/utils/net/http/httpServer"
	"gorm.io/gorm"
	"net/http"
	"strings"
)

// ListMethodGet 是一个通用的列表查询方法，用于根据不同的条件获取数据库中的记录。
// 它使用了泛型 T，允许任何类型的列表查询。
// 参数:
//
//	w http.ResponseWriter: 用于写入HTTP响应。
//	r *http.Request: 包含当前HTTP请求的详细信息。
//	tx *gorm.DB: GORM数据库连接对象，用于执行数据库操作。
//	c userDb.Curd: 查询配置，包含了查询所需的配置信息，如选择查询、条件查询等。
//	p Pager: 分页配置，用于指定查询的分页信息。
func ListMethodGet[T any](w http.ResponseWriter, r *http.Request, tx *gorm.DB, c userDb.Curd, p Pager) {
	if config.Dbg {
		tx = tx.Debug()
	}
	query := r.URL.Query()
	if c.MustField != nil {
		for k, rule := range c.MustField {
			v := query.Get(k)
			if !rule.MatchString(v) {
				httpServer.SetReturnErrorDisableLog(w, fmt.Errorf("参数%s值格式错误", k), http.StatusBadRequest, "参数错误")
				return
			}
		}
	}
	_tx := tx.Model(new(T))
	_tx.Scopes(userDb.FilterWhereStruct(new(T), "", false, r))
	if c.Select.Query != "" {
		_tx.Select(c.Select.Query, c.Select.Args...)
	}
	if c.Where.Query != "" {
		_tx.Where(c.Where.Query, c.Where.Args...)
	}
	if c.Omit != nil && len(c.Omit) > 0 {
		_tx.Omit(c.Omit...)
	}
	for _, item := range c.Preload {
		_tx.Preload(item.Query, item.Args...)
	}
	var list []T
	switch strings.ToLower(c.Pk) {
	case "id":
		RespListsPkId(w, r, _tx, list, p)
	case "row_id":
		RespListsPkRowId(w, r, _tx, list, p)
	default:
		return
	}
}
