package main

import (
	"reflect"
	"time"

	"github.com/radovskyb/watcher"
	"openpitrix.io/openpitrix/pkg/logger"
	"openpitrix.io/watcher/pkg/config"
)

func main() {
	config.LoadConf()
	watch()
}

func watch() {
	w := watcher.New()
	w.SetMaxEvents(1)
	w.FilterOps(watcher.Write, watcher.Create)
	w.Add(config.Global.GlobalConfig.WatchedFile)

	go func() {
		for {
			select {
			case event := <-w.Event:
				logger.Info(nil, "%+v", event)
				reflect.ValueOf()
			case err := <-w.Error:
				panic(err)
			case <-w.Closed:
				logger.Error(nil, "The watcher is closed!")
				return
			}
		}
	}()

	var duration time.Duration = time.Duration(config.Global.GlobalConfig.Duration) * time.Second

	err := w.Start(duration)

	if err != nil {
		logger.Critical(nil, "Failed to start watching: %+v", err)
		panic(err)
	}
}
