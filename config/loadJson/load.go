package loadJson

import (
	"encoding/json"
	"fmt"
	"os"

	"helay.net/go/utils/v3/close/osClose"
	"helay.net/go/utils/v3/config"
	"helay.net/go/utils/v3/tools"
)

func LoadJson(i any) {
	if err := LoadJsonBase(i); err != nil {
		panic(fmt.Errorf("解析配置文件失败 %v", err))
	}
}

func LoadJsonBase(i any) error {
	reader, err := os.Open(tools.Fileabs(config.Cpath))
	defer osClose.CloseFile(reader)
	if err != nil {
		return fmt.Errorf("打开配置文件失败：%s", err.Error())
	}
	y := json.NewDecoder(reader)
	return y.Decode(i)
}
