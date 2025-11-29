package route

import (
	"embed"
	"fmt"
)

// go:embed resources/*
var resourceFS embed.FS

func init() {
	if fs, err := resourceFS.Open("resources/favicon.ico"); err != nil {
		panic(fmt.Errorf("resources/favicon 载入失败 %v", err))
	} else {
		if _, err = fs.Read(favicon); err != nil {
			panic(fmt.Errorf("resources/favicon 读取失败 %v", err))
		}
	}

}
