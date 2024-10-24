package loadIni

import (
	"github.com/xuexila/utils/config"
	"github.com/xuexila/utils/tools"
	"github.com/xuexila/utils/ulogs"
	"gopkg.in/ini.v1"
	"os"
)

func LoadIni(i any) {
	if err := ini.MapTo(i, tools.Fileabs(config.Cpath)); err != nil {
		ulogs.Error("载入配置文件错误", err.Error())
		os.Exit(1)
	}
}