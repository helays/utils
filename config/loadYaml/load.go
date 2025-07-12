package loadYaml

import (
	"fmt"
	"github.com/helays/utils/v2/close/osClose"
	"github.com/helays/utils/v2/config"
	"github.com/helays/utils/v2/logger/ulogs"
	"github.com/helays/utils/v2/tools"
	"github.com/helays/utils/v2/tools/fileinclude"
	"gopkg.in/yaml.v3"
	"io"
	"os"
	"path/filepath"
)

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
