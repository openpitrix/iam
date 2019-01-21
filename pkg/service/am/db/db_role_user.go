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
	logger.Infof(ctx, funcutil.CallerName(1)+":11")
	logger.Infof(ctx, "req: %v", req)

	tx := p.DB.Begin()

	switch {
	case len(req.UserId) == len(req.RoleId):
		for i := 0; i < len(req.RoleId); i++ {
			tx.Exec(
				`INSERT INTO user_role_binding (id, user_id, role_id) VALUES (?,?,?)`,
				genXid(), req.UserId[i], req.RoleId[i],
			)
			if err := tx.Error; err != nil {
				tx.Rollback()
				return nil, err
			}
		}
	case len(req.UserId) == 1:
		logger.Infof(ctx, "debug: req: %v", req)
		for i := 0; i < len(req.RoleId); i++ {
			tx.Exec(
				`INSERT INTO user_role_binding (id, user_id, role_id) VALUES (?,?,?)`,
				genXid(), req.UserId[0], req.RoleId[i],
			)
			if err := tx.Error; err != nil {
				tx.Rollback()
				return nil, err
			}
		}
	case len(req.RoleId) == 1:
		for i := 0; i < len(req.UserId); i++ {
			tx.Exec(
				`INSERT INTO user_role_binding (id, user_id, role_id) VALUES (?,?,?)`,
				genXid(), req.UserId[i], req.RoleId[0],
			)
			if err := tx.Error; err != nil {
				tx.Rollback()
				return nil, err
			}
		}
	}

	if err := tx.Commit().Error; err != nil {
		logger.Warnf(ctx, "%+v", err)
		return nil, err
	}
	logger.Infof(ctx, funcutil.CallerName(1)+":22")

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

	tx := p.DB.Begin()

	switch {
	case len(req.UserId) == len(req.RoleId):
		for i := 0; i < len(req.RoleId); i++ {
			logger.Infof(ctx, "req: %v", req)

			tx.Exec(
				`delete from user_role_binding where user_id=? and role_id=?`,
				req.UserId[i],
				req.RoleId[i],
			)
			if err := tx.Error; err != nil {
				tx.Rollback()
				return nil, err
			}
		}
	case len(req.UserId) == 1:
		for i := 0; i < len(req.RoleId); i++ {
			tx.Exec(
				`delete from user_role_binding where user_id=? and role_id=?`,
				req.UserId[0],
				req.RoleId[i],
			)
			if err := tx.Error; err != nil {
				tx.Rollback()
				return nil, err
			}
		}
	case len(req.RoleId) == 1:
		for i := 0; i < len(req.UserId); i++ {
			tx.Exec(
				`delete from user_role_binding where user_id=? and role_id=?`,
				req.UserId[i],
				req.RoleId[0],
			)
			if err := tx.Error; err != nil {
				tx.Rollback()
				return nil, err
			}
		}
	}

	if err := tx.Commit().Error; err != nil {
		logger.Warnf(ctx, "%+v", err)
		return nil, err
	}

	return &pbam.Empty{}, nil
}
