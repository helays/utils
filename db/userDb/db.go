package userDb

import (
	"github.com/helays/utils/config"
	"github.com/helays/utils/logger/ulogs"
	"github.com/helays/utils/tools"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"net/http"
	"reflect"
	"strconv"
	"strings"
)

//
// ━━━━━━神兽出没━━━━━━
// 　　 ┏┓     ┏┓
// 　　┏┛┻━━━━━┛┻┓
// 　　┃　　　　　 ┃
// 　　┃　　━　　　┃
// 　　┃　┳┛　┗┳  ┃
// 　　┃　　　　　 ┃
// 　　┃　　┻　　　┃
// 　　┃　　　　　 ┃
// 　　┗━┓　　　┏━┛　Code is far away from bug with the animal protecting
// 　　　 ┃　　　┃    神兽保佑,代码无bug
// 　　　　┃　　　┃
// 　　　　┃　　　┗━━━┓
// 　　　　┃　　　　　　┣┓
// 　　　　┃　　　　　　┏┛
// 　　　　┗┓┓┏━┳┓┏┛
// 　　　　 ┃┫┫ ┃┫┫
// 　　　　 ┗┻┛ ┗┻┛
//
// ━━━━━━感觉萌萌哒━━━━━━
//
// User helei
// Date: 2023/9/1 11:19
//

// Paginate 分页通用部分
func Paginate(r *http.Request, pageField, pageSizeField string, pageSize int) func(db *gorm.DB) *gorm.DB {
	if pageField == "" {
		pageField = "pageNo"
	}
	if pageSizeField == "" {
		pageSizeField = "pageSize"
	}
	return func(db *gorm.DB) *gorm.DB {
		page, _ := strconv.Atoi(r.URL.Query().Get(pageField))
		if page < 1 {
			page = 1
		}
		limit, _ := strconv.Atoi(r.URL.Query().Get(pageSizeField))
		if limit < 1 {
			limit = pageSize
		}
		limit = tools.Ternary(limit < 1, 30, limit)
		tx := db
		if r.URL.Query().Get("rall") != "1" {
			tx.Offset((page - 1) * limit).Limit(limit)
		}
		_sort := r.URL.Query().Get("sort")
		if _sort != "" && !config.SpecialChartPreg.MatchString(_sort) {
			if _sort[0] == '-' {
				tx.Order(_sort[1:] + " desc")
			} else {
				tx.Order(_sort)
			}
		}
		return tx
	}
}

// FilterWhereString 过滤string 条件
func FilterWhereString(r *http.Request, query string, field string, like bool) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		value := r.URL.Query().Get(field)
		if value == "" {
			return db
		}
		if like {
			return db.Where(query, "%"+value+"%")
		}
		return db.Where(query, value)
	}
}

type stackItem struct {
	s     any    // model struct
	alias string // 设置的别名
}

// FilterWhereByDbModel 通过DB 实例设置的model 自动映射查询字段
func FilterWhereByDbModel(alias string, enableDefault bool, r *http.Request, likes ...map[string]string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		stack := []stackItem{{db.Statement.Model, alias}}
		query := r.URL.Query()
		conditions := make([]clause.Expression, 0) // 收集所有查询条件
		for len(stack) > 0 {
			item := stack[len(stack)-1]
			stack = stack[:len(stack)-1]

			t := reflect.TypeOf(item.s)
			if t.Kind() == reflect.Ptr {
				t = t.Elem()
			}
			if t.Kind() != reflect.Struct {
				continue
			}

			tableName := item.alias
			v := reflect.ValueOf(item.s)
			if tableName == "" {
				tbName := v.MethodByName("TableName")
				if tbName.IsValid() {
					tableName = tbName.Call([]reflect.Value{})[0].String()
				} else {
					tableName = tools.SnakeString(t.Name())
				}

				item.alias = tableName
			}
			if v.Kind() == reflect.Ptr {
				v = v.Elem()
			}

			for i := 0; i < t.NumField(); i++ {
				field := t.Field(i)
				tagName := field.Tag.Get("json")
				// 处理嵌套结构体
				if field.Type.Kind() == reflect.Struct && field.Tag.Get("gorm") == "" && tagName == "" {
					stack = append(stack, stackItem{v.Field(i).Interface(), item.alias})
					continue
				}
				// 检查字段类型
				fieldType := field.Type.String()
				if fieldType != "int" && fieldType != "string" {
					continue
				}
				if tagName == "" {
					continue
				}
				val := query.Get(strings.Split(tagName, ",")[0])
				if val == "" {
					if !enableDefault {
						continue
					}
					// 如果没有传值，判断是否有默认值
					if val = field.Tag.Get("default"); val == "" {
						continue
					}
				}
				// 这里还需要解析出字段本身的名字，去数据库进行查询，通过将结构体转成蛇形方式。
				fieldName := tableName + "." + tools.SnakeString(field.Name)

				if fieldType == "int" {
					valList := strings.Split(val, ",")
					if len(valList) > 1 {
						conditions = append(conditions, clause.IN{Column: fieldName, Values: toInterfaceSlice(valList)})
					} else {
						conditions = append(conditions, clause.Eq{Column: fieldName, Value: val})
					}
				} else {
					lastVal := applyLikes(val, field.Tag.Get("dblike"), likes, fieldName)
					conditions = append(conditions, clause.Like{Column: fieldName, Value: lastVal})
				}
			}
		}
		// 一次性应用所有查询条件
		if len(conditions) > 0 {
			db.Clauses(conditions...)
		}
		return db
	}
}

