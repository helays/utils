package loadAuto

import (
	"fmt"
	"log"
	"path/filepath"

	"helay.net/go/utils/v3/config"
	loadIni2 "helay.net/go/utils/v3/config/loadIni"
	loadJson2 "helay.net/go/utils/v3/config/loadJson"
	"helay.net/go/utils/v3/config/loadYaml"
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
		log.Printf("默认配置解析器解析失败，尝试其他解析器 %v\n", err)
	}
	for _, v := range loadFunc {
		err = v(i)
		if err == nil {
			return
		}
	}
	panic(fmt.Errorf("解析配置文件失败 %v", err))
}
