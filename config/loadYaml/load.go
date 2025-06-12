package loadYaml

import (
	"fmt"
	"github.com/helays/utils/close/osClose"
	"github.com/helays/utils/config"
	"github.com/helays/utils/logger/ulogs"
	"github.com/helays/utils/tools"
	"gopkg.in/yaml.v3"
	"os"
)

// TODO 后续可使用 标准库里面的 模板引擎预处理 进行嵌套引用。

func LoadYaml(i any) {
	ulogs.DieCheckerr(LoadYamlBase(i), "解析配置文件失败")
}

func LoadYamlBase(i any) error {
	reader, err := os.Open(tools.Fileabs(config.Cpath))
	defer osClose.CloseFile(reader)
	if err != nil {
		return fmt.Errorf("打开配置文件失败：%s", err.Error())
	}
	y := yaml.NewDecoder(reader)
	return y.Decode(i)
}
