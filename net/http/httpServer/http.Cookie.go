package httpServer

import (
	"github.com/helays/utils/tools"
	"net/http"
	"time"
)

func (router Router) SetCookie(w http.ResponseWriter, k, value, path string) {
	path = tools.Ternary(path == "", "/", path)
	cookie := http.Cookie{
		Name:       k,
		Value:      value,
		Path:       path,
		Domain:     router.CookieDomain,
		Expires:    time.Time{},
		RawExpires: "",
		MaxAge:     0,
		Secure:     router.CookieSecure,
		HttpOnly:   router.CookieHttpOnly,
		SameSite:   0,
		Raw:        "",
		Unparsed:   nil,
	}
	http.SetCookie(w, &cookie)
}

func (router Router) DelCookie(w http.ResponseWriter, k, path string) {

	cookie := http.Cookie{
		Name:       k,
		Value:      "",
		Path:       path,
		Domain:     router.CookieDomain,
		RawExpires: "",
		MaxAge:     -1,
		Secure:     router.CookieSecure,
		HttpOnly:   router.CookieHttpOnly,
		SameSite:   0,
		Raw:        "",
		Unparsed:   nil,
	}
	http.SetCookie(w, &cookie)
}
