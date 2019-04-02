package config

import (
	"github.com/koding/multiconfig"
	"openpitrix.io/openpitrix/pkg/etcd"
	"openpitrix.io/openpitrix/pkg/logger"
)

const (
	UPDATE_OP_ETCD = "UPDATE_OP_ETCD"
)

var HANDLERS = map[string]string{
	UPDATE_OP_ETCD: "updateOpenpitrixEtcd",
}

type Pi struct {
	GlobalConfig *Config
	EtcdClient   *etcd.Etcd
}

type Config struct {
	WatchedFile string `default:"global_config.yaml"`   //The file that need to be watched
	Duration    int64    `default:10`                     //The duration for polling cycle which repeats
	Handler     string `default:"updateOpenpitrixEtcd"` //The action func name to run when files change
	LogLevel    string `default:"info"`
	Etcd        *EtcdConfig
}

var Global *Pi

func LoadConf() {
	m := multiconfig.New()
	//get config
	config := &Config{}
	err := m.Load(config)

	if err != nil {
		logger.Critical(nil, "Failed to load config: %+v", err)
		panic(err)
	}
	logger.SetLevelByString(config.LogLevel)
	logger.Debug(nil, "LoadConf: %+v", config)

	Global = &Pi{
		GlobalConfig: config,
		EtcdClient:   openEtcd(*config.Etcd),
	}

}
