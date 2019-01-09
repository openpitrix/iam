// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package im

import (
	"context"

	"openpitrix.io/iam/pkg/pb/im"
)

var _ pbim.AccountManagerServer = (*Server)(nil)

func (p *Server) CreateGroup(context.Context, *pbim.Group) (*pbim.Group, error) {
	panic("todo")
}
func (p *Server) DeleteGroups(context.Context, *pbim.GroupIdList) (*pbim.Empty, error) {
	panic("todo")
}

func (p *Server) CreateUser(context.Context, *pbim.User) (*pbim.User, error) {
	panic("todo")
}
func (p *Server) DeleteUsers(context.Context, *pbim.UserIdList) (*pbim.Empty, error) {
	panic("todo")
}

func (p *Server) ListUsers(context.Context, *pbim.Range) (*pbim.ListUesrsResponse, error) {
	panic("todo")
}
func (p *Server) ListGroups(context.Context, *pbim.Range) (*pbim.ListGroupsResponse, error) {
	panic("todo")
}

func (p *Server) GetUser(context.Context, *pbim.UserId) (*pbim.User, error) {
	panic("todo")
}
func (p *Server) GetUsersByGroupId(context.Context, *pbim.GroupId) (*pbim.UserList, error) {
	panic("todo")
}
func (p *Server) ModifyUser(context.Context, *pbim.User) (*pbim.User, error) {
	panic("todo")
}

func (p *Server) ComparePassword(context.Context, *pbim.Password) (*pbim.Empty, error) {
	panic("todo")
}
func (p *Server) ModifyPassword(context.Context, *pbim.Password) (*pbim.Empty, error) {
	panic("todo")
}

func (p *Server) GetGroup(context.Context, *pbim.GroupId) (*pbim.Group, error) {
	panic("todo")
}
func (p *Server) ModifyGroup(context.Context, *pbim.Group) (*pbim.Group, error) {
	panic("todo")
}

func (p *Server) JoinGroup(context.Context, *pbim.JoinGroupRequest) (*pbim.Empty, error) {
	panic("todo")
}
func (p *Server) LeaveGroup(context.Context, *pbim.LeaveGroupRequest) (*pbim.Empty, error) {
	panic("todo")
}
