package server

import (
	"io"
	"log"
	"os"
	"strings"
	"sync"
)

// FilteringLogger 实现日志过滤功能
type FilteringLogger struct {
	out    io.Writer
	config LogFilterConfig
	mu     sync.Mutex
	logger *log.Logger
}

// NewFilteringLogger 创建一个新的过滤日志器
func NewFilteringLogger(config LogFilterConfig) *FilteringLogger {
	fl := &FilteringLogger{
		out:    os.Stderr,
		config: config,
	}
	fl.logger = log.New(fl, "", log.LstdFlags)
	return fl
}

// Write 实现 io.Writer 接口，用于过滤日志
func (fl *FilteringLogger) Write(p []byte) (n int, err error) {
	// 如果配置了屏蔽所有服务器错误
	if fl.config.SuppressAllServerErrors {
		return len(p), nil
	}

	msg := string(p)

	// 检查是否包含 TLS handshake error
	if fl.config.SuppressTLSHandshakeError {
		if strings.Contains(msg, "TLS handshake error") {
			// 过滤掉特定的 TLS 握手错误
			if strings.Contains(msg, "client sent an HTTP request to an HTTPS server") ||
				strings.Contains(msg, "no common certificate") ||
				strings.Contains(msg, "bad certificate") {
				return len(p), nil
			}
		}
	}

	// 检查是否包含客户端断开连接的错误
	if fl.config.SuppressClientDisconnect {
		if strings.Contains(msg, "read from client: EOF") ||
			strings.Contains(msg, "write: broken pipe") ||
			strings.Contains(msg, "connection reset by peer") {
			return len(p), nil
		}
	}

	// 其他日志正常输出
	return fl.out.Write(p)
}

// GetErrorLog 返回一个 *log.Logger 供 http.Server 使用
func (fl *FilteringLogger) GetErrorLog() *log.Logger {
	return fl.logger
}
