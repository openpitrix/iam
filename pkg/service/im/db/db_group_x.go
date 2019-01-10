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
		db_spec.DBSpec.UserGroupBindingTableName,
	)

	_, err := p.DB.ExecContext(ctx, sql, genXid(), req.GetUid(), req.GetGid())
	if err != nil {
		logger.Warnf(ctx, "%v", sql)
		logger.Warnf(ctx, "%+v", err)
		return nil, err
	}

	return &pbim.Empty{}, nil
}

func (p *Database) LeaveGroup(ctx context.Context, req *pbim.LeaveGroupRequest) (*pbim.Empty, error) {
	logger.Infof(ctx, funcutil.CallerName(1))

	sql := fmt.Sprintf(
		`delete from %s where user_id=? AND group_id=?`,
		db_spec.DBSpec.UserGroupBindingTableName,
	)

	_, err := p.DB.ExecContext(ctx, sql, req.GetUid(), req.GetGid())
	if err != nil {
		logger.Warnf(ctx, "%v", sql)
		logger.Warnf(ctx, "%+v", err)
		return nil, err
	}

	return &pbim.Empty{}, nil
}
