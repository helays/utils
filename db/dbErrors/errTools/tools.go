package errTools

import (
	"errors"
	"fmt"
	"github.com/go-sql-driver/mysql"
	"github.com/helays/utils/v2/db/dbErrors"
	"github.com/helays/utils/v2/db/dbErrors/errPostgres"
	"github.com/jackc/pgx/v5/pgconn"
	"strings"
)

func Error(err error) error {
	switch _err := err.(type) {
	case *pgconn.PgError:
		info, ok := errPostgres.PgErrorMap[_err.Code]
		if !ok {
			return err
		}
		return doPostgres(_err, info)
	}
	return err
}

func doPostgres(err *pgconn.PgError, info dbErrors.DbError) error {
	switch info.Code {
	case "42703":
		return errors.New(err.Message)
	case "22001":
		return fmt.Errorf("超长导致%s", info.ZH)
	}
	return fmt.Errorf("%s：%s", info.ZH, err.Message)
}

// IsTableNotExist 检查错误是否为表不存在
func IsTableNotExist(err error) bool {
	if err == nil {
		return false
	}
	// 检查 PostgresSQL 错误
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		return pgErr.Code == "42P01" // PostgreSQL 表不存在的错误码
	}

	// 检查 MySQL 错误
	var myErr *mysql.MySQLError
	if errors.As(err, &myErr) {
		return myErr.Number == 1146 // MySQL 表不存在的错误码 (ER_NO_SUCH_TABLE)
	}
	errStr := err.Error()

	if strings.Contains(errStr, "no such table") {
		return true // 检查 SQLite 错误（通常以字符串匹配）
	} else if strings.Contains(errStr, "Invalid object name") {
		return true // 检查 SQL Server 错误
	} else if strings.Contains(errStr, "ORA-00942") { // Oracle 表或视图不存在
		return true // 检查 Oracle 错误
	}
	return false
}

// IsColumnNotExist 检查错误是否为字段(列)不存在
func IsColumnNotExist(err error) bool {
	if err == nil {
		return false
	}

	// 检查 PostgreSQL 错误
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		return pgErr.Code == "42703" // PostgreSQL 列不存在的错误码
	}

	// 检查 MySQL 错误
	var myErr *mysql.MySQLError
	if errors.As(err, &myErr) {
		return myErr.Number == 1054 // MySQL 列不存在的错误码 (ER_BAD_FIELD_ERROR)
	}

	errStr := err.Error()

	if strings.Contains(errStr, "no such column") {
		return true // 检查 SQLite 错误
	} else if strings.Contains(errStr, "Invalid column name") {
		return true // 检查 SQL Server 错误
	} else if strings.Contains(errStr, "ORA-00904") { // Oracle 无效标识符(列不存在)
		return true // 检查 Oracle 错误
	}
	return false
}

// IsDuplicateKeyError 检查错误是否为主键或唯一约束冲突，并返回约束类型
// 返回值: 0-非重复错误, 1-主键重复, 2-唯一键重复
func IsDuplicateKeyError(err error) int {
	if err == nil {
		return 0
	}

	errStr := err.Error()

	// 检查 PostgreSQL 错误
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		if pgErr.Code == "23505" { // 唯一约束违反
			if strings.Contains(pgErr.ConstraintName, "pkey") || strings.Contains(pgErr.Message, "primary key") {
				return 1 // 主键重复
			}
			return 2 // 唯一键重复
		}
		return 0
	}

	// 检查 MySQL 错误
	var myErr *mysql.MySQLError
	if errors.As(err, &myErr) {
		if myErr.Number == 1062 { // ER_DUP_ENTRY
			if strings.Contains(myErr.Message, "PRIMARY") || strings.Contains(myErr.Message, "primary key") {
				return 1 // 主键重复
			}
			return 2 // 唯一键重复
		}
		return 0
	}

	// SQLite 错误
	if strings.Contains(errStr, "UNIQUE constraint failed") {
		if strings.Contains(errStr, "PRIMARY") || strings.Contains(errStr, "primary key") {
			return 1
		}
		return 2
	}

	// SQL Server 错误
	if strings.Contains(errStr, "Violation of PRIMARY KEY constraint") {
		return 1
	} else if strings.Contains(errStr, "Violation of UNIQUE KEY constraint") {
		return 2
	}

	// Oracle 错误
	if strings.Contains(errStr, "ORA-00001") {
		if strings.Contains(errStr, "PRIMARY") || strings.Contains(errStr, "primary key") {
			return 1
		}
		return 2
	}
	return 0
}
