package dataType

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"io"
	"testing"
)

func TestGzipRawRoundTrip(t *testing.T) {
	raw := bytes.Repeat([]byte("hello gzip raw "), 100)
	g := NewGzipRaw(raw)

	// Value: 写入数据库时应为合法 gzip 压缩字节
	val, err := g.Value()
	if err != nil {
		t.Fatalf("Value error: %v", err)
	}
	compressed, ok := val.([]byte)
	if !ok {
		t.Fatalf("Value expected []byte, got %T", val)
	}
	if !isValidGzip(compressed) {
		t.Fatalf("Value output is not valid gzip data")
	}

	// Scan: 读取时原样返回压缩字节，不解压
	var dst GzipRaw
	if err := dst.Scan(compressed); err != nil {
		t.Fatalf("Scan error: %v", err)
	}
	if !bytes.Equal(dst.GetValue(), compressed) {
		t.Fatalf("GzipRaw.Scan should keep compressed bytes as-is, got %v want %v", dst.GetValue(), compressed)
	}
	if bytes.Equal(dst.GetValue(), raw) {
		t.Fatalf("GzipRaw should NOT decompress on Scan, but result equals raw data")
	}
}

func TestGzipAutoRoundTrip(t *testing.T) {
	raw := bytes.Repeat([]byte("hello gzip auto "), 100)
	g := NewGzipAuto(raw)

	val, err := g.Value()
	if err != nil {
		t.Fatalf("Value error: %v", err)
	}
	compressed, ok := val.([]byte)
	if !ok {
		t.Fatalf("Value expected []byte, got %T", val)
	}
	if !isValidGzip(compressed) {
		t.Fatalf("Value output is not valid gzip data")
	}

	// Scan: 读取时自动解压，还原原始数据
	var dst GzipAuto
	if err := dst.Scan(compressed); err != nil {
		t.Fatalf("Scan error: %v", err)
	}
	if !bytes.Equal(dst.GetValue(), raw) {
		t.Fatalf("GzipAuto.Scan should decompress to original, got %v want %v", dst.GetValue(), raw)
	}
}

func TestGzipNilHandling(t *testing.T) {
	var raw GzipRaw
	v, err := raw.Value()
	if err != nil {
		t.Fatalf("GzipRaw.Value nil error: %v", err)
	}
	if v != nil {
		t.Fatalf("GzipRaw nil Value expected nil, got %v", v)
	}
	if err := raw.Scan(nil); err != nil {
		t.Fatalf("GzipRaw.Scan nil error: %v", err)
	}
	if !raw.IsNil() {
		t.Fatalf("GzipRaw nil IsNil should be true")
	}

	var auto GzipAuto
	v, err = auto.Value()
	if err != nil {
		t.Fatalf("GzipAuto.Value nil error: %v", err)
	}
	if v != nil {
		t.Fatalf("GzipAuto nil Value expected nil, got %v", v)
	}
	if err := auto.Scan(nil); err != nil {
		t.Fatalf("GzipAuto.Scan nil error: %v", err)
	}
	if !auto.IsNil() {
		t.Fatalf("GzipAuto nil IsNil should be true")
	}
}

func TestGzipStringScan(t *testing.T) {
	// Scan 支持 string 类型输入
	var raw GzipRaw
	str := "not-gzipped-plain-string"
	if err := raw.Scan(str); err != nil {
		t.Fatalf("GzipRaw.Scan string error: %v", err)
	}
	if raw.GetValue() == nil || string(raw.GetValue()) != str {
		t.Fatalf("GzipRaw Scan string mismatch, got %v", raw.GetValue())
	}
}

func TestGzipAutoToleratesUncompressedData(t *testing.T) {
	// GzipAuto 读取到未压缩数据时应原样返回（容错）
	var auto GzipAuto
	str := []byte("plain uncompressed data")
	if err := auto.Scan(str); err != nil {
		t.Fatalf("GzipAuto.Scan plain error: %v", err)
	}
	if !bytes.Equal(auto.GetValue(), str) {
		t.Fatalf("GzipAuto should tolerate uncompressed data, got %v want %v", auto.GetValue(), str)
	}
}

