
# Watcher

[![Build Status](https://travis-ci.org/openpitrix/openpitrix.svg)](https://travis-ci.org/openpitrix/openpitrix)
[![Docker Build Status](https://img.shields.io/docker/build/openpitrix/openpitrix.svg)](https://hub.docker.com/r/openpitrix/openpitrix/)
[![Go Report Card](https://goreportcard.com/badge/openpitrix.io/openpitrix)](https://goreportcard.com/report/openpitrix.io/openpitrix)
[![GoDoc](https://godoc.org/openpitrix.io/openpitrix?status.svg)](https://godoc.org/openpitrix.io/openpitrix)
[![License](http://img.shields.io/badge/license-apache%20v2-blue.svg)](https://github.com/openpitrix/openpitrix/blob/master/LICENSE)

----

**watcher** is a simple service to trigger a handler when **WATCHER_WATCHED_FILE** watched are updated.

In openpitrix, the global config is load into ConfigMap, watcher watches mounted volume dirs from ConfigMap and notifies the target process(UpdateOpenPitrixEtcd) that the config map has been changed.
It currently only supports UpdateOpenPitrixEtcd.


----

## Usage

#### Local

```$xslt
$ export WATCHER_HANDLER=UpdateOpenPitrixEtcd \
  WATCHER_WATCHED_FILE=./test/global_config.yaml \
  WATCHER_DURATION=10 \
  WATCHER_LOG_LEVEL=debug \
  WATCHER_ETCD_PREFIX=openpitrix \
  WATCHER_ETCD_ENDPOINTS=127.0.0.1:2379 \
  IGNORE_KEYS='{"runtime": true, "cluster": {"registry_mirror": true}}'

$ cd watcher
$ go run cmd/watcher/main.go
```

#### Docker

```$xslt
docker run -it -d \
       -v ${CONFIG_DIR}:/opt/config.yaml \
       -e WATCHER_HANDLER=UpdateOpenPitrixEtcd \
       -e WATCHER_WATCHED_FILE=/opt/config.yaml \
       -e WATCHER_DURATION=10 \
       -e WATCHER_LOG_LEVEL=debug \
       -e WATCHER_ETCD_PREFIX=openpitrix \
       -e WATCHER_ETCD_ENDPOINTS=127.0.0.1:2379 \
       -e IGNORE_KEYS='{"runtime": true, "cluster": {"registry_mirror": true}}' \
       --name watcher openpitrix.io/watcher
```

#### Helm

refrence: https://github.com/openpitrix/helm-chart