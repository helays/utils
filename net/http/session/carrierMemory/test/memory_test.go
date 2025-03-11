package test

import (
	session2 "github.com/helays/utils/net/http/session"
	"github.com/helays/utils/net/http/session/carrierMemory"
	"time"
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
// Date: 2024/12/8 17:07
//

var (
	store = &session2.Store{}
)

func run() {
	defer session2.Close(store)
	store = session2.Init(carrierMemory.New(), &session2.Options{
		CookieName:    "vsclub.ltd",
		CheckInterval: time.Hour,
		Carrier:       "cookie",
		Path:          "",
		Domain:        "",
		MaxAge:        0,
		Secure:        false,
		HttpOnly:      false,
		SameSite:      0,
	})
}
