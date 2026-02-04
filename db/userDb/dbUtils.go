package userDb

import (
	"fmt"
	"strings"

	"github.com/helays/utils/v2/config"
	"github.com/helays/utils/v2/db/dbErrors/errTools"
	"github.com/helays/utils/v2/logger/ulogs"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/logger"
)

var nullSqlConn map[string]*gorm.DB

// GetRawSql 生成sql的通用函数
// pg postgres.Open("")
// mysql mysql.Open("")
func GetRawSql(f func(dTx *gorm.DB) *gorm.DB, dialector gorm.Dialector) (string, []any) {
	db, _ := gorm.Open(dialector, &gorm.Config{
		DryRun:                                   true,
		Logger:                                   logger.Default.LogMode(logger.Silent),
		DisableForeignKeyConstraintWhenMigrating: true,
		SkipDefaultTransaction:                   true,
		DisableAutomaticPing:                     true,
	})

	query := f(db).Statement
	return query.SQL.String(), query.Vars
}

// GetRawSqlByDb 生成sql的通用函数,db为数据库连接
func GetRawSqlByDb(f func(dTx *gorm.DB) *gorm.DB, db *gorm.DB) (string, []any) {
	query := f(db).Statement
	return query.SQL.String(), query.Vars
}

func RenameConstraint(tx *gorm.DB, tableName string, oldName, newName string) error {
	if tx == nil || tx.Dialector == nil || tx.Dialector.Name() != config.DbTypePostgres {
		return nil
	}
	return tx.Exec("ALTER TABLE ? RENAME CONSTRAINT ? TO ?", clause.Table{Name: tableName}, clause.Column{Name: oldName}, clause.Column{Name: newName}).Error
}

// ClearSequenceFieldDefaultValue 清除自增序列字段的默认值
func ClearSequenceFieldDefaultValue(tx *gorm.DB, tableName string, seqFields []string) error {
	if tx == nil || tx.Dialector == nil || tx.Dialector.Name() != config.DbTypePostgres {
		return nil
	}
	for _, seqField := range seqFields {
		err := tx.Exec("ALTER TABLE ? ALTER COLUMN ? DROP DEFAULT", clause.Table{Name: tableName}, clause.Column{Name: seqField}).Error
		// 如果报错，需要非表不存在才行
		if err != nil && !errTools.IsTableNotExist(err) && !errTools.IsColumnNotExist(err) {
			return fmt.Errorf("清除表%s自增序列字段%s失败:%s", tableName, seqField, err.Error())
		}
	}
	return nil
}

// DropSequence 删除postgreSQL数据库中指定表的序列。
func DropSequence(tx *gorm.DB, tableName string, seqFields []string) error {
	if tx == nil || tx.Dialector == nil || tx.Dialector.Name() != config.DbTypePostgres {
		return nil
	}
	for _, seqField := range seqFields {
		var seqName string
		tx.Raw("SELECT pg_get_serial_sequence(?, ?)", tableName, seqField).Scan(&seqName)
		if seqName != "" {
			tx.Exec("DROP SEQUENCE  IF EXISTS ?", clause.Table{Name: seqName})
		}
	}
	return nil
}

// UpdateSeq 更新postgreSQL数据库中指定表的序列值。
// 该函数用于确保序列值在插入新记录时不会产生间隙，通常在删除记录或导入数据后调用。
// 参数:
// tableName - 表名，序列所属的表。
// field - 序列字段名，需要重置的序列对应的字段名
func UpdateSeq(utx *gorm.DB, tableName string) {
	defer func() {
		if err := recover(); err != nil {
			ulogs.Error("【UpdateSeq】更新表[%s]自增序列值失败", tableName, err)
		}
	}()
	if utx == nil || utx.Dialector == nil || utx.Dialector.Name() != config.DbTypePostgres {
		return
	}
	ulogs.Infof("【UpdateSeq】开始更新表【%s】自增序列", tableName)
	// 如果自增字段 autoIncrementField不为空，那么再插入完成后，需要使用这句话 SELECT setval(pg_get_serial_sequence('test', 'id'), COALESCE((SELECT MAX(id)+1 FROM test), 1), false) 重置自增字段的值
	// 如果是pg数据库，这里需要获取当前表的主键字段，并判断其是否是自增主键，如果是自增主键，就将字段查询出来，放入autoIncrementField []string 变量中
	var autoIncrementField []string
	const querySeqColumnsSQL = "SELECT column_name FROM information_schema.columns WHERE table_name = ? AND (column_default LIKE 'nextval%' OR is_identity='YES')"
	if err := utx.Raw(querySeqColumnsSQL, tableName).Scan(&autoIncrementField).Error; err != nil {
		ulogs.Errorf("【UpdateSeq】pg数据库查询表[%s]自增字段失败 %v", tableName, err)
	}
	ulogs.Infof("【UpdateSeq】表【%s】自增字段为：%v", tableName, autoIncrementField)

	const updateSeqValueSQL = "SELECT setval(pg_get_serial_sequence(?, ?), COALESCE((SELECT MAX(?)+1 FROM ?), 1), false)"
	for _, field := range autoIncrementField {
		if err := utx.Debug().Exec(
			updateSeqValueSQL,
			tableName,
			field,
			clause.Column{Name: field},
			clause.Table{Name: tableName},
		).Error; err != nil {
			ulogs.Errorf("【UpdateSeq】pg数据库重置表[%s]自增字段[%s]失败 %v", tableName, field, err)
		}
	}
}

