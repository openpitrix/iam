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

func (p *Database) CreateRole(ctx context.Context, req *pb.CreateRoleRequest) (*pb.CreateRoleResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}

	sql, values := pkgBuildSql_InsertInto(
		dbSpec.RoleTableName,
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

	reply := &pb.CreateRoleResponse{
		Head: &pb.ResponseHeader{
			UserId:     req.GetHead().GetUserId(),
			OwnerPath:  "", // TODO
			AccessPath: "", // TODO
		},
		RoleId: req.GetValue().GetRoleId(),
	}

	return reply, nil
}

func (p *Database) DeleteRoles(ctx context.Context, req *pb.DeleteRolesRequest) (*pb.DeleteRolesResponse, error) {
	sql := pkgBuildSql_Delete(
		dbSpec.RoleTableName, dbSpec.RolePrimaryKeyName,
		req.RoleId...,
	)

	_, err := p.DB.ExecContext(ctx, sql)
	if err != nil {
		return nil, err
	}

	reply := &pb.DeleteRolesResponse{
		Head: &pb.ResponseHeader{
			UserId:     req.GetHead().GetUserId(),
			OwnerPath:  "", // TODO
			AccessPath: "", // TODO
		},
		RoleId: req.RoleId,
	}

	return reply, nil
}
func (p *Database) ModifyRole(ctx context.Context, req *pb.ModifyRoleRequest) (*pb.ModifyRoleResponse, error) {
	sql, values := pkgBuildSql_Update(
		dbSpec.RoleTableName, req.GetValue(),
		dbSpec.RolePrimaryKeyName,
	)

	_, err := p.DB.ExecContext(ctx, sql, values...)
	if err != nil {
		return nil, err
	}

	reply := &pb.ModifyRoleResponse{
		Head: &pb.ResponseHeader{
			UserId:     req.GetHead().GetUserId(),
			OwnerPath:  "", // TODO
			AccessPath: "", // TODO
		},
		RoleId: req.GetValue().GetRoleId(),
	}

	return reply, nil
}
func (p *Database) GetRole(ctx context.Context, req *pb.GetRoleRequest) (*pb.GetRoleResponse, error) {
	var sql = fmt.Sprintf(
		"SELECT * FROM %s WHERE %s=$1",
		dbSpec.RoleTableName,
		dbSpec.RolePrimaryKeyName,
	)
	var value pb.Role
	err := p.DB.Get(&value, sql, req.GetRoleId())
	if err != nil {
		return nil, err
	}

	reply := &pb.GetRoleResponse{
		Head: &pb.ResponseHeader{
			UserId:     req.GetHead().GetUserId(),
			OwnerPath:  "", // TODO
			AccessPath: "", // TODO
		},
		Value: &value,
	}

	return reply, nil
}
func (p *Database) DescribeRoles(context.Context, *pb.DescribeRolesRequest) (*pb.DescribeRolesResponse, error) {
	panic("TODO")
}
