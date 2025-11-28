package middleware

import (
	"context"
	"net/http"
	"time"

	"github.com/helays/utils/v2/logger/ulogs"
	"github.com/helays/utils/v2/net/http/httpServer/request"
	"go.uber.org/zap"
)

func (c *ResponseProcessor) metrics(w *writer) {
	ua := w.r.Header.Get("User-Agent")
	elapsed := time.Since(w.createAt).Milliseconds() // 耗时
	status := w.status
	ip := request.Getip(w.r)
	method := w.r.Method
	contentSize := w.bytesWritten
	uri := w.r.URL.String()
	if c.logger == nil {
		if status >= http.StatusBadRequest {
			ulogs.Errorf("[%s] %s %s %d %d %s [%dms]", ip, method, uri, status, contentSize, ua, elapsed)
		} else {
			ulogs.Infof("[%s] %s %s %d %d %s [%dms]", ip, method, uri, status, contentSize, ua, elapsed)
		}
		return
	} else {
		if status >= http.StatusBadRequest {
			c.logger.Error(
				context.Background(),
				ip,
				zap.String(method, uri),
				zap.Int("status", status),
				zap.Int64("bytes_send", contentSize),
				zap.String("http_user_agent", ua),
				zap.Int64("elapsed", elapsed),
			)
		} else {
			c.logger.Debug(
				context.Background(),
				ip,
				zap.String(method, uri),
				zap.Int("status", status),
				zap.Int64("bytes_send", contentSize),
				zap.String("http_user_agent", ua),
				zap.Int64("elapsed", elapsed),
			)
		}
	}
}
