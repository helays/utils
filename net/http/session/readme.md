# Session 模块使用指南

## 使用内存存储



```go
package test

import (
	"github.com/helays/utils/http/session"
	"github.com/helays/utils/http/session/memory"
	"time"
)

var (
	store = &session.Store{}
)

func run() {

	store = session.Init(memory.New(), &session.Options{
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

```



## 使用文件存储

```go
package test

import (
	"github.com/helays/utils/http/session"
	"github.com/helays/utils/http/session/file"
	"time"
)

var (
	store = &session.Store{}
)

type User struct{}

func run() {
	engine, _ := file.New(file.Instance{Path: "runtime/session"})
	// 在session中需要存储User 结构体数据，需要将结构体注册进去
	engine.Register(User{})
	store = session.Init(engine, &session.Options{
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

```

