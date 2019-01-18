// Copyright 2019 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package db

import (
	"context"

	"openpitrix.io/iam/pkg/internal/funcutil"
	pbam "openpitrix.io/iam/pkg/pb/am"
	"openpitrix.io/logger"
)

func (p *Database) DescribeActions(ctx context.Context, req *pbam.DescribeActionsRequest) (*pbam.ActionList, error) {
	logger.Infof(ctx, funcutil.CallerName(1))

	var query = sqlDescribeActionsBy_RoleId_Protal

	var rows = []DBAction{}
	err := p.DB.Raw(query, req.RoleId, req.Portal).Scan(&rows).Error
	if err != nil {
		logger.Warnf(ctx, "%v", query)
		logger.Warnf(ctx, "%+v", err)
		return nil, err
	}

	var sets []*pbam.Action
	for _, v := range rows {
		sets = append(sets, v.ToPB())
	}

	reply := &pbam.ActionList{
		Value: sets,
	}

	return reply, nil
}
