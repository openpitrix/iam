// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package db

import (
	"context"
	"fmt"
	"strings"

	"openpitrix.io/iam/pkg/pb/im"
	"openpitrix.io/iam/pkg/service/im/db_spec"
	"openpitrix.io/logger"
)

func (p *Database) listUsers_no_gid(ctx context.Context, req *pbim.ListUsersRequest) (*pbim.ListUsersResponse, error) {
	if len(req.Gid) > 0 {
		panic("should use listUsers_with_gid")
	}

	// WHERE name IN (name1,name2) AND name LIKE '%search_word%'
	var whereCondition = func() string {
		ss := genWhereCondition(
			map[string][]string{
				"group_id":     req.Gid,
				"user_id":      req.Uid,
				"name":         req.Name,
				"email":        req.Email,
				"phone_number": req.PhoneNumber,
				"status":       req.Status,
			},
			[]string{
				"user_id",
				"user_name",
				"email",
				"phone_number",
				"description",
				"status",
			},
			req.SearchWord,
		)
		if len(ss) >= 0 {
			return "WHERE " + strings.Join(ss, " AND ")
		}
		return ""
	}()

	// ORDER BY column ASC|DESC;
	var orderBy = genOrderBy(req.SortKey, req.Reverse)

	// LIMIT 20 OFFSET 0
	var limitOffset = genLimitOffset(req.GetLimit(), req.GetOffset())

	// SELECT COUNT(*) FROM user {WHERE ...}
	query := fmt.Sprintf(
		"SELECT COUNT(*) FROM user %s",
		whereCondition,
	)
	total, err := p.getCountByQuery(ctx, query)
	if err != nil {
		logger.Warnf(ctx, "%v", query)
		logger.Warnf(ctx, "%+v", err)
		return nil, err
	}

	// SELECT * FROM user {WHERE ...} {ORDER BY ...} {LIMIT ...}
	query = fmt.Sprintf(
		"SELECT * FROM user %s %s %s",
		whereCondition, orderBy, limitOffset,
	)
	reply, err := p.listUsersByQuery(ctx, query)
	if err != nil {
		logger.Warnf(ctx, "%v", query)
		logger.Warnf(ctx, "%+v", err)
		return nil, err
	}

	reply.Total = int32(total)
	return reply, nil
}

func (p *Database) listUsers_with_gid(ctx context.Context, req *pbim.ListUsersRequest) (*pbim.ListUsersResponse, error) {
	if len(req.Gid) == 0 {
		panic("should use listUsers_no_gid")
	}

	// select user.* from
	//     user, user_group, user_group_binding
	// where
	//     user_group_binding.user_id=user.user_id and
	//     user_group_binding.user_id=user_group.group_id and
	//     ...

	// WHERE name IN (name1,name2) AND name LIKE '%search_word%'
	var whereCondition = func() string {
		ss := genWhereCondition(
			map[string][]string{
				"group_id":     req.Gid,
				"user_id":      req.Uid,
				"name":         req.Name,
				"email":        req.Email,
				"phone_number": req.PhoneNumber,
				"status":       req.Status,
			},
			[]string{
				"user_id",
				"user_name",
				"email",
				"phone_number",
				"description",
				"status",
			},
			req.SearchWord,
		)

		ss = append(
			[]string{
				"user_group_binding.user_id=user.user_id",
				"user_group_binding.user_id=user_group.group_id",
			},
			ss...,
		)

		return "WHERE " + strings.Join(ss, " AND ")
	}()

	// ORDER BY column ASC|DESC;
	var orderBy = genOrderBy(req.SortKey, req.Reverse)

	// LIMIT 20 OFFSET 0
	var limitOffset = genLimitOffset(req.GetLimit(), req.GetOffset())

	// SELECT COUNT(user.*) FROM user, user_group, user_group_binding {WHERE ...}
	query := fmt.Sprintf(
		"SELECT COUNT(user.*) FROM user, user_group, user_group_binding %s",
		whereCondition,
	)
	total, err := p.getCountByQuery(ctx, query)
	if err != nil {
		logger.Warnf(ctx, "%v", query)
		logger.Warnf(ctx, "%+v", err)
		return nil, err
	}

	// SELECT user.* FROM user, user_group, user_group_binding
	// {WHERE ...} {ORDER BY ...} {LIMIT ...}
	query = fmt.Sprintf(
		"SELECT user.* FROM user, user_group, user_group_binding %s %s %s",
		whereCondition, orderBy, limitOffset,
	)
	reply, err := p.listUsersByQuery(ctx, query)
	if err != nil {
		logger.Warnf(ctx, "%v", query)
		logger.Warnf(ctx, "%+v", err)
		return nil, err
	}

	reply.Total = int32(total)
	return reply, nil
}

func (p *Database) listUsersByQuery(ctx context.Context, query string) (*pbim.ListUsersResponse, error) {
	var rows = []db_spec.DBUser{}
	err := p.DB.SelectContext(ctx, &rows, query)
	if err != nil {
		logger.Warnf(ctx, "%v", query)
		logger.Warnf(ctx, "%+v", err)
		return nil, err
	}

	var sets []*pbim.User
	for _, v := range rows {
		v.Password = ""
		sets = append(sets, v.ToPB())
	}

	reply := &pbim.ListUsersResponse{
		User: sets,
	}

	return reply, nil
}
