package httpServer

import (
	"github.com/helays/utils/config"
	"net/http"
	"regexp"
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
	SetReturnData(w, 0, "成功")
}

// 辅助函数，用于匹配正则表达式
func (this *Router) matchRegexp(path string, rules []*regexp.Regexp) bool {
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
func (this *Router) validMustLogin(path string) bool {
	if this.UnLoginPath[path] {
		return false // 不用登录
	}
	if this.UnLoginPathRegexp != nil && this.matchRegexp(path, this.UnLoginPathRegexp) {
		return false
	}
	if this.MustLoginPath[path] {
		return true
	}
	if this.MustLoginPathRegexp != nil && this.matchRegexp(path, this.MustLoginPathRegexp) {
		return true
	}
	return false
}

// 无授权的响应
func (this *Router) unAuthorizedResp(w http.ResponseWriter, r *http.Request) bool {
	if this.UnauthorizedRespMethod == 401 {
		SetReturnData(w, this.UnauthorizedRespMethod, "未登录，请先登录！！")
		return false
	}
	http.Redirect(w, r, this.LoginPath, 302)
	return false
}

// 验证是否是登录后就禁止访问的页面
func (this *Router) validDisableLoginRequestPath(path string) bool {
	if this.DisableLoginPath[path] {
		return true
	}
	if this.DisableLoginPathRegexp != nil && this.matchRegexp(path, this.DisableLoginPathRegexp) {
		return true
	}
	return false
}

// 验证是否是管理页面
func (this *Router) validManagePage(path string) bool {
	if this.ManagePage[path] {
		return true
	}
	if this.ManagePageRegexp != nil && this.matchRegexp(path, this.ManagePageRegexp) {
		return true
	}
	return false
}
