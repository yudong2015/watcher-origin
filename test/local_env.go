package test

import (
	"os"
)

const (
	Handler       = "WATCHER_HANDLER"
	WatchedFile   = "WATCHER_WATCHED_FILE"
	Duration      = "WATCHER_DURATION"
	LogLevel      = "WATCHER_LOG_LEVEL"
	EtcdPrefix    = "WATCHER_ETCD_PREFIX"
	EtcdEndpoints = "WATCHER_ETCD_ENDPOINTS"
)

var Envs = map[string]string{
	Handler:       "UpdateOpenPitrixEtcd",
	WatchedFile:   "./test/global_config.yaml",
	Duration:      "5",
	LogLevel:      "debug",
	EtcdPrefix:    "openpitrix",
	EtcdEndpoints: "127.0.0.1:2379",
}

func LocalEnv() {
	for k, v := range Envs {
		os.Setenv(k, v)
	}
}