func TestGzipJSONRoundTrip(t *testing.T) {
	raw := bytes.Repeat([]byte("json round trip "), 20)

	// GzipRaw JSON
	var rawG GzipRaw
	rawG = NewGzipRaw(raw)
	b, err := json.Marshal(rawG)
	if err != nil {
		t.Fatalf("GzipRaw MarshalJSON error: %v", err)
	}
	var rawG2 GzipRaw
	if err := json.Unmarshal(b, &rawG2); err != nil {
		t.Fatalf("GzipRaw UnmarshalJSON error: %v", err)
	}
	if !bytes.Equal(rawG2.GetValue(), raw) {
		t.Fatalf("GzipRaw JSON round trip mismatch")
	}

	// GzipAuto JSON
	var autoG GzipAuto
	autoG = NewGzipAuto(raw)
	b, err = json.Marshal(autoG)
	if err != nil {
		t.Fatalf("GzipAuto MarshalJSON error: %v", err)
	}
	var autoG2 GzipAuto
	if err := json.Unmarshal(b, &autoG2); err != nil {
		t.Fatalf("GzipAuto UnmarshalJSON error: %v", err)
	}
	if !bytes.Equal(autoG2.GetValue(), raw) {
		t.Fatalf("GzipAuto JSON round trip mismatch")
	}
}

func TestGzipJSONNull(t *testing.T) {
	// null 应还原为空值
	b, _ := json.Marshal(GzipRaw{})
	var raw GzipRaw
	if err := json.Unmarshal(b, &raw); err != nil {
		t.Fatalf("GzipRaw UnmarshalJSON null error: %v", err)
	}
	if !raw.IsNil() {
		t.Fatalf("GzipRaw JSON null should be nil")
	}
}

func TestGzipCompressLevel(t *testing.T) {
	raw := bytes.Repeat([]byte("compress level test data "), 200)

	// GzipRaw 设置压缩级别后往返仍正确
	var rawG GzipRaw
	rawG = NewGzipRaw(raw)
	rawG.SetCompressLevel(gzip.BestCompression)
	if rawG.GetCompressLevel() != gzip.BestCompression {
		t.Fatalf("GzipRaw.GetCompressLevel mismatch, got %d", rawG.GetCompressLevel())
	}
	val, err := rawG.Value()
	if err != nil {
		t.Fatalf("GzipRaw compressed Value error: %v", err)
	}
	if !isValidGzip(val.([]byte)) {
		t.Fatalf("GzipRaw compressed output invalid")
	}

	// GzipAuto 设置压缩级别后往返还原原始数据
	var autoG GzipAuto
	autoG = NewGzipAuto(raw)
	autoG.SetCompressLevel(gzip.BestSpeed)
	if autoG.GetCompressLevel() != gzip.BestSpeed {
		t.Fatalf("GzipAuto.GetCompressLevel mismatch, got %d", autoG.GetCompressLevel())
	}
	val, err = autoG.Value()
	if err != nil {
		t.Fatalf("GzipAuto compressed Value error: %v", err)
	}
	var back GzipAuto
	if err := back.Scan(val); err != nil {
		t.Fatalf("GzipAuto compressed Scan error: %v", err)
	}
	if !bytes.Equal(back.GetValue(), raw) {
		t.Fatalf("GzipAuto compressed round trip mismatch")
	}

	// 默认压缩级别
	d := NewGzipAuto(nil)
	if d.GetCompressLevel() != gzip.DefaultCompression {
		t.Fatalf("default compress level mismatch, got %d", d.GetCompressLevel())
	}
}

func TestGzipInvalidCompressLevel(t *testing.T) {
	raw := []byte("invalid level check")
	var autoG GzipAuto
	autoG = NewGzipAuto(raw)
	autoG.SetCompressLevel(12345) // 非法级别，应导致 NewWriterLevel 报错
	if _, err := autoG.Value(); err == nil {
		t.Fatalf("invalid compress level should error")
	}
}

func isValidGzip(data []byte) bool {
	r, err := gzip.NewReader(bytes.NewReader(data))
	if err != nil {
		return false
	}
	defer func() { _ = r.Close() }()
	out, err := io.ReadAll(r)
	return err == nil && len(out) > 0
}