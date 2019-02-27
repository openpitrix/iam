// Copyright 2019 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package constants

const (
	ColumnUserId                       = "user_id"
	ColumnRoleId                       = "role_id"
	ColumnRoleName                     = "role_name"
	ColumnPortal                       = "portal"
	ColumnCreateTime                   = "create_time"
	ColumnUpdateTime                   = "update_time"
	ColumnStatusTime                   = "status_time"
	ColumnStatus                       = "status"
	ColumnDescription                  = "description"
	ColumnModuleId                     = "module_id"
	ColumnId                           = "Id"
	ColumnOwnerPath                    = "owner_path"
	ColumnController                   = "controller"
	ColumnActionBundleVisibilitySuffix = "_action_bundle_visibility"
)

const (
	TableUserRoleBinding   = "user_role_binding"
	TableRole              = "role"
	TableEnableAction      = "enable_action_bundle"
	TableModuleApi         = "module_api"
	TableRoleModuleBinding = "role_module_binding"
)

// columns that can be search through sql '=' operator
var IndexedColumns = map[string][]string{
	TableRole: {
		ColumnRoleId, ColumnPortal, ColumnStatus,
	},
}

var SearchWordColumnTable = []string{
	TableRole,
}

// columns that can be search through sql 'like' operator
var SearchColumns = map[string][]string{
	TableRole: {
		ColumnRoleId, ColumnRoleName, ColumnPortal, ColumnStatus,
	},
}
