// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package db

import (
	"context"
	"fmt"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"openpitrix.io/iam/pkg/internal/funcutil"
	"openpitrix.io/iam/pkg/pb/im"
	"openpitrix.io/iam/pkg/service/im/db_spec"
	"openpitrix.io/logger"
)

func (p *Database) JoinGroup(ctx context.Context, req *pbim.JoinGroupRequest) (*pbim.Empty, error) {
	logger.Infof(ctx, funcutil.CallerName(1))

	if len(req.Uid) == 0 || len(req.Gid) == 0 {
		err := status.Errorf(codes.InvalidArgument, "empty uid or gid")
		logger.Warnf(ctx, "%+v", err)
		return nil, err
	}
	if !(len(req.Uid) == 1 || len(req.Gid) == 1 || len(req.Uid) == len(req.Gid)) {
		err := status.Errorf(codes.InvalidArgument,
			"uid and gid donot math: gid = %v, uid = %v",
			req.Gid, req.Uid,
		)
		logger.Warnf(ctx, "%+v", err)
		return nil, err
	}

	sql := fmt.Sprintf(
		`insert into %s (id, user_id, group_id) values(?,?,?)`,
		db_spec.UserGroupBindingTableName,
	)

	switch {
	case len(req.Uid) == len(req.Gid):
		for i := 0; i < len(req.Gid); i++ {
			xid := genXid()
			gid := req.Gid[i]
			uid := req.Uid[i]

			_, err := p.DB.ExecContext(ctx, sql, xid, uid, gid)
			if err != nil {
				logger.Warnf(ctx, "%v", sql)
				logger.Warnf(ctx, "%v, %v, %v", xid, uid, gid)
				logger.Warnf(ctx, "%+v", err)
				return nil, err
			}
		}
	case len(req.Uid) == 1:
		for i := 0; i < len(req.Gid); i++ {
			xid := genXid()
			gid := req.Gid[i]
			uid := req.Uid[0]

			_, err := p.DB.ExecContext(ctx, sql, xid, uid, gid)
			if err != nil {
				logger.Warnf(ctx, "%v", sql)
				logger.Warnf(ctx, "%v, %v, %v", xid, uid, gid)
				logger.Warnf(ctx, "%+v", err)
				return nil, err
			}
		}
	case len(req.Gid) == 1:
		for i := 0; i < len(req.Uid); i++ {
			xid := genXid()
			gid := req.Gid[0]
			uid := req.Uid[i]

			_, err := p.DB.ExecContext(ctx, sql, xid, uid, gid)
			if err != nil {
				logger.Warnf(ctx, "%v", sql)
				logger.Warnf(ctx, "%v, %v, %v", xid, uid, gid)
				logger.Warnf(ctx, "%+v", err)
				return nil, err
			}
		}
	}

	return &pbim.Empty{}, nil
}

func (p *Database) LeaveGroup(ctx context.Context, req *pbim.LeaveGroupRequest) (*pbim.Empty, error) {
	logger.Infof(ctx, funcutil.CallerName(1))

	if len(req.Uid) == 0 || len(req.Gid) == 0 {
		err := status.Errorf(codes.InvalidArgument, "empty uid or gid")
		logger.Warnf(ctx, "%+v", err)
		return nil, err
	}
	if !(len(req.Uid) == 1 || len(req.Gid) == 1 || len(req.Uid) == len(req.Gid)) {
		err := status.Errorf(codes.InvalidArgument,
			"uid and gid donot math: gid = %v, uid = %v",
			req.Gid, req.Uid,
		)
		logger.Warnf(ctx, "%+v", err)
		return nil, err
	}

	sql := fmt.Sprintf(
		`delete from %s where user_id=? AND group_id=?`,
		db_spec.UserGroupBindingTableName,
	)

	switch {
	case len(req.Uid) == len(req.Gid):
		for i := 0; i < len(req.Gid); i++ {
			xid := genXid()
			gid := req.Gid[i]
			uid := req.Uid[i]

			_, err := p.DB.ExecContext(ctx, sql, xid, uid, gid)
			if err != nil {
				logger.Warnf(ctx, "%v", sql)
				logger.Warnf(ctx, "%v, %v, %v", xid, uid, gid)
				logger.Warnf(ctx, "%+v", err)
				return nil, err
			}
		}
	case len(req.Uid) == 1:
		for i := 0; i < len(req.Gid); i++ {
			xid := genXid()
			gid := req.Gid[i]
			uid := req.Uid[0]

			_, err := p.DB.ExecContext(ctx, sql, xid, uid, gid)
			if err != nil {
				logger.Warnf(ctx, "%v", sql)
				logger.Warnf(ctx, "%v, %v, %v", xid, uid, gid)
				logger.Warnf(ctx, "%+v", err)
				return nil, err
			}
		}
	case len(req.Gid) == 1:
		for i := 0; i < len(req.Uid); i++ {
			xid := genXid()
			gid := req.Gid[0]
			uid := req.Uid[i]

			_, err := p.DB.ExecContext(ctx, sql, xid, uid, gid)
			if err != nil {
				logger.Warnf(ctx, "%v", sql)
				logger.Warnf(ctx, "%v, %v, %v", xid, uid, gid)
				logger.Warnf(ctx, "%+v", err)
				return nil, err
			}
		}
	}

	return &pbim.Empty{}, nil
}

func (p *Database) GetGroupsByUserId(ctx context.Context, req *pbim.UserId) (*pbim.GroupList, error) {
	logger.Infof(ctx, funcutil.CallerName(1))

	const sql = `
		select user_group.* from
			user, user_group, user_group_binding 
		where
			user_group_binding.user_id=user.user_id and
			user_group_binding.user_id=user_group.group_id and
			user.user_id=?
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
		select user.* from
			user, user_group, user_group_binding 
		where
			user_group_binding.user_id=user.user_id and
			user_group_binding.user_id=user_group.group_id and
			user_group.group_id=?
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
