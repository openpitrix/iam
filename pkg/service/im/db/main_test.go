// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package db

import (
	"flag"
	"os"
	"testing"
)

var (
	flagConfigFile = flag.String("config", "config-sqlite3.json", "set config type")
	flagEnableDB   = flag.Bool("enable-db", false, "enable database")
)

func TestMain(m *testing.M) {
	flag.Parse()
	os.Exit(m.Run())
}
