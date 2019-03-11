// Copyright 2019 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package resource

import (
	"context"
	"math"
	"strings"

	pbim "kubesphere.io/im/pkg/pb"

	imclient "openpitrix.io/iam/pkg/client/im"
	"openpitrix.io/iam/pkg/constants"
	"openpitrix.io/iam/pkg/gerr"
	"openpitrix.io/iam/pkg/pb"
	"openpitrix.io/iam/pkg/util/stringutil"
	"openpitrix.io/logger"
)

func GetDataLevelNum(dataLevel string) int {
	switch dataLevel {
	case constants.DataLevelSelf:
		return 1
	case constants.DataLevelGroup:
		return 2
	case constants.DataLevelAll:
		return 3
	default:
		return 1
	}
}

func CanDo(ctx context.Context, req *pb.CanDoRequest) (*pb.CanDoResponse, error) {
	userId := stringutil.SimplifyString(req.UserId)
	url := stringutil.SimplifyString(req.Url)
	urlMethod := strings.ToLower(stringutil.SimplifyString(req.UrlMethod))
	apiMethod := stringutil.SimplifyString(req.ApiMethod)

	var roleIds []string
	if userId == constants.UserSystem {
		return &pb.CanDoResponse{
			UserId:     userId,
			OwnerPath:  ":" + userId,
			AccessPath: "",
		}, nil
	} else {
		var err error
		roleIds, err = GetRoleIdsByUserIds(ctx, []string{userId})
		if err != nil {
			return nil, err
		}
	}

	enableModuleApis, err := GetEnableModuleApis(ctx, roleIds)
	if err != nil {
		return nil, err
	}

	canDo := false
	var moduleIds []string

	for _, enableModuleApi := range enableModuleApis {
		if (enableModuleApi.Url == url && enableModuleApi.UrlMethod == urlMethod) ||
			(enableModuleApi.ApiMethod == apiMethod) {
			canDo = true
			moduleIds = append(moduleIds, enableModuleApi.ModuleId)
		}
	}

	if !canDo {
		logger.Errorf(ctx, "Permission denied for user_id [%s], url [%s], url_method [%s], api_method [%s]",
			userId, url, urlMethod, apiMethod)
		return nil, gerr.New(ctx, gerr.PermissionDenied, gerr.ErrorPermissionDenied)
	}

	groupPath, err := GetUserGroupPath(ctx, req.UserId)
	if err != nil {
		return nil, err
	}
	ownerPath := groupPath + ":" + req.UserId

	roleModuleBindings, err := GetRoleModuleBindingsByRoleIdsAndModuleIds(ctx, roleIds, moduleIds)
	if err != nil {
		return nil, err
	}

	dataLevel := constants.DataLevelSelf
	for _, roleModuleBinding := range roleModuleBindings {
		if GetDataLevelNum(roleModuleBinding.DataLevel) > GetDataLevelNum(dataLevel) {
			dataLevel = roleModuleBinding.DataLevel
		}
	}

	var accessPath string
	switch dataLevel {
	case constants.DataLevelAll:
		accessPath = ""
	case constants.DataLevelGroup:
		accessPath = groupPath
	case constants.DataLevelSelf:
		accessPath = ownerPath
	}

	reply := &pb.CanDoResponse{
		UserId:     userId,
		OwnerPath:  ownerPath,
		AccessPath: accessPath,
	}

	return reply, nil
}

func GetUserGroupPath(ctx context.Context, userId string) (string, error) {
	var userGroupPath string

	if userId == constants.UserSystem {
		return "", nil
	}

	imClient, err := imclient.NewClient()
	if err != nil {
		logger.Errorf(ctx, "Connect to im service failed: %+v", err)
		return "", gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorInternalError)
	}

	response, err := imClient.GetUserWithGroup(ctx, &pbim.GetUserRequest{UserId: userId})
	if err != nil {
		logger.Errorf(ctx, "Get user with group failed: %+v", err)
		return "", gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorInternalError)
	}

	groups := response.User.GroupSet

	//If one user under different groups, get the highest group path.
	minLevel := math.MaxInt32
	for _, group := range groups {
		level := len(strings.Split(group.GroupPath, "."))
		if level < minLevel {
			minLevel = level
			userGroupPath = group.GetGroupPath()
		}
	}

	return userGroupPath, nil
}
