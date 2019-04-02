package test

import (
	"os"
	"testing"
    "openpitrix.io/watcher/pkg/config"
)

func TestLoadConf(t *testing.T) {
	t.Logf("Test load configurations...")

	os.Setenv("HANDLER", "updateOpenpitrixName")
	os.Setenv("FILE", "global_config.yaml")
	os.Setenv("ETCDCONFIG_PREFIX", "openpitrix")
	os.Setenv("ETCDCONFIG_ENDPOINTS", "openpitrix-etcd:2379")

    config.LoadConf()

	t.Log("Test LoadConf successfully!")

}
