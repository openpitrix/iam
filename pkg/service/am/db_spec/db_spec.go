// Copyright 2019 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package db_spec

import "regexp"

var (
	reUid       = regexp.MustCompile(`^[a-zA-Z0-9-_]{2,64}$`)
	reGid       = regexp.MustCompile(`^[a-zA-Z0-9-_]{2,64}$`)
	reGroupPath = regexp.MustCompile(`^[a-zA-Z0-9_.-]{2,255}$`)
)

const (
	ActionTableName      = "action"
	ActionPrimaryKeyName = "action_id"

	FeatureTableName      = "feature"
	FeaturePrimaryKeyName = "feature_id"

	ModuleTableName      = "module"
	ModulePrimaryKeyName = "module_id"

	RoleTableName      = "role"
	RolePrimaryKeyName = "role_id"

	RoleModuleBindingTableName      = "role_module_binding"
	RoleModuleBindingPrimaryKeyName = "id"

	UserRoleBindingTableName      = "user_role_binding"
	UserRoleBindingPrimaryKeyName = "id"
)

var TableMap = map[string]struct{ Name, PrimaryKey, Sql string }{
	ActionTableName: {
		Name:       ActionTableName,
		PrimaryKey: ActionPrimaryKeyName,
		Sql: `CREATE TABLE IF NOT EXISTS ` + ActionTableName + ` (
			action_id   varchar(50) not null,
			action_name varchar(50),
			feature_id  varchar(50),
			method      varchar(50),
			description varchar(1000),
			url         varchar(500),
			url_method  varchar(20),

			primary key(` + ActionPrimaryKeyName + `)
		);`,
	},
	FeatureTableName: {
		Name:       FeatureTableName,
		PrimaryKey: FeaturePrimaryKeyName,
		Sql: `CREATE TABLE IF NOT EXISTS ` + FeatureTableName + ` (
			feature_id   varchar(50) not null,
			module_id    varchar(50),
			feature_name varchar(50),

			primary key(` + FeaturePrimaryKeyName + `)
		);`,
	},
	ModuleTableName: {
		Name:       ModuleTableName,
		PrimaryKey: ModulePrimaryKeyName,
		Sql: `CREATE TABLE IF NOT EXISTS ` + ModuleTableName + ` (
			module_id   varchar(50) not null,
			module_name varchar(50),

			primary key(` + ModulePrimaryKeyName + `)
		);`,
	},
	RoleTableName: {
		Name:       RoleTableName,
		PrimaryKey: RolePrimaryKeyName,
		Sql: `CREATE TABLE IF NOT EXISTS ` + RoleTableName + ` (
			role_id     varchar(50) not null,
			role_name   varchar(50),
			description varchar(255),
			portal      varchar(50), -- 'admin,isv,dev,normal',
			create_time timestamp,
			update_time timestamp,
			owner       varchar(50),
			owner_path  varchar(50),

			primary key(` + RolePrimaryKeyName + `)
		);`,
	},

	// binding

	RoleModuleBindingTableName: {
		Name:       RoleModuleBindingTableName,
		PrimaryKey: RoleModuleBindingPrimaryKeyName,
		Sql: `CREATE TABLE IF NOT EXISTS ` + RoleModuleBindingTableName + ` (
			id              varchar(50) not null,
			role_id         varchar(50),
			module_id       varchar(50),
			data_level      varchar(50),
			enabled_actions text,
			create_time     timestamp,
			update_time     timestamp,
			owner           varchar(50),

			primary key(` + RoleModuleBindingPrimaryKeyName + `)
		);`,
	},
	UserRoleBindingTableName: {
		Name:       UserRoleBindingTableName,
		PrimaryKey: UserRoleBindingPrimaryKeyName,
		Sql: `CREATE TABLE IF NOT EXISTS ` + UserRoleBindingTableName + ` (
			id                   varchar(50) not null,
			user_id              varchar(50),
			role_id              varchar(50)

			primary key(` + UserRoleBindingPrimaryKeyName + `)
		);`,
	},
}
