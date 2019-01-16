// Copyright 2019 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package db

import (
	"context"

	pbam "openpitrix.io/iam/pkg/pb/am"
)

func (p *Database) CreateRole(context.Context, *pbam.Role) (*pbam.Role, error) {
	panic("todo")
}
func (p *Database) DeleteRoles(context.Context, *pbam.RoleIdList) (*pbam.Empty, error) {
	panic("todo")
}
func (p *Database) ModifyRole(context.Context, *pbam.Role) (*pbam.Role, error) {
	panic("todo")
}
func (p *Database) DescribeRoles(context.Context, *pbam.DescribeRolesRequest) (*pbam.RoleList, error) {
	panic("todo")
}