// applyLikes 处理 like 查询的值
func applyLikes(val, dblikeTag string, likes []map[string]string, fieldName string) string {
	lastVal := val
	if dblikeTag == "%" {
		lastVal = "%" + val + "%"
	}
	if len(likes) > 0 {
		if custom, ok := likes[0][fieldName]; ok {
			switch custom {
			case "%%":
				lastVal = "%" + val + "%"
			case "-%":
				lastVal = val + "%"
			case "%-":
				lastVal = "%" + val
			default:
				lastVal = val
			}
		}
	}
	return lastVal
}

// toInterfaceSlice 将字符串切片转换为 interface{} 切片
func toInterfaceSlice(strs []string) []any {
	result := make([]any, len(strs))
	for i, v := range strs {
		result[i] = v
	}
	return result
}

// FilterWhereStruct 通过结构体 自动映射查询字段
// Deprecated: 在FilterWhereByDbModel出来后，尽量通过这个函数来实现通过结构体 自动处理query 参数 转换到 sql where里面
func FilterWhereStruct(s any, alias string, enableDefault bool, r *http.Request, likes ...map[string]string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		t := reflect.TypeOf(s)
		if t.Kind() == reflect.Ptr {
			t = t.Elem()
		}
		if t.Kind() != reflect.Struct {
			return db
		}
		tableName := alias
		v := reflect.ValueOf(s)
		if tableName == "" {
			tbName := v.MethodByName("TableName")
			if tbName.IsValid() {
				tableName = tbName.Call([]reflect.Value{})[0].String()
			} else {
				tableName = tools.SnakeString(t.Name())
			}

			alias = tableName
		}
		if v.Kind() == reflect.Ptr {
			v = v.Elem()
		}
		query := r.URL.Query()
		for i := 0; i < t.NumField(); i++ {
			if t.Field(i).Type.Kind() == reflect.Struct && t.Field(i).Tag.Get("gorm") == "" && t.Field(i).Tag.Get("json") == "" {
				db.Scopes(FilterWhereStruct(v.Field(i).Interface(), alias, enableDefault, r, likes...))
				continue
			}
			if t.Field(i).Type.String() != "int" && t.Field(i).Type.String() != "string" {
				continue
			}
			tagName := t.Field(i).Tag.Get("json")

			if tagName == "" {
				continue
			}
			val := query.Get(strings.Split(tagName, ",")[0])
			if val == "" {
				if !enableDefault {
					continue
				}
				// 如果没有传值，判断是否有默认值
				val = t.Field(i).Tag.Get("default")
				if val == "" {
					continue
				}
			}
			// 这里还需要解析出字段本身的名字，去数据库进行查询，通过将结构体转成蛇形方式。
			fieldName := tableName + "." + tools.SnakeString(t.Field(i).Name)
			if t.Field(i).Type.String() == "int" {
				valList := strings.Split(val, ",")
				if len(valList) > 1 {
					db.Where(fieldName+" in ?", valList)
				} else {
					db.Where(fieldName+" = ?", val)
				}
			} else {
				lastVal := val
				if t.Field(i).Tag.Get("dblike") == "%" {
					lastVal = "%" + val + "%"
				}
				if len(likes) > 0 {
					if custom, ok := likes[0][fieldName]; ok {
						switch custom {
						case "%%":
							lastVal = "%" + val + "%"
						case "-%":
							lastVal = val + "%"
						case "%-":
							lastVal = "%" + val
						default:
							lastVal = val
						}
					}
				}
				db.Where(fieldName+" like ?", lastVal)
			}
		}
		return db
	}
}

