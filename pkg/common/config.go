package common

import (
	"openpitrix.io/openpitrix/pkg/logger"
    "fmt"
    "github.com/koding/multiconfig"
    "os"
    "openpitrix.io/watcher/pkg/test"
)

const CONFIG_PREFIX = "WATCHER"

type Config struct {
	WatchedFile string `default:"/opt/global_config.yaml"`   //The file that need to be watched
	Duration    int64  `default:"10"`                     //The duration for polling cycle which repeats
	Handler     string `default:"UpdateOpenpitrixEtcd"` //The action func name to run when files change
	LogLevel    string `default:"info"`
	Etcd        *Etcd
}

var Global = new(Config)

func LoadConf() {
    if os.Getenv("LOCAL") == "1" { //if run at local, export envs
        test.LocalEnv()
    }

	loader := multiconfig.MultiLoader(
	   &multiconfig.TagLoader{},
	   &multiconfig.EnvironmentLoader{Prefix: CONFIG_PREFIX, CamelCase: true},
     )
	//get config from env
	Global.Etcd = &Etcd{}
	err := loader.Load(Global)
	if err != nil {
		errMsg := fmt.Sprintf("Failed to load config: %+v", err)
		panic(errMsg)
	}

	logger.SetLevelByString(Global.LogLevel)
	logger.Debug(nil, "LoadConf: %+v", Global)
}

type NilError struct {
    msg string
}

func (e NilError) Error() string {
    return e.msg
}