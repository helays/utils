package dataType

import (
	"bytes"
	"database/sql/driver"
	"encoding/json"
	"strconv"
	"strings"

	"gorm.io/gorm"
	"gorm.io/gorm/schema"
	"helay.net/go/utils/v3/config"
	"helay.net/go/utils/v3/tools"
)

// Bool 注意当使用这个类型时，在定义模型时，默认值需要带上括号。不然pg数据库会报错。
type Bool struct {
	bool
}

func NewBool(b bool) Bool {
	return Bool{bool: b}
}

// noinspection all
func (b Bool) Value() (driver.Value, error) {
	if b.bool {
		return int64(1), nil
	}
	return int64(0), nil
}

// noinspection all
func (b *Bool) Scan(value any) error {
	if value == nil {
		b.bool = false
		return nil
	}
	ok, e := tools.Any2bool(value)
	if e != nil {
		return e
	}
	b.bool = ok
	return nil
}

// noinspection all
func (b Bool) GormDataType() string {
	return "int"
}

// noinspection all
func (Bool) GormDBDataType(db *gorm.DB, _ *schema.Field) string {
	switch db.Dialector.Name() {
	case config.DbTypeSqlite:
		return "integer"
	case config.DbTypeMysql:
		return "tinyint(1)"
	case config.DbTypePostgres:
		return "int2"
	case config.DbTypeSqlserver:
		return "bit"
	}
	return "int"
}

// noinspection all
func (b Bool) Bool() bool {
	return b.bool
}

// noinspection all
func (b Bool) Int() int {
	if b.bool {
		return 1
	}
	return 0
}

// Resverse 反转
// noinspection all
func (b *Bool) Resverse() {
	b.bool = !b.bool
}

func (b *Bool) Set(b2 Bool) {
	b.bool = b2.bool
}

func (b *Bool) SetBool(b2 bool) {
	b.bool = b2
}

func (b *Bool) SetInt(i int) {
	b.bool = i != 0
}

func (b *Bool) SetString(s string) {
	b.bool = s != ""
}

func (b *Bool) Equals(b2 Bool) bool {
	return b.bool == b2.bool
}

func (b Bool) MarshalJSON() ([]byte, error) {
	return json.Marshal(b.bool)
}

func (b *Bool) UnmarshalJSON(data []byte) error {
	ds := bytes.TrimSpace(data)
	if len(ds) == 0 || bytes.Equal(ds, []byte("null")) {
		b.bool = false
		return nil
	}
	switch ds[0] {
	case '"': // JSON 字符串
		s, err := strconv.Unquote(string(ds))
		if err != nil {
			return err
		}
		ok, err := tools.Any2bool(s)
		if err != nil {
			return err
		}
		b.bool = ok
	case 't': // true
		b.bool = true
	case 'f': // false
		b.bool = false
	default: // JSON 数字
		f, err := strconv.ParseFloat(string(ds), 64)
		if err != nil {
			return err
		}
		ok, err := tools.Any2bool(f)
		if err != nil {
			return err
		}
		b.bool = ok
	}
	return nil
}

func (b Bool) GobEncode() ([]byte, error) {
	// bool类型可以直接用1个字节表示
	if b.bool {
		return []byte{1}, nil
	}
	return []byte{0}, nil
}

func (b *Bool) GobDecode(data []byte) error {
	if len(data) == 0 {
		b.bool = false
		return nil
	}
	// 任何非零值都视为 true
	b.bool = data[0] != 0
	return nil
}

func (b Bool) ToPtr() *Bool {
	return &b
}

// UnmarshalParam 实现 gin binding.BindUnmarshaler（query/form 参数绑定）。
//
// 目的：让 query/form 里的布尔原生解析（"true"/"false"/"1"/"0"/"yes"/"no"…），
// 不再落到 gin 的 Struct 分支 json.Unmarshal（其遇空串会报 unexpected end of JSON input）。
//
// 注意：gin 对**显式空值**（?x=）仍判定为"已设置"（trySetCustom 恒返回 isSet=true），
// 因此 *Bool（指针）字段会被赋成零值 false；若某筛选参数要求"空 = 未传 = nil"，
// 请由前端不发送空串（现状如此），或该 DTO 改用值类型 dataType.Bool。
func (b *Bool) UnmarshalParam(src string) error {
	s := strings.TrimSpace(src)
	if s == "" {
		return nil
	}
	ok, err := tools.Any2bool(s)
	if err != nil {
		return err
	}
	b.bool = ok
	return nil
}

// UnmarshalText 兼容 encoding.TextUnmarshaler（部分三方解码器/框架偏好它）
func (b *Bool) UnmarshalText(text []byte) error { return b.UnmarshalParam(string(text)) }
