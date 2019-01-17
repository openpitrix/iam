// Copyright 2019 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package db

import (
	"context"
	"fmt"
	"strings"

	"openpitrix.io/iam/pkg/internal/funcutil"
	pbam "openpitrix.io/iam/pkg/pb/am"
	"openpitrix.io/iam/pkg/service/am/db_spec"
	"openpitrix.io/logger"
)

func (p *Database) DescribeUsersWithRole(ctx context.Context, req *pbam.DescribeUsersWithRoleRequest) (*pbam.DescribeUsersWithRoleResponse, error) {
	logger.Infof(ctx, funcutil.CallerName(1))

	// fix repeated fileds
	if len(req.RoleId) == 1 && strings.Contains(req.RoleId[0], ",") {
		req.RoleId = strings.Split(req.RoleId[0], ",")
	}
	if len(req.UserId) == 1 && strings.Contains(req.UserId[0], ",") {
		req.UserId = strings.Split(req.UserId[0], ",")
	}

	// WHERE name IN (name1,name2) AND name LIKE '%search_word%'
	var whereCondition = func() string {
		ss := genWhereCondition(
			map[string][]string{
				"role_id": req.RoleId,
				"user_id": req.UserId,
			},
			nil, "",
		)

		ss = append(
			[]string{
				"user_role_binding.user_id=user.user_id",
				"user_role_binding.group_id=user_group.group_id",
			},
			ss...,
		)

		//
		if len(ss) > 0 {
			return "WHERE " + strings.Join(ss, " AND ")
		}
		return ""
	}()

	// LIMIT 20 OFFSET 0
	var limitOffset = genLimitOffset(req.GetLimit(), req.GetOffset())

	// SELECT COUNT(user.*) FROM user, role, user_role_binding {WHERE ...}
	query := fmt.Sprintf(
		"SELECT COUNT(user.*) FROM user, role, user_role_binding %s",
		whereCondition,
	)
	total, err := p.getCountByQuery(ctx, query)
	if err != nil {
		logger.Warnf(ctx, "%v", query)
		logger.Warnf(ctx, "%+v", err)
		return nil, err
	}

	query = fmt.Sprintf(
		"SELECT COUNT(user.*) FROM user, role, user_role_binding %s %s",
		whereCondition, limitOffset,
	)

	var rows = []db_spec.DBUserWithRole{}
	err = p.DB.SelectContext(ctx, &rows, query)
	if err != nil {
		logger.Warnf(ctx, "%v", query)
		logger.Warnf(ctx, "%+v", err)
		return nil, err
	}

	var sets []*pbam.UserWithRole
	for _, v := range rows {
		sets = append(sets, v.ToPB())
	}

	reply := &pbam.DescribeUsersWithRoleResponse{
		User:  sets,
		Total: int32(total),
	}

	return reply, nil
}
