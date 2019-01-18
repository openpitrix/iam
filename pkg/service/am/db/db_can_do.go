// Copyright 2019 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package db

import (
	"context"
	"regexp"
	"strings"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"openpitrix.io/iam/pkg/internal/funcutil"
	pbam "openpitrix.io/iam/pkg/pb/am"
	"openpitrix.io/logger"
)

func (p *Database) CanDo(ctx context.Context, req *pbam.CanDoRequest) (*pbam.CanDoResponse, error) {
	logger.Infof(ctx, funcutil.CallerName(1))

	// skip version /v123/
	if matched, _ := regexp.MatchString(`/v\d+`, req.UrlMethod); matched {
		if idx := strings.Index(req.UrlMethod[2:], "/"); idx > 0 {
			req.UrlMethod = req.UrlMethod[idx+1:]
		}
	}

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

	// 1. get role list by user_id
	// 2. get action list by role_id
	// 3. check action rule

	// try get OwnerPath from user_path

	reply := &pbam.CanDoResponse{
		UserId:     req.UserId,
		AccessPath: "AccessPath-todo",
		OwnerPath:  "OwnerPath-todo",
	}

	return reply, nil
}
