package httpServer

import (
	"bytes"
	"github.com/dchest/captcha"
	"github.com/helays/utils/logger/ulogs"
	"io"
	"net/http"
	"time"
)

// 显示 favicon
func (ro Router) favicon(w http.ResponseWriter) {
	w.WriteHeader(200)
	rd := bytes.NewReader(favicon[:])
	_, _ = io.Copy(w, rd)
}

// BeforeAction 所有应用前置操作
func (ro *Router) BeforeAction(w http.ResponseWriter, r *http.Request) bool {
	if ro.CookiePath == "" {
		ro.CookiePath = "/"
	}
	// 在判断登录前，应该判断当前接口是否需要鉴权，否则就不读取下方的session
	// 登录这里不应该使用 GetUp更新session，
	// 不用登录的接口，这里就直接返回继续访问
	if !ro.validMustLogin(r.URL.Path) {
		return true
	}

	// 这里改用session 系统
	var loginInfo LoginInfo
	// 如果session 存在，那么当session 剩余24小时的时候，更新session。
	err := ro.Store.GetUpByTimeLeft(w, r, ro.SessionLoginName, &loginInfo, time.Hour*24)
	if err != nil || !loginInfo.IsLogin {
		ulogs.Checkerr(err, "session 获取失败")
		// 未登录的，终止请求，响应401 或者302
		return ro.unAuthorizedResp(w, r)
	}

	// 登录禁止访问的页面
	if ro.validDisableLoginRequestPath(r.URL.Path) {
		http.Redirect(w, r, ro.HomePage, 302)
		return false
	}
	// 控制管理员访问的
	if ro.validManagePage(r.URL.Path) && !loginInfo.IsManage {
		SetReturnCode(w, r, http.StatusForbidden, "无权限访问")
		return false
	}
	return true
}

func (ro Router) Captcha(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	w.Header().Set("Pragma", "no-cache")
	w.Header().Set("Expires", "0")
	var content bytes.Buffer

	// 验证码存储在session中
	captchaId := captcha.NewLen(4)
	if err := ro.Store.Set(w, r, "captcha", captchaId, 4*time.Minute); err != nil {
		InternalServerError(w)
		return
	}

	if err := captcha.WriteImage(&content, captchaId, 106, 40); err != nil {
		InternalServerError(w)
		ulogs.Error(err, "captcha writeImage")
		return
	}
	w.Header().Set("Content-Type", "image/png")
	http.ServeContent(w, r, "", time.Time{}, bytes.NewReader(content.Bytes()))
}

// Middleware 中间件
func (ro *Router) Middleware(w http.ResponseWriter, r *http.Request, f func(w http.ResponseWriter, r *http.Request, ro *Router)) {
	f(w, r, ro)
}
