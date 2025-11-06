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

func (this CustomDate) String() string {
	return time.Time(this).Format(time.DateOnly)
}

func (this CustomDate) Format(layout string) string {
	return time.Time(this).Format(layout)
}

func (this CustomDate) After(u time.Time) bool {
	return time.Time(this).After(u)
}

func (this CustomDate) Before(u time.Time) bool {
	return time.Time(this).Before(u)
}

func (this CustomDate) Sub(u time.Time) time.Duration {
	return time.Time(this).Sub(u)
}

func (this CustomDate) Unix() int64 {
	return time.Time(this).Unix()
}

func (this *CustomDate) Scan(value interface{}) (err error) {
	nullTime := &sql.NullTime{}
	err = nullTime.Scan(value)
	*this = CustomDate(nullTime.Time)
	return
}

func (this CustomDate) Value() (driver.Value, error) {
	return time.Time(this), nil
}

func (this CustomDate) GormDataType() string {
	return "time"
}

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
func (this CustomDate) GobEncode() ([]byte, error) {
	return time.Time(this).GobEncode()
}
func (this *CustomDate) GobDecode(b []byte) error {
	return (*time.Time)(this).GobDecode(b)
}
func (this CustomDate) MarshalJSON() ([]byte, error) {
	t := time.Time(this)
	//if t.IsZero() {
	//	return []byte("null"), nil
	//}
	return []byte(`"` + t.Format(time.DateOnly) + `"`), nil
}
func (this *CustomDate) UnmarshalJSON(b []byte) (err error) {
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
	*this = CustomDate(_t)
	return err
}

type CustomTime time.Time

func NewCustomTimeNow() CustomTime {
	return CustomTime(time.Now())
}

func NewCustomTime(t time.Time) CustomTime {
	return CustomTime(t)
}

func (c CustomTime) String() string {
	return time.Time(c).Format(time.DateTime)
}

func (c CustomTime) Format(layout string) string {
	return time.Time(c).Format(layout)
}

func (c CustomTime) After(u time.Time) bool {
	return time.Time(c).After(u)
}

func (c CustomTime) Before(u time.Time) bool {
	return time.Time(c).Before(u)
}

func (c CustomTime) Sub(u time.Time) time.Duration {
	return time.Time(c).Sub(u)
}

func (c CustomTime) Unix() int64 {
	return time.Time(c).Unix()
}

func (c *CustomTime) Scan(value interface{}) (err error) {
	nullTime := &sql.NullTime{}
	err = nullTime.Scan(value)
	*c = CustomTime(nullTime.Time)
	return
}

func (c CustomTime) Value() (driver.Value, error) {
	return time.Time(c).Format(time.DateTime), nil
}

// GormDataType 这个时间自定义字段，只能用time类型，不然当表中设置有 autoUpdateTime控制属性时会有问题。
// 使用Update方法，只有DataType是time的时候，才会生成time.Time类型，
// 其他类型都会转成时间戳。
func (c CustomTime) GormDataType() string {
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

func (c CustomTime) GobEncode() ([]byte, error) {
	return time.Time(c).GobEncode()
}

func (c *CustomTime) GobDecode(b []byte) error {
	return (*time.Time)(c).GobDecode(b)
}

// MarshalJSON 序列化至json字符串
func (c CustomTime) MarshalJSON() ([]byte, error) {
	t := time.Time(c)
	return []byte(`"` + t.Format(time.DateTime) + `"`), nil
}

// UnmarshalJSON 反序列化
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

func (c CustomTime) ToPtr() *CustomTime {
	return &c
}

func (c CustomTime) ToTime() time.Time {
	return time.Time(c)
}

func (c CustomTime) IsZero() bool {
	return c.ToTime().IsZero()
}

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
