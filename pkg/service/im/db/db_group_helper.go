// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package db

import (
	"context"
	"sort"

	"github.com/chai2010/template"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"openpitrix.io/iam/pkg/internal/strutil"
	"openpitrix.io/iam/pkg/pb/im"
	"openpitrix.io/iam/pkg/service/im/db_spec"
	"openpitrix.io/logger"
)

func (p *Database) getUserGroupBindList(userId, groupId []string) (
	s []db_spec.UserGroupBinding,
) {
	if len(userId) == 0 || len(groupId) == 0 {
		return nil
	}
	if !(len(userId) == 1 || len(groupId) == 1 || len(userId) == len(groupId)) {
		return nil
	}

	if len(userId) == len(groupId) {
		for i := 0; i < len(userId); i++ {
			s = append(s, db_spec.UserGroupBinding{
				UserId:  userId[i],
				GroupId: groupId[i],
			})
		}
		return
	}
	if len(userId) == 1 {
		for i := 0; i < len(groupId); i++ {
			s = append(s, db_spec.UserGroupBinding{
				UserId:  userId[0],
				GroupId: groupId[i],
			})
		}
		return
	}
	if len(groupId) == 1 {
		for i := 0; i < len(userId); i++ {
			s = append(s, db_spec.UserGroupBinding{
				UserId:  userId[i],
				GroupId: groupId[0],
			})
		}
		return
	}

	return
}

// UserGroupBinding

func (p *Database) getAllSubGroupIds(ctx context.Context, req *pbim.GroupIdList) ([]string, error) {
	const sqlTmpl = `
		SELECT * FROM user_group WHERE 1=0
			{{range $i, $v := .GroupId}}
				OR group_path LINK '%{{$v}}%'
				OR group_id='{{$v}}'
			{{end}}
	`
	var query, err = template.Render(sqlTmpl, req)
	if err != nil {
		err := status.Errorf(codes.Internal, "%v", err)
		logger.Errorf(ctx, "%+v", err)
		return nil, err
	}

	query = strutil.SimplifyString(query)
	logger.Infof(ctx, "query: %s", query)

	var rows []db_spec.UserGroup
	p.DB.Raw(query).Find(&rows)
	if err := p.DB.Error; err != nil {
		logger.Warnf(ctx, "%+v", err)
		return nil, err
	}

	var allGroupId []string
	for _, v := range rows {
		allGroupId = append(allGroupId, v.GroupId)
	}

	sort.Strings(allGroupId)
	return allGroupId, nil
}
