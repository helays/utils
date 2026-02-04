package captcha

import (
	"bytes"
	"net/http"
	"time"

	"github.com/dchest/captcha"
	"helay.net/go/utils/v3/logger/ulogs"
	"helay.net/go/utils/v3/net/http/response"
	"helay.net/go/utils/v3/net/http/session"
)

// Text 文字验证码
func (c *Captcha) Text() http.HandlerFunc {
	return c.TextWithStore(func(w http.ResponseWriter, r *http.Request, key, captchaId string) error {
		sv := session.Value{Field: key, Value: captchaId, TTL: c.opt.Text.ExpireTime}
		return session.GetSession().Set(w, r, &sv)
	})
}

type TextStoreFunc func(w http.ResponseWriter, r *http.Request, key, captchaId string) error

func (c *Captcha) TextWithStore(store TextStoreFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		c.TextHandFunc(w, r, store)
	}
}

func (c *Captcha) TextHandFunc(w http.ResponseWriter, r *http.Request, store TextStoreFunc) {
	w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	w.Header().Set("Pragma", "no-cache")
	w.Header().Set("Expires", "0")

	captchaId := captcha.NewLen(c.opt.Text.Length)
	if err := store(w, r, CaptchaTextKey, captchaId); err != nil {
		ulogs.Error(err)
		response.InternalServerError(w)
		return
	}
	var content bytes.Buffer
	if err := captcha.WriteImage(&content, captchaId, c.opt.Text.Width, c.opt.Text.Height); err != nil {
		ulogs.Errorf("验证码生成失败 %v", err)
		response.InternalServerError(w)
		return
	}
	w.Header().Set("Content-Type", "image/png")
	http.ServeContent(w, r, "", time.Time{}, bytes.NewReader(content.Bytes()))
}
