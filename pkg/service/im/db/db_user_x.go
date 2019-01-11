// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package db

import (
	"context"
	"fmt"

	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"openpitrix.io/iam/pkg/internal/funcutil"
	"openpitrix.io/iam/pkg/pb/im"
	"openpitrix.io/iam/pkg/service/im/db_spec"
	"openpitrix.io/logger"
)

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

func (p *Database) ComparePassword(ctx context.Context, req *pbim.Password) (*pbim.Bool, error) {
	logger.Infof(ctx, funcutil.CallerName(1))

	var user db_spec.DBUser
	err := p.DB.Get(&user, "select * from user where user_id=?", req.GetUid())
	if err != nil {
		logger.Warnf(ctx, "uid = %s, err = %+v", req.GetUid(), err)
		return nil, err
	}

	err = bcrypt.CompareHashAndPassword(
		[]byte(user.Password), []byte(req.GetPassword()),
	)
	if err != nil {
		return &pbim.Bool{Value: false}, nil
	}

	// OK
	return &pbim.Bool{Value: true}, nil
}
func (p *Database) ModifyPassword(ctx context.Context, req *pbim.Password) (*pbim.Empty, error) {
	logger.Infof(ctx, funcutil.CallerName(1))

	var user db_spec.DBUser
	err := p.DB.Get(&user, "select * from user where user_id=?", req.GetUid())
	if err != nil {
		logger.Warnf(ctx, "uid = %s, err = %+v", req.GetUid(), err)
		return nil, err
	}

	hashedPass, err := bcrypt.GenerateFromPassword(
		[]byte(req.GetPassword()), bcrypt.DefaultCost,
	)
	if err != nil {
		err := status.Errorf(codes.Internal, "bcrypt failed")
		logger.Warnf(ctx, "%+v", err)
		return nil, err
	}

	sql := fmt.Sprintf(
		`update %s set password="%s" where user_id="%s"`,
		db_spec.UserTableName,
		string(hashedPass),
		req.GetUid(),
	)

	_, err = p.DB.ExecContext(ctx, sql)
	if err != nil {
		logger.Warnf(ctx, "%v", sql)
		logger.Warnf(ctx, "%+v", err)
		return nil, err
	}

	return &pbim.Empty{}, nil
}
