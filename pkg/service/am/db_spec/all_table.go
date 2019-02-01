// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package db_spec

const (
	EnableActionTableName      = "enable_action_bundle"
	ModuleApiTableName         = "module_api"
	RoleTableName              = "role"
	RoleModuleBindingTableName = "role_module_binding"
	UserRoleBindingTableName   = "user_role_binding"
)

var TableNameList = []string{
	EnableActionTableName,
	ModuleApiTableName,
	RoleTableName,
	RoleModuleBindingTableName,
	UserRoleBindingTableName,
}
