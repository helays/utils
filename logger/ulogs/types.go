/*
Package ulogs 提供一个高性能的异步日志库

# 设计哲学

1. 简单可靠：不添加多余的检查，相信开发者会正确初始化
2. 高性能：使用channel缓冲 + 对象池，减少GC压力
3. 双输出：stdout/stderr分离，符合Unix哲学
4. Fail Fast：未初始化直接panic，开发阶段暴露问题

# 快速开始

	package main

	import (
		"time"
		"yourmodule/ulogs"
	)

	func main() {
		// 1. 必须初始化，否则panic
		ulogs.New(&ulogs.Config{
			BufferSize:   20000, // stdout缓冲区大小
			ErrorBufSize: 10000, // stderr缓冲区大小
		})

		// 2. 设置日志级别（可选）
		ulogs.SetLevel(ulogs.LogLevelDebug)

		// 3. 启动处理器
		ulogs.Start()
		defer ulogs.Shutdown()

		// 4. 正常使用
		ulogs.Info("Server started")
		ulogs.Debugf("Config loaded: %v", config)
		ulogs.Errorf("Failed to connect: %v", err)
	}

# 性能特性

- 无锁设计：每个级别独立channel
- 对象池复用：减少内存分配和GC压力
- 异步处理：日志写入不阻塞业务逻辑
- 批量缓冲：log.Logger内置缓冲，无需额外批量

# 许可证

Copyright (c) 2024 helays
MIT License
*/
package ulogs

import (
	"log"
	"os"
	"sync"
)

const (
	defaultBufferSize   = 10000 // 标准输出缓冲区
	defaultErrorBufSize = 5000  // 错误输出缓冲区（可以小一些）
)

type Config struct {
	BufferSize   int `json:"buffer_size" yaml:"buffer_size" ini:"buffer_size"`          // 标准输出缓冲区
	ErrorBufSize int `json:"error_buf_size" yaml:"error_buf_size" ini:"error_buf_size"` // 错误输出缓冲区
}
type logEntry struct {
	logger      *log.Logger
	isFormatted bool
	format      string
	args        []any
}

func (l *logEntry) reset() {
	l.logger = nil
	l.isFormatted = false
	l.format = ""
	l.args = nil
}

type FastLogger struct {
	stdoutChan chan *logEntry
	stderrChan chan *logEntry
	done       chan struct{}
	wg         sync.WaitGroup

	traceLogger *log.Logger
	debugLogger *log.Logger
	infoLogger  *log.Logger
	warnLogger  *log.Logger
	errorLogger *log.Logger
	fatalLogger *log.Logger

	// 对象池
	entryPool *sync.Pool
}

var (
	globalLogger *FastLogger
	once         sync.Once
)

func init() {
	globalLogger = &FastLogger{
		stdoutChan:  make(chan *logEntry, defaultBufferSize),
		stderrChan:  make(chan *logEntry, defaultErrorBufSize),
		done:        make(chan struct{}),
		traceLogger: log.New(os.Stdout, "【TRACE】", log.LstdFlags),
		debugLogger: log.New(os.Stdout, "【DEBUG】", log.LstdFlags),
		infoLogger:  log.New(os.Stdout, "【INFO】", log.LstdFlags),
		warnLogger:  log.New(os.Stderr, "【WARN】", log.LstdFlags),
		errorLogger: log.New(os.Stderr, "【ERROR】", log.LstdFlags),
		fatalLogger: log.New(os.Stderr, "【FATAL】", log.LstdFlags),
	}
	// 初始化对象池
	globalLogger.entryPool = &sync.Pool{
		New: func() any {
			return &logEntry{}
		},
	}
}

// SetBufferSize 设置缓冲区大小
// 注意：必须在 Start 之前调用，否则不会生效
func SetBufferSize(stdoutSize, stderrSize int) {
	if stdoutSize > 0 {
		globalLogger.stdoutChan = make(chan *logEntry, stdoutSize)
	}

	if stderrSize > 0 {
		globalLogger.stderrChan = make(chan *logEntry, stderrSize)
	}
}

func SetLogger(level int, l *log.Logger) {
	switch level {
	case LogLevelTrace:
		globalLogger.traceLogger = l
	case LogLevelDebug:
		globalLogger.debugLogger = l
	case LogLevelInfo:
		globalLogger.infoLogger = l
	case LogLevelWarn:
		globalLogger.warnLogger = l
	case LogLevelError:
		globalLogger.errorLogger = l
	case LogLevelFatal:
		globalLogger.fatalLogger = l

	}
}
func Start() {
	globalLogger.wg.Add(2)
	go globalLogger.processor(globalLogger.stdoutChan) // 标准输出日志处理
	go globalLogger.processor(globalLogger.stderrChan) // 错误输出日志处理
}

func Shutdown() {
	close(globalLogger.done)
	globalLogger.wg.Wait()
}

func (f *FastLogger) processor(ch <-chan *logEntry) {
	defer f.wg.Done()
	for {
		select {
		case entry := <-ch:
			f.writeEntry(entry)
		case <-f.done:
			// 退出前处理完所有日志
			for {
				select {
				case entry := <-ch:
					f.writeEntry(entry)
				default:
					return
				}
			}
		}
	}
}

func (f *FastLogger) writeEntry(entry *logEntry) {
	defer f.putEntry(entry)
	if entry.isFormatted {
		entry.logger.Printf(entry.format, entry.args...)
	} else {
		entry.logger.Println(entry.args...)
	}
}

// getEntry 从对象池获取entry
func (f *FastLogger) getEntry() *logEntry {
	return f.entryPool.Get().(*logEntry)
}

// putEntry 归还entry到池中（会自动reset）
func (f *FastLogger) putEntry(entry *logEntry) {
	// 重要：先reset再归还
	entry.reset()
	f.entryPool.Put(entry)
}
