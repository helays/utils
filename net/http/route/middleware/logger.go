package middleware

import (
	"context"
	"net/http"
	"time"

	"github.com/helays/utils/v2/logger/ulogs"
	"github.com/helays/utils/v2/logger/zaploger"
	"github.com/helays/utils/v2/net/http/request"
	"go.uber.org/zap"
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

// StdLogger 日志标准输出器

type StdLogger struct{}

func NewStdLogger() *StdLogger {
	return &StdLogger{}
}

func (s *StdLogger) Write(l *Logs) {
	if l.Status >= http.StatusBadRequest {
		ulogs.Errorf("[%s] %s %s %d %d %s [%s]", l.Ip, l.Method, l.Uri, l.Status, l.ContentSize, l.UserAgent, l.Elapsed)
	} else {
		ulogs.Infof("[%s] %s %s %d %d %s [%s]", l.Ip, l.Method, l.Uri, l.Status, l.ContentSize, l.UserAgent, l.Elapsed)
	}
}

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
		return
	}
	z.logger.Debug(context.Background(), l.Ip, msg...)
}
