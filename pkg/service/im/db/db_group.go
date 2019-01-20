// Copyright 2018 The OpenPitrix Authors. All rights reserved.
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
	"openpitrix.io/iam/pkg/pb/im"
	"openpitrix.io/iam/pkg/service/im/db_spec"
	"openpitrix.io/logger"
)

func (p *Database) CreateGroup(ctx context.Context, req *pbim.Group) (*pbim.Group, error) {
	logger.Infof(ctx, funcutil.CallerName(1))

	if req == nil {
		err := status.Errorf(codes.InvalidArgument, "empty field")
		logger.Warnf(ctx, "%+v", err)
		return nil, err
	}
	if req != nil {
		if req.Gid == "" {
			req.Gid = genGid()
		}

		if isZeroTimestamp(req.CreateTime) {
			req.CreateTime = ptypes.TimestampNow()
		}
		if isZeroTimestamp(req.UpdateTime) {
			req.UpdateTime = ptypes.TimestampNow()
		}
		if isZeroTimestamp(req.StatusTime) {
			req.StatusTime = ptypes.TimestampNow()
		}
	}

	var dbGroup = db_spec.PBGroupToDB(req)
	if err := dbGroup.ValidateForInsert(); err != nil {
		logger.Warnf(ctx, "%+v", err)
		return nil, err
	}

	// check group_path valid
	switch {
	case dbGroup.GroupPath == dbGroup.Gid+".":
		// skip root
	case strings.HasSuffix(dbGroup.GroupPath, "."+dbGroup.Gid+"."):
		idx := len(dbGroup.GroupPath) - len(dbGroup.Gid)
		parentGroupPath := dbGroup.GroupPath[:idx-1]

		if parentGroupPath == "" {
			err := status.Errorf(codes.InvalidArgument, "invalid parent group path: %s", parentGroupPath)
			logger.Warnf(ctx, "%+v", err)
			return nil, err
		}

		var query = fmt.Sprintf(
			"SELECT COUNT(*) FROM %s WHERE group_path = '%s'",
			db_spec.UserGroupTableName, parentGroupPath,
		)
		total, err := p.getCountByQuery(ctx, query)
		if err != nil {
			logger.Warnf(ctx, "%v", query)
			logger.Warnf(ctx, "%+v", err)
			return nil, err
		}
		if total != 1 {
			err := status.Errorf(codes.InvalidArgument, "invalid group path: %s", dbGroup.GroupPath)
			logger.Warnf(ctx, "%+v", err)
			return nil, err
		}

		dbGroup.ParentGid = parentGroupPath
		dbGroup.GroupPathLevel = strings.Count(dbGroup.GroupPath, ".") + 1

		if idx = strings.LastIndex(parentGroupPath, "."); idx >= 0 {
			dbGroup.ParentGid = parentGroupPath[idx:]
		}
	}

	sql, values := pkgBuildSql_InsertInto(
		db_spec.UserGroupTableName,
		dbGroup,
	)
	if len(values) == 0 {
		err := status.Errorf(codes.InvalidArgument, "empty field")
		logger.Warnf(ctx, "%v, %v", sql, values)
		logger.Warnf(ctx, "%+v", err)
		return nil, err
	}

	_, err := p.dbx.ExecContext(ctx, sql, values...)
	if err != nil {
		logger.Warnf(ctx, "%v, %v", sql, values)
		logger.Warnf(ctx, "%+v", err)
		return nil, err
	}

	return req, nil
}

