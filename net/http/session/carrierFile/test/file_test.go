package test

import (
	session2 "github.com/helays/utils/v2/net/http/session"
	"github.com/helays/utils/v2/net/http/session/carrierFile"
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
// Date: 2024/12/8 17:11
//

var (
	store = &session2.Store{}
)

type User struct{}

func run() {
	defer session2.Close(store)
	engine, _ := carrierFile.New(carrierFile.Instance{Path: "runtime/session"})
	// 在session中需要存储User 结构体数据，需要将结构体注册进去
	// 需要在session 初始化之前进行注册
	engine.Register(User{})

	store = session2.Init(engine, &session2.Options{
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
