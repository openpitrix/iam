// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package db

import (
	"context"
	"fmt"
	"sort"
	"strings"

	"openpitrix.io/iam/pkg/pb/im"
	"openpitrix.io/iam/pkg/service/im/db_spec"
	"openpitrix.io/logger"
)

func (p *Database) listUsers(ctx context.Context, req *pbim.ListUsersRequest) (*pbim.ListUsersResponse, error) {
	// WHERE name IN (name1,name2) AND name LIKE '%search_word%'
	var whereCondition = func() string {
		m := make(map[string]string)

		// name IN(name1,name2,...)
		var keyFileds = []struct {
			Name  string
			Value []string
		}{
			{Name: "gid", Value: req.Gid},
			{Name: "uid", Value: req.Uid},
			{Name: "name", Value: req.Name},
			{Name: "email", Value: req.Email},
			{Name: "phone_number", Value: req.PhoneNumber},
			{Name: "status", Value: req.Status},
		}
		for _, v := range keyFileds {
			if len(v.Value) > 0 {
				m[v.Name] = fmt.Sprintf(
					"%s IN(%s)", v.Name, strings.Join(req.Gid, ","),
				)
			}
		}

		// name LIKE '%search_word%'
		if req.SearchWord != "" {
			var searchWordFieldNames = []string{
				"user_id",
				"user_name",
				"email",
				"phone_number",
				"description",
				"status",
			}
			for _, name := range searchWordFieldNames {
				if _, exists := m[name]; !exists {
					m[name] = fmt.Sprintf(
						"%s LIKE '%%%s%%'", name, req.SearchWord,
					)
				}
			}
		}
		if len(m) == 0 {
			return "" // no where condition
		}

		var ss []string
		for _, v := range m {
			ss = append(ss, v)
		}
		sort.Strings(ss)

		// WHERE condition1 AND condition2 AND ...
		return "WHERE " + strings.Join(ss, " AND ")
	}()

	// ORDER BY column ASC|DESC;
	var orderBy = func() string {
		if req.SortKey == "" {
			return ""
		}
		if req.Reverse {
			return "ORDER BY " + req.SortKey + " DESC"
		} else {
			return "ORDER BY " + req.SortKey + " ASC"
		}
	}()

	// LIMIT 20 OFFSET 0
	var limitOffset = func() string {
		if offset, limit := req.GetOffset(), req.GetLimit(); offset > 0 || limit > 0 {
			return fmt.Sprintf("LIMIT %d OFFSET %d", limit, offset)
		} else {
			return fmt.Sprintf("LIMIT %d OFFSET %d", 20, 0)
		}
	}()

	// SELECT COUNT(*) FROM {name} {WHERE ...}
	query := fmt.Sprintf(
		"SELECT COUNT(*) FROM %s %s",
		"user", whereCondition,
	)
	total, err := p.getCountByQuery(ctx, query)
	if err != nil {
		logger.Warnf(ctx, "%v", query)
		logger.Warnf(ctx, "%+v", err)
		return nil, err
	}

	// SELECT * FROM {name} {WHERE ...} {ORDER BY ...} {LIMIT ...}
	query = fmt.Sprintf(
		"SELECT COUNT(*) FROM %s %s %s %s",
		"user", whereCondition,
		orderBy, limitOffset,
	)

	var rows = []db_spec.DBUser{}
	err = p.DB.SelectContext(ctx, &rows, query)
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
		User:  sets,
		Total: int32(total),
	}

	return reply, nil
}

func (p *Database) getCountByQuery(ctx context.Context, query string) (total int, err error) {
	rows, err := p.DB.QueryContext(ctx, query)
	if err != nil {
		logger.Warnf(ctx, "%v", query)
		logger.Warnf(ctx, "%+v", err)
		return 0, err
	}
	defer rows.Close()

	if rows.Next() {
		if err := rows.Scan(&total); err != nil {
			logger.Warnf(ctx, "%v", query)
			logger.Warnf(ctx, "%+v", err)
			return 0, err
		}
	}
	return
}
