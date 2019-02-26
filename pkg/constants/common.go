// Copyright 2019 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package constants

const (
	PrefixRoleId               = "rid-"
	PrefixUserRoleBindingId    = "bid-"
	PrefixRoleModuleBindingId  = "bid-"
	PrefixEnableActionBundleId = "eid-"
)

const (
	ModuleIdM0 = "m0"
)

const (
	PortalGlobalAdmin = "global_admin"
	PortalIsv         = "isv"
	PortalUser        = "user"
)

var PortalSet = []string{PortalGlobalAdmin, PortalIsv, PortalUser}

const (
	RoleGlobalAdmin = "global_admin"
	RoleDeveloper   = "developer"
	RoleIsv         = "isv"
	RoleUser        = "user"
)

const (
	UserSystem = "system"
)

const (
	DataLevelAll   = "all"
	DataLevelGroup = "group"
	DataLevelSelf  = "self"
)

const (
	StatusActive  = "active"
	StatusDeleted = "deleted"
)

const (
	ControllerSelf   = "self"
	ControllerPitrix = "pitrix"
)

const (
	ActionCreate   = "create"
	ActionModify   = "modify"
	ActionDelete   = "delete"
	ActionDescribe = "describe"
)
