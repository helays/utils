package middleware

import (
	"bufio"
	"compress/flate"
	"compress/gzip"
	"errors"
	"io"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/helays/utils/v2/close/vclose"
	"github.com/helays/utils/v2/logger/zaploger"
)

type CompressionConfig struct {
	Enabled             bool     `json:"enabled" yaml:"enabled" ini:"enabled"`                                           // 是否启用压缩
	Level               int      `json:"level" yaml:"level" ini:"level"`                                                 // 压缩级别
	ExcludeContentTypes []string `json:"exclude_content_types" yaml:"exclude_content_types" ini:"exclude_content_types"` // 不需要压缩的 MIME 类型
}

type ResponseProcessor struct {
	compressOpt CompressionConfig
	loggerOpt   zaploger.Config
	logEvent    []Logger // 日志事件
}

func NewResponseProcessor() *ResponseProcessor {
	c := &ResponseProcessor{}
	return c
}

// SetCompressionConfig 设置压缩配置
func (c *ResponseProcessor) SetCompressionConfig(opt CompressionConfig) {
	c.compressOpt = opt
	if c.compressOpt.Enabled {
		if c.compressOpt.Level < gzip.DefaultCompression || c.compressOpt.Level > gzip.BestCompression {
			c.compressOpt.Level = gzip.DefaultCompression
		}
		for _, contentType := range c.compressOpt.ExcludeContentTypes {
			excludeContentTypes[strings.ToLower(contentType)] = struct{}{}
		}
	}
}

// AddLogHandler 追加日志事件处理器
func (c *ResponseProcessor) AddLogHandler(le ...Logger) {
	for _, lg := range le {
		if lg != nil {
			c.logEvent = append(c.logEvent, lg)
		}
	}
}

func (c *ResponseProcessor) Handler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 创建一个响应处理程序
		enhancedWriter := &writer{w: w, r: r, status: http.StatusOK, createAt: time.Now()}
		defer c.metrics(enhancedWriter)
		if !c.compressOpt.Enabled {
			next.ServeHTTP(enhancedWriter, r)
			return
		}
		w.Header().Add("Vary", acceptEncoding)
		if r.Header.Get("Upgrade") != "" {
			next.ServeHTTP(w, r)
			return
		}
		encoding := parseAcceptEncoding(r.Header.Get(acceptEncoding))
		var encWriter io.WriteCloser
		switch encoding {
		case Gzip:
			encWriter, _ = gzip.NewWriterLevel(w, c.compressOpt.Level)
		case Deflate:
			encWriter, _ = flate.NewWriter(w, c.compressOpt.Level)
		default:
			next.ServeHTTP(enhancedWriter, r)
			return
		}
		defer vclose.Close(encWriter)
		w.Header().Set("Transfer-Encoding", "chunked")
		w.Header().Set("Content-Encoding", encoding.String())
		r.Header.Del(acceptEncoding)
		enhancedWriter.compressor = encWriter
		next.ServeHTTP(enhancedWriter, r)
	})
}

type writer struct {
	w            http.ResponseWriter
	r            *http.Request
	compressor   io.Writer
	status       int       // 状态码（如200、404）
	bytesWritten int64     // 响应体字节数
	createAt     time.Time // 创建时间
}

func (c *writer) WriteHeader(status int) {
	c.status = status

	h := c.w.Header()
	if c.compressor != nil {
		if !shouldCompress(h.Get("Content-Type")) {
			c.closeCompressor()
		} else {
			c.Header().Del("Content-Length")
		}
	}

	c.w.WriteHeader(status)
}

func (c *writer) Header() http.Header {
	return c.w.Header()
}

func (c *writer) Write(b []byte) (n int, err error) {
	h := c.w.Header()
	if c.compressor == nil || !shouldCompress(h.Get("Content-Type")) {
		c.closeCompressor()
		n, err = c.w.Write(b)
	} else {
		h.Del("Content-Length")
		n, err = c.compressor.Write(b)
	}
	c.bytesWritten += int64(n)
	return n, err
}

func (c *writer) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	hijacker, ok := c.w.(http.Hijacker)
	if !ok {
		return nil, nil, errors.New("ResponseWriter 不支持 Hijacker 接口")
	}
	// 重要：关闭并清理压缩器
	if c.compressor != nil {
		// 如果有压缩器，需要关闭它
		if closer, ok := c.compressor.(io.Closer); ok {
			vclose.Close(closer)
		}
		c.closeCompressor()
	}

	// 移除压缩相关的头部，避免客户端误解
	c.w.Header().Del("Content-Encoding")
	c.w.Header().Del("Vary")

	return hijacker.Hijack()
}

func (c *writer) ReadFrom(r io.Reader) (n int64, err error) {
	h := c.w.Header()
	if c.compressor == nil || !shouldCompress(h.Get("Content-Type")) {
		c.closeCompressor()
		return io.Copy(c.w, r)
	}
	h.Del("Content-Length")
	return io.Copy(c.compressor, r)
}

type flusher interface {
	Flush() error
}

func (c *writer) Flush() {
	// Flush compressed data if compressor supports it.
	if f, ok := c.compressor.(flusher); ok {
		_ = f.Flush()
	}
	// Flush HTTP response.
	if f, ok := c.w.(http.Flusher); ok {
		f.Flush()
	}
}

func (c *writer) closeCompressor() {
	c.w.Header().Del("Content-Encoding")
	c.w.Header().Del("Vary")
	c.compressor = nil
}
