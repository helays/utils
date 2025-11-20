package csrf_std

import (
	"net/http"

	"github.com/helays/utils/v2/security/csrf"
)

func WrapHandler(handler http.HandlerFunc, config *csrf.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if config.ShouldValidate(r.Method) {
			switch config.Strategy {
			case csrf.StrategyToken:
			case csrf.StrategyDoubleTap:
			case csrf.StrategyReferer:

			}
		}
		handler(w, r)
	}
}