func (p *Database) DeleteGroups(ctx context.Context, req *pbim.GroupIdList) (*pbim.Empty, error) {
	logger.Infof(ctx, funcutil.CallerName(1))

	if len(req.Gid) == 1 && strings.Contains(req.Gid[0], ",") {
		req.Gid = strings.Split(req.Gid[0], ",")
	}

	if req == nil || len(req.Gid) == 0 || !isValidGids(req.Gid...) {
		err := status.Errorf(codes.InvalidArgument, "empty field")
		logger.Warnf(ctx, "%+v", err)
		return nil, err
	}

	sql := pkgBuildSql_Delete(
		db_spec.UserGroupTableName,
		db_spec.UserGroupPrimaryKeyName,
		req.Gid...,
	)

	tx, err := p.dbx.Beginx()
	if err != nil {
		logger.Warnf(ctx, "%+v", err)
		return nil, err
	}

	_, err = tx.ExecContext(ctx, sql)
	if err != nil {
		logger.Warnf(ctx, "%v", sql)
		logger.Warnf(ctx, "%+v", err)
		return nil, err
	}

	// delete binding
	for _, gid := range req.Gid {
		sql := fmt.Sprintf(
			`delete from %s where group_id=?`,
			db_spec.UserGroupBindingTableName,
		)

		_, err := tx.ExecContext(ctx, sql, gid)
		if err != nil {
			logger.Warnf(ctx, "%v", sql)
			logger.Warnf(ctx, "%+v", err)
			return nil, err
		}
		if err != nil {
			logger.Warnf(ctx, "gid = %v, err = %+v", gid, err)
		}
	}

	err = tx.Commit()
	if err != nil {
		logger.Warnf(ctx, "%+v", err)
		return nil, err
	}

	reply := &pbim.Empty{}
	return reply, nil
}

func (p *Database) ModifyGroup(ctx context.Context, req *pbim.Group) (*pbim.Group, error) {
	logger.Infof(ctx, funcutil.CallerName(1))

	if req == nil || req.Gid == "" {
		err := status.Errorf(codes.InvalidArgument, "empty field")
		logger.Warnf(ctx, "%+v", err)
		return nil, err
	}

	var dbGroup = db_spec.PBGroupToDB(req)

	// ignore read only field
	{
		dbGroup.GroupPath = ""

		dbGroup.CreateTime = time.Time{}
		dbGroup.UpdateTime = time.Now()

		switch {
		case dbGroup.Status == "":
			dbGroup.StatusTime = time.Time{}
		default:
			dbGroup.StatusTime = time.Now()
		}
	}

	if err := dbGroup.ValidateForUpdate(); err != nil {
		logger.Warnf(ctx, "%+v", err)
		return nil, err
	}

	sql, values := pkgBuildSql_Update(
		db_spec.UserGroupTableName, dbGroup,
		db_spec.UserGroupPrimaryKeyName,
	)
	if len(values) == 0 {
		return p.GetGroup(ctx, &pbim.GroupId{Gid: req.Gid})
	}

	_, err := p.dbx.ExecContext(ctx, sql, values...)
	if err != nil {
		logger.Warnf(ctx, "%v, %v", sql, values)
		logger.Warnf(ctx, "%+v", err)
		return nil, err
	}

	return p.GetGroup(ctx, &pbim.GroupId{Gid: req.Gid})
}

func (p *Database) GetGroup(ctx context.Context, req *pbim.GroupId) (*pbim.Group, error) {
	logger.Infof(ctx, funcutil.CallerName(1))

	var query = fmt.Sprintf(
		"SELECT * FROM %s WHERE %s=? LIMIT 1 OFFSET 0;",
		db_spec.UserGroupTableName,
		db_spec.UserGroupPrimaryKeyName,
	)

	var v = db_spec.DBGroup{}
	err := p.dbx.GetContext(ctx, &v, query, req.GetGid())
	if err != nil {
		logger.Warnf(ctx, "%v", query)
		logger.Warnf(ctx, "%+v", err)
		return nil, err
	}

	return v.ToPB(), nil
}

func (p *Database) ListGroups(ctx context.Context, req *pbim.ListGroupsRequest) (*pbim.ListGroupsResponse, error) {
	logger.Infof(ctx, funcutil.CallerName(1))

	// fix repeated fileds
	if len(req.Gid) == 1 && strings.Contains(req.Gid[0], ",") {
		req.Gid = strings.Split(req.Gid[0], ",")
	}
	if len(req.Uid) == 1 && strings.Contains(req.Uid[0], ",") {
		req.Uid = strings.Split(req.Uid[0], ",")
	}
	if len(req.Name) == 1 && strings.Contains(req.Name[0], ",") {
		req.Name = strings.Split(req.Name[0], ",")
	}
	if len(req.Status) == 1 && strings.Contains(req.Status[0], ",") {
		req.Status = strings.Split(req.Status[0], ",")
	}

	if err := p.validateListGroupsReq(req); err != nil {
		return nil, err
	}

	if len(req.Uid) > 0 {
		return p.listGroups_with_uid(ctx, req)
	} else {
		return p.listGroups_no_uid(ctx, req)
	}
}
