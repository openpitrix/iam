// Copyright 2019 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package db_spec

type DBModule struct {
	ModuleId   string `db:"module_id"`
	ModuleName string `db:"module_name"`
}
