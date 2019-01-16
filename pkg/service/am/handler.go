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

func (p *Server) DescribeActions(context.Context, *pbam.DescribeActionsRequest) (*pbam.ActionList, error) {
	panic("todo")
}

func (p *Server) CanDo(context.Context, *pbam.CanDoRequest) (*pbam.CanDoResponse, error) {
	panic("todo")
}
func (p *Server) CreateRole(context.Context, *pbam.Role) (*pbam.Role, error) {
	panic("todo")
}
func (p *Server) DeleteRoles(context.Context, *pbam.RoleIdList) (*pbam.Empty, error) {
	panic("todo")
}
func (p *Server) ModifyRole(context.Context, *pbam.Role) (*pbam.Role, error) {
	panic("todo")
}
func (p *Server) DescribeRoles(context.Context, *pbam.DescribeRolesRequest) (*pbam.RoleList, error) {
	panic("todo")
}

func (p *Server) DescribeUsersWithRole(context.Context, *pbam.DescribeUsersWithRoleRequest) (*pbam.DescribeUsersWithRoleResponse, error) {
	panic("todo")
}

func (p *Server) GetRoleModule(context.Context, *pbam.RoleId) (*pbam.RoleModule, error) {
	panic("todo")
}
func (p *Server) ModifyRoleModule(context.Context, *pbam.RoleModule) (*pbam.RoleModule, error) {
	panic("todo")
}

func (p *Server) BindUserRole(context.Context, *pbam.BindUserRoleRequest) (*pbam.Empty, error) {
	panic("todo")
}
func (p *Server) UnbindUserRole(context.Context, *pbam.UnbindUserRoleRequest) (*pbam.Empty, error) {
	panic("todo")
}
