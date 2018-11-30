// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package db

const SqlTableSchema_Role = `
	CREATE TABLE IF NOT EXISTS role (
		name VARCHAR(50) NOT NULL,
		rule JSON        NOT NULL,

		PRIMARY KEY (name)
	);
`

const SqlTableSchema_RoleBinding = `
	CREATE TABLE IF NOT EXISTS role_binding (
		role_name VARCHAR(50) NOT NULL,
		xid       VARCHAR(50) NOT NULL,

		PRIMARY KEY (role_name, xid)
	);
`
