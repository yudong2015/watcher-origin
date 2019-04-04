package main

import (
	"context"
	"openpitrix.io/openpitrix/pkg/logger"
	"openpitrix.io/watcher/pkg/common"
)

func aa() {

	etcdConfig := common.EtcdConfig{
		Prefix:    "open",
		Endpoints: "172.31.140.135:2379",
	}
	key := "global_config"
	etcd := etcdConfig.OpenEtcd()
	//get old config from etcd, and compare with global_config
	ctx := context.Background()
	err2 := etcd.Dlock(ctx, common.DlockKey, func() error {
		defer etcd.Close()
		//_, err := etcd.Put(ctx, key, "mmmmmmmm")
		get, err := etcd.Get(ctx, key)
		if err != nil {
			return err
		}

		logger.Info(nil, "count: %d", get.Count)
		if get.Count > 0 {
			logger.Info(nil, "get: %s", get.Kvs[0].Value)
		}

		if err != nil {
			logger.Critical(nil, "Failed to put data into etcd: %+v", err)
		}
		return nil
	})
	if err2 != nil {
		logger.Critical(nil, "Failed to update etcd: %+v", err2)
	}
}
