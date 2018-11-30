// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package db_spec

import (
	"openpitrix.io/iam/pkg/pb/am"
)

const (
	RoleBindingTableName = "iam_role_binding"

	RoleBindingTableSchema = `
		CREATE TABLE IF NOT EXISTS ` + RoleBindingTableName + ` (
			role_name VARCHAR(50)  NOT NULL,
			xid       VARCHAR(50)  NOT NULL,

			PRIMARY KEY (role_name, xid)
		);
	`
)

type RoleBinding struct {
	RoleName string
	Xid      string
}

func (RoleBinding) GetTableName() string {
	return RoleBindingTableName
}
func (RoleBinding) GetTableSchema(dbtype string) string {
	return RoleBindingTableSchema
}

func NewRoleBindingFrom(x *pbam.RoleBinding) *RoleBinding {
	return &RoleBinding{
		RoleName: x.RoleName,
		Xid:      x.Xid,
	}
}

func (p *RoleBinding) ToPbRoleBinding() *pbam.RoleBinding {
	return &pbam.RoleBinding{
		RoleName: p.RoleName,
		Xid:      p.Xid,
	}
}
