package httpServerWithDb

import (
	"github.com/helays/utils/db/userDb"
	"github.com/helays/utils/logger/ulogs"
	"github.com/helays/utils/net/http/httpServer"
	"github.com/helays/utils/tools"
	"gorm.io/gorm"
	"net/http"
)

const (
	PageField     = "p"
	PageSizeField = "pn"
	PageSize      = 30
)

type Pager struct {
	PageSize      int    `ini:"page_size" yaml:"page_size" json:"page_size"` // 系统默认查询数量
	PageSizeField string `ini:"page_size_field" yaml:"page_size_field" json:"page_size_field"`
	PageField     string `ini:"page_field" yaml:"page_field" json:"page_field"`
	Order         string
}

// RespFilter 是一个接口，定义了 FiltterDatas 方法
type RespFilter interface {
	RespFilter()
}

type RespDataStruct[T any] struct {
	Lists T     `json:"lists"`
	Total int64 `json:"total"`
}

// RespListsWithFilter 是一个通用的查询列表函数，有默认的数据过滤函数
func RespListsWithFilter[T RespFilter](w http.ResponseWriter, r *http.Request, tx *gorm.DB, respData T, pager Pager) {
	var totals int64
	tx.Scopes(userDb.QueryDateTimeRange(r))
	tx.Count(&totals)
	tx.Order(pager.Order)
	if err := tx.Scopes(userDb.Paginate(r, pager.PageField, pager.PageSizeField, pager.PageSize)).Find(&respData).Error; err != nil {
		httpServer.SetReturn(w, 1, "数据查询失败")
		ulogs.Error(err, r.URL.Path, r.URL.RawQuery, "respLists", "tx.Find")
		return
	}
	// 调用 FiltterDatas 方法
	respData.RespFilter()
	httpServer.SetReturnData(w, 0, "成功", RespDataStruct[T]{Lists: respData, Total: totals})
}

// respLists 通用查询列表
func respLists[T any](w http.ResponseWriter, r *http.Request, tx *gorm.DB, respData T, pager Pager) {
	var totals int64
	tx.Scopes(userDb.QueryDateTimeRange(r))
	tx.Count(&totals)
	tx.Order(pager.Order)
	if err := tx.Scopes(userDb.Paginate(r, pager.PageField, pager.PageSizeField, pager.PageSize)).Find(&respData).Error; err != nil {
		httpServer.SetReturn(w, 1, "数据查询失败")
		ulogs.Error(err, r.URL.Path, r.URL.RawQuery, "respLists", "tx.Find")
		return
	}

	httpServer.SetReturnData(w, 0, "成功", RespDataStruct[T]{Lists: respData, Total: totals})
}

// RespListsPkRowId 通用查询列表 主键 row_id
func RespListsPkRowId[T any](w http.ResponseWriter, r *http.Request, tx *gorm.DB, respData T, pager ...Pager) {
	var (
		pageField     = PageField     // 页面默认字段
		pageSizeField = PageSizeField // 页面呈现数量默认字段
		pageSize      = PageSize      // 每页默认数量
		order         = "row_id desc"
	)
	if len(pager) > 0 {
		_pager := pager[0]
		pageField = tools.Ternary(_pager.PageField == "", pageField, _pager.PageField)
		pageSizeField = tools.Ternary(_pager.PageSizeField == "", pageSizeField, _pager.PageSizeField)
		pageSize = tools.Ternary(_pager.PageSize < 1, pageSize, _pager.PageSize)
		order = tools.Ternary(_pager.Order == "", order, _pager.Order)
	}
	respLists(w, r, tx, respData, Pager{
		PageSize:      pageSize,
		PageSizeField: pageSizeField,
		PageField:     pageField,
		Order:         order,
	})
}

// RespListsPkId 根据查询参数分页获取数据列表，并按指定字段排序。
// 该函数是一个泛型函数，可以处理任何类型的响应数据。
// 参数:
//
//	w: http.ResponseWriter，用于写入HTTP响应。
//	r: *http.Request，当前的HTTP请求。
//	tx: *gorm.DB，数据库事务对象，用于执行数据库查询。
//	respData: T，响应数据的结构体，用于存储查询结果。
//	pager: ...Pager，可变参数，用于自定义分页和排序行为。
func RespListsPkId[T any](w http.ResponseWriter, r *http.Request, tx *gorm.DB, respData T, pager ...Pager) {
	var (
		pageField     = PageField     // 页面默认字段
		pageSizeField = PageSizeField // 页面呈现数量默认字段
		pageSize      = PageSize      // 每页默认数量
		order         = "id desc"
	)
	if len(pager) > 0 {
		_pager := pager[0]
		pageField = tools.Ternary(_pager.PageField == "", pageField, _pager.PageField)
		pageSizeField = tools.Ternary(_pager.PageSizeField == "", pageSizeField, _pager.PageSizeField)
		pageSize = tools.Ternary(_pager.PageSize < 1, pageSize, _pager.PageSize)
		order = tools.Ternary(_pager.Order == "", order, _pager.Order)
	}
	respLists(w, r, tx, respData, Pager{
		PageSize:      pageSize,
		PageSizeField: pageSizeField,
		PageField:     pageField,
		Order:         order,
	})
}
