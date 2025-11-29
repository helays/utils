package dbresponse

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/helays/utils/v2/config"
	"github.com/helays/utils/v2/db/userDb"
	"github.com/helays/utils/v2/logger/ulogs"
	"github.com/helays/utils/v2/net/http/response"
	"github.com/helays/utils/v2/tools"
	"gorm.io/gorm"
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

type pageListResp struct {
	Lists any   `json:"lists"`
	Total int64 `json:"total"`
}

// RespListsWithFilter 是一个通用的查询列表函数，有默认的数据过滤函数
func RespListsWithFilter[B any, T RespFilter](w http.ResponseWriter, r *http.Request, tx *gorm.DB, c userDb.Curd, p Pager) {
	_tx, ok := queryBase(w, r, tx, new(B), c)
	if !ok {
		return
	}
	_tx.Scopes(userDb.QueryDateTimeRange(r))
	var totals int64
	_tx.Count(&totals)
	_tx.Order(p.Order)
	var respData T
	if err := _tx.Scopes(userDb.Paginate(r, p.PageField, p.PageSizeField, p.PageSize)).Find(&respData).Error; err != nil {
		response.SetReturn(w, 1, "数据查询失败")
		ulogs.Error(err, r.URL.Path, r.URL.RawQuery, "respLists", "tx.Find")
		return
	}
	// 调用 FiltterDatas 方法
	respData.RespFilter()
	response.SetReturnData(w, 0, "成功", pageListResp{Lists: respData, Total: totals})
}

// 查询返回数据列 base
func queryBase(w http.ResponseWriter, r *http.Request, tx *gorm.DB, model any, c userDb.Curd) (*gorm.DB, bool) {
	session := tx.Session(&gorm.Session{})
	if config.Dbg {
		session = tx.Debug()
	}
	query := r.URL.Query()
	if c.MustField != nil {
		for k, rule := range c.MustField {
			v := query.Get(k)
			if !rule.MatchString(v) {
				response.SetReturnErrorDisableLog(w, fmt.Errorf("参数%s值格式错误", k), http.StatusBadRequest, "参数错误")
				return nil, false
			}
		}
	}
	_tx := session.Model(model)
	_tx.Scopes(userDb.FilterWhereByDbModel(c.Alias, c.EnableDefault, r))
	if c.Select.Query != "" {
		_tx.Select(c.Select.Query, c.Select.Args...)
	}
	for _, join := range c.Joins {
		_tx.Joins(join.Query, join.Args...)
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
	return _tx, true
}

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
	_tx, ok := queryBase(w, r, tx, new(T), c)
	if !ok {
		return
	}
	resp := make([]T, 0)
	switch strings.ToLower(c.Pk) {
	case "id":
		RespListsPkId(w, r, _tx, resp, p)
	case "row_id":
		RespListsPkRowId(w, r, _tx, resp, p)
	default:
		return
	}
}

// RespListsPkRowId 通用查询列表 主键 row_id
func RespListsPkRowId(w http.ResponseWriter, r *http.Request, tx *gorm.DB, resp any, pager ...Pager) {
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
	respLists(w, r, tx, resp, Pager{
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
func RespListsPkId(w http.ResponseWriter, r *http.Request, tx *gorm.DB, resp any, pager ...Pager) {
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
	respLists(w, r, tx, resp, Pager{
		PageSize:      pageSize,
		PageSizeField: pageSizeField,
		PageField:     pageField,
		Order:         order,
	})
}

// respLists 通用查询列表
func respLists(w http.ResponseWriter, r *http.Request, tx *gorm.DB, resp any, pager Pager) {
	var totals int64
	tx.Scopes(userDb.QueryDateTimeRange(r))
	tx.Count(&totals)
	tx.Scopes(userDb.AutoSetSort(r, pager.Order, true)) // 通过 请求中的 get参数、tx 自动解析

	// 下面用反射创建 slice ，貌似开销较大
	//modelType := reflect.TypeOf(tx.Statement.Model) // 直接获取 tx 里面指向的模型
	//// 创建一个与 model 类型相同的切片
	//lst := reflect.MakeSlice(reflect.SliceOf(modelType), 0, 0).Interface()
	if err := tx.Scopes(userDb.Paginate(r, pager.PageField, pager.PageSizeField, pager.PageSize)).Find(&resp).Error; err != nil {
		response.SetReturn(w, 1, "数据查询失败")
		ulogs.Error(err, r.URL.Path, r.URL.RawQuery, "respLists", "tx.Find")
		return
	}
	response.SetReturnData(w, 0, "成功", pageListResp{Lists: resp, Total: totals})
}
