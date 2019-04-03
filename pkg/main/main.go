package main

import (
	"time"

	"github.com/radovskyb/watcher"
	"openpitrix.io/openpitrix/pkg/logger"
	"openpitrix.io/watcher/pkg/common"
    "fmt"
)

const (
    UPDATE_OP_ETCD = "UpdateOpenpitrixEtcd"
)

func main() {
	common.LoadConf()
	watch()
}

func watch() {
	w := watcher.New()
	w.SetMaxEvents(1)
	w.FilterOps(watcher.Write, watcher.Create)
	global := common.Global()
	w.Add(global.WatchedFile)

	go func() {
		for {
			select {
			case event := <-w.Event:
				logger.Info(nil, "%+v", event)
                switch global.Handler {
                case UPDATE_OP_ETCD:
                    UpdateOpenpitrixEtcd()
                default:
                    msg := fmt.Sprintf( "The func %s not exist!", global.Handler)
                    panic(msg)
                }
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
