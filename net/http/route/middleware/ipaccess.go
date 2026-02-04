package middleware

import (
	"net"
	"net/http"

	"github.com/helays/utils/v2/net/http/request"
	"github.com/helays/utils/v2/net/ipmatch"
)

type IPAccessMiddleware struct {
	allowIPMatch *ipmatch.IPMatcher
	denyIPMatch  *ipmatch.IPMatcher
	debugIPMatch *ipmatch.IPMatcher
}

func NewIPAccessMiddleware() *IPAccessMiddleware {
	return &IPAccessMiddleware{}
}

func (i *IPAccessMiddleware) SetAllow(allowIPMatch *ipmatch.IPMatcher) {
	i.allowIPMatch = allowIPMatch
	i.allowIPMatch.Build()
}

func (i *IPAccessMiddleware) SetDeny(denyIPMatch *ipmatch.IPMatcher) {
	i.denyIPMatch = denyIPMatch
	i.denyIPMatch.Build()
}

func (i *IPAccessMiddleware) SetDebug(debugIPMatch *ipmatch.IPMatcher) {
	i.debugIPMatch = debugIPMatch
	i.debugIPMatch.Build()
}

func (i *IPAccessMiddleware) Close() {

}

func (i *IPAccessMiddleware) DenyIpAccess(next http.Handler) http.Handler {
	if i.denyIPMatch == nil {
		return next
	}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if i.denyIPMatch.Contains(filterIPAndPort(r)) {
			if w != nil {
				w.WriteHeader(http.StatusForbidden)
			}
			return
		}
		next.ServeHTTP(w, r)
	})
}

func (i *IPAccessMiddleware) AllowIpAccess(next http.Handler) http.Handler {
	if i.allowIPMatch == nil {
		return next
	}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !i.allowIPMatch.Contains(filterIPAndPort(r)) {
			if w != nil {
				w.WriteHeader(http.StatusForbidden)
			}
			return
		}
		next.ServeHTTP(w, r)
	})
}

func (i *IPAccessMiddleware) DebugIPAccess(next http.Handler) http.Handler {
	if i.debugIPMatch == nil {
		return next
	}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if debug, ok := r.Context().Value(DebugCtxKey).(bool); ok && debug {
			if !i.debugIPMatch.Contains(filterIPAndPort(r)) {
				if w != nil {
					w.WriteHeader(http.StatusForbidden)
				}
				return
			}
		}

		next.ServeHTTP(w, r)
	})
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
