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
)

type AnyMap map[interface{}]interface{}

var IGNORE_KEYS map[string]interface{}

func init() {
	IGNORE_KEYS = map[string]interface{}{
		"runtime": true,
	}
}

func UpdateOpenpitrixEtcd() {
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
	logger.Debug(nil, "global_config_yaml: %s", content)
	err = yaml.Unmarshal(content, newConfigMap)
	if err != nil {
		logger.Critical(nil, "Failed to Unmarshal to newConfigMap: %+v", err)
	}
	logger.Debug(nil, "global_config_map: %v", newConfigMap)

	//get old config from etcd, and compare with global_config
	ctx, cancel := context.WithTimeout(context.Background(), common.EtcdDlockTimeOut)
	defer cancel()
	err = etcd.Dlock(ctx, func() error {
		logger.Info(nil, "Updating openpitrix etcd...")
		get, err := etcd.Client.Get(ctx, pi.GlobalConfigKey)
		if err != nil {
			return err
		}
		var modifyed = new(bool)
		logger.Debug(nil, "get-count: %d", get.Count)
		logger.Debug(nil, "get: %+v", get.Kvs)
		if get.Count == 0 {
			//init global_config if empty in etcd
			oldConfig = content
		} else {
			//update old config from new config
			oldConfig = get.Kvs[0].Value
			err := yaml.Unmarshal(oldConfig, oldConfigMap)
			if err != nil {
				return err
			}
			compareOpenpitrixConfig(newConfigMap, oldConfigMap, IGNORE_KEYS, modifyed)
			logger.Debug(nil, "modifyed: %t, Config updated: %v", *modifyed, oldConfigMap)
		}

		if *modifyed { //put updated config to etcd
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

//Base old config in etcd, update that from new config.
//return if there is diffrence from new and old
func compareOpenpitrixConfig(new, old AnyMap, ignoreKeys map[string]interface{}, modifyed *bool) {
	for k, v := range old {
		kStr := k.(string)

		logger.Debug(nil, "key: %s", kStr)
		logger.Debug(nil, "oldValue: %+v", v)
		logger.Debug(nil, "newValue: %+v", new[k])

		//check if k is in ignore updating map
		var t interface{}
		if ignoreKeys == nil || ignoreKeys[kStr] == nil {
			t = nil
		} else {
			t = reflect.TypeOf(ignoreKeys[kStr]).Kind()
		}

		if t == reflect.Bool && ignoreKeys[kStr].(bool) {
			continue //olny in this condition, ignore update old config
		} else if t == reflect.Map {
			//get sub-ignore-keys
			ignoreKeys = ignoreKeys[kStr].(map[string]interface{})
		}

		if v == nil { //check if old value and new value are nil
			if new == nil || new[k] == nil {
				continue
			} else {
				logger.Info(nil, "Updating, key: %s, oldValue: %v, newValue: %v", k, v, new[k])
				old[k] = new[k]
				continue
			}
		}

		//update old config from new config
		switch reflect.TypeOf(v).Kind() {
		case reflect.Map:
			compareOpenpitrixConfig(new[k].(AnyMap), v.(AnyMap), ignoreKeys, modifyed)
		default:
			if new[k] != v {
				logger.Info(nil, "Updating, key: %s, oldValue: %v, newValue: %v", k, v, new[k])
				old[k] = new[k]
				*modifyed = true
			}
		}
	}
}
