package router

import (
	"net/http"
	"time"

	"github.com/helays/utils/v2/logger/ulogs"
	"github.com/helays/utils/v2/net/http/httpServer/response"
)

// BeforeAction 所有应用前置操作
// Deprecated: 请使用 http server 中间件模式
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
		response.SetReturnCode(w, r, http.StatusForbidden, "无权限访问")
		return false
	}
	return true
}

// Middleware 中间件
// Deprecated: 请使用 http server 中间件模式
func (ro *Router) Middleware(w http.ResponseWriter, r *http.Request, f func(w http.ResponseWriter, r *http.Request, ro *Router)) {
	f(w, r, ro)
}
