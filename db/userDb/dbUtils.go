package userDb

import (
	"database/sql"
	"fmt"
	"github.com/helays/utils/config"
	"github.com/helays/utils/db/dbErrors/errTools"
	"github.com/helays/utils/logger/ulogs"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/logger"
)

var nullSqlConn map[string]*gorm.DB

// GetRawSql 生成sql的通用函数
func GetRawSql(f func(dTx *gorm.DB) *gorm.DB, dbTypes ...string) (string, []any) {
	dbType := "pg"
	if len(dbTypes) > 1 {
		dbType = dbTypes[0]
	}
	db, ok := nullSqlConn[dbType]
	if !ok {
		var dialector gorm.Dialector
		switch dbType {
		case "pg":
			dialector = postgres.Open("")
		case "mysql":
			dialector = mysql.Open("")
			//case "sqlite":
			//	dialector = sqlite.Open("")
			//case "sqlserver":
			//	dialector = sqlserver.Open("")
			//case "tidb":
			//	dialector = mysql.Open("")
			//case "clickhouse":
			//	dialector = clickhouse.Open("")

		}
		db, _ = gorm.Open(dialector, &gorm.Config{
			DryRun:                                   true,
			Logger:                                   logger.Default.LogMode(logger.Silent),
			DisableForeignKeyConstraintWhenMigrating: true,
			SkipDefaultTransaction:                   true,
			DisableAutomaticPing:                     true,
		})
	}

	query := f(db).Statement
	return query.SQL.String(), query.Vars
}

// GetRawSqlByDb 生成sql的通用函数,db为数据库连接
func GetRawSqlByDb(f func(dTx *gorm.DB) *gorm.DB, db *gorm.DB) (string, []any) {
	query := f(db).Statement
	return query.SQL.String(), query.Vars
}

func CloseDb(conn *sql.DB) {
	if conn != nil {
		_ = conn.Close()
	}
}

func CloseMysqlRows(rows *sql.Rows) {
	CloseRows(rows)
}

func CloseRows(rows *sql.Rows) {
	if rows != nil {
		_ = rows.Close()
	}
}

// Deprecated: As of utils v1.1.0, this value is simply [tools.CloseDb].
func CloseMysql(conn *sql.DB) {
	if conn != nil {
		_ = conn.Close()
	}
}

func CloseStmt(stmt *sql.Stmt) {
	if stmt != nil {
		_ = stmt.Close()
	}
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
			ulogs.Error("更新自增序列值失败", err)
		}
	}()
	if utx == nil || utx.Dialector == nil || utx.Dialector.Name() != config.DbTypePostgres {
		return
	}
	// 如果自增字段 autoIncrementField不为空，那么再插入完成后，需要使用这句话 SELECT setval(pg_get_serial_sequence('test', 'id'), COALESCE((SELECT MAX(id)+1 FROM test), 1), false) 重置自增字段的值
	var autoIncrementField []string
	// 如果是pg数据库，这里需要获取当前表的主键字段，并判断其是否是自增主键，如果是自增主键，就将字段查询出来，放入autoIncrementField []string 变量中
	if err := utx.Raw("SELECT column_name FROM information_schema.columns WHERE table_name = ? AND column_default LIKE 'nextval%'", tableName).Scan(&autoIncrementField).Error; err != nil {
		ulogs.Error(err, "pg数据库查询自增字段失败")
	}
	for _, field := range autoIncrementField {
		if err := utx.Debug().Exec(
			"SELECT setval(pg_get_serial_sequence(?, ?), COALESCE((SELECT MAX(?)+1 FROM ?), 1), false)",
			tableName,
			field,
			clause.Column{Name: field},
			clause.Table{Name: tableName},
		).Error; err != nil {
			ulogs.Error(err, "pg数据库重置自增字段失败", tableName)
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
			ulogs.Error(err, "pg数据库重置自增字段失败", tableName)
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
