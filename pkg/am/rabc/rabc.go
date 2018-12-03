// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

// Package rabc implements rabc with local file and Database.
package rabc

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
	ListRoles(filter *pbam.RoleNameFilter) (*pbam.RoleList, error)

	CreateRoleBinding(x *pbam.RoleBindingList) error
	DeleteRoleBinding(xid ...string) error

	GetRoleBindingByRoleName(name string) (*pbam.RoleBindingList, error)
	GetRoleBindingByXidList(xid ...string) (*pbam.RoleBindingList, error)
	ListRoleBindings(filter *pbam.RoleNameFilter) (*pbam.RoleBindingList, error)

	CanDo(x *pbam.Action) bool
}

func OpenDatabase(dbtype, dbpath string) (Interface, error) {
	return openDatabase(dbtype, dbpath)
}
