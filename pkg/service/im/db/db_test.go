// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package db

import (
	"testing"

	"openpitrix.io/iam/pkg/service/im/config"
)

func TestDB(t *testing.T) {
	if !*flagEnableDB {
		t.Skip()
	}

	cfg, err := config.Load(*flagConfigFile)
	Assert(t, err == nil, err)

	db, err := OpenDatabase(cfg, nil)
	Assert(t, err == nil, err)
	defer db.Close()

	_ = db
}
