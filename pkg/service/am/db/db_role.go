// Copyright 2019 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package db

import (
	"context"

	pbam "openpitrix.io/iam/pkg/pb/am"
	"openpitrix.io/iam/pkg/util/funcutil"
	"openpitrix.io/logger"
)

func (p *Database) CreateRole(ctx context.Context, req *pbam.Role) (*pbam.Role, error) {
	logger.Infof(ctx, funcutil.CallerName(1))

	panic("todo")
}
func (p *Database) DeleteRoles(ctx context.Context, req *pbam.RoleIdList) (*pbam.Empty, error) {
	logger.Infof(ctx, funcutil.CallerName(1))

	panic("todo")
}
func (p *Database) ModifyRole(ctx context.Context, req *pbam.Role) (*pbam.Role, error) {
	logger.Infof(ctx, funcutil.CallerName(1))

	panic("todo")
}
func (p *Database) DescribeRoles(ctx context.Context, req *pbam.DescribeRolesRequest) (*pbam.RoleList, error) {
	logger.Infof(ctx, funcutil.CallerName(1))

	panic("todo")
}
