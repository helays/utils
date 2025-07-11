package router

import (
	"github.com/helays/utils/tools"
	"net/http"
	"time"
)

func (ro Router) SetCookie(w http.ResponseWriter, k, value, path string) {
	path = tools.Ternary(path == "", "/", path)
	cookie := http.Cookie{
		Name:       k,
		Value:      value,
		Path:       path,
		Domain:     ro.CookieDomain,
		Expires:    time.Time{},
		RawExpires: "",
		MaxAge:     0,
		Secure:     ro.CookieSecure,
		HttpOnly:   ro.CookieHttpOnly,
		SameSite:   0,
		Raw:        "",
		Unparsed:   nil,
	}
	http.SetCookie(w, &cookie)
}

func (ro Router) DelCookie(w http.ResponseWriter, k, path string) {

	cookie := http.Cookie{
		Name:       k,
		Value:      "",
		Path:       path,
		Domain:     ro.CookieDomain,
		RawExpires: "",
		MaxAge:     -1,
		Secure:     ro.CookieSecure,
		HttpOnly:   ro.CookieHttpOnly,
		SameSite:   0,
		Raw:        "",
		Unparsed:   nil,
	}
	http.SetCookie(w, &cookie)
}
