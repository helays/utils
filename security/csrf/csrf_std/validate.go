package csrf_std

import (
	"fmt"
	"net/http"

	"helay.net/go/utils/v3/security/csrf"
)

// getClientToken 获取客户端的token
// 顺序是先从最安全的位置获取，然后是中间位置，最后是最不安全的位置
func (s *Std) getClientToken(r *http.Request) string {
	// 从Header中获取
	if token := r.Header.Get(csrf.DefaultTokenName); token != "" {
		return token
	}
	// 从Form中获取
	if token := r.FormValue(csrf.DefaultCookieName); token != "" {
		return token
	}

	// 从 Query中获取
	return r.URL.Query().Get(csrf.DefaultCookieName)
}

func (s *Std) getCookieToken(r *http.Request) string {
	// 从Cookie中获取
	if cookie, err := r.Cookie(csrf.DefaultCookieName); err == nil {
		return cookie.Value
	}
	return ""
}

// 双重提交Cookie
func (s *Std) validateDoubleTapStrategy(r *http.Request) (string, bool) {
	// 获取客户端Token
	clientToken := s.getClientToken(r)
	if clientToken == "" {
		return "", false
	}
	cookieToken := s.getCookieToken(r)

	return cookieToken, clientToken == cookieToken
}

func (s *Std) validateToken(w http.ResponseWriter, r *http.Request, pattern, clientToken string, config *csrf.Config) error {
	var token string
	tokenField := config.GetTokenBinding(pattern)
	if config.TokenMode == csrf.TokenModePerRequest {
		err := s.sessionManager.Flashes(w, r, tokenField, &token)
		if err != nil {
			return err
		}
	} else {
		err := s.sessionManager.Get(w, r, tokenField, &token)
		if err != nil {
			return err
		}
	}
	if clientToken != token {
		return fmt.Errorf("请求验证失败")
	}
	return nil
}
