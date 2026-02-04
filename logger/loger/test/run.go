package main

import (
	"helay.net/go/utils/v3/config/loadIni"
	"helay.net/go/utils/v3/config/parseCmd"
	"helay.net/go/utils/v3/logger/loger"
	"strings"
	"time"
)

type config struct {
	Log loger.Loger
}

func main() {
	var log = new(config)
	parseCmd.Parseparams(nil)
	loadIni.LoadIni(log)

	loger.Init(log.Log)
	go func() {
		for {
			log.Log.Error(time.Now().Unix())
		}
	}()
	for {
		log.Log.Log(strings.Repeat(time.Now().String(), 2))
	}
}
