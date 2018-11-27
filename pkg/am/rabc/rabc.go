// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package rabc

import (
	"openpitrix.io/iam/pkg/pb/am"
)

type Interface interface {
	CanDo(x pbam.Action) bool

	AllRoles() []pbam.Role
	AllRoleBindings() []pbam.RoleBinding

	GetRoleByName(name string) (role pbam.Role, ok bool)
	GetRoleByXid(xid []string) pbam.RoleList

	CreateRole(role pbam.Role) error
	ModifyRole(role pbam.Role) error
	DeleteRole(name string) error

	CreateRoleBinding(x []pbam.RoleBinding) error
	DeleteRoleBinding(xid []string) error

	Close() error
}

type DBOptions struct {
	DBType     string // mysql/sqlite3
	DBEngine   string // InnoDB/...
	DBEncoding string // utf8/...
}

func OpenDatabase(dbpath string, opt *DBOptions) (Interface, error) {
	panic("TODO")
}

func OpenFile(jsonFile string) (Interface, error) {
	panic("TODO")
}
