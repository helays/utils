package middleware

import (
	"context"
	"net/http"
	"time"

	"go.uber.org/zap"
	"helay.net/go/utils/v3/logger/ulogs"
	"helay.net/go/utils/v3/logger/zaploger"
	"helay.net/go/utils/v3/net/http/request"
)

func (c *ResponseProcessor) metrics(w *writer) {
	l := Logs{
		Status:      w.status,
		Ip:          request.Getip(w.r),
		Method:      w.r.Method,
		Uri:         w.r.URL.String(),
		ContentSize: w.bytesWritten,
		UserAgent:   w.r.Header.Get("User-Agent"),
		Elapsed:     time.Since(w.createAt),
	}

	for _, lg := range c.logEvent {
		if lg != nil {
			lg.Write(&l)
		}
	}

}

// ================= StdLogger 日志标准输出器 ================

// StdLogger 日志标准输出器
// http code >= 500 输出到Stderr,日志等级Fatal
// http code >= 400 输出的Stderr，日志等级Error
// http code >= 300 输出到Stdout，日志等级Warn
// 否则输出到 Stdout，日志等级 Info
type StdLogger struct{}

func NewStdLogger() *StdLogger {
	return &StdLogger{}
}

func (s *StdLogger) Write(l *Logs) {
	if l.Status >= http.StatusInternalServerError {
		ulogs.Fatalf("[%s] %s %s %d %d %s [%s]", l.Ip, l.Method, l.Uri, l.Status, l.ContentSize, l.UserAgent, l.Elapsed)
	} else if l.Status >= http.StatusBadRequest {
		ulogs.Errorf("[%s] %s %s %d %d %s [%s]", l.Ip, l.Method, l.Uri, l.Status, l.ContentSize, l.UserAgent, l.Elapsed)
	} else if l.Status >= http.StatusMultipleChoices {
		ulogs.Warnf("[%s] %s %s %d %d %s [%s]", l.Ip, l.Method, l.Uri, l.Status, l.ContentSize, l.UserAgent, l.Elapsed)
	} else {
		ulogs.Infof("[%s] %s %s %d %d %s [%s]", l.Ip, l.Method, l.Uri, l.Status, l.ContentSize, l.UserAgent, l.Elapsed)
	}
}

// ================= DebugStdLogger 日志标准输出器 ================

// DebugStdLogger 调试日志输出器（将正常日志输出到debug级别）
// http code >= 500 输出到Stderr,日志等级Fatal
// http code >= 400 输出的Stderr，日志等级Error
// http code >= 300 输出到Stdout，日志等级Warn
// 否则输出到 Stdout，日志等级 Debug（正常日志使用debug级别）
type DebugStdLogger struct{}

func NewDebugStdLogger() *DebugStdLogger {
	return &DebugStdLogger{}
}

func (d *DebugStdLogger) Write(l *Logs) {
	if l.Status >= http.StatusInternalServerError {
		ulogs.Fatalf("[%s] %s %s %d %d %s [%s]", l.Ip, l.Method, l.Uri, l.Status, l.ContentSize, l.UserAgent, l.Elapsed)
	} else if l.Status >= http.StatusBadRequest {
		ulogs.Errorf("[%s] %s %s %d %d %s [%s]", l.Ip, l.Method, l.Uri, l.Status, l.ContentSize, l.UserAgent, l.Elapsed)
	} else if l.Status >= http.StatusMultipleChoices {
		ulogs.Warnf("[%s] %s %s %d %d %s [%s]", l.Ip, l.Method, l.Uri, l.Status, l.ContentSize, l.UserAgent, l.Elapsed)
	} else {
		ulogs.Debugf("[%s] %s %s %d %d %s [%s]", l.Ip, l.Method, l.Uri, l.Status, l.ContentSize, l.UserAgent, l.Elapsed)
	}
}

// ================= ZapLogger 日志输出器 ================

type ZapLogger struct {
	logger *zaploger.Logger
}

func NewZapLogger(c *zaploger.Config) (*ZapLogger, error) {
	logger, err := zaploger.New(c)
	if err != nil {
		return nil, err
	}
	return &ZapLogger{logger}, nil
}

func (z *ZapLogger) Write(l *Logs) {
	msg := []any{
		zap.String(l.Method, l.Uri),
		zap.Int("status", l.Status),
		zap.Int64("bytes_send", l.ContentSize),
		zap.String("http_user_agent", l.UserAgent),
		zap.Any("elapsed", l.Elapsed),
	}

	if l.Status >= http.StatusBadRequest {
		z.logger.Error(context.Background(), l.Ip, msg...)
	} else if l.Status >= http.StatusMultipleChoices {
		z.logger.Warn(context.Background(), l.Ip, msg...)
	} else {
		z.logger.Debug(context.Background(), l.Ip, msg...)
	}

}
