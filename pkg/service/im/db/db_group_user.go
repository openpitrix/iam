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

func (p *Database) JoinGroup(ctx context.Context, req *pbim.JoinGroupRequest) (*pbim.Empty, error) {
	logger.Infof(ctx, funcutil.CallerName(1))

	sql := fmt.Sprintf(
		`insert into %s (id, user_id, group_id) values(?,?,?)`,
		db_spec.UserGroupBindingTableName,
	)

	xid := genXid()
	_, err := p.DB.ExecContext(ctx, sql, xid, req.GetUid(), req.GetGid())
	if err != nil {
		logger.Warnf(ctx, "%v", sql)
		logger.Warnf(ctx, "%v, %v, %v", xid, req.GetUid(), req.GetGid())
		logger.Warnf(ctx, "%+v", err)
		return nil, err
	}

	return &pbim.Empty{}, nil
}

func (p *Database) LeaveGroup(ctx context.Context, req *pbim.LeaveGroupRequest) (*pbim.Empty, error) {
	logger.Infof(ctx, funcutil.CallerName(1))

	sql := fmt.Sprintf(
		`delete from %s where user_id=? AND group_id=?`,
		db_spec.UserGroupBindingTableName,
	)

	_, err := p.DB.ExecContext(ctx, sql, req.GetUid(), req.GetGid())
	if err != nil {
		logger.Warnf(ctx, "%v", sql)
		logger.Warnf(ctx, "%+v", err)
		return nil, err
	}

	return &pbim.Empty{}, nil
}

func (p *Database) GetGroupsByUserId(ctx context.Context, req *pbim.UserId) (*pbim.GroupList, error) {
	logger.Infof(ctx, funcutil.CallerName(1))

	const sql = `
		select t2.* from
		user t1, user_group t2, user_group_binding t3 
		where t1.user_id=t3.user_id and t2.group_id=t3.group_id
		and t1.user_id=?
	`
	var rows []db_spec.DBGroup
	err := p.DB.Select(&rows, sql, req.GetUid())
	if err != nil {
		logger.Warnf(ctx, "%v", sql)
		logger.Warnf(ctx, "%+v", err)
		return nil, err
	}

	var sets []*pbim.Group
	for _, v := range rows {
		sets = append(sets, v.ToPB())
	}

	reply := &pbim.GroupList{
		Value: sets,
	}
	return reply, nil
}

func (p *Database) GetUsersByGroupId(ctx context.Context, req *pbim.GroupId) (*pbim.UserList, error) {
	logger.Infof(ctx, funcutil.CallerName(1))

	const sql = `
		select t1.* from
		user t1, user_group t2, user_group_binding t3 
		where t1.user_id=t3.user_id and t2.group_id=t3.group_id
		and t2.group_id=?
	`
	var rows []db_spec.DBUser
	err := p.DB.Select(&rows, sql, req.GetGid())
	if err != nil {
		logger.Warnf(ctx, "%v", sql)
		logger.Warnf(ctx, "%+v", err)
		return nil, err
	}

	var sets []*pbim.User
	for _, v := range rows {
		sets = append(sets, v.ToPB())
	}

	reply := &pbim.UserList{
		Value: sets,
	}
	return reply, nil
}
