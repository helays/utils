package loadIni

import (
	"gopkg.in/ini.v1"
	"helay.net/go/utils/v3/config"
	"helay.net/go/utils/v3/logger/ulogs"
	"helay.net/go/utils/v3/tools"
)

func LoadIni(i any) {
	ulogs.DieCheckerr(LoadIniBase(i), "载入配置文件失败")
}

// LoadIniBase 载入配置基础功能
func LoadIniBase(i any) error {
	return ini.MapTo(i, tools.Fileabs(config.Cpath))
}
