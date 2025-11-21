package httpServer

import (
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/helays/utils/v2/net/http/httpServer/request"
	"github.com/helays/utils/v2/net/http/httpServer/response"
	"github.com/helays/utils/v2/net/http/httpServer/responsewriter"
)

// 默认验证中间件
func (h *HttpServer) defaultValid(next http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		recorder := responsewriter.New(w) // 创建包装器
		if len(h.serverNameMap) > 0 {
			// 提取并转换为小写的host（忽略端口部分）
			host := strings.ToLower(strings.SplitN(r.Host, ":", 2)[0])
			if _, ok := h.serverNameMap[host]; !ok {
				recorder.WriteHeader(http.StatusBadGateway)
				return
			}
		}
		start := time.Now()
		defer h.httpLogger(recorder, r, start)
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

// 检测ip白名单和黑名单
func (h *HttpServer) checkIpAccess(w http.ResponseWriter, r *http.Request) bool {
	addr := request.Getip(r)
	ip, _, err := net.SplitHostPort(addr)
	if err != nil {
		ip = addr // 如果无端口，则直接使用原地址
	}
	if h.denyIpList != nil && h.denyIpList.Contains(ip) {
		response.Forbidden(w, "你的IP已被监管")
		return false
	}
	if h.allowIpList != nil {
		if !h.allowIpList.Contains(ip) {
			response.Forbidden(w, "你的IP不在系统白名单内")
			return false
		}
	}
	return true
}

func (h *HttpServer) debugIpAccess(w http.ResponseWriter, r *http.Request) bool {
	// 判断path是否以debug开头
	if !strings.HasPrefix(r.URL.Path, "/debug/") {
		return true
	}
	addr := request.Getip(r)
	ip, _, err := net.SplitHostPort(addr)
	if err != nil {
		ip = addr // 如果无端口，则直接使用原地址
	}
	if !h.debugAllowIpList.Contains(ip) {
		response.Forbidden(w, http.StatusText(http.StatusForbidden))
		return false
	}
	return true
}

// 跨域中间件
func (h *HttpServer) cors(next http.Handler) http.HandlerFunc {
	if h.CORS == nil || !h.CORS.Enabled {
		return next.ServeHTTP
	}
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodOptions {
			h.CORS.HandlePreflight(w, r.Header.Get("Origin"))
			return
		}
		if shouldContinue := h.CORS.Apply(w, r.Header.Get("Origin")); !shouldContinue {
			return // 严格模式下被拒绝
		}
		next.ServeHTTP(w, r)
	}
}
