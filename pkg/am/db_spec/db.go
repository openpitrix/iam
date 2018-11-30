// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

// OpenPitrix Access Management service database spec package.
package db_spec

type TableSpec interface {
	GetTableName() string
	GetTableSchema(dbtype string) string
}

var AllTableSpecList = []TableSpec{
	Role{},
	RoleBinding{},
}
