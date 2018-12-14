// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package client

import (
	"openpitrix.io/iam/pkg/pb/am"
	"openpitrix.io/iam/pkg/pb/im"
)

type Interface interface {
	CreateGroup(group ...*pbim.Group) error
	CreateUser(user ...*pbim.User) error
	CreateRole(role ...*pbam.Role) error
	CreateActionRule(rule ...*pbam.ActionRule) error
	CreateUserRoleBinding(bind ...*pbam.RoleXidBinding) error

	DeleteGroup(gid ...string) error
	DeleteUser(uid ...string) error
	DeleteRole(name ...string) error
	DeleteActionRule(name ...string) error
	DeleteUserRoleBinding(bind ...*pbam.RoleXidBinding) error

	ModifyGroup(group *pbim.Group) error
	ModifyUser(user *pbim.User) error
	MidifyPassword(uid, password string) error
	ModyfyRole(rule ...*pbam.ActionRule) error
	ModyfyActionRule(rule ...*pbam.ActionRule) error

	GetGroup(gid string) (*pbim.Group, error)
	GetUser(uid string) (*pbim.User, error)
	GetRole(name string) (*pbam.Role, error)
	GetActionRule(name string) (*pbam.ActionRule, error)
	GetUserRoleBinding(xid string) ([]*pbam.RoleXidBinding, error)

	ListGroups(filter *pbim.Range) ([]*pbim.Group, error)
	ListUsers(filter *pbim.Range) ([]*pbim.User, error)
	ListRoles(filter *pbam.NameFilter) ([]*pbam.Role, error)
	ListActionRules(filter *pbam.NameFilter) ([]*pbam.ActionRule, error)
	ListUserRoleBinding(filter *pbam.RoleBindingFilter) ([]*pbam.RoleXidBinding, error)

	ComparePassword(uid, password string) (ok bool, err error)

	CanDoAction(uid, method string, namespace ...string)
	CanDoActionWithUserNamespace(uid, method string)
	CanDoActionNoNamespace(uid, method string)
}
