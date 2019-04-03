package config

import (
	"strings"

	"openpitrix.io/openpitrix/pkg/etcd"
	"openpitrix.io/openpitrix/pkg/logger"
)


const (
	GlobalConfigKey = "global_config"
	DlockKey        = "dlock_" + GlobalConfigKey
)

type EtcdConfig struct {
	Prefix string `default:"openpitrix"`
	Endpoints  string `default:"openpitrix-etcd:2379"` // Example: "localhost:2379,localhost:22379,localhost:32379"
}


func (config EtcdConfig)openEtcd() *etcd.Etcd {
	endpoints := strings.Split(config.Endpoints, ",")
	etcd, err := etcd.Connect(endpoints, config.Prefix)
	if err != nil {
		logger.Critical(nil, "failed to connect etcd")
		panic(err)
	}
	return etcd
}
