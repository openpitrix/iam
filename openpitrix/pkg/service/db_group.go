// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package service

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"openpitrix.io/iam/openpitrix/pkg/pb"
)

func (p *Database) CreateGroup(ctx context.Context, req *pb.CreateGroupRequest) (*pb.CreateGroupResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}

	sql, values := pkgBuildSql_InsertInto(req.GetValue())
	if len(values) == 0 {
		err := status.Errorf(codes.InvalidArgument, "empty field")
		return nil, err
	}
	_, err := p.DB.ExecContext(ctx, sql, values...)
	if err != nil {
		return nil, err
	}

	reply := &pb.CreateGroupResponse{
		Head: &pb.ResponseHeader{
			UserId:     req.GetHead().GetUserId(),
			OwnerPath:  "", // TODO
			AccessPath: "", // TODO
		},
		GroupId: req.GetValue().GetGroupId(),
	}

	return reply, nil
}

func (p *Database) DeleteGroups(context.Context, *pb.DeleteGroupsRequest) (*pb.DeleteGroupsResponse, error) {
	panic("TODO")
}
func (p *Database) ModifyGroup(context.Context, *pb.ModifyGroupRequest) (*pb.ModifyGroupResponse, error) {
	panic("TODO")
}
func (p *Database) GetGroup(context.Context, *pb.GetGroupRequest) (*pb.GetGroupResponse, error) {
	panic("TODO")
}
func (p *Database) DescribeGroups(context.Context, *pb.DescribeGroupsRequest) (*pb.DescribeGroupsResponse, error) {
	panic("TODO")
}
