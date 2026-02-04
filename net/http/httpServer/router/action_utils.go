package router

import (
	"net/http"
	"regexp"

	"github.com/pkg/errors"
	"helay.net/go/utils/v3/config"
	"helay.net/go/utils/v3/logger/ulogs"
	"helay.net/go/utils/v3/net/http/response"
	"helay.net/go/utils/v3/net/http/route"
)

//
// ━━━━━━神兽出没━━━━━━
// 　　 ┏┓     ┏┓
// 　　┏┛┻━━━━━┛┻┓
// 　　┃　　　　　 ┃
// 　　┃　　━　　　┃
// 　　┃　┳┛　┗┳  ┃
// 　　┃　　　　　 ┃
// 　　┃　　┻　　　┃
// 　　┃　　　　　 ┃
// 　　┗━┓　　　┏━┛　Code is far away from bug with the animal protecting
// 　　　 ┃　　　┃    神兽保佑,代码无bug
// 　　　　┃　　　┃
// 　　　　┃　　　┗━━━┓
// 　　　　┃　　　　　　┣┓
// 　　　　┃　　　　　　┏┛
// 　　　　┗┓┓┏━┳┓┏┛
// 　　　　 ┃┫┫ ┃┫┫
// 　　　　 ┗┻┛ ┗┻┛
//
// ━━━━━━感觉萌萌哒━━━━━━
//
//
// User helay
// Date: 2024/11/30 15:00
//

// SwitchDebug 切换调试模式开关
func SwitchDebug(w http.ResponseWriter, r *http.Request) {
	config.Dbg = !config.Dbg
	response.SetReturnData(w, 0, "成功")
}

// 辅助函数，用于匹配正则表达式
func (ro *Router) matchRegexp(path string, rules []*regexp.Regexp) bool {
	for _, rule := range rules {
		if rule.MatchString(path) {
			return true
		}
	}
	return false
}

// 验证是否需要登录才能访问，
// 首先要判断 免登录权限，优先级高
// 然后判断 是否要登录
func (ro *Router) validMustLogin(path string) bool {
	if ro.UnLoginPath[path] {
		return false // 不用登录
	}
	if ro.UnLoginPathRegexp != nil && ro.matchRegexp(path, ro.UnLoginPathRegexp) {
		return false
	}
	if ro.MustLoginPath[path] {
		return true
	}
	if ro.MustLoginPathRegexp != nil && ro.matchRegexp(path, ro.MustLoginPathRegexp) {
		return true
	}
	return false
}

// 无授权的响应
func (ro *Router) unAuthorizedResp(w http.ResponseWriter, r *http.Request) bool {
	if ro.UnauthorizedRespMethod == 401 {
		response.SetReturnData(w, ro.UnauthorizedRespMethod, "未登录，请先登录！！")
		return false
	}
	http.Redirect(w, r, ro.LoginPath, 302)
	return false
}

// 验证是否是登录后就禁止访问的页面
func (ro *Router) validDisableLoginRequestPath(path string) bool {
	if ro.DisableLoginPath[path] {
		return true
	}
	if ro.DisableLoginPathRegexp != nil && ro.matchRegexp(path, ro.DisableLoginPathRegexp) {
		return true
	}
	return false
}

// 验证是否是管理页面
func (ro *Router) validManagePage(path string) bool {
	if ro.ManagePage[path] {
		return true
	}
	if ro.ManagePageRegexp != nil && ro.matchRegexp(path, ro.ManagePageRegexp) {
		return true
	}
	return false
}

func (ro *Router) errorWithLog(w http.ResponseWriter, resp route.ErrorResp) {
	if ro.ErrorWithLog != nil {

		ro.ErrorWithLog(w, &resp)
		return
	}
	ulogs.Error("errorWithLog:", errors.WithStack(resp.Error))
	route.RenderErrorText(w, &resp)
}
