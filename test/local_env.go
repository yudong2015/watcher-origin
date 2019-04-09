package test

import (
	"os"
)

const (
	HANDLER        = "WATCHER_HANDLER"
	WATCHED_FILE   = "WATCHER_WATCHED_FILE"
	DURATION       = "WATCHER_DURATION"
	LOG_LEVEL      = "WATCHER_LOG_LEVEL"
	ETCD_PREFIX    = "WATCHER_ETCD_PREFIX"
	ETCD_ENDPOINTS = "WATCHER_ETCD_ENDPOINTS"
)

var Envs = map[string]string{
	HANDLER:        "UpdateOpenpitrixEtcd",
	WATCHED_FILE:   "./test/global_config.yaml",
	DURATION:       "5",
	LOG_LEVEL:      "debug",
	ETCD_PREFIX:    "openpitrix",
	ETCD_ENDPOINTS: "127.0.0.1:2379",
}

func LocalEnv() {
	for k, v := range Envs {
		os.Setenv(k, v)
	}
}
