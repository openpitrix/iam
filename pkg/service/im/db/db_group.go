// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package db

import (
	"context"
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

	// gen group_path from parent_id
	if req.ParentGroupId != "" {
		if parent, err := p.GetGroup(ctx, &pbim.GroupId{GroupId: req.ParentGroupId}); err == nil {
			if dbGroup.GroupPath == "" {
				dbGroup.GroupPath = parent.GroupPath + "." + dbGroup.GroupId
			}
		} else {
			err := status.Errorf(codes.InvalidArgument, "invalid parent_id: %v", req.ParentGroupId)
			logger.Warnf(ctx, "%+v", err)
			return nil, err
		}
	}

	// check group_path valid
	switch {
	case dbGroup.GroupPath == "":
		dbGroup.GroupPath = dbGroup.GroupId

	case dbGroup.GroupPath == dbGroup.GroupId:
		// skip root

	case strings.HasSuffix(dbGroup.GroupPath, "."+dbGroup.GroupId):
		// check parent path
		ids := strings.Split(dbGroup.GroupPath, ".")
		if len(ids) > 1 {
			ids = ids[:len(ids)-1]
		}

		if len(ids) > 0 {
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
		}

	default:
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

	if req == nil || len(req.GroupId) == 0 || !isValidIds(req.GroupId...) {
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

	var v = UserGroup{GroupId: req.GroupId}
	if err := p.DB.Model(User{}).Take(&v).Error; err != nil {
		logger.Warnf(ctx, "%+v", err)
		return nil, err
	}

	return v.ToPB(), nil
}

func (p *Database) ListGroups(ctx context.Context, req *pbim.ListGroupsRequest) (*pbim.ListGroupsResponse, error) {
	logger.Infof(ctx, funcutil.CallerName(1))

	// fix repeated fileds
	if len(req.UserId) == 1 && strings.Contains(req.UserId[0], ",") {
		req.UserId = strings.Split(req.UserId[0], ",")
	}
	if len(req.GroupId) == 1 && strings.Contains(req.GroupId[0], ",") {
		req.GroupId = strings.Split(req.GroupId[0], ",")
	}
	if len(req.GroupName) == 1 && strings.Contains(req.GroupName[0], ",") {
		req.GroupName = strings.Split(req.GroupName[0], ",")
	}
	if len(req.Status) == 1 && strings.Contains(req.Status[0], ",") {
		req.Status = strings.Split(req.Status[0], ",")
	}

	req.UserId = trimEmptyString(req.UserId)
	req.GroupId = trimEmptyString(req.GroupId)
	req.GroupName = trimEmptyString(req.GroupName)
	req.Status = trimEmptyString(req.Status)

	// limit & offset
	if req.Limit == 0 && req.Offset == 0 {
		req.Limit = 20
		req.Offset = 0
	}
	if req.Limit < 0 {
		req.Limit = 0
	}
	if req.Limit > 200 {
		req.Limit = 200
	}
	if req.Offset < 0 {
		req.Offset = 0
	}

	if !isValidIds(req.GroupId...) {
		err := status.Errorf(codes.InvalidArgument, "invalid gid: %v", req.GroupId)
		logger.Warnf(ctx, "%+v", err)
		return nil, err
	}
	if !isValidIds(req.UserId...) {
		err := status.Errorf(codes.InvalidArgument, "invalid uid: %v", req.UserId)
		logger.Warnf(ctx, "%+v", err)
		return nil, err
	}
	if !isValidSearchWord(req.SearchWord) {
		err := status.Errorf(codes.InvalidArgument, "invalid search_word: %v", req.SearchWord)
		logger.Warnf(ctx, "%+v", err)
		return nil, err
	}
	if !isValidSortKey(req.SortKey) {
		err := status.Errorf(codes.InvalidArgument, "invalid sort_key: %v", req.SortKey)
		logger.Warnf(ctx, "%+v", err)
		return nil, err
	}

	var (
		inKeys   []string
		inValues []interface{}
	)

	if len(req.UserId) > 0 {
		inKeys = append(inKeys, "user.user_id in(?)")
		inValues = append(inValues, req.UserId)
	}
	if len(req.GroupId) > 0 {
		inKeys = append(inKeys, "user_group.group_id in(?)")
		inValues = append(inValues, req.GroupId)
	}
	if len(req.GroupName) > 0 {
		inKeys = append(inKeys, "user_group.group_name in(?)")
		inValues = append(inValues, req.GroupName)
	}
	if len(req.Status) > 0 {
		inKeys = append(inKeys, "user_group.status in(?)")
		inValues = append(inValues, req.Status)
	}

	if req.SearchWord != "" {
		var likeKey = "%" + req.SearchWord + "%"

		var likeSql = `(1=0`
		likeSql += " OR user_group.group_id LIKE ?"
		likeSql += " OR user_group.group_name LIKE ?"
		likeSql += " OR user_group.description LIKE ?"
		likeSql += " OR user_group.status LIKE ?"
		likeSql += ")"

		inKeys = append(inKeys, likeSql)
		inValues = append(inValues,
			likeKey, // group_id
			likeKey, // group_name
			likeKey, // description
			likeKey, // status
		)
	}

	var query = ""
	if len(inKeys) > 0 {
		if len(req.UserId) > 0 {
			query += "SELECT user_group.* from user, user_group, user_group_binding"
			query += " WHERE "
			query += " user_group_binding.user_id=user.user_id AND"
			query += " user_group_binding.group_id=user_group.group_id AND"
			query += strings.Join(inKeys, " AND ")
		} else {
			query += "SELECT user_group.* from user_group"
			query += " WHERE "
			query += strings.Join(inKeys, " AND ")
		}
	} else {
		query = "SELECT user_group.* from user_group WHERE 1=1"
		inValues = nil
	}

	logger.Infof(ctx, "query: %s", query)
	logger.Infof(ctx, "inValues: %v", inValues)

	var rows []UserGroup
	p.DB.Limit(req.Limit).Offset(req.Offset).Where(query, inValues...).Find(&rows)
	if err := p.DB.Error; err != nil {
		logger.Warnf(ctx, "%+v", err)
		return nil, err
	}

	var sets []*pbim.Group
	for _, v := range rows {
		sets = append(sets, v.ToPB())
	}

	reply := &pbim.ListGroupsResponse{
		Group: sets,
	}

	return reply, nil
}
