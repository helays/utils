package dataType

import (
	"database/sql/driver"
	"fmt"
	"strings"
	"time"

	"github.com/helays/utils/v2/config"
	"github.com/helays/utils/v2/rule-engine/formatter"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

type Timestamp struct {
	timestamp time.Time
	format    string
}

func NewTimestamp(t time.Time, formats ...string) *Timestamp {
	format := time.DateTime
	if len(formats) > 0 {
		format = formats[0]
	}
	return &Timestamp{
		timestamp: t,
		format:    format,
	}
}

// noinspection all
func (t *Timestamp) Scan(value interface{}) error {
	if value == nil {
		t.timestamp = time.Time{}
		return nil
	}
	switch v := value.(type) {
	case int64:
		// 假设存储的是微秒级时间戳
		t.timestamp = time.Unix(v/1e6, (v%1e6)*1e3)
		return nil
	case int32:
		// 假设存储的是秒级时间戳
		t.timestamp = time.Unix(int64(v), 0)
		return nil
	case int:
		// 假设存储的是秒级时间戳
		t.timestamp = time.Unix(int64(v), 0)
		return nil
	case time.Time:
		t.timestamp = v
		return nil
	case []byte:
		// 尝试解析为时间戳数字
		var ts int64
		_, err := fmt.Sscanf(string(v), "%d", &ts)
		if err == nil {
			t.timestamp = time.Unix(ts/1e6, (ts%1e6)*1e3)
			return nil
		}
		// 如果解析数字失败，尝试解析为时间字符串
		tf := formatter.FormatRule[time.Time]{FormatType: "output_date"}
		parsedTime, err := tf.Format(v)
		if err != nil {
			return fmt.Errorf("无法解析时间戳: %v", err)
		}
		t.timestamp = parsedTime
		return nil
	case string:
		// 尝试解析为时间戳数字
		var ts int64
		_, err := fmt.Sscanf(v, "%d", &ts)
		if err == nil {
			t.timestamp = time.Unix(ts/1e6, (ts%1e6)*1e3)
			return nil
		}

		// 如果解析数字失败，尝试解析为时间字符串
		tf := formatter.FormatRule[time.Time]{FormatType: "output_date"}
		parsedTime, err := tf.Format(v)
		if err != nil {
			return fmt.Errorf("无法解析时间戳: %v", err)
		}
		t.timestamp = parsedTime
		return nil
	default:
		return fmt.Errorf("不支持的扫描类型: %T", value)
	}
}

// Value 实现 driver.Valuer 接口，将时间 转成微秒级时间戳
// noinspection all
func (t Timestamp) Value() (driver.Value, error) {
	if t.timestamp.IsZero() {
		return nil, nil
	}
	return t.timestamp.UnixMicro(), nil
}

// GormDataType gorm db data type
// noinspection all
func (t Timestamp) GormDataType() string {
	return "timestamp"
}

// GormDBDataType gorm db data type
// noinspection all
func (t Timestamp) GormDBDataType(db *gorm.DB, field *schema.Field) string {
	switch db.Dialector.Name() {
	case config.DbTypeSqlite:
		return "INTEGER"
	case config.DbTypeMysql:
		return "BIGINT"
	case config.DbTypePostgres:
		return "BIGINT"
	case config.DbTypeSqlserver:
		return "BIGINT"
	}
	return "BIGINT"
}

// Time 返回内部的 time.Time
// noinspection all
func (t *Timestamp) Time() time.Time {
	return t.timestamp
}

// String 实现 Stringer 接口
// noinspection all
func (t *Timestamp) String() string {
	if t.timestamp.IsZero() {
		return ""
	}
	return t.timestamp.Format(t.format)
}

// Format 使用指定格式格式化时间
// noinspection all
func (t *Timestamp) Format(layout string) string {
	if t.timestamp.IsZero() {
		return ""
	}
	return t.timestamp.Format(layout)
}

// Unix 返回秒级时间戳
// noinspection all
func (t *Timestamp) Unix() int64 {
	return t.timestamp.Unix()
}

// UnixMilli 返回毫秒级时间戳
// noinspection all
func (t *Timestamp) UnixMilli() int64 {
	return t.timestamp.UnixMilli()
}

// UnixMicro 返回微秒级时间戳
// noinspection all
func (t *Timestamp) UnixMicro() int64 {
	return t.timestamp.UnixMicro()
}

// UnixNano 返回纳秒级时间戳
// noinspection all
func (t *Timestamp) UnixNano() int64 {
	return t.timestamp.UnixNano()
}

// IsZero 判断时间是否为零值
// noinspection all
func (t *Timestamp) IsZero() bool {
	return t.timestamp.IsZero()
}

// SetTime 设置时间
// noinspection all
func (t *Timestamp) SetTime(time time.Time) {
	t.timestamp = time
}

// SetFormat 设置格式化字符串
// noinspection all
func (t *Timestamp) SetFormat(format string) {
	t.format = format
}

// MarshalJSON 实现 JSON 序列化接口
// noinspection all
func (t *Timestamp) MarshalJSON() ([]byte, error) {
	if t.timestamp.IsZero() {
		return []byte("null"), nil
	}
	return []byte(`"` + t.String() + `"`), nil
}

// UnmarshalJSON 实现 JSON 反序列化接口
// noinspection all
func (t *Timestamp) UnmarshalJSON(data []byte) error {
	s := strings.Trim(string(data), "\"")
	if s == "null" || s == "" {
		return nil
	}
	tf := formatter.FormatRule[time.Time]{FormatType: "output_date"}
	_t, err := tf.Format(s)
	if err != nil {
		return err
	}
	t.timestamp = _t
	return err
}
