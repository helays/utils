package session

import (
	"context"
	"net/http"

	"helay.net/go/utils/v3/logger/ulogs"
)

// Middleware session 中间件
// noinspection all
func Middleware() func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			if sessionId, err := session.GetSessionId(w, r); err != nil {
				ulogs.Errorf("解析Session ID失败 path=%s, method=%s %v", r.URL.Path, r.Method, err)
			} else {
				ctx = context.WithValue(ctx, SessionID, sessionId)
			}
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func GetSessionID(ctx context.Context) string {
	if val := ctx.Value(SessionID); val != nil {
		if sessionId, ok := val.(string); ok {
			return sessionId
		}
	}
	return ""
}
