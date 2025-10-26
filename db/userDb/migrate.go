package userDb

import (
	"reflect"
	"strings"

	"github.com/helays/utils/v2/config"
	"github.com/helays/utils/v2/db/tablename"
	"github.com/helays/utils/v2/logger/ulogs"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

// AutoMigrate 根据结构体自动创建表
func AutoMigrate(db *gorm.DB, c tablename.TableName, model any) {
	tx := db.Session(&gorm.Session{})
	switch tx.Dialector.Name() {
	case config.DbTypeMysql:
		AutoCreateTableWithStruct(tx.Set(c.MigrateComment()), model, c.MigrateError())
	default:
		AutoCreateTableWithStruct(tx, model, c.MigrateError())
	}

}

// AutoCreateTableWithStruct 根据结构体判断是否需要创建表
func AutoCreateTableWithStruct(db *gorm.DB, tb any, errmsg string) {
	t := reflect.TypeOf(tb)
	if t.Kind() != reflect.Struct {
		return
	}
	if !db.Migrator().HasTable(tb) {
		ulogs.DieCheckerr(db.Debug().Migrator().CreateTable(tb), errmsg)
		return
	}
	// 如果表存在，在判断结构体中是否有新增字段，如果有，就自动改变表
	AutoCreateTableWithColumn(db, tb, errmsg)
}

// AutoCreateTableWithColumn 根据表字段判断是否有数据缺失
// 一次性拿到数据库的字段结构以及索引结构
// 根据结构体信息，生成本地模型的字段结构和索引结构
// 首先比对数据库与本地索引结构，将多余的索引删除调
// 比对字段信息变化，如果有,就进行alter操作
// 最后再次比对数据库与本地模型的索引结构，将缺失的索引补充上。
func AutoCreateTableWithColumn(db *gorm.DB, tb any, errmsg string) {
	stmt := db.Session(&gorm.Session{Initialized: true}).Statement
	if err := stmt.Parse(tb); err != nil {
		ulogs.DieCheckerr(err, "表模型解析失败", errmsg)
		return
	}
	// 查询数据库表字段元数据
	columnTypes, err := db.Session(&gorm.Session{Initialized: true}).Migrator().ColumnTypes(tb)
	if err != nil {
		ulogs.DieCheckerr(err, "查询数据库表字段元数据失败", errmsg)
		return
	}
	// 查询数据库表索引元数据
	dstIndexTypes, err := db.Session(&gorm.Session{Initialized: true}).Migrator().GetIndexes(tb)
	if err != nil {
		ulogs.DieCheckerr(err, "查询数据库表索引元数据失败", errmsg)
		return
	}
	// 根据结构体，解析本地模型生成索引元数据
	srcIndexTypes := stmt.Schema.ParseIndexes()
	var (
		dstIndexTypesMap  = make(map[string]gorm.Index)
		srcIndexTypesMap  = make(map[string]*schema.Index)
		dstColumnTypesMap = make(map[string]gorm.ColumnType)
	)
	for _, ct := range columnTypes {
		dstColumnTypesMap[ct.Name()] = ct
	}
	for _, index := range srcIndexTypes {
		srcIndexTypesMap[index.Name] = index
	}
	// 删除多余的索引
	for _, index := range dstIndexTypes {
		if isPk, ok := index.PrimaryKey(); ok && isPk {
			continue
		}
		// 判断数据库里面的索引是否需要删除
		idxName := index.Name()
		if _, ok := srcIndexTypesMap[idxName]; !ok {
			ulogs.Infof("表【%s】字段[%s]索引需要删除", stmt.Schema.Table, idxName)
			if err = db.Debug().Migrator().DropIndex(tb, idxName); err != nil {
				if !strings.Contains(err.Error(), "check that it exists") {
					ulogs.DieCheckerr(err, "删除数据库索引失败", errmsg)
				}
				ulogs.Errorf("表【%s】字段[%s]索引删除失败，可能是索引名有特殊字符，请人工删除 %v", stmt.Schema.Table, idxName, err)
			}
			continue
		}
		dstIndexTypesMap[index.Name()] = index
	}
	// 判断字段是否有变化
	for _, item := range stmt.Schema.Fields {
		if item.IgnoreMigration {
			continue
		}
		dstColumn, _ok := dstColumnTypesMap[item.DBName]
		// 判断字段缺失
		if !_ok {
			ulogs.Infof("表【%s】字段[%s]缺失，正在自动创建表字段", stmt.Schema.Table, item.DBName)
			ulogs.DieCheckerr(db.Debug().AutoMigrate(tb), errmsg)
			return
		}
		// 主键无相关方法，暂不处理
		// 自增无相关方法，

		// 判断字段说明是否改变
		if v, ok := dstColumn.Comment(); ok && v != item.Comment {
			ulogs.Infof("表【%s】字段[%s]字段说明不一致，正在自动重建 %s %s", stmt.Schema.Table, item.DBName, v, item.Comment)
			if err = db.Debug().Migrator().AlterColumn(tb, item.DBName); err != nil {
				ulogs.Errorf("表【%s】字段[%s]字段说明修改失败 %v", stmt.Schema.Table, item.DBName, err)
			}
			continue
		}
		// 判断允许null 是否改变
		// 这个也不能用，有的数据库时间字段设置的允许null,但是会自动改成not null。
	}
	// 判断索引是否有新增
	for _, item := range srcIndexTypes {
		idxName := item.Name
		if _, ok := dstIndexTypesMap[idxName]; !ok {
			err = db.Debug().Migrator().CreateIndex(tb, idxName)
		}
	}
}
