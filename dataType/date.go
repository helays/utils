package dataType

import (
	"database/sql"
	"database/sql/driver"
	"github.com/helays/utils/config"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
	"strings"
	"time"
)

type CustomTime time.Time

func (this CustomTime) String() string {
	return time.Time(this).Format(time.DateTime)
}

func (this CustomTime) After(u time.Time) bool {
	return time.Time(this).After(u)
}

func (this CustomTime) Before(u time.Time) bool {
	return time.Time(this).Before(u)
}

func (this CustomTime) Sub(u time.Time) time.Duration {
	return time.Time(this).Sub(u)
}

func (this CustomTime) Unix() int64 {
	return time.Time(this).Unix()
}

func (this *CustomTime) Scan(value interface{}) (err error) {
	nullTime := &sql.NullTime{}
	err = nullTime.Scan(value)
	*this = CustomTime(nullTime.Time)
	return
}

func (this CustomTime) Value() (driver.Value, error) {
	return time.Time(this), nil
}

func (this CustomTime) GormDataType() string {
	return "time"
}

// GormDBDataType gorm db data type
func (CustomTime) GormDBDataType(db *gorm.DB, field *schema.Field) string {
	switch db.Dialector.Name() {
	case "sqlite":
		return "timestamp"
	case "mysql":
		return "timestamp"
	case "postgres":
		return "timestamp with time zone"
	case "sqlserver":
		return "timestamp"
	}
	return ""
}

func (this CustomTime) GobEncode() ([]byte, error) {
	return time.Time(this).GobEncode()
}

func (this *CustomTime) GobDecode(b []byte) error {
	return (*time.Time)(this).GobDecode(b)
}

func (this CustomTime) MarshalJSON() ([]byte, error) {
	b := []byte{'"'}
	b = append(b, []byte(time.Time(this).Format(time.DateTime))...)
	b = append(b, '"')
	return b, nil
}
func (this *CustomTime) UnmarshalJSON(b []byte) (err error) {
	s := strings.Trim(string(b), "\"")
	if s == "null" {
		*this = CustomTime{}
		return nil
	}
	_t, err := time.ParseInLocation(time.DateTime, s, config.CstSh)
	if err != nil {
		_t, err = time.ParseInLocation(time.RFC3339Nano, s, config.CstSh)
	}
	*this = CustomTime(_t)
	return err
}
