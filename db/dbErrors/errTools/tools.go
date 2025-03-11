package errTools

import (
	"errors"
	"fmt"
	"github.com/helays/utils/db/dbErrors"
	"github.com/helays/utils/db/dbErrors/errPostgres"
	"github.com/jackc/pgx/v5/pgconn"
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
	}
	return fmt.Errorf("%s：%s", info.ZH, err.Message)
}
