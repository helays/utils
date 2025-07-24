package httpServer

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
	"github.com/gorilla/handlers"
	"github.com/helays/utils/v2/config"
	"github.com/helays/utils/v2/crypto/md5"
	"github.com/helays/utils/v2/logger/ulogs"
	"github.com/helays/utils/v2/logger/zaploger"
	"github.com/helays/utils/v2/net/http/httpServer/request"
	"github.com/helays/utils/v2/net/http/httpServer/response"
	"github.com/helays/utils/v2/net/http/mime"
	"github.com/helays/utils/v2/net/ipAccess"
	"github.com/helays/utils/v2/tools"
	"net"
	"net/http"
	"os"
	"strings"
	"time"
)

// HttpServerStart 公功 http server 启动函数
func (h *HttpServer) HttpServerStart() {
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

	server := &http.Server{Addr: h.ListenAddr}
	if h.EnableGzip {
		server.Handler = handlers.CompressHandler(h.mux)
	} else {
		server.Handler = h.mux
	}
	defer Closehttpserver(server)
	ulogs.Log("启动Http(s) Server", h.ListenAddr)
	var err error
	if h.Ssl {
		server.TLSConfig = &tls.Config{
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
			server.TLSConfig.ClientCAs = pool
			server.TLSConfig.ClientAuth = tls.RequireAndVerifyClientCert
		}
		if err = server.ListenAndServeTLS(tools.Fileabs(h.Crt), tools.Fileabs(h.Key)); err != nil {
			ulogs.Error("HTTPS Service 服务启动失败", server.Addr, err)
			os.Exit(1)
		}
		return
	}
	go h.hotUpdate(server)
	var isQuit bool
	go h.stopServer(server, &isQuit)
	if err = server.ListenAndServe(); err != nil {
		if !errors.Is(err, http.ErrServerClosed) {
			ulogs.Error("HTTP Service 启动失败", server.Addr, err)
			os.Exit(1)
		}
	}
	if h.Hotupdate && !isQuit {
		h.HttpServerStart()
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

// 检测ip白名单和黑名单
func (this *HttpServer) checkIpAccess(w http.ResponseWriter, r *http.Request) bool {
	addr := request.Getip(r)
	ip, _, err := net.SplitHostPort(addr)
	if err != nil {
		ip = addr // 如果无端口，则直接使用原地址
	}
	if this.allowIpList != nil {
		if !this.allowIpList.Contains(ip) {
			response.Forbidden(w, "你的IP不在系统白名单内")
			return false
		}
		return true
	}
	if this.denyIpList != nil && this.denyIpList.Contains(ip) {
		response.Forbidden(w, "你的IP已被监管")
		return false
	}
	return true
}

func (h *HttpServer) debugIpAccess(w http.ResponseWriter, r *http.Request) bool {
	// 判断path是否以debug开头
	if !strings.HasPrefix(r.URL.Path, "/debug/") {
		return true
	}
	addr := request.Getip(r)
	ip, _, err := net.SplitHostPort(addr)
	if err != nil {
		ip = addr // 如果无端口，则直接使用原地址
	}
	if !h.debugAllowIpList.Contains(ip) {
		response.Forbidden(w, http.StatusText(http.StatusForbidden))
		return false
	}
	return true
}

// Closehttpserver 关闭http server
func Closehttpserver(s *http.Server) {
	if s != nil {
		_ = s.Close()
	}
}

// 用于检测参数变更，然后热更新。
func (this *HttpServer) hotUpdate(server *http.Server) {
	if this.Hotupdate {
		go func() {
			hash := this.hash()
			tck := time.NewTicker(1 * time.Second)
			defer tck.Stop()
			for range tck.C {
				if hash == this.hash() {
					continue
				}
				tck.Stop()
				break
			}
			Closehttpserver(server)
		}()
	}
}

// 计算 httpserver 模块摘要
func (this HttpServer) hash() string {
	strArr := append([]string{
		this.ListenAddr,
		this.Auth,
		tools.Booltostring(this.Ssl),
		this.Ca,
		this.Crt,
		this.Key,
		this.SocketTimeout.String(),
	}, append(this.Allowip, this.Denyip...)...)
	return md5.Md5string(strings.Join(strArr, ""))
}

func (this HttpServer) stopServer(server *http.Server, isQuit *bool) {
	config.SetEnableHttpServer(true)
	_ = <-config.CloseHttpserverSig
	*isQuit = true
	ulogs.Log("http server已关闭")
	Closehttpserver(server)
	config.CloseHttpserverSig <- 1
}
