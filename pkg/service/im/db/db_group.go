// Copyright 2018 The OpenPitrix Authors. All rights reserved.
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
	"openpitrix.io/iam/pkg/pb/im"
	"openpitrix.io/logger"
)

func (p *Database) CreateGroup(ctx context.Context, req *pbim.Group) (*pbim.Group, error) {
	logger.Infof(ctx, funcutil.CallerName(1))

	var dbGroup = NewUserGroupFromPB(req)
	if err := dbGroup.ValidateForInsert(); err != nil {
		logger.Warnf(ctx, "%+v", err)
		return nil, err
	}

	// check group_path valid
	switch {
	case dbGroup.GroupPath == "":
		dbGroup.GroupPath = dbGroup.GroupId + "."

	case dbGroup.GroupPath == dbGroup.GroupId+".":
		// skip root

	case strings.HasSuffix(dbGroup.GroupPath, "."+dbGroup.GroupId+"."):
		prefixPath := dbGroup.GroupPath[:len(dbGroup.GroupPath)-len("."+dbGroup.GroupId+".")]
		ids := strings.Split(prefixPath, ".")

		var count int
		p.DB.Model(&UserGroup{}).Where("group_id in (?)", ids).Count(&count)
		if err := p.DB.Error; err != nil {
			logger.Warnf(ctx, "%+v", err)
			return nil, err
		}
		if count != len(ids) {
			err := status.Errorf(codes.InvalidArgument, "invalid parent group path: %s", dbGroup.GroupPath)
			logger.Warnf(ctx, "%+v", err)
			return nil, err
		}

	default: // check parent path
		err := status.Errorf(codes.InvalidArgument, "invalid parent group path: %s", dbGroup.GroupPath)
		logger.Warnf(ctx, "%+v", err)
		return nil, err
	}

	if err := p.DB.Create(dbGroup).Error; err != nil {
		logger.Warnf(ctx, "%+v, %v", err, dbGroup)
		return nil, err
	}

	return req, nil
}

func (p *Database) DeleteGroups(ctx context.Context, req *pbim.GroupIdList) (*pbim.Empty, error) {
	logger.Infof(ctx, funcutil.CallerName(1))

	if req == nil || len(req.GroupId) == 0 || !isValidGids(req.GroupId...) {
		err := status.Errorf(codes.InvalidArgument, "empty field")
		logger.Warnf(ctx, "%+v", err)
		return nil, err
	}

	tx := p.DB.Begin()
	{
		if err := tx.Delete(UserGroup{}, "group_id in (?)", req.GroupId).Error; err != nil {
			tx.Rollback()
			return nil, err
		}
		if err := tx.Delete(UserGroupBinding{}, "group_id in (?)", req.GroupId).Error; err != nil {
			tx.Rollback()
			return nil, err
		}
	}

	if err := tx.Commit().Error; err != nil {
		logger.Warnf(ctx, "%+v", err)
		return nil, err
	}

	reply := &pbim.Empty{}
	return reply, nil
}

func (p *Database) ModifyGroup(ctx context.Context, req *pbim.Group) (*pbim.Group, error) {
	logger.Infof(ctx, funcutil.CallerName(1))

	if req.GroupId == "" {
		err := status.Errorf(codes.InvalidArgument, "empty field")
		logger.Warnf(ctx, "%+v", err)
		return nil, err
	}

	var dbGroup = NewUserGroupFromPB(req)
	if err := p.DB.Model(dbGroup).Updates(dbGroup).Error; err != nil {
		logger.Warnf(ctx, "%+v", err)
		return nil, err
	}

	return p.GetGroup(ctx, &pbim.GroupId{GroupId: req.GroupId})
}

func (p *Database) GetGroup(ctx context.Context, req *pbim.GroupId) (*pbim.Group, error) {
	logger.Infof(ctx, funcutil.CallerName(1))

	var query = fmt.Sprintf(
		"SELECT * FROM user_group WHERE group_id=? LIMIT 1 OFFSET 0;")

	var v = UserGroup{GroupId: req.GroupId}
	if err := p.DB.Model(User{}).Take(&v).Error; err != nil {
		logger.Warnf(ctx, "%v", query)
		logger.Warnf(ctx, "%+v", err)
		return nil, err
	}

	return v.ToPB(), nil
}

func (p *Database) ListGroups(ctx context.Context, req *pbim.ListGroupsRequest) (*pbim.ListGroupsResponse, error) {
	logger.Infof(ctx, funcutil.CallerName(1))

	// fix repeated fileds
	if len(req.GroupId) == 1 && strings.Contains(req.GroupId[0], ",") {
		req.GroupId = strings.Split(req.GroupId[0], ",")
	}
	if len(req.UserId) == 1 && strings.Contains(req.UserId[0], ",") {
		req.UserId = strings.Split(req.UserId[0], ",")
	}
	if len(req.GroupName) == 1 && strings.Contains(req.GroupName[0], ",") {
		req.GroupName = strings.Split(req.GroupName[0], ",")
	}
	if len(req.Status) == 1 && strings.Contains(req.Status[0], ",") {
		req.Status = strings.Split(req.Status[0], ",")
	}

	if err := p.validateListGroupsReq(req); err != nil {
		return nil, err
	}

	if len(req.UserId) > 0 {
		return p.listGroups_with_uid(ctx, req)
	} else {
		return p.listGroups_no_uid(ctx, req)
	}
}
