// Copyright 2019 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package db

import (
	"context"
	"strings"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"openpitrix.io/iam/pkg/internal/funcutil"
	pbam "openpitrix.io/iam/pkg/pb/am"
	"openpitrix.io/logger"
)

func (p *Database) BindUserRole(ctx context.Context, req *pbam.BindUserRoleRequest) (*pbam.Empty, error) {
	logger.Infof(ctx, funcutil.CallerName(1))

	if len(req.RoleId) == 1 && strings.Contains(req.RoleId[0], ",") {
		req.RoleId = strings.Split(req.RoleId[0], ",")
	}
	if len(req.UserId) == 1 && strings.Contains(req.UserId[0], ",") {
		req.UserId = strings.Split(req.UserId[0], ",")
	}

	if len(req.UserId) == 0 || len(req.RoleId) == 0 {
		err := status.Errorf(codes.InvalidArgument, "empty user_id or role_id")
		logger.Warnf(ctx, "%+v", err)
		return nil, err
	}
	if !(len(req.UserId) == 1 || len(req.RoleId) == 1 || len(req.UserId) == len(req.RoleId)) {
		err := status.Errorf(codes.InvalidArgument,
			"user_id and role_id donot math: user_id = %v, role_id = %v",
			req.UserId, req.RoleId,
		)
		logger.Warnf(ctx, "%+v", err)
		return nil, err
	}

	type UserRoleBinding struct {
		Id     string
		RoleId string
		UserId string
	}

	switch {
	case len(req.UserId) == len(req.RoleId):
		for i := 0; i < len(req.RoleId); i++ {
			p.DB.NewRecord(UserRoleBinding{
				Id:     genXid(),
				RoleId: req.RoleId[i],
				UserId: req.UserId[i],
			})
		}
	case len(req.UserId) == 1:
		for i := 0; i < len(req.RoleId); i++ {
			p.DB.NewRecord(UserRoleBinding{
				Id:     genXid(),
				RoleId: req.RoleId[i],
				UserId: req.UserId[0],
			})
		}
	case len(req.RoleId) == 1:
		for i := 0; i < len(req.UserId); i++ {
			p.DB.NewRecord(UserRoleBinding{
				Id:     genXid(),
				RoleId: req.RoleId[0],
				UserId: req.UserId[i],
			})
		}
	}

	return &pbam.Empty{}, nil
}

func (p *Database) UnbindUserRole(ctx context.Context, req *pbam.UnbindUserRoleRequest) (*pbam.Empty, error) {
	logger.Infof(ctx, funcutil.CallerName(1))

	if len(req.RoleId) == 1 && strings.Contains(req.RoleId[0], ",") {
		req.RoleId = strings.Split(req.RoleId[0], ",")
	}
	if len(req.UserId) == 1 && strings.Contains(req.UserId[0], ",") {
		req.UserId = strings.Split(req.UserId[0], ",")
	}

	if len(req.UserId) == 0 || len(req.RoleId) == 0 {
		err := status.Errorf(codes.InvalidArgument, "empty user_id or role_id")
		logger.Warnf(ctx, "%+v", err)
		return nil, err
	}
	if !(len(req.UserId) == 1 || len(req.RoleId) == 1 || len(req.UserId) == len(req.RoleId)) {
		err := status.Errorf(codes.InvalidArgument,
			"user_id and role_id donot math: user_id = %v, role_id = %v",
			req.UserId, req.RoleId,
		)
		logger.Warnf(ctx, "%+v", err)
		return nil, err
	}

	type UserRoleBinding struct {
		Id     string
		RoleId string
		UserId string
	}
	switch {
	case len(req.UserId) == len(req.RoleId):
		for i := 0; i < len(req.RoleId); i++ {
			p.DB.Delete(UserRoleBinding{
				RoleId: req.RoleId[i],
				UserId: req.UserId[i],
			})
		}
	case len(req.UserId) == 1:
		for i := 0; i < len(req.RoleId); i++ {
			p.DB.Delete(UserRoleBinding{
				RoleId: req.RoleId[1],
				UserId: req.UserId[0],
			})
		}
	case len(req.RoleId) == 1:
		for i := 0; i < len(req.UserId); i++ {
			p.DB.Delete(UserRoleBinding{
				RoleId: req.RoleId[0],
				UserId: req.UserId[i],
			})
		}
	}

	return &pbam.Empty{}, nil
}
