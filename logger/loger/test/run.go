package main

import (
	"github.com/helays/utils/v2/config/loadIni"
	"github.com/helays/utils/v2/config/parseCmd"
	"github.com/helays/utils/v2/logger/loger"
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
