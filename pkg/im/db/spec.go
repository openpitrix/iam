// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package db

// http://www.mysqltutorial.org/mysql-adjacency-list-tree/
// http://mikehillyer.com/articles/managing-hierarchical-data-in-mysql/

var _TableSchemaMap = []struct {
	Name   string
	Schema string
	Value  interface{}
}{
	{TableName_User, SqlTableSchema_User, User{}},
	{TableName_Group, SqlTableSchema_Group, Group{}},
}

const (
	TableName_User  = `user`
	TableName_Group = `group`
)

const SqlTableSchema_User = `
	CREATE TABLE IF NOT EXISTS user (
		uid VARCHAR(50) NOT NULL,
		gid VARCHAR(50) NOT NULL,

		name        TEXT NOT NULL,
		email       TEXT NOT NULL,
		description TEXT NOT NULL,
		password    TEXT NOT NULL,
		status      TEXT NOT NULL,
		extra       JSON NOT NULL,

		create_time TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
		update_time TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
		status_time TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,

		PRIMARY KEY (uid),
		FOREIGN KEY (gid) REFERENCES group (gid)
	);
`

const SqlTableSchema_Group = `
	CREATE TABLE IF NOT EXISTS group (
		gid        VARCHAR(50) NOT NULL,
		gid_parent VARCHAR(50) NOT NULL, -- root.gid_parent == root.gid

		name        TEXT NOT NULL,
		email       TEXT NOT NULL,
		description TEXT NOT NULL,
		status      TEXT NOT NULL,
		extra       JSON NOT NULL,

		create_time TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
		update_time TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
		status_time TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,

		PRIMARY KEY (gid),
		FOREIGN KEY (gid_parent) REFERENCES group (gid)
			ON DELETE CASCADE
			ON UPDATE CASCADE
	);
`