// ResetSequence 重置postgreSQL数据库中指定表的序列值。
func ResetSequence(tx *gorm.DB, tableName string, seqFields []string) error {
	if tx == nil || tx.Dialector == nil || tx.Dialector.Name() != config.DbTypePostgres {
		return nil
	}
	for _, field := range seqFields {
		err := tx.Exec("SELECT setval(pg_get_serial_sequence(?, ?), 1, false)", tableName, field).Error
		if err != nil && !errTools.IsTableNotExist(err) && !errTools.IsColumnNotExist(err) {
			ulogs.Errorf("【ResetSequence】pg数据库重置表[%s]自增字段[%s]失败 %v", tableName, field, err)
		}
	}
	return nil
}

// UnionAllScope 可复用的 UNION ALL 查询 Scope
func UnionAllScope(queries ...*gorm.DB) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if len(queries) == 0 {
			return db
		}
		// 单表查询直接返回
		if len(queries) == 1 {
			return db.Table("(?) AS union_table", queries[0])
		}
		// 多表查询构建 UNION ALL
		unionQuery := queries[0]
		for i := 1; i < len(queries); i++ {
			unionQuery = db.Session(&gorm.Session{}).Raw("? UNION ALL ?", unionQuery, queries[i])
		}
		return db.Session(&gorm.Session{}).Table("(?) AS union_table", unionQuery)
	}
}

func FindTableWithPrefix(tx *gorm.DB, prefix string) ([]string, error) {
	var tables []string
	var err error
	switch tx.Dialector.Name() {
	case config.DbTypePostgres:
		// 还要查询当前的搜索模式
		// 获取当前搜索模式
		var searchPath string
		if err = tx.Raw("SHOW search_path").Scan(&searchPath).Error; err != nil {
			return nil, fmt.Errorf("获取当前连接的搜索模式失败：%s", err.Error())
		}
		// 默认搜索模式是第一个模式
		currentSchema := "public" // 默认值
		if len(searchPath) > 0 {
			currentSchema = strings.Split(searchPath, ",")[0] // 取第一个模式
			currentSchema = strings.TrimSpace(currentSchema)  // 去除空格
		}
		_tx := tx.Raw("SELECT table_name FROM information_schema.tables WHERE table_schema like ? and table_name LIKE ?", currentSchema, prefix+"%")
		err = _tx.Scan(&tables).Error
	case config.DbTypeMysql:
		currentDataBase := tx.Migrator().CurrentDatabase()
		_tx := tx.Raw("SELECT table_name FROM information_schema.tables WHERE table_schema like ? and table_name LIKE ?", currentDataBase, prefix+"%")
		err = _tx.Scan(&tables).Error
	}
	if err != nil {
		return nil, fmt.Errorf("根据前缀查询表清单失败：%s", err.Error())
	}
	return tables, nil
}

func Truncate(tx *gorm.DB, tableName string) error {
	var sql string
	dialect := tx.Dialector.Name()
	switch dialect {
	case config.DbTypeSqlite:
		sql = fmt.Sprintf("DELETE FROM %s", tableName)
	case config.DbTypeMysql:
		sql = fmt.Sprintf("TRUNCATE TABLE %s", tableName)
	case config.DbTypePostgres:
		sql = fmt.Sprintf("TRUNCATE TABLE %s", tableName)
	default:
		return fmt.Errorf("不支持的数据库类型：%s", dialect)
	}
	err := tx.Exec(sql).Error
	if err != nil {
		return fmt.Errorf("清空表失败：%s", err.Error())
	}
	if dialect == "sqlite" {
		if err = tx.Exec(fmt.Sprintf("DELETE FROM sqlite_sequence WHERE name = '%s'", tableName)).Error; err != nil {
			return fmt.Errorf("sqlite重置序列失败：%s", err.Error())
		}
	}
	return nil
}

// IsTiDB 判断是否是TiDB数据库
func IsTiDB(tx *gorm.DB) bool {
	if tx == nil || tx.Dialector == nil {
		return false
	}
	if driver, ok := tx.Dialector.(*mysql.Dialector); ok {
		return strings.Contains(strings.ToLower(driver.Config.ServerVersion), "tidb")
	}
	return false
}
