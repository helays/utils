package httpServer

import (
	"net/http"
	"strings"
	"time"

	"github.com/helays/utils/v2/close/httpClose"
	"github.com/helays/utils/v2/close/vclose"
	"golang.org/x/net/websocket"
)

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

func (h *HttpServer) middleware(u string, f http.Handler, callback ...MiddlewareFunc) {
	var handler http.Handler
	mid := []MiddlewareFunc{
		h.denyIPAccess,
		h.allowIPAccess,
		h.debugIPAccess,
		h.cors,
	}

	if h.Security.DefaultValidLast {
		mid = append(mid, callback...)
		mid = append(mid, h.defaultValid)
	} else {
		mid = append(mid, h.defaultValid)
		mid = append(mid, callback...)
	}
	handler = Chain(mid...)(f)
	// 最终的处理函数
	finalHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer httpClose.CloseReq(r)
		handler.ServeHTTP(w, r)
	})
	h.mux.Handle(u, finalHandler)
}

// Socket 公共中间件
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
		defer h.socketLogger(ws, start)
		// 白名单验证
		if h.denyIPMatch != nil && h.checkDenyIpAccess(nil, ws.Request()) {
			vclose.Close(ws)
			return
		} else if h.allowIPMatch != nil && !h.checkAllowIPAccess(nil, ws.Request()) {
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
