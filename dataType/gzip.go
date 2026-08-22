package dataType

import (
	"bytes"
	"compress/gzip"
	"database/sql/driver"
	"encoding/base64"
	"encoding/json"
	"io"

	"gorm.io/gorm"
	"gorm.io/gorm/schema"
	"helay.net/go/utils/v3/tools"
)

// GzipRaw 定义一个 Gzip 存储类型：写入数据库时对原始数据自动 Gzip 压缩，
// 读取（Scan）时原样返回数据库中的压缩字节，不解压。
// 底层列类型由 BlobDbDataType 自动判断。
// noinspection all
type GzipRaw struct {
	data          []byte
	compressLevel int
}

// NewGzipRaw 构造一个 GzipRaw，data 为待压缩的原始字节
func NewGzipRaw(data []byte) GzipRaw {
	return GzipRaw{
		data:          data,
		compressLevel: gzip.DefaultCompression,
	}
}

// Value 实现了 driver.Valuer 接口，写入数据库时将原始数据 Gzip 压缩
// noinspection all
func (g GzipRaw) Value() (driver.Value, error) {
	if g.data == nil {
		return nil, nil
	}
	return gzipCompressLevel(g.data, g.compressLevel)
}

// Scan 实现了 sql.Scanner 接口，读取数据库时原样保留压缩字节，不解压
// noinspection all
func (g *GzipRaw) Scan(val any) error {
	if val == nil {
		g.data = nil
		return nil
	}
	b, err := tools.Any2bytes(val)
	if err != nil {
		return err
	}
	g.data = b
	return nil
}

// GormDataType gorm common data type
// noinspection all
func (GzipRaw) GormDataType() string {
	return "blob"
}

// GormDBDataType gorm db data type，由 BlobDbDataType 自动判断
// noinspection all
func (GzipRaw) GormDBDataType(db *gorm.DB, field *schema.Field) string {
	return BlobDbDataType(db, field)
}

// SetCompressLevel 设置压缩级别，取值同 compress/gzip（如 gzip.BestSpeed、
// gzip.DefaultCompression、gzip.BestCompression）。0 表示使用默认压缩级别。
// noinspection all
func (g *GzipRaw) SetCompressLevel(level int) {
	g.compressLevel = level
}

// GetCompressLevel 返回当前压缩级别
// noinspection all
func (g GzipRaw) GetCompressLevel() int {
	return g.compressLevel
}

// SetValue 设置原始数据
// noinspection all
func (g *GzipRaw) SetValue(v []byte) {
	g.data = v
}

// GetValue 返回原始数据（对于 GzipRaw，即调用方设置/扫描进来的字节）
// noinspection all
func (g GzipRaw) GetValue() []byte {
	return g.data
}

// IsNil 判断是否为空值
// noinspection all
func (g GzipRaw) IsNil() bool {
	return g.data == nil
}

// MarshalJSON 以 base64 输出原始数据
// noinspection all
func (g GzipRaw) MarshalJSON() ([]byte, error) {
	if g.data == nil {
		return []byte("null"), nil
	}
	return json.Marshal(base64.StdEncoding.EncodeToString(g.data))
}

// UnmarshalJSON 从 base64 还原原始数据
// noinspection all
func (g *GzipRaw) UnmarshalJSON(b []byte) error {
	if b == nil || len(b) == 0 || string(b) == "null" {
		g.data = nil
		return nil
	}
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}
	data, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		return err
	}
	g.data = data
	return nil
}

// GzipAuto 定义一个 Gzip 存储类型：写入数据库时对原始数据自动 Gzip 压缩，
// 读取（Scan）时自动解压，还原出压缩前的原始数据。
// 底层列类型由 BlobDbDataType 自动判断。
// noinspection all
type GzipAuto struct {
	data          []byte
	compressLevel int
}

// NewGzipAuto 构造一个 GzipAuto，data 为待压缩的原始字节
func NewGzipAuto(data []byte) GzipAuto {
	return GzipAuto{
		data:          data,
		compressLevel: gzip.DefaultCompression,
	}
}

// Value 实现了 driver.Valuer 接口，写入数据库时将原始数据 Gzip 压缩
// noinspection all
func (g GzipAuto) Value() (driver.Value, error) {
	if g.data == nil {
		return nil, nil
	}
	return gzipCompressLevel(g.data, g.compressLevel)
}

// Scan 实现了 sql.Scanner 接口，读取数据库时自动解压还原原始数据
// noinspection all
func (g *GzipAuto) Scan(val any) error {
	if val == nil {
		g.data = nil
		return nil
	}
	b, err := tools.Any2bytes(val)
	if err != nil {
		return err
	}
	decoded, err := gzipDecompress(b)
	if err != nil {
		return err
	}
	g.data = decoded
	return nil
}

// GormDataType gorm common data type
// noinspection all
func (GzipAuto) GormDataType() string {
	return "blob"
}

// GormDBDataType gorm db data type，由 BlobDbDataType 自动判断
// noinspection all
func (GzipAuto) GormDBDataType(db *gorm.DB, field *schema.Field) string {
	return BlobDbDataType(db, field)
}

// SetCompressLevel 设置压缩级别，取值同 compress/gzip（如 gzip.BestSpeed、
// gzip.DefaultCompression、gzip.BestCompression）。0 表示使用默认压缩级别。
// noinspection all
func (g *GzipAuto) SetCompressLevel(level int) {
	g.compressLevel = level
}

// GetCompressLevel 返回当前压缩级别
// noinspection all
func (g GzipAuto) GetCompressLevel() int {
	return g.compressLevel
}

// SetValue 设置原始数据
// noinspection all
func (g *GzipAuto) SetValue(v []byte) {
	g.data = v
}

// GetValue 返回解压后的原始数据
// noinspection all
func (g GzipAuto) GetValue() []byte {
	return g.data
}

// IsNil 判断是否为空值
// noinspection all
func (g GzipAuto) IsNil() bool {
	return g.data == nil
}

// MarshalJSON 以 base64 输出原始数据
// noinspection all
func (g GzipAuto) MarshalJSON() ([]byte, error) {
	if g.data == nil {
		return []byte("null"), nil
	}
	return json.Marshal(base64.StdEncoding.EncodeToString(g.data))
}

// UnmarshalJSON 从 base64 还原原始数据
// noinspection all
func (g *GzipAuto) UnmarshalJSON(b []byte) error {
	if b == nil || len(b) == 0 || string(b) == "null" {
		g.data = nil
		return nil
	}
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}
	data, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		return err
	}
	g.data = data
	return nil
}

// 使用标准库 compress/gzip 压缩数据，level 为压缩级别
func gzipCompressLevel(data []byte, level int) ([]byte, error) {
	var buf bytes.Buffer
	zw, err := gzip.NewWriterLevel(&buf, level)
	if err != nil {
		return nil, err
	}
	if _, err = zw.Write(data); err != nil {
		return nil, err
	}
	if err = zw.Close(); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// 解压 gzip 数据；若数据不是合法 gzip（兼容旧数据/未压缩数据），
// 则回退返回原字节且不报错。
func gzipDecompress(data []byte) ([]byte, error) {
	if len(data) == 0 {
		return data, nil
	}
	r, err := gzip.NewReader(bytes.NewReader(data))
	if err != nil {
		// 兼容未压缩数据，原样返回
		return data, nil
	}
	defer func() {
		_ = r.Close()
	}()
	out, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}
	return out, nil
}
