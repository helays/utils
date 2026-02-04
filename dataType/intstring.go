package dataType

import (
	"bytes"
	"database/sql/driver"
	"encoding/gob"
	"encoding/json"
	"strconv"
	"strings"
	"sync"

	"golang.org/x/exp/constraints"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
	"helay.net/go/utils/v3/config"
	"helay.net/go/utils/v3/tools"
)

type IntString[T constraints.Integer] struct {
	value        T
	jsonAsString bool
	valid        bool // 用于区分零值和未设置值
}

// NewIntString 添加一些便利方法
func NewIntString[T constraints.Integer](v T) IntString[T] {
	return IntString[T]{
		value:        v,
		jsonAsString: true,
		valid:        true,
	}
}

func StringToIntString[T constraints.Integer](v string) (IntString[T], error) {
	if v == "" {
		return ZeroIntString[T](), config.ErrInvalidParam
	}
	dst, err := strconv.ParseInt(v, 10, 64)
	if err != nil {
		return ZeroIntString[T](), err
	}
	return NewIntString[T](T(dst)), nil
}

func NewIntStringAsNumber[T constraints.Integer](v T) IntString[T] {
	return IntString[T]{
		value:        v,
		jsonAsString: false,
		valid:        true,
	}
}

// ZeroIntString 零值（未设置值）
func ZeroIntString[T constraints.Integer]() IntString[T] {
	return IntString[T]{
		valid: false,
	}
}

// Value 注意，在写数据库的时候，貌似只支持 int64
// noinspection all
func (i IntString[T]) Value() (driver.Value, error) {
	if !i.valid {
		return nil, nil
	}
	return int64(i.value), nil
}

// noinspection all
func (i *IntString[T]) Scan(val any) error {
	i.jsonAsString = true // 从数据读取出来也默认为字符串
	if val == nil {
		i.value = 0
		i.valid = false
		return nil
	}
	i.valid = true
	v, err := tools.Any2Int[T](val)
	if err != nil {
		return err
	}
	i.value = T(v)
	return nil
}

// noinspection all
func (i IntString[T]) GormDBDataType(db *gorm.DB, field *schema.Field) string {
	var zero T
	switch db.Dialector.Name() {
	case config.DbTypeSqlite:
		// 使用类型开关检查具体类型
		switch any(zero).(type) {
		case int, int64, int32, int16, uint, uint64, uint32, uint16:
			return "INTEGER"
		case int8, byte:
			return "TINYINT"
		default:
			return "INTEGER"
		}
	case config.DbTypeMysql:
		switch any(zero).(type) {
		case int, int32, int16:
			return "INT"
		case int64:
			return "BIGINT"
		case uint, uint32, uint16:
			return "INT UNSIGNED"
		case uint64:
			return "BIGINT UNSIGNED"
		case int8, byte:
			return "TINYINT"
		default:
			return "INT"
		}
	case config.DbTypePostgres:
		switch any(zero).(type) {
		case int, int32, int16:
			return "INTEGER"
		case int64, uint, uint32, uint16, uint64:
			return "BIGINT" // PostgreSQL 没有 UNSIGNED
		case int8, byte:
			return "SMALLINT"
		default:
			return "INTEGER"
		}
	case config.DbTypeSqlserver:
		switch any(zero).(type) {
		case int, int32:
			return "INT"
		case int64:
			return "BIGINT"
		case uint, uint32:
			return "INT"
		case uint64:
			return "BIGINT"
		case int16:
			return "SMALLINT"
		case int8, byte:
			return "TINYINT"
		default:
			return "INT"
		}
	}
	return "INT"
}

// SetJsonAsString 或者使用方法设置
// noinspection all
func (i *IntString[T]) SetJsonAsString(asString bool) {
	i.jsonAsString = asString
}

func (i IntString[T]) GetJsonAsString() bool {
	return i.jsonAsString
}

