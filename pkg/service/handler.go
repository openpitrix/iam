// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package service

import (
	"context"

	"openpitrix.io/iam/openpitrix/pkg/pb"
	"openpitrix.io/iam/openpitrix/pkg/version"
)

func (p *Server) GetVersion(ctx context.Context, req *pb.Empty) (*pb.String, error) {
	reply := &pb.String{Value: version.GetVersionString()}
	return reply, nil
}

func (p *Server) DescribeActions(context.Context, *pb.DescribeActionsRequest) (*pb.DescribeActionsResponse, error) {
	panic("TODO")
}

func (p *Server) GetOwnerPath(context.Context, *pb.GetOwnerPathRequest) (*pb.String, error) {
	panic("TODO")
}

func (p *Server) GetAccessPath(context.Context, *pb.GetAccessPathRequest) (*pb.String, error) {
	panic("TODO")
}

func (p *Server) CanDoAction(context.Context, *pb.CanDoActionRequest) (*pb.CanDoActionResponse, error) {
	panic("TODO")
}

func (p *Server) CreateGroup(ctx context.Context, req *pb.CreateGroupRequest) (*pb.CreateGroupResponse, error) {
	return p.db.CreateGroup(ctx, req)
}
func (p *Server) DeleteGroups(ctx context.Context, req *pb.DeleteGroupsRequest) (*pb.DeleteGroupsResponse, error) {
	return p.db.DeleteGroups(ctx, req)
}
func (p *Server) ModifyGroup(ctx context.Context, req *pb.ModifyGroupRequest) (*pb.ModifyGroupResponse, error) {
	return p.db.ModifyGroup(ctx, req)
}
func (p *Server) GetGroup(ctx context.Context, req *pb.GetGroupRequest) (*pb.GetGroupResponse, error) {
	return p.db.GetGroup(ctx, req)
}
func (p *Server) DescribeGroups(ctx context.Context, req *pb.DescribeGroupsRequest) (*pb.DescribeGroupsResponse, error) {
	return p.db.DescribeGroups(ctx, req)
}

func (p *Server) ModifyRoleModuleBindings(context.Context, *pb.ModifyRoleModuleBindingsRequest) (*pb.ModifyRoleModuleBindingsResponse, error) {
	panic("TODO")
}

func (p *Server) CreateRole(ctx context.Context, req *pb.CreateRoleRequest) (*pb.CreateRoleResponse, error) {
	return p.db.CreateRole(ctx, req)
}

func (p *Server) DeleteRoles(ctx context.Context, req *pb.DeleteRolesRequest) (*pb.DeleteRolesResponse, error) {
	return p.db.DeleteRoles(ctx, req)
}
func (p *Server) ModifyRole(ctx context.Context, req *pb.ModifyRoleRequest) (*pb.ModifyRoleResponse, error) {
	return p.db.ModifyRole(ctx, req)
}
func (p *Server) GetRole(ctx context.Context, req *pb.GetRoleRequest) (*pb.GetRoleResponse, error) {
	return p.db.GetRole(ctx, req)
}
func (p *Server) DescribeRoles(ctx context.Context, req *pb.DescribeRolesRequest) (*pb.DescribeRolesResponse, error) {
	return p.db.DescribeRoles(ctx, req)
}

func (p *Server) ComparePassword(context.Context, *pb.UserPassword) (*pb.Bool, error) {
	panic("TODO")
}
func (p *Server) ModifyPassword(context.Context, *pb.UserPassword) (*pb.Bool, error) {
	panic("TODO")
}

func (p *Server) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
	return p.db.CreateUser(ctx, req)
}
func (p *Server) DeleteUsers(ctx context.Context, req *pb.DeleteUsersRequest) (*pb.DeleteUsersResponse, error) {
	return p.db.DeleteUsers(ctx, req)
}
func (p *Server) ModifyUser(ctx context.Context, req *pb.ModifyUserRequest) (*pb.ModifyUserResponse, error) {
	return p.db.ModifyUser(ctx, req)
}
func (p *Server) DescribeUsers(ctx context.Context, req *pb.DescribeUsersRequest) (*pb.DescribeUsersResponse, error) {
	return p.db.DescribeUsers(ctx, req)
}
func (p *Server) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.GetUserResponse, error) {
	return p.db.GetUser(ctx, req)
}
