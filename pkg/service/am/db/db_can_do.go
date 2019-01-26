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

	// 2. query RoleModuleList
	query, err = template.Render(`
		select distinct * from role_module_binding where 1=1
			and role_id in (
				{{range $i, $v := .}}
					{{if eq $i 0}} '{{$v.RoleId}}' {{else}} ,'{{$v.RoleId}}' {{end}}
				{{end}}
			)
		`, roleList,
	)
	if err != nil {
		logger.Warnf(ctx, "%+v", err)
		return nil, err
	}

	query = strutil.SimplifyString(query)
	logger.Infof(ctx, "query: %s", query)

	var roleModuleList []db_spec.RoleModuleBinding
	p.DB.Raw(query).Find(&roleModuleList)
	if err := p.DB.Error; err != nil {
		logger.Warnf(ctx, "%+v", err)
		return nil, err
	}
	if len(roleModuleList) == 0 {
		err := status.Errorf(codes.PermissionDenied, "no nodule")
		logger.Warnf(ctx, "%+v", err)
		return nil, err
	}

	// 3. query ModuleApiList
	query, err = template.Render(`
		select distinct * from module_api where 1=1
			and module_id in (
				{{range $i, $v := .}}
					{{if eq $i 0}} '{{$v.ModuleId}}' {{else}} ,'{{$v.ModuleId}}' {{end}}
				{{end}}
			)
		`, roleModuleList,
	)
	if err != nil {
		logger.Warnf(ctx, "%+v", err)
		return nil, err
	}

	query = strutil.SimplifyString(query)
	logger.Infof(ctx, "query: %s", query)

	var moduleApiList []db_spec.ModuleApi
	p.DB.Raw(query).Find(&moduleApiList)
	if err := p.DB.Error; err != nil {
		logger.Warnf(ctx, "%+v", err)
		return nil, err
	}
	if len(moduleApiList) == 0 {
		err := status.Errorf(codes.PermissionDenied, "no nodule api")
		logger.Warnf(ctx, "%+v", err)
		return nil, err
	}

	// can do
	var canDoReuest = false
	for _, v := range roleModuleList {
		if v.IsCheckAll != 0 {
			canDoReuest = true
			break
		}
	}
	if !canDoReuest {
		for _, v := range moduleApiList {
			if v.Url == req.Url && v.UrlMethod == req.UrlMethod {
				canDoReuest = true
				break
			}
		}
	}

	if !canDoReuest {
		err := status.Errorf(codes.PermissionDenied, "disabled")
		logger.Warnf(ctx, "%+v", err)
		return nil, err
	}

	// get groupPath from IM server
	groupPath, err := p.getGrouprPathByUserId(ctx, req.UserId)
	if err != nil {
		logger.Warnf(ctx, "%+v", err)
		return nil, err
	}
	if groupPath == "" {
		logger.Warnf(ctx, "no group, req: %+v", req)
		// ignore err
	}

	// get access path
	accessPath, err := p.getAccessPathBy(ctx, req, groupPath)
	if err != nil {
		logger.Warnf(ctx, "%+v", err)
		return nil, err
	}

	reply := &pbam.CanDoResponse{
		UserId:     req.UserId,
		OwnerPath:  groupPath + ":" + req.UserId, // todo: portal: group or +uid
		AccessPath: accessPath,
	}

	return reply, nil
}

func (p *Database) getAccessPathBy(ctx context.Context, req *pbam.CanDoRequest, ownerPath string) (string, error) {
	query, err := template.Render(sqlGetAccessPath_tmpl, &sqlGetAccessPath_args{
		UserId:    req.UserId,
		OwnerPath: ownerPath,
		Url:       req.Url,
		UrlMethod: req.UrlMethod,
	})
	if err != nil {
		logger.Warnf(ctx, "%+v", err)
		return "", err
	}

	type Result struct {
		AccessPath string
	}
	var rows []Result

	err = p.DB.Raw(query).Scan(&rows).Error
	if err != nil {
		logger.Warnf(ctx, "%v", query)
		logger.Warnf(ctx, "%+v", err)
		return "", err
	}
	if len(rows) == 0 {
		return "", nil
	}

	return rows[0].AccessPath, nil
}

func (p *Database) getGrouprPathByUserId(ctx context.Context, userId string) (string, error) {
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
