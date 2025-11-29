package route

import (
	"errors"
	"io/fs"
	"net/http"
)

func toHTTPError(err error) (msg string, httpStatus int) {
	switch {
	case errors.Is(err, fs.ErrNotExist):
		return "404 page not found", http.StatusNotFound
	case errors.Is(err, fs.ErrPermission):
		return "403 Forbidden", http.StatusForbidden
	default:
		return "500 Internal Server Error", http.StatusInternalServerError
	}
}

func RenderErrorText(w http.ResponseWriter, resp *ErrorResp) {
	if resp == nil {
		return
	}
	desc := resp.Msg
	if desc == "" && resp.Error != nil {
		desc = resp.Error.Error()
	} else {
		desc = http.StatusText(resp.Code)
	}
	http.Error(w, desc, resp.Code)
}
