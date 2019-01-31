// Copyright 2019 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package db

import (
	"context"
	"fmt"
	"sort"
	"strings"

	"github.com/chai2010/template"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"openpitrix.io/iam/pkg/internal/funcutil"
	"openpitrix.io/iam/pkg/internal/strutil"
	pbam "openpitrix.io/iam/pkg/pb/am"
	pbim "openpitrix.io/iam/pkg/pb/im"
	"openpitrix.io/iam/pkg/service/am/db_spec"
	"openpitrix.io/iam/pkg/validator"
	"openpitrix.io/logger"
)

func (p *Database) CanDo(ctx context.Context, req *pbam.CanDoRequest) (*pbam.CanDoResponse, error) {
	logger.Infof(ctx, funcutil.CallerName(1))

	req.UserId = strutil.SimplifyString(req.UserId)
	req.Url = strutil.SimplifyString(req.Url)
	req.UrlMethod = strutil.SimplifyString(req.UrlMethod)

	req.UrlMethod = strings.ToLower(req.UrlMethod)

	if !validator.IsValidId(req.UserId) {
		err := status.Errorf(codes.InvalidArgument, "invalid UserId: %v", req.UserId)
		logger.Warnf(ctx, "%+v", err)
		return nil, err
	}
	if req.Url == "" {
		err := status.Errorf(codes.InvalidArgument, "empty Url")
		logger.Warnf(ctx, "%+v", err)
		return nil, err
	}
	if req.UrlMethod == "" {
		err := status.Errorf(codes.InvalidArgument, "empty UrlMethod")
		logger.Warnf(ctx, "%+v", err)
		return nil, err
	}

	// 1. query RoleList
	query, err := template.Render(`
		select distinct role.* from
			role, user_role_binding
		where 1=1
			and user_role_binding.role_id=role.role_id
			and user_role_binding.user_id='{{.UserId}}'
		`, req,
	)
	if err != nil {
		logger.Warnf(ctx, "%+v", err)
		return nil, err
	}

	query = strutil.SimplifyString(query)
	logger.Infof(ctx, "query: %s", query)

	var roleList []db_spec.Role
	p.DB.Raw(query).Find(&roleList)
	if err := p.DB.Error; err != nil {
		logger.Warnf(ctx, "%+v", err)
		return nil, err
	}
	if len(roleList) == 0 {
		err := status.Errorf(codes.PermissionDenied, "no role")
		logger.Warnf(ctx, "%+v", err)
		return nil, err
	}

	// 2. query RoleModuleBindingList
	query, err = template.Render(`
		select distinct
			role_module_binding.*
		from
			role_module_binding, module_api
		where 1=1
			and role_module_binding.module_id=module_api.module_id
			and module_api.url_method='{{get_url_method}}'
			and module_api.url='{{get_url}}'

			and role_module_binding.role_id in (
				{{range $i, $v := .}}
					{{if eq $i 0}} '{{$v.RoleId}}' {{else}} ,'{{$v.RoleId}}' {{end}}
				{{end}}
			)
		`,
		roleList, template.FuncMap{
			"get_url_method": func() string { return req.UrlMethod },
			"get_url":        func() string { return req.Url },
		},
	)
	if err != nil {
		logger.Warnf(ctx, "%+v", err)
		return nil, err
	}

	query = strutil.SimplifyString(query)
	logger.Infof(ctx, "query: %s", query)

	var roleModuleBindList []db_spec.RoleModuleBinding
	p.DB.Raw(query).Find(&roleModuleBindList)
	if err := p.DB.Error; err != nil {
		logger.Warnf(ctx, "%+v", err)
		return nil, err
	}
	if len(roleModuleBindList) == 0 {
		err := status.Errorf(codes.PermissionDenied, "no module")
		logger.Warnf(ctx, "%+v", err)
		return nil, err
	}

	// check all?
	var isCheckAll = false
	for _, v := range roleModuleBindList {
		if v.IsCheckAll != 0 {
			isCheckAll = true
			break
		}
	}

	// can do?
	// 3. isCheckAll or query enable_action
	var canDoReuest = false
	if isCheckAll {
		canDoReuest = true
	} else {
		query, err = template.Render(`
			select distinct enable_action.* from
				enable_action, role_module_binding, module_api
			where 1=1
				and enable_action.bind_id=role_module_binding.bind_id
				and enable_action.action_id=module_api.action_id
				and module_api.module_id=role_module_binding.module_id

				and role_module_binding.role_id in (
					{{range $i, $v := .}}
						{{if eq $i 0}} '{{$v.RoleId}}' {{else}} ,'{{$v.RoleId}}' {{end}}
					{{end}}
				)
				and module_api.url_method='{{get_url_method}}'
				and module_api.url='{{get_url}}'

				LIMIT 1
			`,
			roleList, template.FuncMap{
				"get_url_method": func() string { return req.UrlMethod },
				"get_url":        func() string { return req.Url },
			},
		)
		if err != nil {
			logger.Warnf(ctx, "%+v", err)
			return nil, err
		}

		query = strutil.SimplifyString(query)
		logger.Infof(ctx, "query: %s", query)

		var enableActionList []db_spec.EnableAction
		p.DB.Raw(query).Find(&enableActionList)
		if err := p.DB.Error; err != nil {
			logger.Warnf(ctx, "%+v", err)
			return nil, err
		}
		if len(enableActionList) > 0 {
			canDoReuest = true
		}
	}

	// disabled
	if !canDoReuest {
		err := status.Errorf(codes.PermissionDenied, "disabled")
		logger.Warnf(ctx, "%+v", err)
		return nil, err
	}

	// get groupPath from IM server
	groupPath, err := p.getShortestGroupPathByUserId(ctx, req.UserId)
	if err != nil {
		logger.Warnf(ctx, "%+v", err)
		return nil, err
	}
	if groupPath == "" {
		logger.Warnf(ctx, "no group, req: %+v", req)
		// ignore err
	}

	// get data_level
	var dataLevel string
	if dataLevel == "" {
		for _, v := range roleModuleBindList {
			if v.DataLevel == db_spec.DataLevel_All {
				dataLevel = db_spec.DataLevel_All
				break
			}
		}
	}
	if dataLevel == "" {
		for _, v := range roleModuleBindList {
			if v.DataLevel == db_spec.DataLevel_Group {
				dataLevel = db_spec.DataLevel_Group
				break
			}
		}
	}
	if dataLevel == "" {
		dataLevel = db_spec.DataLevel_Self
	}

	// get access path
	var accessPath string
	switch dataLevel {
	case db_spec.DataLevel_All:
		accessPath = "" // all path
	case db_spec.DataLevel_Group:
		accessPath = groupPath
	case db_spec.DataLevel_Self:
		accessPath = groupPath + ":" + req.UserId
	default:
		logger.Warnf(ctx, "unreachable, should panic")
		accessPath = "???"
	}

	// Done
	reply := &pbam.CanDoResponse{
		UserId:     req.UserId,
		OwnerPath:  groupPath + ":" + req.UserId,
		AccessPath: accessPath,
	}

	return reply, nil
}

