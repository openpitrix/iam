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

func (p *Database) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}

	sql, values := pkgBuildSql_InsertInto(
		dbTableSchemaMap[&pb.Group{}].TableName,
		req.GetValue(),
	)
	if len(values) == 0 {
		err := status.Errorf(codes.InvalidArgument, "empty field")
		return nil, err
	}
	_, err := p.DB.ExecContext(ctx, sql, values...)
	if err != nil {
		return nil, err
	}

	reply := &pb.CreateUserResponse{
		Head: &pb.ResponseHeader{
			UserId:     req.GetHead().GetUserId(),
			OwnerPath:  "", // TODO
			AccessPath: "", // TODO
		},
		UserId: req.GetValue().GetUserId(),
	}

	return reply, nil
}
func (p *Database) DeleteUsers(context.Context, *pb.DeleteUsersRequest) (*pb.DeleteUsersResponse, error) {
	panic("TODO")
}
func (p *Database) ModifyUser(context.Context, *pb.ModifyUserRequest) (*pb.ModifyUserResponse, error) {
	panic("TODO")
}
func (p *Database) DescribeUsers(context.Context, *pb.DescribeUsersRequest) (*pb.DescribeUsersResponse, error) {
	panic("TODO")
}
func (p *Database) GetUser(context.Context, *pb.GetUserRequest) (*pb.GetUserResponse, error) {
	panic("TODO")
}
