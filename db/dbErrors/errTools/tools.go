package errTools

import (
	"errors"
	"fmt"
	"github.com/go-sql-driver/mysql"
	"github.com/helays/utils/db/dbErrors"
	"github.com/helays/utils/db/dbErrors/errPostgres"
	"github.com/helays/utils/logger/ulogs"
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
	// 其他数据库或未知错误
	ulogs.Errorf("unknown error type: %T, err: %v\n", err, err)
	return false
}
