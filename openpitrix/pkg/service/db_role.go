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

func (p *Database) CreateRole(ctx context.Context, req *pb.CreateRoleRequest) (*pb.CreateRoleResponse, error) {
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

func (p *Database) DeleteRole(context.Context, *pb.DeleteRoleRequest) (*pb.DeleteRoleResponse, error) {
	panic("TODO")
}
func (p *Database) ModifyRole(context.Context, *pb.ModifyRoleRequest) (*pb.ModifyRoleResponse, error) {
	panic("TODO")
}
func (p *Database) GetRole(context.Context, *pb.GetRoleRequest) (*pb.GetRoleResponse, error) {
	panic("TODO")
}
func (p *Database) DescribeRoles(context.Context, *pb.DescribeRolesRequest) (*pb.DescribeRolesResponse, error) {
	panic("TODO")
}