// noinspection all
func (i IntString[T]) MarshalJSON() ([]byte, error) {
	if !i.valid {
		return []byte("null"), nil
	}
	if i.jsonAsString {
		v := tools.Any2string(i.value)
		return []byte(`"` + v + `"`), nil
	}
	return json.Marshal(i.value)

}

// noinspection all
func (i *IntString[T]) UnmarshalJSON(data []byte) error {
	if data == nil || len(data) == 0 {
		i.valid = false
		return nil
	}
	ds := strings.Trim(string(data), `"`)
	if ds == "" || ds == "null" || ds == "nil" || ds == "undefined" {
		i.valid = false
		return nil
	}
	i.jsonAsString = false
	rawDs := string(data)
	if len(rawDs) > 0 && rawDs[0] == '"' {
		i.jsonAsString = true
	}

	dst, err := tools.Any2Int[T](ds)
	if err != nil {
		return err
	}
	i.valid = true
	i.value = T(dst)
	return nil
}

// UnmarshalText 实现 TextMarshaler 接口
// 这个主要是从query参数中进行提取，所以这里统一处理为字符串
// noinspection all
func (i *IntString[T]) UnmarshalText(text []byte) error {
	err := i.UnmarshalJSON(text)
	if err != nil {
		return err
	}
	i.jsonAsString = true
	return nil
}

// noinspection all
func (i *IntString[T]) SetValue(v T) {
	i.value = v
	i.valid = true
}

// noinspection all
func (i IntString[T]) GetValue() T {
	if !i.valid {
		var zero T
		return zero
	}
	return i.value
}

// 实现 Stringer 接口
// noinspection all
func (i IntString[T]) String() string {
	if !i.valid {
		return ""
	}
	return tools.Any2string(i.value)
}

// Equals 实现 Equals 方法用于比较
// noinspection all
func (i IntString[T]) Equals(other IntString[T]) bool {
	if !i.valid && !other.valid {
		return true // 两个都是无效值，认为相等
	}
	if i.valid != other.valid {
		return false // 一个有效一个无效，不相等
	}
	return i.value == other.value
}

// IsZero 添加零值检查
// noinspection all
func (i IntString[T]) IsZero() bool {
	if !i.valid {
		return true // 未设置的值被认为是零值
	}
	var zero T
	return i.value == zero
}

// IsValid 检查值是否有效
// noinspection all
func (i IntString[T]) IsValid() bool {
	return i.valid
}

// 全局缓存编码器/解码器
var gobBufferPool = &sync.Pool{New: func() any { return new(bytes.Buffer) }}

// gobIntString 用于 Gob 序列化的内部结构
type gobIntString[T constraints.Integer] struct {
	Value        T
	JsonAsString bool
	Valid        bool
}

// GobEncode
// noinspection all
func (i IntString[T]) GobEncode() ([]byte, error) {
	// 从池中获取 buffer
	buf := gobBufferPool.Get().(*bytes.Buffer)
	buf.Reset()
	defer gobBufferPool.Put(buf) // 用完后放回池中

	// 创建一次编码器（避免重复分配）
	enc := gob.NewEncoder(buf)

	// 一次编码整个结构体
	if err := enc.Encode(gobIntString[T]{
		Value:        i.value,
		JsonAsString: i.jsonAsString,
		Valid:        i.valid,
	}); err != nil {
		return nil, err
	}
	// 返回数据的副本（buf 会被放回池中重用）
	return bytes.Clone(buf.Bytes()), nil
}

// noinspection all
func (i *IntString[T]) GobDecode(data []byte) error {
	var temp gobIntString[T]

	// 使用传入的数据创建 reader
	dec := gob.NewDecoder(bytes.NewReader(data))
	if err := dec.Decode(&temp); err != nil {
		return err
	}

	// 赋值给接收者
	i.value = temp.Value
	i.jsonAsString = temp.JsonAsString
	i.valid = temp.Valid

	return nil
}

func (i *IntString[T]) GetValid() bool {
	return i.valid
}
