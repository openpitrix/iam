// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package db

import (
	"context"
	"fmt"

	"openpitrix.io/iam/pkg/internal/funcutil"
	"openpitrix.io/iam/pkg/pb/im"
	"openpitrix.io/iam/pkg/service/im/db_spec"
	"openpitrix.io/logger"
)

func (p *Database) _ListUsers_all_count(ctx context.Context, req *pbim.Range) (total int, err error) {
	logger.Infof(ctx, funcutil.CallerName(1))

	var query = fmt.Sprintf("SELECT COUNT(*) FROM %s", db_spec.UserTableName)

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

	return total, nil
}

func (p *Database) _ListUsers_all(ctx context.Context, req *pbim.Range) (*pbim.ListUsersResponse, error) {
	logger.Infof(ctx, funcutil.CallerName(1))

	total, err := p._ListUsers_all_count(ctx, req)
	if err != nil {
		logger.Warnf(ctx, "%+v", err)
		return nil, err
	}

	var query = fmt.Sprintf("SELECT * FROM %s", db_spec.UserTableName)
	if offset, limit := req.GetOffset(), req.GetLimit(); offset > 0 || limit > 0 {
		query += fmt.Sprintf(" LIMIT %d OFFSET %d;", limit, offset)
	} else {
		query += fmt.Sprintf(" LIMIT %d OFFSET %d;", 20, 0)
	}

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

func (p *Database) _ListUsers_bySearchWord(ctx context.Context, req *pbim.Range) (*pbim.ListUsersResponse, error) {
	logger.Infof(ctx, funcutil.CallerName(1))

	var searchWord = req.GetSearchWord()

	if searchWord == "" {
		return p._ListUsers_all(ctx, req)
	}

	if !pkgSearchWordValid(searchWord) {
		return nil, fmt.Errorf("invalid search_word: %q", searchWord)
	}

	var searchWordFieldNames = pkgGetDBTableStringFieldNames(new(db_spec.DBUser))
	if len(searchWordFieldNames) == 0 {
		return p._ListUsers_all(ctx, req)
	}

	var (
		queryHeaer       = fmt.Sprintf("SELECT * FROM %s ", db_spec.UserTableName)
		queryCountHeader = fmt.Sprintf("SELECT COUNT(*) FROM %s ", db_spec.UserTableName)
		queryTail        string
	)

	for i, name := range searchWordFieldNames {
		if i == 0 {
			queryTail += " WHERE `" + name + "` LIKE '%" + searchWord + "%'"
		} else {
			queryTail += " OR `" + name + "` LIKE '%" + searchWord + "%'"
		}
	}

	if offset, limit := req.GetOffset(), req.GetLimit(); offset > 0 || limit > 0 {
		queryTail += fmt.Sprintf(" LIMIT %d OFFSET %d;", limit, offset)
	} else {
		queryTail += fmt.Sprintf(" LIMIT %d OFFSET %d;", 20, 0)
	}

	// total
	var total int
	{
		rows, err := p.DB.QueryContext(ctx, queryCountHeader+queryTail)
		if err != nil {
			logger.Warnf(ctx, "%v", queryCountHeader+queryTail)
			logger.Warnf(ctx, "%+v", err)
			return nil, err
		}
		defer rows.Close()

		if rows.Next() {
			if err := rows.Scan(&total); err != nil {
				logger.Warnf(ctx, "%v", queryCountHeader+queryTail)
				logger.Warnf(ctx, "%+v", err)
				return nil, err
			}
		}
	}

	var rows = []db_spec.DBUser{}
	err := p.DB.SelectContext(ctx, &rows, queryHeaer+queryTail)
	if err != nil {
		logger.Warnf(ctx, "%v", queryCountHeader+queryTail)
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
