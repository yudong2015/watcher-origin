package main

import (
    "openpitrix.io/watcher/pkg/config"
    "io/ioutil"
    yaml "gopkg.in/yaml.v2"
    "openpitrix.io/openpitrix/pkg/logger"
    "context"
    "reflect"
)

type OpenpitrixConfig map[string]interface{}

var IGNORE_KEYS map[string]interface{}

func init(){
    IGNORE_KEYS = map[string]interface{}{
        "runtime": true,
    }
}


func UpdateOpenpitrixEtcd() {
    global := config.Global
    etcdClient := global.EtcdClient

    //read global_config file and convert to map
    content, err := ioutil.ReadFile(global.GlobalConfig.WatchedFile)
    if err != nil {
        logger.Critical(nil, "Failed to read %s: %+v", global.GlobalConfig.WatchedFile, err)
        return
    }
    newConfigMap := OpenpitrixConfig{}
    err = yaml.Unmarshal(content, newConfigMap)
    if err != nil {
        logger.Critical(nil, "Failed to Unmarshal to newConfigMap: %+v", err)
    }

    //get old config from etcd, and compare with global_config
    ctx := context.Background()
    err = etcdClient.Dlock(ctx, config.DlockKey, func() error {
        get, err := etcdClient.Get(ctx, config.GlobalConfigKey)
        if err != nil {
            return err
        }
        var oldConfig []byte
        if get.Count == 0 {
            oldConfig = content
        } else {
            oldConfig = get.Kvs[0].Value
            oldConfigMap := OpenpitrixConfig{}
            err := yaml.Unmarshal(oldConfig, oldConfigMap)
            if err != nil {
                return err
            }
            logger.Debug(nil, "%v", oldConfigMap)
            compareOpenpitrixConfig(newConfigMap, oldConfigMap, IGNORE_KEYS)
            logger.Debug(nil, "%v", oldConfigMap)
        }

        _, err = etcdClient.Put(ctx, config.GlobalConfigKey, string(oldConfig))
        if err != nil {
            logger.Critical(nil, "Failed to put data into etcd: %+v", err)
        }
        return nil
    })
    if err != nil {
        logger.Critical(nil, "Failed to update etcd: %+v", err)
    }
}


//Base old config in etcd, update that from new config.
func compareOpenpitrixConfig(new, old map[string]interface{}, ignoreKeys map[string]interface{}) {
    for k, v := range old {

        //get
        t := reflect.TypeOf(ignoreKeys[k]).Kind()
        if t == reflect.Bool && ignoreKeys[k].(bool) {
            return
        } else if t == reflect.Map {
            ignoreKeys = ignoreKeys[k].(map[string]interface{})
        }

        switch v.(type) {
        case map[string]interface{}:
            compareOpenpitrixConfig(new[k].(map[string]interface{}), v.(map[string]interface{}), ignoreKeys)
        default:
            if new[k] != v {
                old[k] = new[k]
            }
        }
    }
}

