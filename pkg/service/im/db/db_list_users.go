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

func (p *Database) validateListUsersReq(req *pbim.ListUsersRequest) error {
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

	if !isValidGids(req.GroupId...) {
		return fmt.Errorf("invalid gid: %v", req.GroupId)
	}
	if !isValidUids(req.UserId...) {
		return fmt.Errorf("invalid uid: %v", req.UserId)
	}
	if !isValidNames(req.UserName...) {
		return fmt.Errorf("invalid name: %v", req.UserName)
	}
	if !isValidEmails(req.Email...) {
		return fmt.Errorf("invalid email: %v", req.Email)
	}
	if !isValidPhoneNumbers(req.PhoneNumber...) {
		return fmt.Errorf("invalid phone_number: %v", req.PhoneNumber)
	}
	if !isValidStatus(req.Status...) {
		return fmt.Errorf("invalid status: %v", req.Status)
	}

	return nil
}

func (p *Database) listUsers_no_gid(ctx context.Context, req *pbim.ListUsersRequest) (*pbim.ListUsersResponse, error) {
	if len(req.GroupId) > 0 {
		panic("should use listUsers_with_gid")
	}

	// WHERE name IN (name1,name2) AND name LIKE '%search_word%'
	var whereCondition = func() string {
		ss := genWhereCondition(
			map[string][]string{
				"group_id":     req.GroupId,
				"user_id":      req.UserId,
				"user_name":    req.UserName,
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
		if len(ss) > 0 {
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
	var total int
	p.DB.Raw(query).Count(&total)
	if err := p.DB.Error; err != nil {
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
	if len(req.GroupId) == 0 {
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
				"group_id":     req.GroupId,
				"user_id":      req.UserId,
				"user_name":    req.UserName,
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

	// SELECT COUNT(user.*) FROM user, user_group, user_group_binding {WHERE ...}
	query := fmt.Sprintf(
		"SELECT COUNT(user.*) FROM user, user_group, user_group_binding %s",
		whereCondition,
	)
	var total int
	p.DB.Raw(query).Count(&total)
	if err := p.DB.Error; err != nil {
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
	err := p.dbx.SelectContext(ctx, &rows, query)
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
