// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package db

import (
	"context"
	"strings"

	"github.com/chai2010/template"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	idpkg "openpitrix.io/iam/pkg/id"
	"openpitrix.io/iam/pkg/internal/funcutil"
	"openpitrix.io/iam/pkg/internal/strutil"
	"openpitrix.io/iam/pkg/pb/im"
	"openpitrix.io/iam/pkg/service/im/db_spec"
	"openpitrix.io/iam/pkg/validator"
	"openpitrix.io/logger"
)

func (p *Database) CreateGroup(ctx context.Context, req *pbim.Group) (*pbim.Group, error) {
	logger.Infof(ctx, funcutil.CallerName(1))

	// must generate new id
	req.GroupId = idpkg.GenId("gid-")

	if req.GroupPath == "" {
		if req.ParentGroupId == "" {
			req.GroupPath = req.GroupId
		} else {
			req.ParentGroupId = strutil.SimplifyString(req.ParentGroupId)

			// get parent group_path
			parentGroup, err := p.GetGroup(ctx, &pbim.GroupId{
				GroupId: req.ParentGroupId,
			})
			if err != nil {
				err = status.Errorf(codes.InvalidArgument, "get parent info failed: %v", err)
				logger.Warnf(ctx, "%+v", err)
				return nil, err
			}

			req.GroupPath = parentGroup.GroupPath + "." + req.GroupId
		}
	}

	var dbGroup = db_spec.NewUserGroupFromPB(req).AdjustForCreate()
	if err := dbGroup.IsValidForCreate(); err != nil {
		err = status.Errorf(codes.InvalidArgument, "%v", err)
		logger.Warnf(ctx, "%+v", err)
		return nil, err
	}

	var all_parent_group_id []string
	// ParentGroupId must be in GroupPath, skip it
	for _, id := range strings.Split(dbGroup.GroupPath, ".") {
		if id != dbGroup.GroupId {
			all_parent_group_id = append(all_parent_group_id, id)
		}
	}

	// check all parent group_id exists
	if len(all_parent_group_id) > 0 {
		var query, err = template.Render(`
		SELECT COUNT(*) FROM user_group WHERE 1=0
			{{range $i, $v := .}}
				OR group_id='{{$v}}'
			{{end}}
		`, all_parent_group_id,
		)
		if err != nil {
			logger.Warnf(ctx, "%+v", err)
			return nil, err
		}

		query = strutil.SimplifyString(query)
		logger.Infof(ctx, "query: %s", query)

		var total int
		p.DB.Raw(query).Count(&total)
		if err := p.DB.Error; err != nil {
			logger.Warnf(ctx, "%+v", err)
			return nil, err
		}
		if total != len(all_parent_group_id) {
			err := status.Errorf(codes.InvalidArgument,
				"some group_id in all_parent_group_id(%q) donot exists",
				dbGroup.GroupPath,
			)
			logger.Warnf(ctx, "%+v", err)
			return nil, err
		}
	}

	// create new record
	if err := p.DB.Create(dbGroup).Error; err != nil {
		logger.Warnf(ctx, "%+v, %v", err, dbGroup)
		return nil, err
	}

	// get again
	return p.GetGroup(ctx, &pbim.GroupId{GroupId: req.GroupId})
}

func (p *Database) DeleteGroups(ctx context.Context, req *pbim.GroupIdList) (*pbim.Empty, error) {
	logger.Infof(ctx, funcutil.CallerName(1))

	if len(req.GroupId) == 0 {
		err := status.Errorf(codes.InvalidArgument, "empty GroupId")
		logger.Warnf(ctx, "%+v", err)
		return nil, err
	}
	if !validator.IsValidId(req.GroupId...) {
		err := status.Errorf(codes.InvalidArgument, "invalid GroupId: %v", req.GroupId)
		logger.Warnf(ctx, "%+v", err)
		return nil, err
	}

	// 1. get group and sub group id list
	allGroupId, err := p.getAllSubGroupIds(ctx, req)
	if err != nil {
		logger.Warnf(ctx, "%+v", err)
		return nil, err
	}

	// 2. delete user_group & user_group_binding
	tx := p.DB.Begin()
	{
		if err := tx.Delete(db_spec.UserGroup{}, "group_id IN(?)", allGroupId).Error; err != nil {
			tx.Rollback()
			return nil, err
		}
		if err := tx.Delete(db_spec.UserGroupBinding{}, "group_id IN(?)", allGroupId).Error; err != nil {
			tx.Rollback()
			return nil, err
		}
	}

	if err := tx.Commit().Error; err != nil {
		logger.Warnf(ctx, "%+v", err)
		return nil, err
	}

	reply := &pbim.Empty{}
	return reply, nil
}

