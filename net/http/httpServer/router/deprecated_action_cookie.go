package router

import (
	"net/http"
	"time"

	"helay.net/go/utils/v3/tools"
)

// Deprecated: 请使用 helay.net/go/utils/v3/net/http/cookiekit.SetCookie
func (ro *Router) SetCookie(w http.ResponseWriter, k, value, path string) {
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

// Deprecated: 请使用 helay.net/go/utils/v3/net/http/cookiekit.DelCookie
func (ro *Router) DelCookie(w http.ResponseWriter, k, path string) {

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
