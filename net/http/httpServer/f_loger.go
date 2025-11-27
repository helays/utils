package httpServer

import (
	"context"
	"net/http"
	"time"

	"github.com/helays/utils/v2/logger/ulogs"
	"github.com/helays/utils/v2/net/http/httpServer/request"
	"github.com/helays/utils/v2/net/http/responsewriter"
	"go.uber.org/zap"
	"golang.org/x/net/websocket"
)

// httpLogger http 日志
func (h *HttpServer) httpLogger(w *responsewriter.StatusRecorder, r *http.Request, start time.Time) {
	ua := r.Header.Get("User-Agent")
	elapsed := time.Since(start).Milliseconds() // 耗时
	if h.logger != nil {
		// 这里输出info 级别的请求日志
		h.logger.Info(context.Background(),
			request.Getip(r),
			zap.String(r.Method, r.URL.String()),
			zap.Int("status", w.GetStatus()),
			zap.Int64("bytes_send", w.GetBytesWritten()),
			zap.String("http_user_agent", ua),
			zap.Int64("elapsed", elapsed),
		)
	} else {
		ulogs.Debug(request.Getip(r), r.Method, r.URL.String(), w.GetStatus(), w.GetBytesWritten(), ua, elapsed)
	}
}

// socketLogger socket 日志
func (h *HttpServer) socketLogger(ws *websocket.Conn, start time.Time) {
	elapsed := time.Since(start).Milliseconds() // 耗时
	req := ws.Request()
	ua := req.Header.Get("User-Agent")

	if h.logger != nil {
		// 这里输出info级别的请求日志
		h.logger.Info(context.Background(),
			request.Getip(ws.Request()),
			zap.String("method", "WEBSOCKET"),
			zap.String("url", req.URL.Path),
			zap.String("http_user_agent", ua),
			zap.Int64("elapsed", elapsed),
		)
	} else {
		ulogs.Debug(request.Getip(ws.Request()), "WEBSOCKET", req.URL.Path, ua, elapsed)
	}
}
