package loadIni

import (
	"github.com/helays/utils/v2/config"
	"github.com/helays/utils/v2/logger/ulogs"
	"github.com/helays/utils/v2/tools"
	"gopkg.in/ini.v1"
)

func LoadIni(i any) {
	ulogs.DieCheckerr(LoadIniBase(i), "载入配置文件失败")
}

// LoadIniBase 载入配置基础功能
func LoadIniBase(i any) error {
	return ini.MapTo(i, tools.Fileabs(config.Cpath))
}
