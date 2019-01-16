// Copyright 2019 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package db

import (
	"context"
	"fmt"

	"openpitrix.io/logger"
	"openpitrix.io/iam/pkg/pb/am"
	"openpitrix.io/iam/pkg/service/am/db_spec"
)

func (p *Database) DescribeActions(ctx context.Context, req *pbam.DescribeActionsRequest) (*pbam.ActionList, error) {
	// SELECT * FROM name
	var query = fmt.Sprintf(
		"SELECT * FROM %s;", db_spec.ActionTableName,
	)

	var rows = []db_spec.DBAction{}
	err := p.DB.SelectContext(ctx, &rows, query)
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
