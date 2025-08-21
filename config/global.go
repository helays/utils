package config

import (
	"errors"
	"math/rand/v2"
	"net/http"
	"sync"
	"time"
)

var RandPool = sync.Pool{
	New: func() interface{} {
		return rand.New(rand.NewPCG(uint64(time.Now().UnixNano()), uint64(time.Now().UnixNano())))
	},
}
var (
	Help   bool   // 打印显示帮助信息
	Cpath  string // 配置文件路径
	Appath string // 当前路径
	Dbg    bool   // Debug 模式

	CstSh = time.FixedZone("CST", 8*3600) // 东八区

	PublicKeyByt  []byte // 公钥
	PrivateKeyByt []byte // 私钥

	CloseHttpserverSig   = make(chan byte, 1)
	EnableParseParamsLog = true
)

// 用于控制 是否启用http server的
var (
	enableHttpserver   bool
	enableHttpserverMu sync.RWMutex
)

var (
	ErrNotFound = errors.New(http.StatusText(http.StatusNotFound))
)

func GetIsEnableHttpServer() bool {
	enableHttpserverMu.RLock()
	defer enableHttpserverMu.RUnlock()
	return enableHttpserver
}

func SetEnableHttpServer(b bool) {
	enableHttpserverMu.Lock()
	defer enableHttpserverMu.Unlock()
	enableHttpserver = b
}
