package loadIni

import (
	"fmt"

	"gopkg.in/ini.v1"
	"helay.net/go/utils/v3/config"
	"helay.net/go/utils/v3/tools"
)

func LoadIni(i any) {
	if err := LoadIniBase(i); err != nil {
		panic(fmt.Errorf("解析配置文件失败 %v", err))
	}
}

// LoadIniBase 载入配置基础功能
func LoadIniBase(i any) error {
	return ini.MapTo(i, tools.Fileabs(config.Cpath))
}
