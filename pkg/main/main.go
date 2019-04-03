package main

import (
	"time"

	"github.com/radovskyb/watcher"
	"openpitrix.io/openpitrix/pkg/logger"
	"openpitrix.io/watcher/pkg/config"
)

const (
    UPDATE_OP_ETCD = "updateOpenpitrixEtcd"
)

var HANDLERS = map[string]func(){
    UPDATE_OP_ETCD: UpdateOpenpitrixEtcd,
}

func main() {
	config.LoadConf()
	watch()
}

func watch() {
	w := watcher.New()
	w.SetMaxEvents(1)
	w.FilterOps(watcher.Write, watcher.Create)
	global := config.Global
	w.Add(global.GlobalConfig.WatchedFile)

	go func() {
		for {
			select {
			case event := <-w.Event:
				logger.Info(nil, "%+v", event)
			    HANDLERS[global.GlobalConfig.Handler]()
			case err := <-w.Error:
				panic(err)
			case <-w.Closed:
				logger.Error(nil, "The watcher is closed!")
				return
			}
		}
	}()

	duration := time.Duration(global.GlobalConfig.Duration) * time.Second

	err := w.Start(duration)

	if err != nil {
		logger.Critical(nil, "Failed to start watching: %+v", err)
		panic(err)
	}
}
