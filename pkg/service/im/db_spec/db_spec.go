// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package db_spec

const (
	UserTableName      = "user"
	UserPrimaryKeyName = "user_id"

	UserGroupTableName      = "user_group"
	UserGroupPrimaryKeyName = "group_id"

	UserGroupBindingTableName      = "user_group_binding"
	UserGroupBindingPrimaryKeyName = "id"
)

var TableMap = map[string]struct{ Name, PrimaryKey, Sql string }{
	UserTableName: {
		Name:       UserTableName,
		PrimaryKey: UserPrimaryKeyName,
		Sql: `CREATE TABLE IF NOT EXISTS ` + UserTableName + ` (
			user_id      varchar(50) not null,
			user_name    varchar(50),
			email        varchar(50),
			phone_number varchar(50),
			description  varchar(1000),
			password     varchar(128),
			status       varchar(10),
			create_time  timestamp,
			update_time  timestamp,
			status_time  timestamp,
			extra        json,

			primary key(` + UserPrimaryKeyName + `)
		);`,
	},
	UserGroupTableName: {
		Name:       UserGroupTableName,
		PrimaryKey: UserGroupPrimaryKeyName,
		Sql: `CREATE TABLE IF NOT EXISTS ` + UserGroupTableName + ` (
			group_id        varchar(50) not null,
			group_path      varchar(255),
			group_name      varchar(50),
			description     varchar(1000),
			status          varchar(10),
			create_time     timestamp,
			update_time     timestamp,
			status_time     timestamp,
			extra           json,

			parent_group_id varchar(50),
			level           int,

			primary key(` + UserGroupPrimaryKeyName + `)
		);`,
	},
	UserGroupBindingTableName: {
		Name:       UserGroupBindingTableName,
		PrimaryKey: UserGroupBindingPrimaryKeyName,
		Sql: `CREATE TABLE IF NOT EXISTS ` + UserGroupBindingTableName + ` (
			id       varchar(50) not null,
			user_id  varchar(50),
			group_id varchar(50),

			primary key(` + UserGroupBindingPrimaryKeyName + `)
		);`,
	},
}
