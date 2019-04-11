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

	if global.Handler != Envs[Handler] {
		t.Errorf("Handler value is wrong: %s, want: %s", global.Handler, Envs[Handler])
	}

	if global.WatchedFile != Envs[WatchedFile] {
		t.Errorf("WatchedFile value is wrong: %s, want: %s", global.WatchedFile, Envs[WatchedFile])
	}

	d, _ := strconv.ParseInt(Envs[Duration], 10, 64)
	if global.Duration != d {
		t.Errorf("Duration value is wrong: %d, want: %s", global.Duration, Envs[Duration])
	}

	if global.LogLevel != Envs[LogLevel] {
		t.Errorf("LogLevel value is wrong: %s, want: %s", global.LogLevel, Envs[LogLevel])
	}

	t.Log("Test LoadConf successfully!")
}
