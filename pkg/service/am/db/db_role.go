// Copyright 2019 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package db

import (
	"context"
	"strings"

	"github.com/golang/protobuf/ptypes"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"openpitrix.io/iam/pkg/internal/funcutil"
	pbam "openpitrix.io/iam/pkg/pb/am"
	"openpitrix.io/iam/pkg/service/am/db_spec"
	"openpitrix.io/logger"
)

func (p *Database) CreateRole(ctx context.Context, req *pbam.Role) (*pbam.Role, error) {
	logger.Infof(ctx, funcutil.CallerName(1))

	if req == nil {
		err := status.Errorf(codes.InvalidArgument, "empty field")
		logger.Warnf(ctx, "%+v", err)
		return nil, err
	}
	if req != nil {
		if req.RoleId == "" {
			req.RoleId = genRoleId()
		}

		if isZeroTimestamp(req.CreateTime) {
			req.CreateTime = ptypes.TimestampNow()
		}
		if isZeroTimestamp(req.UpdateTime) {
			req.UpdateTime = ptypes.TimestampNow()
		}
	}

	var dbRole = db_spec.PBRoleToDB(req)
	if err := dbRole.ValidateForInsert(); err != nil {
		logger.Warnf(ctx, "%+v", err)
		return nil, err
	}

	var sql = `INSERT INTO role (
		role_id,
		role_name,
		description,
		portal,
		owner,
		owner_path,
		create_time,
		update_time
	) VALUES (
		?, -- role_id,
		?, -- role_name,
		?, -- description,
		?, -- portal,
		?, -- owner,
		?, -- owner_path,
		?, -- create_time,
		?  -- update_time
	);`

	_, err := p.DB.ExecContext(ctx, sql,
		dbRole.RoleId,      // ?, -- role_id,
		dbRole.RoleName,    // ?, -- role_name,
		dbRole.Description, // ?, -- description,
		dbRole.Portal,      // ?, -- portal,
		dbRole.Owner,       // ?, -- owner,
		dbRole.OwnerPath,   // ?, -- owner_path,
		dbRole.CreateTime,  // ?, -- create_time,
		dbRole.UpdateTime,  // ?  -- update_time
	)
	if err != nil {
		logger.Warnf(ctx, "%v, %v", sql, dbRole)
		logger.Warnf(ctx, "%+v", err)
		return nil, err
	}

	return nil, nil
}
func (p *Database) DeleteRoles(ctx context.Context, req *pbam.RoleIdList) (*pbam.Empty, error) {
	logger.Infof(ctx, funcutil.CallerName(1))

	if len(req.RoleId) == 1 && strings.Contains(req.RoleId[0], ",") {
		req.RoleId = strings.Split(req.RoleId[0], ",")
	}

	if req == nil || len(req.RoleId) == 0 || !isValidGids(req.RoleId...) {
		err := status.Errorf(codes.InvalidArgument, "empty field")
		logger.Warnf(ctx, "%+v", err)
		return nil, err
	}

	sql := pkgBuildSql_Delete(
		db_spec.RoleTableName,
		db_spec.RolePrimaryKeyName,
		req.RoleId...,
	)

	_, err := p.DB.ExecContext(ctx, sql)
	if err != nil {
		logger.Warnf(ctx, "%v", sql)
		logger.Warnf(ctx, "%+v", err)
		return nil, err
	}

	// TODO: delete binding

	return &pbam.Empty{}, nil
}
func (p *Database) ModifyRole(ctx context.Context, req *pbam.Role) (*pbam.Role, error) {
	logger.Infof(ctx, funcutil.CallerName(1))

	panic("todo")
}
func (p *Database) DescribeRoles(ctx context.Context, req *pbam.DescribeRolesRequest) (*pbam.RoleList, error) {
	logger.Infof(ctx, funcutil.CallerName(1))

	panic("todo")
}
