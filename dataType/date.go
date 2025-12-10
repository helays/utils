package dataType

import (
	"database/sql"
	"database/sql/driver"
	"strings"
	"time"

	"github.com/helays/utils/v2/config"
	"github.com/helays/utils/v2/rule-engine/formatter"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

type CustomDate time.Time

// noinspection all
func (c CustomDate) String() string {
	return time.Time(c).Format(time.DateOnly)
}

// noinspection all
func (c CustomDate) Format(layout string) string {
	return time.Time(c).Format(layout)
}

// noinspection all
func (c CustomDate) After(u time.Time) bool {
	return time.Time(c).After(u)
}

// noinspection all
func (c CustomDate) Before(u time.Time) bool {
	return time.Time(c).Before(u)
}

// noinspection all
func (c CustomDate) Sub(u time.Time) time.Duration {
	return time.Time(c).Sub(u)
}

// noinspection all
func (c CustomDate) Unix() int64 {
	return time.Time(c).Unix()
}

// noinspection all
func (c *CustomDate) Scan(value interface{}) (err error) {
	nullTime := &sql.NullTime{}
	err = nullTime.Scan(value)
	*c = CustomDate(nullTime.Time)
	return
}

// noinspection all
func (c CustomDate) Value() (driver.Value, error) {
	return time.Time(c), nil
}

// noinspection all
func (c CustomDate) GormDataType() string {
	return "time"
}

// noinspection all
func (CustomDate) GormDBDataType(db *gorm.DB, field *schema.Field) string {
	switch db.Dialector.Name() {
	case config.DbTypeSqlite:
		return "date"
	case config.DbTypeMysql:
		return "date"
	case config.DbTypePostgres:
		return "date"
	case config.DbTypeSqlserver:
		return "date"
	}
	return ""
}

// noinspection all
func (c CustomDate) GobEncode() ([]byte, error) {
	return time.Time(c).GobEncode()
}

// noinspection all
func (c *CustomDate) GobDecode(b []byte) error {
	return (*time.Time)(c).GobDecode(b)
}

// noinspection all
func (c CustomDate) MarshalJSON() ([]byte, error) {
	t := time.Time(c)
	//if t.IsZero() {
	//	return []byte("null"), nil
	//}
	return []byte(`"` + t.Format(time.DateOnly) + `"`), nil
}

// noinspection all
func (c *CustomDate) UnmarshalJSON(b []byte) (err error) {
	s := strings.Trim(string(b), "\"")
	if s == "null" || s == "" {
		//*this = CustomTime{}
		return nil
	}
	tf := formatter.FormatRule[time.Time]{FormatType: "output_date"}
	_t, err := tf.Format(s)
	if err != nil {
		return err
	}
	*c = CustomDate(_t)
	return err
}

type CustomTime time.Time

func NewCustomTimeNow() CustomTime {
	return CustomTime(time.Now())
}

func NewCustomTime(t time.Time) CustomTime {
	return CustomTime(t)
}

// noinspection all
func (c CustomTime) String() string {
	return time.Time(c).Format(time.DateTime)
}

// noinspection all
func (c CustomTime) Format(layout string) string {
	return time.Time(c).Format(layout)
}

// noinspection all
func (c CustomTime) After(u time.Time) bool {
	return time.Time(c).After(u)
}

// noinspection all
func (c CustomTime) Before(u time.Time) bool {
	return time.Time(c).Before(u)
}

// noinspection all
func (c CustomTime) Sub(u time.Time) time.Duration {
	return time.Time(c).Sub(u)
}

// noinspection all
func (c CustomTime) Unix() int64 {
	return time.Time(c).Unix()
}

// noinspection all
func (c *CustomTime) Scan(value interface{}) (err error) {
	nullTime := &sql.NullTime{}
	err = nullTime.Scan(value)
	*c = CustomTime(nullTime.Time)
	return
}

// noinspection all
func (c CustomTime) Value() (driver.Value, error) {
	return time.Time(c).Format(time.DateTime), nil
}

// GormDataType 这个时间自定义字段，只能用time类型，不然当表中设置有 autoUpdateTime控制属性时会有问题。
// 使用Update方法，只有DataType是time的时候，才会生成time.Time类型，
// 其他类型都会转成时间戳。
// noinspection all
func (c CustomTime) GormDataType() string {
	return "time"
}

// GormDBDataType gorm db data type
// noinspection all
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

// noinspection all
func (c CustomTime) GobEncode() ([]byte, error) {
	return time.Time(c).GobEncode()
}

// noinspection all
func (c *CustomTime) GobDecode(b []byte) error {
	return (*time.Time)(c).GobDecode(b)
}

// MarshalJSON 序列化至json字符串
// noinspection all
func (c CustomTime) MarshalJSON() ([]byte, error) {
	t := time.Time(c)
	return []byte(`"` + t.Format(time.DateTime) + `"`), nil
}

// UnmarshalJSON 反序列化
// noinspection all
func (c *CustomTime) UnmarshalJSON(b []byte) (err error) {
	s := strings.Trim(string(b), "\"")
	if s == "null" || s == "" {
		return nil
	}
	tf := formatter.FormatRule[time.Time]{FormatType: "output_date"}
	_t, err := tf.Format(s)
	if err != nil {
		return err
	}
	*c = CustomTime(_t)
	return err
}

// noinspection all
func (c CustomTime) ToPtr() *CustomTime {
	return &c
}

// noinspection all
func (c CustomTime) ToTime() time.Time {
	return time.Time(c)
}

// noinspection all
func (c CustomTime) IsZero() bool {
	return c.ToTime().IsZero()
}

// noinspection all
func (c CustomTime) Equal(t time.Time) bool {
	return c.ToTime().Equal(t)
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
