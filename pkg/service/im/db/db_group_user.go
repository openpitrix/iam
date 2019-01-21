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

func (p *Database) JoinGroup(ctx context.Context, req *pbim.JoinGroupRequest) (*pbim.Empty, error) {
	logger.Infof(ctx, funcutil.CallerName(1))

	if len(req.GroupId) == 1 && strings.Contains(req.GroupId[0], ",") {
		req.GroupId = strings.Split(req.GroupId[0], ",")
	}
	if len(req.UserId) == 1 && strings.Contains(req.UserId[0], ",") {
		req.UserId = strings.Split(req.UserId[0], ",")
	}

	if len(req.UserId) == 0 || len(req.GroupId) == 0 {
		err := status.Errorf(codes.InvalidArgument, "empty uid or gid")
		logger.Warnf(ctx, "%+v", err)
		return nil, err
	}
	if !(len(req.UserId) == 1 || len(req.GroupId) == 1 || len(req.UserId) == len(req.GroupId)) {
		err := status.Errorf(codes.InvalidArgument,
			"uid and gid donot math: gid = %v, uid = %v",
			req.GroupId, req.UserId,
		)
		logger.Warnf(ctx, "%+v", err)
		return nil, err
	}

	tx := p.DB.Begin()

	switch {
	case len(req.UserId) == len(req.GroupId):
		for i := 0; i < len(req.GroupId); i++ {
			xid := genId("xid-", 12)
			uid := req.UserId[i]
			gid := req.GroupId[i]

			tx.Exec(
				`INSERT INTO user_group_binding (id, user_id, group_id) VALUES (?,?,?)`,
				xid, uid, gid,
			)
			if err := tx.Error; err != nil {
				tx.Rollback()
				return nil, err
			}
		}
	case len(req.UserId) == 1:
		for i := 0; i < len(req.GroupId); i++ {
			xid := genId("xid-", 12)
			gid := req.GroupId[i]
			uid := req.UserId[0]

			tx.Exec(
				`INSERT INTO user_group_binding (id, user_id, group_id) VALUES (?,?,?)`,
				xid, uid, gid,
			)
			if err := tx.Error; err != nil {
				tx.Rollback()
				return nil, err
			}
		}
	case len(req.GroupId) == 1:
		for i := 0; i < len(req.UserId); i++ {
			xid := genId("xid-", 12)
			gid := req.GroupId[0]
			uid := req.UserId[i]

			tx.Exec(
				`INSERT INTO user_group_binding (id, user_id, group_id) VALUES (?,?,?)`,
				xid, uid, gid,
			)
			if err := tx.Error; err != nil {
				tx.Rollback()
				return nil, err
			}
		}
	}

	if err := tx.Commit().Error; err != nil {
		logger.Warnf(ctx, "%+v", err)
		return nil, err
	}

	return &pbim.Empty{}, nil
}

func (p *Database) LeaveGroup(ctx context.Context, req *pbim.LeaveGroupRequest) (*pbim.Empty, error) {
	logger.Infof(ctx, funcutil.CallerName(1))

	if len(req.GroupId) == 1 && strings.Contains(req.GroupId[0], ",") {
		req.GroupId = strings.Split(req.GroupId[0], ",")
	}
	if len(req.UserId) == 1 && strings.Contains(req.UserId[0], ",") {
		req.UserId = strings.Split(req.UserId[0], ",")
	}

	if len(req.UserId) == 0 || len(req.GroupId) == 0 {
		err := status.Errorf(codes.InvalidArgument, "empty uid or gid")
		logger.Warnf(ctx, "%+v", err)
		return nil, err
	}
	if !(len(req.UserId) == 1 || len(req.GroupId) == 1 || len(req.UserId) == len(req.GroupId)) {
		err := status.Errorf(codes.InvalidArgument,
			"uid and gid donot math: gid = %v, uid = %v",
			req.GroupId, req.UserId,
		)
		logger.Warnf(ctx, "%+v", err)
		return nil, err
	}

	tx := p.DB.Begin()

	switch {
	case len(req.UserId) == len(req.GroupId):
		for i := 0; i < len(req.GroupId); i++ {
			gid := req.GroupId[i]
			uid := req.UserId[i]

			tx.Exec(
				`delete from user_group_binding where user_id=? and group_id=?`,
				uid, gid,
			)
			if err := tx.Error; err != nil {
				tx.Rollback()
				return nil, err
			}
		}
	case len(req.UserId) == 1:
		for i := 0; i < len(req.GroupId); i++ {
			gid := req.GroupId[i]
			uid := req.UserId[0]

			tx.Exec(
				`delete from user_group_binding where user_id=? and group_id=?`,
				uid, gid,
			)
			if err := tx.Error; err != nil {
				tx.Rollback()
				return nil, err
			}
		}
	case len(req.GroupId) == 1:
		for i := 0; i < len(req.UserId); i++ {
			gid := req.GroupId[0]
			uid := req.UserId[i]

			tx.Exec(
				`delete from user_group_binding where user_id=? and group_id=?`,
				uid, gid,
			)
			if err := tx.Error; err != nil {
				tx.Rollback()
				return nil, err
			}
		}
	}

	if err := tx.Commit().Error; err != nil {
		logger.Warnf(ctx, "%+v", err)
		return nil, err
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
			user_group_binding.group_id=user_group.group_id and
			user.user_id=?
	`
	var rows []UserGroup
	p.DB.Raw(sql, sql, req.UserId).Scan(&rows)
	if err := p.DB.Error; err != nil {
		logger.Warnf(ctx, "%v", sql)
		logger.Warnf(ctx, "%+v", err)
		return nil, err
	}

	logger.Infof(ctx, "rows: %v", rows)

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
	var rows []User
	p.DB.Raw(sql, sql, req.GroupId).Scan(&rows)
	if err := p.DB.Error; err != nil {
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
