package route

import (
	"embed"
	"fmt"
)

//go:embed resources/*
var resourceFS embed.FS

var favicon []byte

func init() {
	var err error
	if favicon, err = resourceFS.ReadFile("resources/favicon.ico"); err != nil {
		panic(fmt.Errorf("resources/favicon 载入失败 %v", err))
	}
}
