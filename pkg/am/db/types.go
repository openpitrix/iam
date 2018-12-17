// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package db

import (
	"openpitrix.io/iam/pkg/pb/am"
)

type Role struct {
	Name string `db:"name"`
	Rule string `db:"rule"` // JSON string
}

type RoleBinding struct {
	RoleName string `db:"role_name"`
	Xid      string `db:"xid"`
}

func NewRoleBindingFrom(x *pbam.RoleXidBinding) *RoleBinding {
	return &RoleBinding{
		RoleName: x.RoleName,
		Xid:      x.Xid,
	}
}

func (p *RoleBinding) ToPbRoleBinding() *pbam.RoleXidBinding {
	return &pbam.RoleXidBinding{
		RoleName: p.RoleName,
		Xid:      p.Xid,
	}
}

func NewRoleFrom(x *pbam.Role) *Role {
	return &Role{
		Name: x.Name,
		Rule: encodeRuleList(x.Rule),
	}
}

func (p *Role) ToPbRole() *pbam.Role {
	return &pbam.Role{
		Name: p.Name,
		Rule: decodeRuleList(p.Rule),
	}
}
