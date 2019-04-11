package test

import (
	"testing"

	"context"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"openpitrix.io/openpitrix/pkg/pi"
	"openpitrix.io/watcher/pkg/common"
	"openpitrix.io/watcher/pkg/handler"
	"os"
)

func TestWatchOpenPitrixConfig(t *testing.T) {
	t.Logf("Test watching openpitrix config...")

	LocalEnv()
	common.LoadConf()
	global := common.Global

	content, contentMap, err := common.ReadYamlFile(
	    "./global_config.yaml")
	if err != nil {
		t.Errorf("Failed to read content: %+v", err)
	}

	//init etcd: put global_config into etcd
	//err = putToEtcd(global, content)
	//if err != nil {
	//    t.Errorf("Failed to put content into etcd: %+v", err)
	//}

	const TmpFile = "./config_tmp.yaml"
	global.WatchedFile = TmpFile
	defer os.Remove(TmpFile)

	//Update config values in content
	contentMap := make(map[string]interface{})
	err = yaml.Unmarshal(content, contentMap)
	if err != nil {
		t.Errorf("Failed to unmarshal content to map: %+v", err)
	}
	var pilot map[string]interface{}
	pilot = common.InterfaceToMap(contentMap["pilot"])
	pilot["port"] = 30114
	t.Logf("contentMap: %+v", contentMap)
	t.Logf("contentMap: %+v", pilot)

	t.Log("Test successfully!")
}

func putToEtcd(global *common.Config, content []byte) error {
	etcd := global.Etcd
	ctx, cancel := context.WithTimeout(context.Background(), common.EtcdDlockTimeOut)
	defer cancel()
	err := etcd.Dlock(ctx, func() error {
		_, err := etcd.Client.Put(ctx, pi.GlobalConfigKey, string(content))
		if err != nil {
			return err
		}
		return nil
	})
	return err
}
