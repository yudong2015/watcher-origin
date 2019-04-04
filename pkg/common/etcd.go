package common

import (
	"strings"

	"openpitrix.io/openpitrix/pkg/etcd"
	"openpitrix.io/openpitrix/pkg/logger"
	"time"
)


const (
	GlobalConfigKey = "global_config"
	DlockKey        = "dlock_" + GlobalConfigKey
	EtcdDlockTimeOut = time.Second * 60
)

type EtcdConfig struct {
	Prefix string `default:"openpitrix"`
	Endpoints  string `default:"openpitrix-etcd:2379"` // Example: "localhost:2379,localhost:22379,localhost:32379"
}


func (config *EtcdConfig)OpenEtcd() *etcd.Etcd {
	endpoints := strings.Split(config.Endpoints, ",")
	logger.Info(nil,"Start to open etcd: %v...", config)
	client, err := etcd.Connect(endpoints, config.Prefix)
	if err != nil {
		logger.Critical(nil, "failed to connect etcd")
		panic(err)
	}
	logger.Info(nil,"Opened etcd.")
	return client
}