func FilterWhereData(data any, tableName ...string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		t := reflect.TypeOf(data)
		if t.Kind() != reflect.Struct {
			return db
		}
		for i := 0; i < t.NumField(); i++ {
			tagName := t.Field(i).Tag.Get("db")
			if !strings.Contains(tagName, "filter") {
				continue
			}
			var tagMap = make(map[string]string)
			for _, v := range strings.Split(tagName, ";") {
				if strings.Contains(v, ":") {
					tagMap[strings.Split(v, ":")[0]] = strings.Split(v, ":")[1]
				}
			}
			jsonTag := strings.Split(t.Field(i).Tag.Get("json"), ",")[0]
			if len(tableName) > 0 {
				jsonTag = tableName[0] + "." + jsonTag
			}
			switch t.Field(i).Type.Kind() {
			case reflect.String:
				v := reflect.ValueOf(data).Field(i).String()
				if v == "" {
					continue
				}
				if strings.Contains(tagName, "%%") {
					db.Where(jsonTag+" like ?", "%"+v+"%")
				} else if strings.Contains(tagName, "%-") {
					db.Where(jsonTag+" like ?", "%"+v)
				} else if strings.Contains(tagName, "-%") {
					db.Where(jsonTag+" like ?", v+"%")
				} else {
					db.Where(jsonTag+"=?", v)
				}
			case reflect.Int:
				v := reflect.ValueOf(data).Field(i).Int()
				if tagMap["ignore"] != "" {
					// 传递值 = 忽略值
					if _tmp, err := strconv.Atoi(tagMap["ignore"]); err == nil && v == int64(_tmp) {
						continue
					}
				}
				db.Where(jsonTag+"=?", v)
			}
		}
		return db
	}
}

// QueryDateTimeRange 时间区间查询
func QueryDateTimeRange(r *http.Request, filed ...string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		sTime := r.URL.Query().Get("begin_time")
		eTime := r.URL.Query().Get("end_time")
		sField := "create_time"
		if len(filed) > 0 {
			sField = filed[0]
		}
		if sTime != "" {
			db.Where(sField+" > ?", sTime)
		}
		if eTime != "" {
			db.Where(sField+" <= ?", eTime)
		}
		return db
	}
}

// AutoCreateTableWithStruct 根据结构体判断是否需要创建表
func AutoCreateTableWithStruct(db *gorm.DB, tb any, errmsg string) {
	t := reflect.TypeOf(tb)
	if t.Kind() != reflect.Struct {
		return
	}
	if !db.Migrator().HasTable(tb) {
		ulogs.DieCheckerr(db.Debug().AutoMigrate(tb), errmsg)
	}
	// 如果表存在，在判断结构体中是否有新增字段，如果有，就自动改变表
	AutoCreateTableWithColumn(db, tb, errmsg, t)
}

// AutoCreateTableWithColumn 根据表字段判断是否有数据缺失
func AutoCreateTableWithColumn(db *gorm.DB, tb any, errmsg string, t reflect.Type) bool {
	// 如果表存在，在判断结构体中是否有新增字段，如果有，就自动改变表
	for i := 0; i < t.NumField(); i++ {
		if t.Field(i).Type.Kind() == reflect.Struct && t.Field(i).Tag.Get("gorm") == "" && t.Field(i).Tag.Get("json") == "" {
			if AutoCreateTableWithColumn(db, tb, errmsg, t.Field(i).Type) {
				return true
			}
			continue
		}
		tag := t.Field(i).Tag.Get("gorm")
		if tag == "" {
			continue
		}
		if tag == "-:all" || tag == "-:migration" || strings.Contains(tag, "-:migration") {
			continue
		}
		column := tools.SnakeString(t.Field(i).Name)
		for _, item := range strings.Split(tag, ";") {
			if !strings.HasPrefix(item, "column") {
				continue
			}
			column = item[7:]
		}

		if !db.Migrator().HasColumn(tb, column) {
			ulogs.Log("表字段有缺失，正在自动创建表字段：", reflect.TypeOf(tb).String(), column)
			ulogs.DieCheckerr(db.Debug().AutoMigrate(tb), errmsg)
			return true // 创建一次就行了
		}
	}
	return false
}
