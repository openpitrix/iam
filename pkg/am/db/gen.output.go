// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

// Code generated. DO NOT EDIT.

package db

type Role struct {
	Name string `db:"name, size:50, primarykey"`
	Rule string `db:"rule"`
}

type RoleBinding struct {
	Name string `db:"name, size:50, primarykey"`
	Xid  string `db:"xid, size:50, primarykey"`
}
