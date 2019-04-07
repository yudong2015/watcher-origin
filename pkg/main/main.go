package main

import (
	"time"

	"fmt"
	"github.com/radovskyb/watcher"
	"openpitrix.io/openpitrix/pkg/logger"
	"openpitrix.io/watcher/pkg/common"
)

const (
	UPDATE_OP_ETCD = "UpdateOpenpitrixEtcd"
)

func watch() {
	w := watcher.New()
	w.SetMaxEvents(1)
	w.FilterOps(watcher.Write, watcher.Create)

	global := common.Global
	w.Add(global.WatchedFile)
	logger.Info(nil, "Watching file: %s", global.WatchedFile)

	go func() {
		for {
			select {
			case event := <-w.Event:
				logger.Info(nil, "%+v", event)
				handle()
			case err := <-w.Error:
				panic(err)
			case <-w.Closed:
				logger.Error(nil, "The watcher is closed!")
				return
			}
		}
	}()

	duration := time.Duration(global.Duration) * time.Second

	err := w.Start(duration)

	if err != nil {
		logger.Critical(nil, "Failed to start watching: %+v", err)
		panic(err)
	}
}

func handle() {
	switch common.Global.Handler {
	case UPDATE_OP_ETCD:
		UpdateOpenpitrixEtcd()
	default:
		msg := fmt.Sprintf("The func %s not exist!", common.Global.Handler)
		panic(msg)
	}
}

func main() {
	common.LoadConf()
	handle()
	watch()
}
