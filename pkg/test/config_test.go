package test

import (
	"strconv"
	"testing"

	"openpitrix.io/watcher/pkg/common"
)

func TestLoadConf(t *testing.T) {
	t.Logf("Test load configurations...")

	LocalEnv()
	common.LoadConf()
	global := common.Global

	if global.Handler != Envs[HANDLER] {
		t.Errorf("Handler value is wrong: %s, want: %s", global.Handler, Envs[HANDLER])
	}

	if global.WatchedFile != Envs[WATCHED_FILE] {
		t.Errorf("WatchedFile value is wrong: %s, want: %s", global.WatchedFile, Envs[WATCHED_FILE])
	}

	d, _ := strconv.ParseInt(Envs[DURATION], 10, 64)
	if global.Duration != d {
		t.Errorf("Duration value is wrong: %d, want: %s", global.Duration, Envs[DURATION])
	}

	if global.LogLevel != Envs[LOG_LEVEL] {
		t.Errorf("LogLevel value is wrong: %s, want: %s", global.LogLevel, Envs[LOG_LEVEL])
	}

	t.Log("Test LoadConf successfully!")
}
