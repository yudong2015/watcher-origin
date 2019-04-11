package test

import (
	"testing"

	"context"
	"io/ioutil"
	"openpitrix.io/openpitrix/pkg/pi"
	"openpitrix.io/watcher/pkg/common"
)

func TestEtcd(t *testing.T) {
	t.Logf("Test etcd...")

	LocalEnv()
	common.LoadConf()
	global := common.Global

	etcd := global.Etcd
	ctx, cancel := context.WithTimeout(context.Background(), common.EtcdDlockTimeOut)
	defer cancel()
	etcd.Dlock(ctx, func() error {
		get, err := etcd.Client.Get(ctx, pi.GlobalConfigKey)
		if err != nil {
			t.Errorf("Failed to get %s from etcd!", pi.GlobalConfigKey)
			return err
		}
		t.Logf("The origin: %+v", get.Kvs)

		content, err := ioutil.ReadFile(global.WatchedFile)
		t.Logf("The content: %s", content)
		_, err = etcd.Client.Put(ctx, pi.GlobalConfigKey, string(content))
		if err != nil {
			t.Error("Failed to put into etcd!")
			return err
		}

		get, err = etcd.Client.Get(ctx, pi.GlobalConfigKey)
		if err != nil {
			t.Errorf("Failed to get %s from etcd!", pi.GlobalConfigKey)
			return err
		}
		t.Logf("The new: %+v", get.Kvs)

		return nil
	})

	t.Log("Test etcd successfully!")
}
