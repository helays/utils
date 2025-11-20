package httpServer

import (
	"net/http"
	"strings"
	"time"

	"github.com/helays/utils/v2/close/httpClose"
	"github.com/helays/utils/v2/close/vclose"
	"github.com/helays/utils/v2/logger/ulogs"
	"github.com/helays/utils/v2/net/http/httpServer/request"
	"github.com/helays/utils/v2/net/http/httpServer/responsewriter"
	"go.uber.org/zap"
	"golang.org/x/net/context"
	"golang.org/x/net/websocket"
)

func (h *HttpServer) middleware(u string, f http.Handler, callback ...MiddlewareFunc) {
	h.initMux()
	// 统一处理 finalHandler 的生成逻辑
	h.mux.HandleFunc(u, func(w http.ResponseWriter, r *http.Request) {
		defer httpClose.CloseReq(r) // 确保请求被正确关闭
		var finalHandler http.Handler
		if h.DefaultValidFirst {
			finalHandler = f
			if len(callback) > 0 {
				finalHandler = Chain(callback...)(finalHandler)
			}
			finalHandler = h.defaultValid(finalHandler)
		} else {
			finalHandler = h.defaultValid(f)
			if len(callback) > 0 {
				finalHandler = Chain(callback...)(finalHandler)
			}
		}
		finalHandler.ServeHTTP(w, r)
	})

}
func (h *HttpServer) defaultValid(next http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 创建包装器
		recorder := responsewriter.New(w) // 默认200
		if len(h.serverNameMap) > 0 {
			// 提取并转换为小写的host（忽略端口部分）
			host := strings.ToLower(strings.SplitN(r.Host, ":", 2)[0])
			if _, ok := h.serverNameMap[host]; !ok {
				recorder.WriteHeader(http.StatusBadGateway)
				return
			}
		}
		start := time.Now()
		defer func() {
			ua := r.Header.Get("User-Agent")
			elapsed := time.Since(start).Milliseconds() // 耗时
			if h.logger != nil {
				// 这里输出info 级别的请求日志
				h.logger.Info(context.Background(),
					request.Getip(r),
					zap.String(r.Method, r.URL.String()),
					zap.Int("status", recorder.GetStatus()),
					zap.Int64("bytes_send", recorder.GetBytesWritten()),
					zap.String("http_user_agent", ua),
					zap.Int64("elapsed", elapsed),
				)
			} else {
				ulogs.Debug(request.Getip(r), r.Method, r.URL.String(), recorder.GetStatus(), recorder.GetBytesWritten(), ua, elapsed)
			}
		}()
		// add header
		recorder.Header().Set("server", "vs/1.0")
		recorder.Header().Set("connection", "keep-alive")
		// 白名单验证
		if h.enableCheckIpAccess && !h.checkIpAccess(recorder, r) {
			return
		}
		if !h.debugIpAccess(recorder, r) {
			return
		}
		if h.CommonCallback != nil && !h.CommonCallback(recorder, r) {
			return
		}
		next.ServeHTTP(recorder, r)
	}
}

func (h *HttpServer) socketMiddleware(u string, f websocket.Handler) {
	handler := websocket.Handler(func(ws *websocket.Conn) {
		defer vclose.Close(ws)
		if len(h.serverNameMap) > 0 {
			// 提取并转换为小写的host（忽略端口部分）
			host := strings.ToLower(strings.SplitN(ws.Request().Host, ":", 2)[0])
			if _, ok := h.serverNameMap[host]; !ok {
				// 对于WebSocket，我们可能不能直接返回HTTP状态码，但可以决定是否关闭连接
				vclose.Close(ws)
				return
			}
		}

		start := time.Now()
		ua := ws.Request().Header.Get("User-Agent")
		defer func() {
			elapsed := time.Since(start).Milliseconds() // 耗时
			if h.logger != nil {
				// 这里输出info级别的请求日志
				h.logger.Info(context.Background(),
					request.Getip(ws.Request()),
					zap.String("method", "WEBSOCKET"),
					zap.String("url", u),
					zap.String("http_user_agent", ua),
					zap.Int64("elapsed", elapsed),
				)
			} else {
				ulogs.Debug(request.Getip(ws.Request()), "WEBSOCKET", u, ua, elapsed)
			}
		}()
		// 白名单验证
		if h.enableCheckIpAccess && !h.checkIpAccess(nil, ws.Request()) {
			vclose.Close(ws)
			return
		}

		// 如果有通用回调并且该回调不允许继续，则不进行后续处理
		if h.CommonCallback != nil && !h.CommonCallback(nil, ws.Request()) {
			vclose.Close(ws)
			return
		}
		f(ws)
	})
	h.mux.Handle(u, handler)
}

// MiddlewareFunc 修改为支持 http.Handler
type MiddlewareFunc func(next http.Handler) http.Handler

// Chain 中间件链（兼容 http.Handler）
func Chain(middlewares ...MiddlewareFunc) MiddlewareFunc {
	return func(final http.Handler) http.Handler {
		for i := len(middlewares) - 1; i >= 0; i-- {
			final = middlewares[i](final)
		}
		return final
	}
}
