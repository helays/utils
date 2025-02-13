package loadAuto

import (
	"github.com/helays/utils/config"
	loadIni2 "github.com/helays/utils/config/loadIni"
	loadJson2 "github.com/helays/utils/config/loadJson"
	"github.com/helays/utils/config/loadYaml"
	"github.com/helays/utils/logger/ulogs"
	"path/filepath"
)

var (
	loadFunc = map[string]func(i any) error{
		".ini":  loadIni2.LoadIniBase,
		".json": loadJson2.LoadJsonBase,
		".yaml": loadYaml.LoadYamlBase,
	}
)

// Load 载入配置文件
func Load[T any](i T) {
	ext := filepath.Ext(config.Cpath)
	var err error
	loadFirst, ok := loadFunc[ext]
	if ok {
		delete(loadFunc, ext)
		if err = loadFirst(i); err == nil {
			return
		}
	}
	for _, v := range loadFunc {
		err = v(i)
		if err == nil {
			return
		}
	}
	ulogs.DieCheckerr(err, "载入配置文件失败")
}
