// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package db

import (
	"context"
	"fmt"

	"openpitrix.io/iam/pkg/pb/im"
	"openpitrix.io/iam/pkg/service/im/db_spec"
	"openpitrix.io/logger"
)

func (p *Database) GetUsersByGroupId(ctx context.Context, req *pbim.GroupId) (*pbim.UserList, error) {
	const sql = `
		select t1.* from
		user t1, user_group t2, user_group_binding t3 
		where t1.user_id=t3.user_id and t2.group_id=t3.group_id
		where t2.group_id=?
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
	var user db_spec.DBUser
	err := p.DB.Get(&user, "select * from user where user_id=?", req.GetUid())
	if err != nil {
		logger.Warnf(ctx, "uid = %s, err = %+v", req.GetUid(), err)
		return nil, err
	}

	if user.Password == req.GetPassword() {
		return &pbim.Bool{Value: true}, nil
	} else {
		return &pbim.Bool{Value: false}, nil
	}
}
func (p *Database) ModifyPassword(ctx context.Context, req *pbim.Password) (*pbim.Empty, error) {
	sql := fmt.Sprintf(
		`update %s set password="%s" where user_id="%s"`,
		db_spec.DBSpec.UserTableName,
		req.GetUid(), req.GetPassword(),
	)

	_, err := p.DB.ExecContext(ctx, sql)
	if err != nil {
		logger.Warnf(ctx, "%v", sql)
		logger.Warnf(ctx, "%+v", err)
		return nil, err
	}

	return &pbim.Empty{}, nil
}
