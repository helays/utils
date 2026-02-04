package httpServer

import (
	"net"
	"net/http"
	"strings"

	"helay.net/go/utils/v3/net/http/request"
)

// 默认验证中间件
func (h *HttpServer) defaultValid(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if len(h.serverNameMap) > 0 {
			// 提取并转换为小写的host（忽略端口部分）
			host := strings.ToLower(strings.SplitN(r.Host, ":", 2)[0])
			if _, ok := h.serverNameMap[host]; !ok {
				w.WriteHeader(http.StatusBadGateway)
				return
			}
		}
		// add header
		w.Header().Set("server", "vs/1.0")
		w.Header().Set("connection", "keep-alive")

		if h.CommonCallback != nil && !h.CommonCallback(w, r) {
			return
		}
		next.ServeHTTP(w, r)
	})
}

// 跨域中间件
func (h *HttpServer) cors(next http.Handler) http.Handler {
	if h.Security.CORS == nil || !h.Security.CORS.Enabled {
		return next
	}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodOptions {
			h.Security.CORS.HandlePreflight(w, r.Header.Get("Origin"))
			return
		}
		if shouldContinue := h.Security.CORS.Apply(w, r.Header.Get("Origin")); !shouldContinue {
			return // 严格模式下被拒绝
		}
		next.ServeHTTP(w, r)
	})
}

func (h *HttpServer) denyIPAccess(next http.Handler) http.Handler {
	if h.denyIPMatch == nil {
		return next
	}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if h.checkDenyIpAccess(w, r) {
			return
		}
		next.ServeHTTP(w, r)
	})
}
func (h *HttpServer) allowIPAccess(next http.Handler) http.Handler {
	if h.allowIPMatch == nil {
		return next
	}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !h.checkAllowIPAccess(w, r) {
			return
		}
		next.ServeHTTP(w, r)
	})
}

// 调试IP访问控制
func (h *HttpServer) debugIPAccess(next http.Handler) http.Handler {
	if h.debugIPMatch == nil {
		return next
	}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 判断path是否以debug开头
		if strings.HasPrefix(r.URL.Path, "/debug/") {
			if !h.debugIPMatch.Contains(filterIPAndPort(r)) {
				if w != nil {
					w.WriteHeader(http.StatusForbidden)
				}
				return
			}
		}
		next.ServeHTTP(w, r)
	})
}

// 检测ip黑名单
func (h *HttpServer) checkDenyIpAccess(w http.ResponseWriter, r *http.Request) bool {
	if h.denyIPMatch.Contains(filterIPAndPort(r)) {
		if w != nil {
			w.WriteHeader(http.StatusForbidden)
		}
		return true
	}
	return false
}

// 检测ip白名单
func (h *HttpServer) checkAllowIPAccess(w http.ResponseWriter, r *http.Request) bool {
	if h.allowIPMatch.Contains(filterIPAndPort(r)) {
		return true
	}
	if w != nil {
		w.WriteHeader(http.StatusForbidden)
	}
	return false
}

// 过滤ip和端口
func filterIPAndPort(r *http.Request) string {
	addr := request.Getip(r)
	ip, _, err := net.SplitHostPort(addr)
	if err != nil {
		ip = addr // 如果无端口，则直接使用原地址
	}
	return ip
}
