// Copyright 2019 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package db

import (
	"context"
	"fmt"
	"strings"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"openpitrix.io/iam/pkg/internal/funcutil"
	pbam "openpitrix.io/iam/pkg/pb/am"
	"openpitrix.io/iam/pkg/service/am/db_spec"
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

	sql := fmt.Sprintf(
		`insert into %s (id, user_id, role_id) values(?,?,?)`,
		db_spec.UserRoleBindingTableName,
	)

	switch {
	case len(req.UserId) == len(req.RoleId):
		for i := 0; i < len(req.RoleId); i++ {
			xid := genXid()
			role_id := req.RoleId[i]
			user_id := req.UserId[i]

			_, err := p.DB.ExecContext(ctx, sql, xid, user_id, role_id)
			if err != nil {
				logger.Warnf(ctx, "%v", sql)
				logger.Warnf(ctx, "%v, %v, %v", xid, user_id, role_id)
				logger.Warnf(ctx, "%+v", err)
				return nil, err
			}
		}
	case len(req.UserId) == 1:
		for i := 0; i < len(req.RoleId); i++ {
			xid := genXid()
			role_id := req.RoleId[i]
			user_id := req.UserId[0]

			_, err := p.DB.ExecContext(ctx, sql, xid, user_id, role_id)
			if err != nil {
				logger.Warnf(ctx, "%v", sql)
				logger.Warnf(ctx, "%v, %v, %v", xid, user_id, role_id)
				logger.Warnf(ctx, "%+v", err)
				return nil, err
			}
		}
	case len(req.RoleId) == 1:
		for i := 0; i < len(req.UserId); i++ {
			xid := genXid()
			role_id := req.RoleId[0]
			user_id := req.UserId[i]

			_, err := p.DB.ExecContext(ctx, sql, xid, user_id, role_id)
			if err != nil {
				logger.Warnf(ctx, "%v", sql)
				logger.Warnf(ctx, "%v, %v, %v", xid, user_id, role_id)
				logger.Warnf(ctx, "%+v", err)
				return nil, err
			}
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

	sql := fmt.Sprintf(
		`delete from %s where user_id=? AND role_id=?`,
		db_spec.UserRoleBindingTableName,
	)

	switch {
	case len(req.UserId) == len(req.RoleId):
		for i := 0; i < len(req.RoleId); i++ {
			xid := genXid()
			role_id := req.RoleId[i]
			user_id := req.UserId[i]

			_, err := p.DB.ExecContext(ctx, sql, xid, user_id, role_id)
			if err != nil {
				logger.Warnf(ctx, "%v", sql)
				logger.Warnf(ctx, "%v, %v, %v", xid, user_id, role_id)
				logger.Warnf(ctx, "%+v", err)
				return nil, err
			}
		}
	case len(req.UserId) == 1:
		for i := 0; i < len(req.RoleId); i++ {
			xid := genXid()
			role_id := req.RoleId[i]
			user_id := req.UserId[0]

			_, err := p.DB.ExecContext(ctx, sql, xid, user_id, role_id)
			if err != nil {
				logger.Warnf(ctx, "%v", sql)
				logger.Warnf(ctx, "%v, %v, %v", xid, user_id, role_id)
				logger.Warnf(ctx, "%+v", err)
				return nil, err
			}
		}
	case len(req.RoleId) == 1:
		for i := 0; i < len(req.UserId); i++ {
			xid := genXid()
			role_id := req.RoleId[0]
			user_id := req.UserId[i]

			_, err := p.DB.ExecContext(ctx, sql, xid, user_id, role_id)
			if err != nil {
				logger.Warnf(ctx, "%v", sql)
				logger.Warnf(ctx, "%v, %v, %v", xid, user_id, role_id)
				logger.Warnf(ctx, "%+v", err)
				return nil, err
			}
		}
	}

	return &pbam.Empty{}, nil
}
