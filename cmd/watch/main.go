// Copyright 2019 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package main

import (
	"openpitrix.io/watcher/pkg/common"
	"openpitrix.io/watcher/pkg/watch"
)

func main() {
	common.LoadConf()
	watch.Watch()
}
