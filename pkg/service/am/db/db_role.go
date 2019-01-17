// Copyright 2019 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package db

import (
	"context"
	"fmt"
	"strings"
	"time"

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

	if req == nil || req.RoleId == "" {
		err := status.Errorf(codes.InvalidArgument, "empty field")
		logger.Warnf(ctx, "%+v", err)
		return nil, err
	}

	var dbRole = db_spec.PBRoleToDB(req)

	// ignore read only field
	{
		dbRole.CreateTime = time.Time{}
		dbRole.UpdateTime = time.Now()
	}

	if err := dbRole.ValidateForUpdate(); err != nil {
		logger.Warnf(ctx, "%+v", err)
		return nil, err
	}

	sql, values := pkgBuildSql_Update(
		db_spec.RoleTableName, dbRole,
		db_spec.RolePrimaryKeyName,
	)
	if len(values) == 0 {
		return p.GetRole(ctx, &pbam.RoleId{RoleId: req.RoleId})
	}

	_, err := p.DB.ExecContext(ctx, sql, values...)
	if err != nil {
		logger.Warnf(ctx, "%v, %v", sql, values)
		logger.Warnf(ctx, "%+v", err)
		return nil, err
	}

	return p.GetRole(ctx, &pbam.RoleId{RoleId: req.RoleId})
}

func (p *Database) GetRole(ctx context.Context, req *pbam.RoleId) (*pbam.Role, error) {
	logger.Infof(ctx, funcutil.CallerName(1))

	var query = fmt.Sprintf(
		"SELECT * FROM %s WHERE %s=? LIMIT 1 OFFSET 0;",
		db_spec.RoleTableName,
		db_spec.RolePrimaryKeyName,
	)

	var v = db_spec.DBRole{}
	err := p.DB.GetContext(ctx, &v, query, req.GetRoleId())
	if err != nil {
		logger.Warnf(ctx, "%v", query)
		logger.Warnf(ctx, "%+v", err)
		return nil, err
	}

	return v.ToPB(), nil
}
func (p *Database) DescribeRoles(ctx context.Context, req *pbam.DescribeRolesRequest) (*pbam.RoleList, error) {
	logger.Infof(ctx, funcutil.CallerName(1))

	// fix repeated fileds
	if len(req.RoleId) == 1 && strings.Contains(req.RoleId[0], ",") {
		req.RoleId = strings.Split(req.RoleId[0], ",")
	}
	if len(req.RoleName) == 1 && strings.Contains(req.RoleName[0], ",") {
		req.RoleName = strings.Split(req.RoleName[0], ",")
	}
	if len(req.Portal) == 1 && strings.Contains(req.Portal[0], ",") {
		req.Portal = strings.Split(req.Portal[0], ",")
	}
	if len(req.UserId) == 1 && strings.Contains(req.UserId[0], ",") {
		req.UserId = strings.Split(req.UserId[0], ",")
	}

	// WHERE name IN (name1,name2) AND name LIKE '%search_word%'
	var whereCondition = func() string {
		ss := genWhereCondition(
			map[string][]string{
				"role_id":   req.RoleId,
				"role_name": req.RoleName,
				"portal":    req.Portal,
				"user_id":   req.UserId,
			},
			nil, "",
		)

		if len(ss) > 0 {
			return "WHERE " + strings.Join(ss, " AND ")
		}
		return ""
	}()

	var query = fmt.Sprintf(
		"SELECT * FROM role %s;",
		whereCondition,
	)

	var rows = []db_spec.DBRole{}
	err := p.DB.SelectContext(ctx, &rows, query)
	if err != nil {
		logger.Warnf(ctx, "%v", query)
		logger.Warnf(ctx, "%+v", err)
		return nil, err
	}

	var sets []*pbam.Role
	for _, v := range rows {
		sets = append(sets, v.ToPB())
	}

	reply := &pbam.RoleList{
		Value: sets,
	}

	return reply, nil
}
