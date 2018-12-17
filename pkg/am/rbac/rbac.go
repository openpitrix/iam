// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

// Package rbac implements rbac with Database.
package rbac

import (
	"openpitrix.io/iam/pkg/pb/am"
)

type Interface interface {
	Close() error

	CreateRole(role *pbam.Role) error
	ModifyRole(role *pbam.Role) error
	DeleteRoleByName(name string) error

	GetRoleByName(name string) (*pbam.Role, error)
	GetRoleByXidList(xid ...string) (*pbam.RoleList, error)
	ListRoles(filter *pbam.NameFilter) (*pbam.RoleList, error)

	CreateRoleBinding(x *pbam.RoleXidBindingList) error
	DeleteRoleBinding(xid ...string) error

	GetRoleBindingByRoleName(name string) (*pbam.RoleXidBindingList, error)
	GetRoleBindingByXidList(xid ...string) (*pbam.RoleXidBindingList, error)
	ListRoleBindings(filter *pbam.NameFilter) (*pbam.RoleXidBindingList, error)

	CanDo(x *pbam.Action) bool
}

func OpenDatabase(dbtype, dbpath string) (Interface, error) {
	return openDatabase(dbtype, dbpath)
}
