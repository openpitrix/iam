// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package db

import (
	"context"
	"fmt"
	"strings"

	"openpitrix.io/iam/pkg/pb/im"
	"openpitrix.io/logger"
)

func (p *Database) validateListGroupsReq(req *pbim.ListGroupsRequest) error {
	if !isValidSearchWord(req.SearchWord) {
		return fmt.Errorf("invalid search_word: %v", req.SearchWord)
	}
	if !isValidSortKey(req.SortKey) {
		return fmt.Errorf("invalid sort_key: %v", req.SortKey)
	}

	if req.Offset < 0 {
		return fmt.Errorf("invalid offset: %v", req.Offset)
	}
	if req.Limit < 0 || req.Limit > 200 {
		return fmt.Errorf("invalid limit: %v", req.Limit)
	}

	// check repeaded fields
	if !isValidIds(req.GroupId...) {
		return fmt.Errorf("invalid gid: %v", req.GroupId)
	}
	if !isValidIds(req.UserId...) {
		return fmt.Errorf("invalid uid: %v", req.UserId)
	}

	return nil
}

func (p *Database) listGroups_no_uid(ctx context.Context, req *pbim.ListGroupsRequest) (*pbim.ListGroupsResponse, error) {
	if len(req.UserId) > 0 {
		panic("uid should nil")
	}

	// WHERE name IN (name1,name2) AND name LIKE '%search_word%'
	var whereCondition = func() string {
		ss := genWhereCondition(
			map[string][]string{
				"group_id":   req.GroupId,
				"group_name": req.GroupName,
				"status":     req.Status,
			},
			[]string{
				"group_id",
				"group_name",
				"description",
				"status",
			},
			req.SearchWord,
		)
		if len(ss) > 0 {
			return "WHERE " + strings.Join(ss, " AND ")
		}

		return ""
	}()

	// ORDER BY column ASC|DESC;
	var orderBy = genOrderBy(req.SortKey, req.Reverse)

	// LIMIT 20 OFFSET 0
	var limitOffset = genLimitOffset(req.GetLimit(), req.GetOffset())

	// SELECT COUNT(*) FROM user_group {WHERE ...}
	query := fmt.Sprintf(
		"SELECT COUNT(*) FROM user_group %s",
		whereCondition,
	)

	logger.Infof(ctx, "4")
	logger.Infof(ctx, "4.1, query: %s", query)
	var total int
	p.DB.Raw(query).Count(&total)
	if err := p.DB.Error; err != nil {
		logger.Warnf(ctx, "%v", query)
		logger.Warnf(ctx, "%+v", err)
		return nil, err
	}

	logger.Infof(ctx, "5")
	// SELECT * FROM user_group {WHERE ...} {ORDER BY ...} {LIMIT ...}
	query = fmt.Sprintf(
		"select * FROM user_group %s %s %s",
		whereCondition, orderBy, limitOffset,
	)
	logger.Infof(ctx, "5.1, query: %s", query)
	reply, err := p.listGroupsByQuery(ctx, query)
	if err != nil {
		logger.Warnf(ctx, "%v", query)
		logger.Warnf(ctx, "%+v", err)
		return nil, err
	}
	logger.Infof(ctx, "6")

	reply.Total = int32(total)
	return reply, nil
}

func (p *Database) listGroups_with_uid(ctx context.Context, req *pbim.ListGroupsRequest) (*pbim.ListGroupsResponse, error) {
	if len(req.UserId) == 0 {
		panic("uid is empty")
	}

	// select user_group.* from
	//     user, user_group, user_group_binding
	// where
	//     user_group_binding.user_id=user.user_id and
	//     user_group_binding.user_id=user_group.group_id and
	//     ...

	// WHERE name IN (name1,name2) AND name LIKE '%search_word%'
	var whereCondition = func() string {
		ss := genWhereCondition(
			map[string][]string{
				"user_group.group_id":   req.GroupId,
				"user_group.group_name": req.GroupName,
				"user.user_id":          req.UserId,
				"user_group.status":     req.Status,
			},
			[]string{
				"user.user_id",
				"user_group.group_name",
				"user_group.description",
				"user_group.status",
			},
			req.SearchWord,
		)

		ss = append(
			[]string{
				"user_group_binding.user_id=user.user_id",
				"user_group_binding.group_id=user_group.group_id",
			},
			ss...,
		)
		if len(ss) > 0 {
			return "WHERE " + strings.Join(ss, " AND ")
		}

		return "WHERE " + strings.Join(ss, " AND ")
	}()

	// ORDER BY column ASC|DESC;
	var orderBy = genOrderBy(req.SortKey, req.Reverse)

	// LIMIT 20 OFFSET 0
	var limitOffset = genLimitOffset(req.GetLimit(), req.GetOffset())

	// SELECT COUNT(user_group.*) FROM user, user_group, user_group_binding {WHERE ...}
	query := fmt.Sprintf(
		"SELECT COUNT(user_group.*) FROM user, user_group, user_group_binding %s",
		whereCondition,
	)

	logger.Infof(ctx, "%v", query)
	var total int
	p.DB.Raw(query).Count(&total)
	if err := p.DB.Error; err != nil {
		logger.Warnf(ctx, "%v", query)
		logger.Warnf(ctx, "%+v", err)
		return nil, err
	}

	logger.Infof(ctx, "%v", query)

	// SELECT user_group.* FROM user, user_group, user_group_binding
	// {WHERE ...} {ORDER BY ...} {LIMIT ...}
	query = fmt.Sprintf(
		"SELECT user_group.* FROM user, user_group, user_group_binding %s %s %s",
		whereCondition, orderBy, limitOffset,
	)
	reply, err := p.listGroupsByQuery(ctx, query)
	if err != nil {
		logger.Warnf(ctx, "%v", query)
		logger.Warnf(ctx, "%+v", err)
		return nil, err
	}

	reply.Total = int32(total)
	return reply, nil
}

func (p *Database) listGroupsByQuery(ctx context.Context, query string) (*pbim.ListGroupsResponse, error) {
	var rows = []UserGroup{}
	err := p.DB.Raw(query).Scan(&rows).Error
	if err != nil {
		logger.Warnf(ctx, "%v", query)
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
