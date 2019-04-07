package main

import (
	"context"
	"openpitrix.io/openpitrix/pkg/logger"
	"openpitrix.io/watcher/pkg/common"
)

func aa() {

	etcd := common.Etcd{
		Prefix:    "open",
		Endpoints: "172.31.140.135:2379",
	}
	key := "global_config"
	etcd.NewEtcdClient()
	//get old config from etcd, and compare with global_config
	ctx := context.Background()
	err2 := etcd.Dlock(ctx, func() error {
		defer etcd.Client.Close()
		//_, err := etcd.Put(ctx, key, "mmmmmmmm")
		get, err := etcd.Client.Get(ctx, key)
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
