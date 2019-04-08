// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package common

import (
	"context"
	"strings"

	"time"

	"go.etcd.io/etcd/clientv3"
	"go.etcd.io/etcd/clientv3/concurrency"
	"go.etcd.io/etcd/clientv3/namespace"

	"openpitrix.io/openpitrix/pkg/logger"
)

const (
	GlobalConfigKey  = "global_config"
	DlockKey         = "dlock_" + GlobalConfigKey
	EtcdDlockTimeOut = time.Second * 60
)

type Etcd struct {
	Prefix    string `default:"openpitrix"`
	Endpoints string `default:"openpitrix-etcd:2379"` // Example: "localhost:2379,localhost:22379,localhost:32379"
	Client    *clientv3.Client
}

func (etcd *Etcd) NewEtcdClient() {
	logger.Info(nil, "New client of etcd: %+v...", etcd)
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   strings.Split(etcd.Endpoints, ","),
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		logger.Critical(nil, "Failed to connect etcd.")
		panic(err)
	}
	cli.KV = namespace.NewKV(cli.KV, etcd.Prefix)
	etcd.Client = cli
}

type Mutex struct {
	*concurrency.Mutex
}

func (etcd *Etcd) NewMutex(ctx context.Context) (*Mutex, error) {
	if etcd.Client == nil {
		return nil, NilError{"The etcd client is nil!"}
	}
	logger.Debug(nil, "Create new session of etcd...")
	session, err := concurrency.NewSession(etcd.Client, concurrency.WithContext(ctx))
	if err != nil {
		return nil, err
	}
	return &Mutex{concurrency.NewMutex(session, DlockKey)}, nil
}

// Lock locks the mutex with a cancelable context. If the context is canceled
// while trying to acquire the lock, the mutex tries to clean its stale lock entry.
func (m *Mutex) Lock(ctx context.Context) error {
	return m.Mutex.Lock(ctx)
}

func (m *Mutex) Unlock(ctx context.Context) error {
	return m.Mutex.Unlock(ctx)
}

type callback func() error

func (etcd *Etcd) Dlock(ctx context.Context, cb callback) error {
	logger.Info(ctx, "Create dlock with key [%s]", DlockKey)
	mutex, err := etcd.NewMutex(ctx)
	if err != nil {
		logger.Critical(ctx, "Dlock lock error, failed to create mutex: %+v", err)
		panic(err)
	}
	logger.Debug(ctx, "Locking...")
	err = mutex.Lock(ctx)
	if err != nil {
		logger.Critical(ctx, "Dlock lock error, failed to lock mutex: %+v", err)
		return err
	}
	defer mutex.Unlock(ctx)
	logger.Debug(ctx, "Running callback...")
	err = cb()
	return err
}
