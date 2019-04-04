package test

import (
	"os"
	"testing"
    "openpitrix.io/watcher/pkg/common"
    "strconv"
)

func TestLoadConf(t *testing.T) {
	t.Logf("Test load configurations...")

    const (
        HANDLER = "WATCHER_HANDLER"
        WATCHED_FILE = "WATCHER_WATCHED_FILE"
        DURATION = "WATCHER_DURATION"
        LOG_LEVEL = "WATCHER_LOG_LEVEL"
        ETCD_PREFIX = "WATCHER_ETCD_PREFIX"
        ETCD_ENDPOINTS = "WATCHER_ETCD_ENDPOINTS"
    )
	ENVs := map[string]string {
        HANDLER: "UpdateOpenpitrixName",
        WATCHED_FILE: "../test/global_config.yaml",
        DURATION: "30",
        LOG_LEVEL: "debug",
        ETCD_PREFIX: "openpitr",
        ETCD_ENDPOINTS: "openpitrix-etcd:237",
    }

    for k, v := range ENVs {
        os.Setenv(k, v)
    }
    common.LoadConf()
    global := common.Global

    if global.Handler!=ENVs[HANDLER] {
        t.Errorf("Handler value is wrong: %s, want: %s", global.Handler, ENVs[HANDLER])
    }

    if global.WatchedFile!=ENVs[WATCHED_FILE] {
        t.Errorf("WatchedFile value is wrong: %s, want: %s", global.WatchedFile, ENVs[WATCHED_FILE])
    }

    d,_ := strconv.ParseInt(ENVs[DURATION], 10, 64)
    if global.Duration!=d {
        t.Errorf("Duration value is wrong: %d, want: %s", global.Duration, ENVs[DURATION])
    }

    if global.LogLevel!=ENVs[LOG_LEVEL] {
        t.Errorf("LogLevel value is wrong: %s, want: %s", global.LogLevel, ENVs[LOG_LEVEL])
    }

	t.Log("Test LoadConf successfully!")
}
