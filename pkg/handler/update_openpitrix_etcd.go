// Copyright 2019 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package handler

import (
	"context"
	"io/ioutil"
	"reflect"

	"gopkg.in/yaml.v2"

	"openpitrix.io/openpitrix/pkg/logger"
	"openpitrix.io/openpitrix/pkg/pi"
	"openpitrix.io/watcher/pkg/common"
	"os"
)

type AnyMap map[interface{}]interface{}

func UpdateOpenPitrixEtcd() {
	global := common.Global
	etcd := global.Etcd

	var content []byte
	var oldConfig []byte
	newConfigMap := AnyMap{}
	oldConfigMap := AnyMap{}

	//read global_config file and convert to map
	content, err := ioutil.ReadFile(global.WatchedFile)
	if err != nil {
		logger.Critical(nil, "Failed to read %s: %+v", global.WatchedFile, err)
		return //just do nothing if failed to read file
	}
	logger.Debug(nil, "global_config yaml: %s", content)
	err = yaml.Unmarshal(content, newConfigMap)
	if err != nil {
		logger.Critical(nil, "Failed to Unmarshal global_config yaml to config map: %+v", err)
	}
	logger.Debug(nil, "global_config map: %v", newConfigMap)

	//get old config from etcd, and compare with global_config
	ctx, cancel := context.WithTimeout(context.Background(), common.EtcdDlockTimeOut)
	defer cancel()
	err = etcd.Dlock(ctx, func() error {
		logger.Info(nil, "Updating openpitrix etcd...")
		get, err := etcd.Client.Get(ctx, pi.GlobalConfigKey)
		if err != nil {
			return err
		}
		var modified = new(bool)
		logger.Debug(nil, "get.count: %d", get.Count)
		logger.Debug(nil, "get: %+v", get.Kvs)
		if get.Count == 0 {
			//init global_config if empty in etcd
			oldConfig = content
		} else {
			//update old config from new config
			oldConfig = get.Kvs[0].Value
			err := yaml.Unmarshal(oldConfig, oldConfigMap)
			if err != nil {
				logger.Error(ctx, "Failed to unmarshal old config to map!")
				return err
			}

			ignoreKeyMap := make(AnyMap)
			err = yaml.Unmarshal([]byte(os.Getenv(common.IgnoreKeys)), ignoreKeyMap)
			if err != nil {
				logger.Error(ctx, "Failed to unmarshal ignore keys to map!")
				return err
			}
			compareOpenPitrixConfig(newConfigMap, oldConfigMap, ignoreKeyMap, modified)
			logger.Debug(nil, "modified: %t, Config updated: %v", *modified, oldConfigMap)
		}

		if *modified { //put updated config to etcd
			oldConfig, err = yaml.Marshal(oldConfigMap)
			if err != nil {
				logger.Critical(nil, "Failed to convert oldConfigMap to oldConfig: %+v", err)

			}
			_, err := etcd.Client.Put(ctx, pi.GlobalConfigKey, string(oldConfig))
			if err != nil {
				logger.Critical(nil, "Failed to put data into etcd: %+v", err)
			}
		}
		return nil
	})
	if err != nil {
		logger.Critical(nil, "Failed to update etcd: %+v", err)
	}
}

//Base old config, update that from new config.
func compareOpenPitrixConfig(new, old AnyMap, ignoreKeys AnyMap, modified *bool) {
	for k, v := range old {
		kStr := k.(string)
		logger.Debug(nil, "key: %s", kStr)


		//check if k is in ignore keys
		var subIgnoreKeys AnyMap
		var t interface{}
		if ignoreKeys == nil || ignoreKeys[kStr] == nil {
			t = nil
		} else {
			t = reflect.TypeOf(ignoreKeys[kStr]).Kind()
		}

		if t == reflect.Bool && ignoreKeys[kStr].(bool) {
			logger.Info(nil,"Ignore to update config: %s", kStr)
			continue //only in this condition, ignore update old config
		} else if t == reflect.Map {
			//get sub-ignore-keys
			subIgnoreKeys = ignoreKeys[kStr].(AnyMap)
		}

		if v == nil { //check if old value and new value are nil
			if new == nil || new[k] == nil {
				continue
			} else {
				logger.Info(nil, "Updating, key: %s, old value: %v, new value: %v", k, v, new[k])
				//update old config from new config
				old[k] = new[k]
				continue
			}
		}

		switch reflect.TypeOf(v).Kind() {
		case reflect.Map:
			compareOpenPitrixConfig(new[k].(AnyMap), v.(AnyMap), subIgnoreKeys, modified)
		default:
			if new[k] != v { //update old config from new config
				logger.Info(nil, "Updating, key: %s, old value: %v, new value: %v", k, v, new[k])
				old[k] = new[k]
				*modified = true
			}
		}
	}
}