func (p *Database) ModifyGroup(ctx context.Context, req *pbim.Group) (*pbim.Group, error) {
	logger.Infof(ctx, funcutil.CallerName(1))

	var dbGroup = db_spec.NewUserGroupFromPB(req).AdjustForUpdate()
	if err := dbGroup.IsValidForUpdate(); err != nil {
		err = status.Errorf(codes.InvalidArgument, "%v", err)
		logger.Warnf(ctx, "%+v", err)
		return nil, err
	}

	if err := p.DB.Model(dbGroup).Updates(dbGroup).Error; err != nil {
		logger.Warnf(ctx, "%+v", err)
		return nil, err
	}

	return p.GetGroup(ctx, &pbim.GroupId{GroupId: req.GroupId})
}

func (p *Database) GetGroup(ctx context.Context, req *pbim.GroupId) (*pbim.Group, error) {
	logger.Infof(ctx, funcutil.CallerName(1))

	var v = db_spec.UserGroup{GroupId: req.GroupId}
	if err := p.DB.Model(db_spec.UserGroup{}).Take(&v).Error; err != nil {
		logger.Warnf(ctx, "%+v", err)
		return nil, err
	}

	return v.ToPB(), nil
}

func (p *Database) ListGroups(ctx context.Context, req *pbim.ListGroupsRequest) (*pbim.ListGroupsResponse, error) {
	logger.Infof(ctx, funcutil.CallerName(1))

	// fix repeated fileds
	if len(req.UserId) == 1 && strings.Contains(req.UserId[0], ",") {
		req.UserId = strings.Split(req.UserId[0], ",")
	}
	if len(req.GroupId) == 1 && strings.Contains(req.GroupId[0], ",") {
		req.GroupId = strings.Split(req.GroupId[0], ",")
	}
	if len(req.GroupName) == 1 && strings.Contains(req.GroupName[0], ",") {
		req.GroupName = strings.Split(req.GroupName[0], ",")
	}
	if len(req.Status) == 1 && strings.Contains(req.Status[0], ",") {
		req.Status = strings.Split(req.Status[0], ",")
	}

	req.UserId = strutil.SimplifyStringList(req.UserId)
	req.GroupId = strutil.SimplifyStringList(req.GroupId)
	req.GroupName = strutil.SimplifyStringList(req.GroupName)
	req.Status = strutil.SimplifyStringList(req.Status)

	// limit & offset
	if req.Limit == 0 && req.Offset == 0 {
		req.Limit = 20
		req.Offset = 0
	}
	if req.Limit < 0 {
		req.Limit = 0
	}
	if req.Limit > 200 {
		req.Limit = 200
	}
	if req.Offset < 0 {
		req.Offset = 0
	}

	if len(req.GroupId) > 0 {
		if !validator.IsValidId(req.GroupId...) {
			err := status.Errorf(codes.InvalidArgument, "invalid gid: %v", req.GroupId)
			logger.Warnf(ctx, "%+v", err)
			return nil, err
		}
	}
	if len(req.UserId) > 0 {
		if !validator.IsValidId(req.UserId...) {
			err := status.Errorf(codes.InvalidArgument, "invalid uid: %v", req.UserId)
			logger.Warnf(ctx, "%+v", err)
			return nil, err
		}
	}
	if req.SearchWord != "" {
		if !validator.IsValidSearchWord(req.SearchWord) {
			err := status.Errorf(codes.InvalidArgument, "invalid search_word: %v", req.SearchWord)
			logger.Warnf(ctx, "%+v", err)
			return nil, err
		}
	}
	if req.SortKey != "" {
		if !validator.IsValidSortKey(req.SortKey) {
			err := status.Errorf(codes.InvalidArgument, "invalid sort_key: %v", req.SortKey)
			logger.Warnf(ctx, "%+v", err)
			return nil, err
		}
		if !((*db_spec.UserGroup)(nil)).IsValidSortKey(req.SortKey) {
			err := status.Errorf(codes.InvalidArgument, "invalid sort_key: %v", req.SortKey)
			logger.Warnf(ctx, "%+v", err)
			return nil, err
		}
	}

	// 1. get group and sub group id list
	if len(req.GroupId) > 0 {
		allGroupId, err := p.getAllSubGroupIds(ctx, &pbim.GroupIdList{
			GroupId: req.GroupId,
		})
		if err != nil {
			logger.Warnf(ctx, "%+v", err)
			return nil, err
		}

		req.GroupId = allGroupId
	}

	const sqlTmpl = `
		{{if not .UserId}}
			select distinct {{if IsCountMode}}COUNT(*){{else}}*{{end}} from user_group where 1=1
				{{if .GroupId}}
					and group_id in (
						{{range $i, $v := .GroupId}}
							{{if eq $i 0}} '{{$v}}' {{else}} ,'{{$v}}' {{end}}
						{{end}}
					)
				{{end}}
				{{if .GroupName}}
					and group_name in (
						{{range $i, $v := .GroupName}}
							{{if eq $i 0}} '{{$v}}' {{else}} ,'{{$v}}' {{end}}
						{{end}}
					)
				{{end}}
				{{if .Status}}
					and status in (
						{{range $i, $v := .Status}}
							{{if eq $i 0}} '{{$v}}' {{else}} ,'{{$v}}' {{end}}
						{{end}}
					)
				{{end}}
				{{if .SearchWord}}
					and (1=0
						OR description LIKE '%{{.SearchWord}}%'
						{{if not .GroupId}}
							OR group_id LIKE '%{{.SearchWord}}%'
						{{end}}
						{{if not .GroupName}}
							OR group_name LIKE '%{{.SearchWord}}%'
						{{end}}
						{{if not .Status}}
							OR status LIKE '%{{.SearchWord}}%'
						{{end}}
					)
				{{end}}
				{{if not IsCountMode}}
					{{if .SortKey}}
						order by {{.SortKey}} {{if .Reverse}} desc {{end}}
					{{end}}
					limit {{.Limit}} offset {{.Offset}}
				{{end}}
		{{else}}
			select distinct {{if IsCountMode}}COUNT(*){{else}}user_group.*{{end}} from
				user, user_group, user_group_binding
			where 1=1
				and user_group_binding.user_id=user.user_id
				and user_group_binding.group_id=user_group.group_id

				{{if .UserId}}
					and user.user_id in (
						{{range $i, $v := .UserId}}
							{{if eq $i 0}} '{{$v}}' {{else}} ,'{{$v}}' {{end}}
						{{end}}
					)
				{{end}}
				{{if .GroupId}}
					and user_group.group_id in (
						{{range $i, $v := .GroupId}}
							{{if eq $i 0}} '{{$v}}' {{else}} ,'{{$v}}' {{end}}
						{{end}}
					)
				{{end}}
				{{if .GroupName}}
					and user_group.group_name in (
						{{range $i, $v := .GroupName}}
							{{if eq $i 0}} '{{$v}}' {{else}} ,'{{$v}}' {{end}}
						{{end}}
					)
				{{end}}
				{{if .Status}}
					and user_group.status in (
						{{range $i, $v := .Status}}
							{{if eq $i 0}} '{{$v}}' {{else}} ,'{{$v}}' {{end}}
						{{end}}
					)
				{{end}}
				{{if .SearchWord}}
					and (1=0
						OR user_group.description LIKE '%{{.SearchWord}}%'
						{{if not .GroupId}}
							OR user_group.group_id LIKE '%{{.SearchWord}}%'
						{{end}}
						{{if not .GroupName}}
							OR user_group.group_name LIKE '%{{.SearchWord}}%'
						{{end}}
						{{if not .Status}}
							OR user_group.status LIKE '%{{.SearchWord}}%'
						{{end}}
					)
				{{end}}
				{{if not IsCountMode}}
					{{if .SortKey}}
						order by user_group.{{.SortKey}} {{if .Reverse}} desc {{end}}
					{{end}}
					limit {{.Limit}} offset {{.Offset}}
				{{end}}
		{{end}}
	`

	// count mode
	var query, err = template.Render(sqlTmpl, req,
		template.FuncMap{
			"IsCountMode": func() bool { return true },
		},
	)
	if err != nil {
		logger.Warnf(ctx, "%+v", err)
		return nil, err
	}

	query = strutil.SimplifyString(query)
	logger.Infof(ctx, "count: %s", query)

	var total int
	p.DB.Raw(query).Count(&total)
	if err := p.DB.Error; err != nil {
		logger.Warnf(ctx, "%+v", err)
		return nil, err
	}
	if total == 0 {
		return &pbim.ListGroupsResponse{}, nil
	}

	// query mode
	query, err = template.Render(sqlTmpl, req,
		template.FuncMap{
			"IsCountMode": func() bool { return false },
		},
	)
	if err != nil {
		logger.Warnf(ctx, "%+v", err)
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
	if len(rows) == 0 {
		return &pbim.ListGroupsResponse{}, nil
	}

	// query user_id
	query, err = template.Render(`
		SELECT * FROM user_group_binding WHERE 1=0
			{{range $i, $v := .}}
				OR group_id='{{$v.GroupId}}'
			{{end}}
		`, rows,
	)
	if err != nil {
		logger.Warnf(ctx, "%+v", err)
		return nil, err
	}

	query = strutil.SimplifyString(query)
	logger.Infof(ctx, "query: %s", query)

	var bindRows []db_spec.UserGroupBinding
	p.DB.Raw(query).Find(&bindRows)
	if err := p.DB.Error; err != nil {
		logger.Warnf(ctx, "%+v", err)
		// ignore err
	}

	// convert to pb type
	var sets []*pbim.Group
	for _, v := range rows {
		sets = append(sets, v.ToPB())
	}

	// save user_id
	for _, v := range bindRows {
		for j, vj := range sets {
			if v.GroupId == vj.GroupId {
				sets[j].UserId = append(sets[j].UserId, v.UserId)
			}
		}
	}

	reply := &pbim.ListGroupsResponse{
		Group: sets,
		Total: int32(total),
	}

	return reply, nil
}

func (p *Database) ListGroupsWithUser(ctx context.Context, req *pbim.ListGroupsWithUserRequest) (*pbim.ListGroupsWithUserResponse, error) {
	err := status.Errorf(codes.Unimplemented, "TODO")
	logger.Warnf(ctx, "%+v", err)
	return nil, err
}
