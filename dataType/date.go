package dataType

import (
	"database/sql"
	"database/sql/driver"
	"github.com/helays/utils/v2/config"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
	"strings"
	"time"
)

type CustomTime time.Time

func (this CustomTime) String() string {
	return time.Time(this).Format(time.DateTime)
}

func (this CustomTime) Format(layout string) string {
	return time.Time(this).Format(layout)
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
	case config.DbTypeSqlite:
		return "timestamp"
	case config.DbTypeMysql:
		return "timestamp"
	case config.DbTypePostgres:
		return "timestamp with time zone"
	case config.DbTypeSqlserver:
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
	t := time.Time(this)
	//if t.IsZero() {
	//	return []byte("null"), nil
	//}
	return []byte(`"` + t.Format(time.DateTime) + `"`), nil
}
func (this *CustomTime) UnmarshalJSON(b []byte) (err error) {
	s := strings.Trim(string(b), "\"")
	if s == "null" || s == "" {
		//*this = CustomTime{}
		return nil
	}
	_t, err := time.ParseInLocation(time.DateTime, s, config.CstSh)
	if err != nil {
		_t, err = time.ParseInLocation(time.RFC3339Nano, s, config.CstSh)
	}
	*this = CustomTime(_t)
	return err
}

type DynamicTime struct {
	time.Time
	Format string
}

func (dt DynamicTime) MarshalJSON() ([]byte, error) {
	if dt.Time.IsZero() {
		return []byte(`null`), nil
	}
	return []byte(`"` + dt.Time.Format(dt.Format) + `"`), nil
}
