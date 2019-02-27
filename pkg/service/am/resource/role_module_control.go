// Copyright 2019 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package resource

import (
	"context"

	"openpitrix.io/iam/pkg/gerr"
	"openpitrix.io/iam/pkg/global"
	"openpitrix.io/iam/pkg/models"
	"openpitrix.io/iam/pkg/pb"
	"openpitrix.io/iam/pkg/util/strutil"
)

func getModuleTree(moduleApis []*models.ModuleApi) map[string]map[string]map[string][]*models.ModuleApi {
	moduleApiMap := make(map[string]map[string]map[string][]*models.ModuleApi)
	for _, moduleApi := range moduleApis {
		_, ok := moduleApiMap[moduleApi.ModuleId]
		if !ok {
			moduleApiMap[moduleApi.ModuleId] = make(map[string]map[string][]*models.ModuleApi)
		}
		_, ok = moduleApiMap[moduleApi.ModuleId][moduleApi.FeatureId]
		if !ok {
			moduleApiMap[moduleApi.ModuleId][moduleApi.FeatureId] = make(map[string][]*models.ModuleApi)
		}
		_, ok = moduleApiMap[moduleApi.ModuleId][moduleApi.FeatureId][moduleApi.ActionBundleId]
		if !ok {
			moduleApiMap[moduleApi.ModuleId][moduleApi.FeatureId][moduleApi.ActionBundleId] = []*models.ModuleApi{}
		}
		moduleApiMap[moduleApi.ModuleId][moduleApi.FeatureId][moduleApi.ActionBundleId] = append(
			moduleApiMap[moduleApi.ModuleId][moduleApi.FeatureId][moduleApi.ActionBundleId], moduleApi)
	}
	return moduleApiMap
}

func buildModuleTree(
	roleModuleBindings []*models.RoleModuleBinding,
	visibilityModuleApis []*models.ModuleApi,
	enableModuleApis []*models.ModuleApi,
) *pb.Module {

	roleModuleBindingMap := make(map[string]*models.RoleModuleBinding)
	for _, roleModuleBinding := range roleModuleBindings {
		roleModuleBindingMap[roleModuleBinding.ModuleId] = roleModuleBinding
	}

	visibilityModuleTree := getModuleTree(visibilityModuleApis)
	enableModuleTree := getModuleTree(enableModuleApis)

	var moduleElems []*pb.ModuleElem
	for moduleId, featureTree := range visibilityModuleTree {
		enableFeatureTree, ok := enableModuleTree[moduleId]
		if !ok {
			enableFeatureTree = make(map[string]map[string][]*models.ModuleApi)
		}

		var featureSet []*pb.Feature
		var moduleName string
		for featureId, actionBundleTree := range featureTree {
			enableActionBundleTree, ok := enableFeatureTree[featureId]
			if !ok {
				enableActionBundleTree = make(map[string][]*models.ModuleApi)
			}

			var enableActionBundleIds []string
			for enableActionBundleId := range enableActionBundleTree {
				enableActionBundleIds = append(enableActionBundleIds, enableActionBundleId)
			}

			var actionBundleSet []*pb.ActionBundle
			var featureName string
			for actionBundleId, apis := range actionBundleTree {
				if len(apis) == 0 {
					continue
				}
				if featureName == "" {
					featureName = apis[0].FeatureName
				}
				if moduleName == "" {
					moduleName = apis[0].ModuleName
				}
				actionBundleSet = append(actionBundleSet, &pb.ActionBundle{
					ActionBundleId:   actionBundleId,
					ActionBundleName: apis[0].ActionBundleName,
				})
			}

			featureSet = append(featureSet, &pb.Feature{
				FeatureId:                featureId,
				FeatureName:              featureName,
				ActionBundleSet:          actionBundleSet,
				CheckedActionBundleIdSet: enableActionBundleIds,
			})
		}

		moduleElem := &pb.ModuleElem{
			ModuleId:   moduleId,
			ModuleName: moduleName,
			DataLevel:  roleModuleBindingMap[moduleId].DataLevel,
			IsCheckAll: roleModuleBindingMap[moduleId].IsCheckAll,
			FeatureSet: featureSet,
		}
		moduleElems = append(moduleElems, moduleElem)
	}

	return &pb.Module{
		ModuleElemSet: moduleElems,
	}
}

