// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package service

import (
	"context"

	"openpitrix.io/iam/openpitrix/pkg/pb"
)

func (p *Server) CreateGroup(ctx context.Context, req *pb.CreateGroupRequest) (*pb.CreateGroupResponse, error) {
	return p.db.CreateGroup(ctx, req)
}
func (p *Server) DeleteGroups(ctx context.Context, req *pb.DeleteGroupsRequest) (*pb.DeleteGroupsResponse, error) {
	return p.db.DeleteGroups(ctx, req)
}
func (p *Server) ModifyGroup(context.Context, *pb.ModifyGroupRequest) (*pb.ModifyGroupResponse, error) {
	panic("TODO")
}
func (p *Server) GetGroup(context.Context, *pb.GetGroupRequest) (*pb.GetGroupResponse, error) {
	panic("TODO")
}
func (p *Server) DescribeGroups(context.Context, *pb.DescribeGroupsRequest) (*pb.DescribeGroupsResponse, error) {
	panic("TODO")
}
