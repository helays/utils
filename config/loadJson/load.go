package loadJson

import (
	"encoding/json"
	"fmt"
	"github.com/helays/utils/v2/close/osClose"
	"github.com/helays/utils/v2/config"
	"github.com/helays/utils/v2/logger/ulogs"
	"github.com/helays/utils/v2/tools"
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