func GetRoleModule(ctx context.Context, req *pb.GetRoleModuleRequest) (*pb.GetRoleModuleResponse, error) {
	roleId := req.RoleId

	visibilityModuleApis, err := GetVisibilityModuleApis(ctx, roleId)
	if err != nil {
		return nil, err
	}

	roleModuleBindings, err := GetRoleModuleBindingsByRoleIds(ctx, []string{roleId})
	if err != nil {
		return nil, err
	}

	enableModuleApis, err := GetEnableModuleApis(ctx, []string{roleId})
	if err != nil {
		return nil, err
	}

	roleModule := buildModuleTree(
		roleModuleBindings,
		visibilityModuleApis,
		enableModuleApis,
	)
	return &pb.GetRoleModuleResponse{
		RoleId: roleId,
		Module: roleModule,
	}, nil
}

func ModifyRoleModule(ctx context.Context, req *pb.ModifyRoleModuleRequest) (*pb.ModifyRoleModuleResponse, error) {
	roleId := req.RoleId
	module := req.Module

	visibilityModuleIds, err := GetVisibilityModuleIds(ctx, roleId)
	if err != nil {
		return nil, err
	}

	roleModuleBindings, err := GetRoleModuleBindingsByRoleIdsAndModuleIds(ctx, []string{roleId}, visibilityModuleIds)
	if err != nil {
		return nil, err
	}

	var bindIds []string
	for _, roleModuleBinding := range roleModuleBindings {
		bindIds = append(bindIds, roleModuleBinding.BindId)
	}

	var newRoleModuleBindings []*models.RoleModuleBinding
	var newEnableActionBundles []*models.EnableActionBundle
	for _, moduleElem := range module.ModuleElemSet {
		moduleId := moduleElem.ModuleId
		if !strutil.Contains(visibilityModuleIds, moduleId) {
			return nil, gerr.New(ctx, gerr.NotFound, gerr.ErrorModuleNotFound, moduleId)
		}

		newRoleModuleBinding := models.NewRoleModuleBinding(
			roleId,
			moduleId,
			moduleElem.DataLevel,
			moduleElem.IsCheckAll,
		)

		if !moduleElem.IsCheckAll {
			for _, feature := range moduleElem.FeatureSet {
				for _, actionBundleId := range feature.CheckedActionBundleIdSet {
					newEnableActionBundle := models.NewEnableActionBundle(newRoleModuleBinding.BindId, actionBundleId)
					newEnableActionBundles = append(newEnableActionBundles, newEnableActionBundle)
				}
			}
		}
		newRoleModuleBindings = append(newRoleModuleBindings, newRoleModuleBinding)
	}

	tx := global.Global().Database.Begin()
	{
		// delete old role_module_binding
		if err := tx.Exec("DELETE from role_module_binding where bind_id in (?)", bindIds).Error; err != nil {
			tx.Rollback()
			return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorInternalError)
		}

		// delete old enable_action_bundle
		if err := tx.Exec("DELETE from enable_action_bundle where bind_id in (?)", bindIds).Error; err != nil {
			tx.Rollback()
			return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorInternalError)
		}

		// insert new role_module_binding
		for _, roleModuleBinding := range newRoleModuleBindings {
			if err := tx.Create(roleModuleBinding).Error; err != nil {
				tx.Rollback()
				return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorInternalError)
			}
		}
		// insert new enable_action_bundle
		for _, enableActionBundle := range newEnableActionBundles {
			if err := tx.Create(enableActionBundle).Error; err != nil {
				tx.Rollback()
				return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorInternalError)
			}
		}
	}

	if err := tx.Commit().Error; err != nil {
		return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorInternalError)
	}

	return &pb.ModifyRoleModuleResponse{
		RoleId: roleId,
	}, nil
}
