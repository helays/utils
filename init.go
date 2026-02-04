package utils

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/helays/utils/v2/config"
)

var once sync.Once

func init() {
	var err error
	once.Do(func() {
		config.Appath, err = filepath.Abs(filepath.Dir(os.Args[0]))
		if err != nil {
			panic(fmt.Errorf("获取系统运行目录失败 %v", err))
		}
	})

}
