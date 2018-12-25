// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package service

import (
	"context"
	"fmt"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"openpitrix.io/iam/openpitrix/pkg/pb"
)

func (p *Database) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}

	sql, values := pkgBuildSql_InsertInto(
		dbSpec.UserTableName,
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
func (p *Database) DeleteUsers(ctx context.Context, req *pb.DeleteUsersRequest) (*pb.DeleteUsersResponse, error) {
	sql := pkgBuildSql_Delete(
		dbSpec.UserTableName, dbSpec.UserPrimaryKeyName,
		req.UserId...,
	)

	_, err := p.DB.ExecContext(ctx, sql)
	if err != nil {
		return nil, err
	}

	reply := &pb.DeleteUsersResponse{
		Head: &pb.ResponseHeader{
			UserId:     req.GetHead().GetUserId(),
			OwnerPath:  "", // TODO
			AccessPath: "", // TODO
		},
		UserId: req.UserId,
	}

	return reply, nil
}
func (p *Database) ModifyUser(ctx context.Context, req *pb.ModifyUserRequest) (*pb.ModifyUserResponse, error) {
	sql, values := pkgBuildSql_Update(
		dbSpec.UserTableName, req.GetValue(),
		dbSpec.UserPrimaryKeyName,
	)

	_, err := p.DB.ExecContext(ctx, sql, values...)
	if err != nil {
		return nil, err
	}

	reply := &pb.ModifyUserResponse{
		Head: &pb.ResponseHeader{
			UserId:     req.GetHead().GetUserId(),
			OwnerPath:  "", // TODO
			AccessPath: "", // TODO
		},
		UserId: req.GetValue().GetUserId(),
	}

	return reply, nil
}
func (p *Database) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.GetUserResponse, error) {
	var sql = fmt.Sprintf(
		"SELECT * FROM %s WHERE %s=$1",
		dbSpec.UserTableName,
		dbSpec.UserPrimaryKeyName,
	)
	var value pb.User
	err := p.DB.Get(&value, sql, req.GetUserId())
	if err != nil {
		return nil, err
	}

	reply := &pb.GetUserResponse{
		Head: &pb.ResponseHeader{
			UserId:     req.GetHead().GetUserId(),
			OwnerPath:  "", // TODO
			AccessPath: "", // TODO
		},
		Value: &value,
	}

	return reply, nil
}
func (p *Database) DescribeUsers(context.Context, *pb.DescribeUsersRequest) (*pb.DescribeUsersResponse, error) {
	panic("TODO")
}
