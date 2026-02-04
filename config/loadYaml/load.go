package loadYaml

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
	"helay.net/go/utils/v3/close/osClose"
	"helay.net/go/utils/v3/config"
	"helay.net/go/utils/v3/logger/ulogs"
	"helay.net/go/utils/v3/tools"
	"helay.net/go/utils/v3/tools/fileinclude"
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
