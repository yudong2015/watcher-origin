package common

import (
	"openpitrix.io/openpitrix/pkg/logger"
    "fmt"
    "github.com/koding/multiconfig"
)

const CONFIG_PREFIX = "WATCHER"

type Config struct {
	WatchedFile string `default:"/opt/global_config.yaml"`   //The file that need to be watched
	Duration    int64  `default:"10"`                     //The duration for polling cycle which repeats
	Handler     string `default:"UpdateOpenpitrixEtcd"` //The action func name to run when files change
	LogLevel    string `default:"info"`
	Etcd        EtcdConfig
}

var global *Config

func Global() *Config {
    return global
}

func LoadConf() *Config {
	loader := multiconfig.MultiLoader(
	   &multiconfig.TagLoader{},
	   &multiconfig.EnvironmentLoader{Prefix: CONFIG_PREFIX, CamelCase: true},
     )
	//get config from env
	Global := &Config{}
	err := loader.Load(Global)
	if err != nil {
		errMsg := fmt.Sprintf("Failed to load config: %+v", err)
		panic(errMsg)
	}

	logger.SetLevelByString(Global.LogLevel)
	logger.Debug(nil, "LoadConf: %+v", Global)

	return Global
}
