// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package im

import (
	"context"

	"openpitrix.io/iam/pkg/pb/im"
	"openpitrix.io/iam/pkg/version"
)

var _ pbim.AccountManagerServer = (*Server)(nil)

func (p *Server) GetVersion(ctx context.Context, req *pbim.Empty) (*pbim.String, error) {
	reply := &pbim.String{Value: version.GetVersionString()}
	return reply, nil
}

func (p *Server) CreateGroup(ctx context.Context, req *pbim.Group) (*pbim.Group, error) {
	return p.db.CreateGroup(ctx, req)
}
func (p *Server) DeleteGroups(ctx context.Context, req *pbim.GroupIdList) (*pbim.Empty, error) {
	return p.db.DeleteGroups(ctx, req)
}

func (p *Server) CreateUser(ctx context.Context, req *pbim.User) (*pbim.User, error) {
	return p.db.CreateUser(ctx, req)
}
func (p *Server) DeleteUsers(ctx context.Context, req *pbim.UserIdList) (*pbim.Empty, error) {
	return p.db.DeleteUsers(ctx, req)
}

func (p *Server) ListUsers(ctx context.Context, req *pbim.Range) (*pbim.ListUsersResponse, error) {
	return p.db.ListUsers(ctx, req)
}
func (p *Server) ListGroups(ctx context.Context, req *pbim.Range) (*pbim.ListGroupsResponse, error) {
	return p.db.ListGroups(ctx, req)
}

func (p *Server) GetUser(ctx context.Context, req *pbim.UserId) (*pbim.User, error) {
	return p.db.GetUser(ctx, req)
}
func (p *Server) GetUsersByGroupId(ctx context.Context, req *pbim.GroupId) (*pbim.UserList, error) {
	return p.db.GetUsersByGroupId(ctx, req)
}
func (p *Server) ModifyUser(ctx context.Context, req *pbim.User) (*pbim.User, error) {
	return p.db.ModifyUser(ctx, req)
}

func (p *Server) ComparePassword(ctx context.Context, req *pbim.Password) (*pbim.Empty, error) {
	return p.db.ComparePassword(ctx, req)
}
func (p *Server) ModifyPassword(ctx context.Context, req *pbim.Password) (*pbim.Empty, error) {
	return p.db.ModifyPassword(ctx, req)
}

func (p *Server) GetGroup(ctx context.Context, req *pbim.GroupId) (*pbim.Group, error) {
	return p.db.GetGroup(ctx, req)
}
func (p *Server) ModifyGroup(ctx context.Context, req *pbim.Group) (*pbim.Group, error) {
	return p.db.ModifyGroup(ctx, req)
}

func (p *Server) JoinGroup(ctx context.Context, req *pbim.JoinGroupRequest) (*pbim.Empty, error) {
	return p.db.JoinGroup(ctx, req)
}
func (p *Server) LeaveGroup(ctx context.Context, req *pbim.LeaveGroupRequest) (*pbim.Empty, error) {
	return p.db.LeaveGroup(ctx, req)
}
