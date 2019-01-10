// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package db_spec

var DBSpec = struct {
	TableNames []string

	UserTableName      string
	UserPrimaryKeyName string

	UserGroupTableName      string
	UserGroupPrimaryKeyName string

	UserGroupBindingTableName      string
	UserGroupBindingPrimaryKeyName string
}{
	TableNames: []string{
		"user",
		"user_group",
		"user_group_binding",
	},

	UserTableName:      "user",
	UserPrimaryKeyName: "user_id",

	UserGroupTableName:      "user_group",
	UserGroupPrimaryKeyName: "group_id",

	UserGroupBindingTableName:      "user_group_binding",
	UserGroupBindingPrimaryKeyName: "id",
}

var DBInitSqlList = []struct{ Name, Sql string }{
	{
		Name: "user",
		Sql: `CREATE TABLE IF NOT EXISTS user (
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

			primary key(user_id)
		);`,
	},
	{
		Name: "user_group",
		Sql: `CREATE TABLE IF NOT EXISTS user_group (
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

			primary key(group_id)
		);`,
	},
	{
		Name: "user_group_binding",
		Sql: `CREATE TABLE IF NOT EXISTS user_group_binding (
			id       varchar(50) not null,
			user_id  varchar(50),
			group_id varchar(50),

			primary key(id)
		);`,
	},
}
