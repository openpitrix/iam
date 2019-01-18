// Copyright 2019 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package am

import (
	"context"

	"openpitrix.io/iam/pkg/pb/am"
	"openpitrix.io/iam/pkg/version"
)

var _ pbam.AccessManagerServer = (*Server)(nil)

func (p *Server) GetVersion(context.Context, *pbam.Empty) (*pbam.String, error) {
	reply := &pbam.String{Value: version.GetVersionString()}
	return reply, nil
}

func (p *Server) DescribeActions(ctx context.Context, req *pbam.DescribeActionsRequest) (*pbam.ActionList, error) {
	return p.db.DescribeActions(ctx, req)
}

func (p *Server) CanDo(ctx context.Context, req *pbam.CanDoRequest) (*pbam.CanDoResponse, error) {
	return p.db.CanDo(ctx, req)
}

func (p *Server) CreateRole(ctx context.Context, req *pbam.Role) (*pbam.Role, error) {
	return p.db.CreateRole(ctx, req)
}
func (p *Server) DeleteRoles(ctx context.Context, req *pbam.RoleIdList) (*pbam.Empty, error) {
	return p.db.DeleteRoles(ctx, req)
}
func (p *Server) ModifyRole(ctx context.Context, req *pbam.Role) (*pbam.Role, error) {
	return p.db.ModifyRole(ctx, req)
}
func (p *Server) GetRole(ctx context.Context, req *pbam.RoleId) (*pbam.Role, error) {
	return p.db.GetRole(ctx, req)
}
func (p *Server) DescribeRoles(ctx context.Context, req *pbam.DescribeRolesRequest) (*pbam.RoleList, error) {
	return p.db.DescribeRoles(ctx, req)
}

func (p *Server) GetUserWithRole(ctx context.Context, req *pbam.UserId) (*pbam.UserWithRole, error) {
	return p.db.GetUserWithRole(ctx, req)
}
func (p *Server) DescribeUsersWithRole(ctx context.Context, req *pbam.DescribeUsersWithRoleRequest) (*pbam.DescribeUsersWithRoleResponse, error) {
	return p.db.DescribeUsersWithRole(ctx, req)
}

func (p *Server) ModifyRoleModule(ctx context.Context, req *pbam.RoleModule) (*pbam.RoleModule, error) {
	return p.db.ModifyRoleModule(ctx, req)
}

func (p *Server) BindUserRole(ctx context.Context, req *pbam.BindUserRoleRequest) (*pbam.Empty, error) {
	return p.db.BindUserRole(ctx, req)
}
func (p *Server) UnbindUserRole(ctx context.Context, req *pbam.UnbindUserRoleRequest) (*pbam.Empty, error) {
	return p.db.UnbindUserRole(ctx, req)
}
