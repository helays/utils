package middleware

import "time"

const DebugCtxKey = "route-debug"

// CompressionAlgorithm 压缩算法类型
type CompressionAlgorithm string

func (c CompressionAlgorithm) String() string {
	return string(c)
}

// 支持的压缩算法常量
const (
	Gzip    CompressionAlgorithm = "gzip"    // 最广泛兼容
	Deflate CompressionAlgorithm = "deflate" // 基本支持
	Brotli  CompressionAlgorithm = "br"      // 最佳压缩比和性能平衡
)

var supportedCompressionAlgorithms = map[CompressionAlgorithm]struct{}{
	Gzip:    struct{}{},
	Deflate: struct{}{},
}

const acceptEncoding string = "Accept-Encoding"

// 需要排除的压缩类型
var excludeContentTypes = map[string]struct{}{
	"image/": {},
	"video/": {},
	"audio/": {},
	"font/":  {},

	"application/zip":         {},
	"application/gzip":        {},
	"application/x-gzip":      {},
	"application/x-bzip2":     {},
	"application/x-tar":       {},
	"application/pdf":         {},
	"application/x-font-woff": {},
	// 添加更多需要排除的压缩类型
}

type (
	// Logs 日志内容结构体
	Logs struct {
		Status      int           // 响应状态码
		Ip          string        // 请求 IP
		Method      string        // 请求方法
		Uri         string        // 请求 URI
		ContentSize int64         // 响应内容大小
		UserAgent   string        // 用户代理
		Elapsed     time.Duration // 响应耗时,微秒级别
	}
	// Logger 日志接口
	Logger interface {
		Write(l *Logs)
	}
)
