package cors_std

import (
	"net/http"

	"github.com/helays/utils/v2/security/cors"
)

func Cors(opt cors.Config) func(next http.Handler) http.Handler {
	if !opt.Enabled {
		return func(next http.Handler) http.Handler {
			return next
		}
	}
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method == http.MethodOptions {
				opt.HandlePreflight(w, r.Header.Get("Origin"))
				return
			}
			if shouldContinue := opt.Apply(w, r.Header.Get("Origin")); !shouldContinue {
				return // 严格模式下被拒绝
			}
			next.ServeHTTP(w, r)
		})
	}
}
