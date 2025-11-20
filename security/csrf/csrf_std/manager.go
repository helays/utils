package csrf_std

import (
	"fmt"
	"net/http"

	"github.com/helays/utils/v2/map/safemap"
	"github.com/helays/utils/v2/net/http/httpServer/response"
	"github.com/helays/utils/v2/net/http/sessionmgr"
	"github.com/helays/utils/v2/security/csrf"
)

type Std struct {
	configs        safemap.SyncMap[string, *csrf.Config]
	sessionManager *sessionmgr.Manager
}

func NewStd(sm *sessionmgr.Manager) *Std {
	s := &Std{
		sessionManager: sm,
	}

	return s
}

// SetConfig 设置路径的CSRF配置
func (s *Std) SetConfig(pattern string, config *csrf.Config) {
	if config == nil {
		return
	}
	s.configs.Store(pattern, config)
}

func (s *Std) GetConfig(pattern string) (*csrf.Config, bool) {
	return s.configs.Load(pattern)
}

// Middleware 中间件
// 在接口运行过程中 实时查询csrf配置信息，并进行验证
func (s *Std) Middleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return s.WrapHandler(next.ServeHTTP)
	}
}

// WrapHandler 包装单个处理器
// csrf配置信息，一般时固定好了后不会再改变。
func (s *Std) WrapHandler(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		if config, exists := s.GetConfig(path); exists && config.ShouldValidate(r.Method) {
			clientToken, ok := s.validateDoubleTapStrategy(r)
			if !ok {
				response.SetReturnErrorDisableLog(w, fmt.Errorf("请求验证失败"), http.StatusForbidden)
				return
			}
			if config.Strategy == csrf.StrategyToken {
				if err := s.validateToken(w, r, path, clientToken, config); err != nil {
					response.SetReturnErrorDisableLog(w, err, http.StatusForbidden)
					return
				}
			}
		}
		handler(w, r)
	}
}

// TokenHandler 设置csrf token
// 前端调用的时候，切记不能每个接口都访问前都访问这个接口，否则会重复设置cookie，导致csrf验证失败。
// 根据接口的安全程度来设置调用频率即可。
func (s *Std) TokenHandler(w http.ResponseWriter, r *http.Request, path string, config *csrf.Config) (string, error) {
	if config == nil || !config.ShouldValidate(r.Method) {
		return "", nil
	}
	token := csrf.GenerateCSRFToken()
	tokenField := config.GetTokenBinding(path)

	// 需要将token在服务端进行存储
	if config.Strategy == csrf.StrategyToken {
		if err := s.sessionManager.Set(w, r, tokenField, token, config.Timeout); err != nil {
			return "", err
		}
	}
	// 设置cookie
	s.setDualCookies(w, config, path, token)
	return token, nil
}

// 设置cookie
func (s *Std) setDualCookies(w http.ResponseWriter, config *csrf.Config, pattern, token string) {
	path := "/"
	// 带有效期或者是每次访问前都需要获取一次token的，就使用独立path
	if config.TokenMode == csrf.TokenModePerRequest || config.TokenMode == csrf.TokenModeTimed {
		path = pattern
	}
	maxAge := int(config.Timeout.Seconds())
	http.SetCookie(w, &http.Cookie{
		Name:     csrf.DefaultCookieName,
		Value:    token,
		Path:     path,
		MaxAge:   maxAge,
		Secure:   config.Secure,
		HttpOnly: true,
		SameSite: config.SameSite,
	})

	http.SetCookie(w, &http.Cookie{
		Name:     csrf.DefaultStatusName,
		Value:    "valid",
		Path:     path,
		MaxAge:   maxAge,
		Secure:   config.Secure,
		HttpOnly: false,
		SameSite: config.SameSite,
	})
}
