package httpServer

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gorilla/handlers"
	"github.com/helays/utils/v2/close/httpClose"
	"github.com/helays/utils/v2/crypto/xxhashkit"
	"github.com/helays/utils/v2/logger/ulogs"
	"github.com/helays/utils/v2/net/http/mime"
	"github.com/helays/utils/v2/net/http/route/middleware"
	"github.com/helays/utils/v2/net/ipkit"
	"github.com/helays/utils/v2/net/ipmatch"
	"github.com/helays/utils/v2/safe"
	"github.com/helays/utils/v2/tools"
	"golang.org/x/net/websocket"
)

// HttpServerStart 公功 http server 启动函数
func (h *HttpServer) HttpServerStart(ctx context.Context) {
	var stop = safe.NewResourceRWMutex(false)
	go func() {
		<-ctx.Done()
		stop.Write(true)
		h.close()
	}()
	for {
		h.initParams(ctx)   // 初始化参数
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

func (h *HttpServer) close() {
	httpClose.Server(h.server)
	ulogs.Log("http server已关闭")
}

func (h *HttpServer) initParams(ctx context.Context) {
	mime.InitMimeTypes()
	h.serverNameMap = make(map[string]byte)
	for _, dom := range h.ServerName {
		h.serverNameMap[strings.ToLower(dom)] = 0
	}
	h.logger = middleware.NewResponseProcessor()
	if len(h.Logger.LogLevelConfigs) > 0 {
		logger, err := middleware.NewZapLogger(&h.Logger)
		if err != nil {
			panic(fmt.Errorf("日志模块初始化失败 %v", err))
		}
		h.logger.AddLogHandler(logger)
	}
	h.iptablesInit(ctx)
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

// ip 黑白名单初始化
func (h *HttpServer) iptablesInit(ctx context.Context) {
	now := time.Now()
	ulogs.Infof("开始初始化http server IP防火墙")
	defer ulogs.Infof("http server IP防火墙初始化完成，耗时：%v", time.Since(now))
	if !h.Security.IPAccess.Enable {
		return
	}
	var err error
	if h.Security.IPAccess.Allow != nil || len(h.Allowip) > 0 {
		h.allowIPMatch, err = ipmatch.NewIPMatcher(ctx, h.Security.IPAccess.Allow)
		ulogs.DieCheckerr(err, "http server ip白名单初始化失败")
		for _, ip := range h.Allowip {
			if ipkit.ISIPv4OrIPv6(ip) == "ipv4" {
				err = h.allowIPMatch.AddIPv4Rule(ip)
			} else {
				err = h.allowIPMatch.AddIPv6Rule(ip)
			}
			ulogs.DieCheckerr(err, "http server ip白名单初始化失败")
		}
		h.allowIPMatch.Build()
	}

	if h.Security.IPAccess.Deny != nil || len(h.Denyip) > 0 {
		h.denyIPMatch, err = ipmatch.NewIPMatcher(ctx, h.Security.IPAccess.Deny)
		ulogs.DieCheckerr(err, "http server ip黑名单初始化失败")
		for _, ip := range h.Denyip {
			if ipkit.ISIPv4OrIPv6(ip) == "ipv4" {
				err = h.denyIPMatch.AddIPv4Rule(ip)
			} else {
				err = h.denyIPMatch.AddIPv6Rule(ip)
			}
			ulogs.DieCheckerr(err, "http server ip黑名单初始化失败")
		}
		h.denyIPMatch.Build()
	}
	if h.Security.IPAccess.Debug != nil {
		h.debugIPMatch, err = ipmatch.NewIPMatcher(ctx, h.Security.IPAccess.Debug)
		ulogs.DieCheckerr(err, "http server ip调试名单初始化失败")
		h.debugIPMatch.Build()
	}

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
	strArr := []string{h.ListenAddr, tools.Booltostring(h.Ssl), h.Ca, h.Crt, h.Key, h.SocketTimeout.String()}
	return xxhashkit.XXHashString(strings.Join(strArr, ""))
}

func (h *HttpServer) addRoute(path string, handle http.Handler, cb ...MiddlewareFunc) {
	if h.route == nil {
		h.route = make(map[string]*routerRule)
	}
	h.route[path] = &routerRule{
		routeType: RouteTypeHTTP,
		cb:        cb,
		handle:    handle,
		path:      path,
	}
}

func (h *HttpServer) addWSRoute(path string, handle websocket.Handler) {
	if h.route == nil {
		h.route = make(map[string]*routerRule)
	}
	h.route[path] = &routerRule{
		routeType: RouteTypeWebSocket,
		wsHandle:  handle,
		path:      path,
	}
}
