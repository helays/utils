package loadYaml

import (
	"fmt"
	"github.com/helays/utils/close/osClose"
	"github.com/helays/utils/config"
	"github.com/helays/utils/logger/ulogs"
	"github.com/helays/utils/tools"
	"github.com/helays/utils/tools/fileinclude"
	"gopkg.in/yaml.v3"
	"io"
	"os"
	"path/filepath"
)

// TODO 后续可使用 标准库里面的 模板引擎预处理 进行嵌套引用。

func LoadYaml(i any) {
	ulogs.DieCheckerr(LoadYamlBase(i), "解析配置文件失败")
}

func LoadYamlBase(i any) error {
	reader, err := fileInclude()
	if err != nil {
		return err
	}

	y := yaml.NewDecoder(reader)
	return y.Decode(i)
}

func fileInclude() (io.Reader, error) {
	p := tools.Fileabs(config.Cpath)
	reader, err := os.Open(p)
	defer osClose.CloseFile(reader)
	if err != nil {
		return nil, fmt.Errorf("打开配置文件失败：%s", err.Error())
	}
	inc := fileinclude.NewProcessor()
	inc.SetPrefix("#include ")
	inc.FromReader(reader, filepath.Dir(p))
	return inc.ToReader()
}
