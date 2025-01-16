package parseCmd

import (
	"flag"
	"fmt"
	"github.com/helays/utils"
	"github.com/helays/utils/config"
	"github.com/helays/utils/logger/ulogs"
	"os"
)

// Parseparams 解析启动参数
func Parseparams(f ...func()) {
	// 解析参数
	var (
		vers     bool
		logLevel string
	)
	flag.BoolVar(&config.Help, "h", false, "参数说明")
	flag.StringVar(&config.Cpath, "c", "conf.ini", "配置文件")
	flag.BoolVar(&config.Dbg, "debug", false, "Debug 模式")
	flag.StringVar(&logLevel, "log-level", "info", "日志级别:\ndebug info warn error fatal")
	flag.BoolVar(&vers, "version", false, "查看版本")
	if len(f) > 0 {
		for _, v := range f {
			if v != nil {
				v()
			}
		}
	}
	flag.Parse()
	if vers {
		fmt.Println(os.Args[0], utils.Version, utils.BuildTime)
		os.Exit(1)
	}
	if config.Help {
		flag.Usage()
		os.Exit(0)
	}
	if config.Dbg {
		logLevel = "debug"
	}
	// 控制日志等级
	switch logLevel {
	case "debug":
		ulogs.Level = ulogs.LogLevelDebug
	case "info":
		ulogs.Level = ulogs.LogLevelInfo
	case "warn":
		ulogs.Level = ulogs.LogLevelWarn
	case "error":
		ulogs.Level = ulogs.LogLevelError
	case "fatal":
		ulogs.Level = ulogs.LogLevelFatal
	}

	if config.EnableParseParamsLog {
		fmt.Println("日志级别", logLevel, ulogs.Level)
		ulogs.Log("运行参数解析完成...")
	}
}
