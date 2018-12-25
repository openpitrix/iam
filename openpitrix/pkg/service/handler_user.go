// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package service

import (
	"context"

	"openpitrix.io/iam/openpitrix/pkg/pb"
)

func (p *Server) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
	return p.db.CreateUser(ctx, req)
}
func (p *Server) DeleteUsers(ctx context.Context, req *pb.DeleteUsersRequest) (*pb.DeleteUsersResponse, error) {
	return p.db.DeleteUsers(ctx, req)
}
func (p *Server) ModifyUser(ctx context.Context, req *pb.ModifyUserRequest) (*pb.ModifyUserResponse, error) {
	return p.db.ModifyUser(ctx, req)
}
func (p *Server) DescribeUsers(context.Context, *pb.DescribeUsersRequest) (*pb.DescribeUsersResponse, error) {
	panic("TODO")
}
func (p *Server) GetUser(context.Context, *pb.GetUserRequest) (*pb.GetUserResponse, error) {
	panic("TODO")
}
