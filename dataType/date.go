package dataType

import (
	"database/sql"
	"database/sql/driver"
	"strings"
	"time"

	"github.com/helays/utils/v2/config"
	"github.com/helays/utils/v2/rule-engine/formatter"
	"github.com/helays/utils/v2/tools"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

type CustomDate struct {
	time.Time
}

// noinspection all
func (c CustomDate) String() string {
	return c.Format(time.DateOnly)
}

// noinspection all
func (c *CustomDate) Scan(value interface{}) (err error) {
	nullTime := &sql.NullTime{}
	err = nullTime.Scan(value)
	c.Time = nullTime.Time
	return
}

// noinspection all
func (c CustomDate) Value() (driver.Value, error) {
	return c.Time, nil
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
func (c CustomDate) MarshalJSON() ([]byte, error) {
	//if t.IsZero() {
	//	return []byte("null"), nil
	//}
	return []byte(`"` + c.Format(time.DateOnly) + `"`), nil
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
	c.Time = _t
	return err
}

type CustomTime struct {
	time.Time
}

func NewCustomTimeNow() CustomTime {
	return CustomTime{Time: time.Now()}
}

func NewCustomTime(t time.Time) CustomTime {
	return CustomTime{Time: t}
}

// noinspection all
func (c CustomTime) String() string {
	return c.Format(time.DateTime)
}

// noinspection all
func (c *CustomTime) Scan(value interface{}) (err error) {
	nullTime := &sql.NullTime{}
	err = nullTime.Scan(value)
	c.Time = nullTime.Time
	return
}

// noinspection all
func (c CustomTime) Value() (driver.Value, error) {
	return c.Format(time.DateTime), nil
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

// MarshalJSON 序列化至json字符串
// noinspection all
func (c CustomTime) MarshalJSON() ([]byte, error) {
	return []byte(`"` + c.Format(time.DateTime) + `"`), nil
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
	c.Time = _t
	return err
}

// noinspection all
func (c CustomTime) ToPtr() *CustomTime {
	return &c
}

// noinspection all
func (c CustomTime) ToTime() time.Time {
	return c.Time
}

// noinspection all
func (c CustomTime) IsZero() bool {
	return c.Time.IsZero()
}

// AdjustTimezoneIfNeeded 调整时区
// noinspection all
func (c *CustomTime) AdjustTimezoneIfNeeded() {
	c.Time = tools.AdjustTimezoneIfNeeded(c.Time)
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