func (p *Database) getShortestGroupPathByUserId(ctx context.Context, userId string) (string, error) {
	conn, err := grpc.Dial(fmt.Sprintf("%s:%d", p.cfg.ImHost, p.cfg.ImPort), grpc.WithInsecure())
	if err != nil {
		logger.Warnf(ctx, "%+v", err)
		return "", err
	}
	defer conn.Close()

	client := pbim.NewAccountManagerClient(conn)
	reply, err := client.GetGroupsByUserId(ctx, &pbim.UserId{UserId: userId})
	if err != nil {
		logger.Warnf(ctx, "%+v", err)
		return "", err
	}

	// no group
	if len(reply.Value) == 0 {
		logger.Warnf(ctx, "no group")
		return "", nil
	}

	// sort group
	sort.Slice(reply.Value, func(i, j int) bool {
		iGroupPathLevel := strings.Count(reply.Value[i].GroupPath, ".") + 1
		jGroupPathLevel := strings.Count(reply.Value[j].GroupPath, ".") + 1

		if iGroupPathLevel != jGroupPathLevel {
			return iGroupPathLevel < jGroupPathLevel
		}

		return reply.Value[i].GroupPath < reply.Value[j].GroupPath
	})

	// take first group
	ownerPath := reply.Value[0].GroupPath

	// OK
	return ownerPath, nil
}
