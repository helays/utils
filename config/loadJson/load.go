package loadJson

import (
	"encoding/json"
	"fmt"
	"helay.net/go/utils/v3/close/osClose"
	"helay.net/go/utils/v3/config"
	"helay.net/go/utils/v3/logger/ulogs"
	"helay.net/go/utils/v3/tools"
	"os"
)

func LoadJson(i any) {
	ulogs.DieCheckerr(LoadJsonBase(i), "解析配置文件失败")
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
