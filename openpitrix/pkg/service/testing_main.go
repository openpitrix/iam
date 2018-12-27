// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package service

import (
	"flag"
	"os"
	"testing"

	"openpitrix.io/iam/openpitrix/pkg/config"
)

var (
	tDatabase   = flag.String("test-database", "auto", "enable database (auto|off|on)")
	tConfigFile = flag.String("test-config", "config.json", "config file")

	tDatabaseEnabled = false
	tConfig          *config.Config
)

func init() {
	flag.Parse()

	switch *tDatabase {
	case "auto":
		tDatabaseEnabled = false
	case "on":
		tDatabaseEnabled = true
	case "off":
		tDatabaseEnabled = false
	}
}

func TestMain(m *testing.M) {
	// call flag.Parse() here if TestMain uses flags
	os.Exit(m.Run())
}
