// Copyright 2019 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package db

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	"github.com/chai2010/template"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"openpitrix.io/iam/pkg/internal/funcutil"
	pbam "openpitrix.io/iam/pkg/pb/am"
	pbim "openpitrix.io/iam/pkg/pb/im"
	"openpitrix.io/logger"
)

func (p *Database) CanDo(ctx context.Context, req *pbam.CanDoRequest) (*pbam.CanDoResponse, error) {
	logger.Infof(ctx, funcutil.CallerName(1))

	if matched, _ := regexp.MatchString(`^/v\d+`, req.Url); matched {
		if idx := strings.Index(req.Url[2:], "/"); idx >= 0 {
			req.Url = req.Url[2:][idx+1:]
		}
	}

	logger.Infof(nil, "req.Url: %v", req.Url)

	type DBCanDo struct {
		Url       string
		UrlMethod string
	}

	var query = sqlCanDo
	var rows = []DBCanDo{}

	err := p.DB.Raw(query, req.UserId, req.Url, req.UrlMethod).Scan(&rows).Error
	if err != nil {
		logger.Warnf(ctx, "%v", query)
		logger.Warnf(ctx, "%+v", err)
		return nil, err
	}
	if len(rows) == 0 {
		err := status.Errorf(codes.PermissionDenied, "disable")
		logger.Warnf(ctx, "%+v", err)
		return nil, err
	}

	// get owner path from IM server
	ownerPath, err := p.getOwnerPathByUserId(ctx, req.UserId)
	if err != nil {
		logger.Warnf(ctx, "%+v", err)
		return nil, err
	}
	if ownerPath == "" {
		//
	}

	// get access path
	accessPath, err := p.getAccessPathBy(ctx, req, ownerPath)
	if len(rows) == 0 {
		logger.Warnf(ctx, "%+v", err)
		return nil, err
	}

	reply := &pbam.CanDoResponse{
		UserId:     req.UserId,
		OwnerPath:  ownerPath,
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

func (p *Database) getOwnerPathByUserId(ctx context.Context, userId string) (string, error) {
	conn, err := grpc.Dial(fmt.Sprintf("%s:%d", p.cfg.ImHost, p.cfg.ImPort), grpc.WithInsecure())
	if err != nil {
		return "", err
	}
	defer conn.Close()

	client := pbim.NewAccountManagerClient(conn)
	reply, err := client.GetGroupsByUserId(ctx, &pbim.UserId{UserId: userId})
	if err != nil {
		return "", err
	}

	logger.Infof(ctx, "getOwnerPathByUserId: %v", reply.GetValue())

	// no group
	if len(reply.Value) == 0 {
		return "", nil
	}

	// take first group
	ownerPath := reply.Value[0].GroupPath

	// OK
	return ownerPath, nil
}
