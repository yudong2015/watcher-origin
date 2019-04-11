package common

import (
	"fmt"

	"github.com/koding/multiconfig"

	"openpitrix.io/openpitrix/pkg/logger"
)

const ConfigPrefix = "WATCHER"

type Config struct {
	WatchedFile string `default:"/opt/global_config.yaml"` //The file that need to be watched
	Duration    int64  `default:"10"`                      //The duration for polling cycle which repeats
	Handler     string `default:"UpdateOpenPitrixEtcd"`    //The action func name to run when files change
	LogLevel    string `default:"info"`
	Etcd        *Etcd
}

var Global = new(Config)

func LoadConf() {
	loader := multiconfig.MultiLoader(
		&multiconfig.TagLoader{},
		&multiconfig.EnvironmentLoader{Prefix: ConfigPrefix, CamelCase: true},
	)
	//get config from env
	Global.Etcd = &Etcd{}
	err := loader.Load(Global.Etcd)
	if err != nil {
		errMsg := fmt.Sprintf("Failed to load etcd config: %+v", err)
		panic(errMsg)
	}
	err = loader.Load(Global)
	if err != nil {
		errMsg := fmt.Sprintf("Failed to load global config: %+v", err)
		panic(errMsg)
	}

	logger.SetLevelByString(Global.LogLevel)
	logger.Debug(nil, "Etcd config: %+v", Global.Etcd)
	logger.Debug(nil, "Global config: %+v", Global)
}

type NilError struct {
	msg string
}

func (e NilError) Error() string {
	return e.msg
}
