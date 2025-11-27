package middleware

import (
	"context"
	"net/http"
	"time"

	"github.com/helays/utils/v2/logger/ulogs"
	"github.com/helays/utils/v2/logger/zaploger"
	"github.com/helays/utils/v2/net/http/httpServer/request"
	"github.com/helays/utils/v2/net/http/responsewriter"
	"go.uber.org/zap"
)

type LoggerMiddleware struct {
	opt    zaploger.Config
	logger *zaploger.Logger
}

func NewLoggerMiddleware(opt zaploger.Config) (*LoggerMiddleware, error) {
	l := &LoggerMiddleware{
		opt: opt,
	}
	if len(l.opt.LogLevelConfigs) > 0 {
		var err error
		l.logger, err = zaploger.New(&l.opt)
		if err != nil {
			return nil, err
		}
	}
	return l, nil
}

func (l *LoggerMiddleware) Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		recorder := responsewriter.New(w)
		start := time.Now()
		defer func() {
			ua := r.Header.Get("User-Agent")
			elapsed := time.Since(start).Milliseconds() // 耗时
			status := recorder.GetStatus()

			if l.logger == nil {
				if status >= 400 || status < 200 {
					ulogs.Error(request.Getip(r), r.Method, r.URL.String(), status, recorder.GetBytesWritten(), ua, elapsed)
				} else {
					ulogs.Debug(request.Getip(r), r.Method, r.URL.String(), status, recorder.GetBytesWritten(), ua, elapsed)
				}
				return
			}
			if status >= 200 && status < 400 {
				l.logger.Info(
					context.Background(),
					request.Getip(r),
					zap.String(r.Method, r.URL.String()),
					zap.Int("status", status),
					zap.Int64("bytes_send", recorder.GetBytesWritten()),
					zap.String("http_user_agent", ua),
					zap.Int64("elapsed", elapsed),
				)
			} else {
				l.logger.Error(
					context.Background(),
					request.Getip(r),
					zap.String(r.Method, r.URL.String()),
					zap.Int("status", status),
					zap.Int64("bytes_send", recorder.GetBytesWritten()),
					zap.String("http_user_agent", ua),
					zap.Int64("elapsed", elapsed),
				)
			}
		}()
		next.ServeHTTP(recorder, r)
	})
}
