package httpServer

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"errors"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gorilla/handlers"
	"github.com/helays/utils/v2/close/httpClose"
	"github.com/helays/utils/v2/crypto/xxhashkit"
	"github.com/helays/utils/v2/logger/ulogs"
	"github.com/helays/utils/v2/logger/zaploger"
	"github.com/helays/utils/v2/net/http/mime"
	"github.com/helays/utils/v2/net/ipAccess"
	"github.com/helays/utils/v2/tools"
	"github.com/helays/utils/v2/tools/mutex"
)

// HttpServerStart 公功 http server 启动函数
func (h *HttpServer) HttpServerStart(ctx context.Context) {
	var stop = mutex.NewSafeResourceRWMutex(false)
	go func() {
		<-ctx.Done()
		stop.Write(true)
		httpClose.Server(h.server)
		ulogs.Log("http server已关闭")
	}()
	for {
		h.initParams()      // 初始化参数
		go h.hotUpdate(ctx) // 启用热更新检测模块
		var err error
		ulogs.Log("启动Http(s) Server", h.ListenAddr)
		if h.Ssl {
			err = h.server.ListenAndServeTLS(tools.Fileabs(h.Crt), tools.Fileabs(h.Key))
		} else {
			err = h.server.ListenAndServe()
		}
		if err != nil {
			if !errors.Is(err, http.ErrServerClosed) {
				ulogs.Error("HTTP Service 启动失败", h.server.Addr, err)
				os.Exit(1)
			}
		}
		// 未启用热更新，或者检测到退出信号 就退出
		if !h.Hotupdate || stop.Read() {
			break
		}
	}

}

func (h *HttpServer) initParams() {
	mime.InitMimeTypes()
	h.serverNameMap = make(map[string]byte)
	for _, dom := range h.ServerName {
		h.serverNameMap[strings.ToLower(dom)] = 0
	}
	if h.Logger.LogLevelConfigs != nil {
		var err error
		h.logger, err = zaploger.New(&h.Logger)
		ulogs.DieCheckerr(err, "http server 日志模块初始化失败")
	}
	h.iptablesInit()
	h.initMux()
	h.initRouter()
	h.server = &http.Server{Addr: h.ListenAddr}
	if h.EnableGzip {
		h.server.Handler = handlers.CompressHandler(h.mux)
	} else {
		h.server.Handler = h.mux
	}

	if h.Ssl {
		h.server.TLSConfig = &tls.Config{
			CipherSuites: []uint16{
				tls.TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305,
				tls.TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305,
				tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
				tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
				tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
				tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
			},
		}
		// 如果包含ca证书，就需要做强制双向https 验证
		if h.Ca != "" {
			caCrt, err := os.ReadFile(tools.Fileabs(h.Ca))
			if err != nil {
				ulogs.Error("HTTPS Service Load Ca error", err)
				os.Exit(1)
			}
			pool := x509.NewCertPool()
			pool.AppendCertsFromPEM(caCrt)
			h.server.TLSConfig.ClientCAs = pool
			h.server.TLSConfig.ClientAuth = tls.RequireAndVerifyClientCert
		}
	}
}

func (h *HttpServer) initMux() {
	if h.mux == nil {
		h.mux = http.NewServeMux()
	}
}

// ip 黑白名单初始化
func (h *HttpServer) iptablesInit() {
	var err error
	if len(h.Allowip) > 0 {
		h.allowIpList, err = ipAccess.NewIPList(h.Allowip...)
		ulogs.DieCheckerr(err, "http server ip白名单初始化失败")
		h.enableCheckIpAccess = true
	}
	if len(h.Denyip) > 0 {
		h.denyIpList, err = ipAccess.NewIPList(h.Denyip...)
		ulogs.DieCheckerr(err, "http server ip黑名单初始化失败")
		h.enableCheckIpAccess = true
	}
	debugAllowIps := []string{"127.0.0.1"}
	if len(h.DebugAllowIp) > 0 {
		debugAllowIps = append(debugAllowIps, h.DebugAllowIp...)
	}
	h.debugAllowIpList, err = ipAccess.NewIPList(debugAllowIps...)
	ulogs.DieCheckerr(err, "http server debug ip白名单初始化失败")
}

// 用于检测参数变更，然后热更新。
func (h *HttpServer) hotUpdate(ctx context.Context) {
	if !h.Hotupdate {
		return
	}
	hash := h.hash()
	tck := time.NewTicker(1 * time.Second)
	defer tck.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-tck.C:
			if hash == h.hash() {
				continue
			}
			// 关闭client,重启程序
			httpClose.Server(h.server)
			return
		}
	}
}

// 计算 httpserver 模块摘要
func (h *HttpServer) hash() string {
	strArr := []string{h.ListenAddr, h.Auth, tools.Booltostring(h.Ssl), h.Ca, h.Crt, h.Key, h.SocketTimeout.String()}
	strArr = append(strArr, h.Allowip...)
	strArr = append(strArr, h.Denyip...)
	return xxhashkit.XXHashString(strings.Join(strArr, ""))
}
