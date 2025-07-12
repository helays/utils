package signalTools

import (
	"github.com/helays/utils/v2/config"
	"github.com/helays/utils/v2/logger/ulogs"
	"os"
	"os/signal"
	"syscall"
)

// SignalHandle 系统信号
// @var funds 结束服务前，需要执行的操作
func SignalHandle(funds ...func()) {
	exitsin := make(chan os.Signal)
	//  syscall.SIGHUP 用作重载信号
	signal.Notify(exitsin, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGTERM) // 注意，syscall.SIGKILL 不能被捕获
	ulogs.Log("退出信号", <-exitsin)                                             // 日志记录
	exit()

}

func exit(funds ...func()) {
	ulogs.Log("执行回调", "数量", len(funds))
	for _, f := range funds {
		if f == nil {
			continue
		}
		f()
	}
	ulogs.Log("各个组件关闭完成，系统即将自动关闭", os.Getpid())
	if config.GetIsEnableHttpServer() {
		config.CloseHttpserverSig <- 1
		_ = <-config.CloseHttpserverSig
	}
	os.Exit(0)
}
