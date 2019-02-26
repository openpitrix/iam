// Copyright 2019 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package am

import (
	"context"

	"openpitrix.io/iam/pkg/constants"
	"openpitrix.io/iam/pkg/pb"
	"openpitrix.io/iam/pkg/service/am/resource"
	"openpitrix.io/iam/pkg/version"
)

var _ pb.AccessManagerServer = (*Server)(nil)

func (p *Server) GetVersion(context.Context, *pb.GetVersionRequest) (*pb.GetVersionResponse, error) {
	return &pb.GetVersionResponse{
		Version: version.GetVersionString(),
	}, nil
}

func (p *Server) CanDo(ctx context.Context, req *pb.CanDoRequest) (*pb.CanDoResponse, error) {
	return resource.CanDo(ctx, req)
}

func (p *Server) CreateRole(ctx context.Context, req *pb.CreateRoleRequest) (*pb.CreateRoleResponse, error) {
	return resource.CreateRole(ctx, req)
}

func (p *Server) DeleteRoles(ctx context.Context, req *pb.DeleteRolesRequest) (*pb.DeleteRolesResponse, error) {
	_, err := CheckRolesPermission(ctx, req.RoleId, constants.ActionDelete)
	if err != nil {
		return nil, err
	}
	return resource.DeleteRoles(ctx, req)
}

func (p *Server) ModifyRole(ctx context.Context, req *pb.ModifyRoleRequest) (*pb.ModifyRoleResponse, error) {
	_, err := CheckRolesPermission(ctx, []string{req.RoleId}, constants.ActionModify)
	if err != nil {
		return nil, err
	}
	return resource.ModifyRole(ctx, req)
}

func (p *Server) GetRole(ctx context.Context, req *pb.GetRoleRequest) (*pb.GetRoleResponse, error) {
	roles, err := CheckRolesPermission(ctx, []string{req.RoleId}, constants.ActionDescribe)
	if err != nil {
		return nil, err
	}
	return &pb.GetRoleResponse{
		Role: roles[0].ToPB(),
	}, nil
}

func (p *Server) GetRoleWithUser(ctx context.Context, req *pb.GetRoleRequest) (*pb.GetRoleWithUserResponse, error) {
	_, err := CheckRolesPermission(ctx, []string{req.RoleId}, constants.ActionDescribe)
	if err != nil {
		return nil, err
	}
	return resource.GetRoleWithUser(ctx, req)
}

func (p *Server) DescribeRoles(ctx context.Context, req *pb.DescribeRolesRequest) (*pb.DescribeRolesResponse, error) {
	return resource.DescribeRoles(ctx, req)
}

func (p *Server) DescribeRolesWithUser(ctx context.Context, req *pb.DescribeRolesRequest) (*pb.DescribeRolesWithUserResponse, error) {
	return resource.DescribeRolesWithUser(ctx, req)
}

func (p *Server) GetRoleModule(ctx context.Context, req *pb.GetRoleModuleRequest) (*pb.GetRoleModuleResponse, error) {
	_, err := CheckRolesPermission(ctx, []string{req.RoleId}, constants.ActionDescribe)
	if err != nil {
		return nil, err
	}
	return resource.GetRoleModule(ctx, req)
}

func (p *Server) ModifyRoleModule(ctx context.Context, req *pb.ModifyRoleModuleRequest) (*pb.ModifyRoleModuleResponse, error) {
	_, err := CheckRolesPermission(ctx, []string{req.RoleId}, constants.ActionModify)
	if err != nil {
		return nil, err
	}
	return resource.ModifyRoleModule(ctx, req)
}

func (p *Server) BindUserRole(ctx context.Context, req *pb.BindUserRoleRequest) (*pb.BindUserRoleResponse, error) {
	_, err := CheckRolesPermission(ctx, req.RoleId, constants.ActionCreate)
	if err != nil {
		return nil, err
	}
	return resource.BindUserRole(ctx, req)
}

func (p *Server) UnbindUserRole(ctx context.Context, req *pb.UnbindUserRoleRequest) (*pb.UnbindUserRoleResponse, error) {
	_, err := CheckRolesPermission(ctx, req.RoleId, constants.ActionCreate)
	if err != nil {
		return nil, err
	}
	return resource.UnbindUserRole(ctx, req)
}
