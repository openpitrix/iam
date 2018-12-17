// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package service

import (
	"context"

	"openpitrix.io/iam/pkg/internal/funcutil"
	"openpitrix.io/iam/pkg/pb/am"
	"openpitrix.io/logger"
)

func (p *Server) CreateActionRule(ctx context.Context, req *pbam.Role) (*pbam.Role, error) {
	logger.Infof(ctx, funcutil.CallerName(1))

	panic("TODO")
}
func (p *Server) DeleteActionRuleByName(ctx context.Context, req *pbam.String) (*pbam.Bool, error) {
	logger.Infof(ctx, funcutil.CallerName(1))

	panic("TODO")
}

func (p *Server) GetActionRuleByName(ctx context.Context, req *pbam.String) (*pbam.Role, error) {
	logger.Infof(ctx, funcutil.CallerName(1))

	panic("TODO")
}
func (p *Server) GetActionRuleByRoleName(ctx context.Context, req *pbam.String) (*pbam.Role, error) {
	logger.Infof(ctx, funcutil.CallerName(1))

	panic("TODO")
}
func (p *Server) GetActionRuleByXid(ctx context.Context, req *pbam.XidList) (*pbam.RoleList, error) {
	logger.Infof(ctx, funcutil.CallerName(1))

	panic("TODO")
}

func (p *Server) ListActionRules(ctx context.Context, req *pbam.NameFilter) (*pbam.ActionRuleList, error) {
	logger.Infof(ctx, funcutil.CallerName(1))

	panic("TODO")
}
func (p *Server) GetXidListByRoleName(ctx context.Context, req *pbam.String) (*pbam.XidList, error) {
	logger.Infof(ctx, funcutil.CallerName(1))

	panic("TODO")
}
func (p *Server) GetRoleListByByXid(ctx context.Context, req *pbam.XidList) (*pbam.RoleList, error) {
	logger.Infof(ctx, funcutil.CallerName(1))

	panic("TODO")
}

func (p *Server) CreateRole(ctx context.Context, req *pbam.Role) (*pbam.Role, error) {
	logger.Infof(ctx, funcutil.CallerName(1))

	panic("TODO")
}

func (p *Server) ModifyRole(ctx context.Context, req *pbam.Role) (*pbam.Role, error) {
	logger.Infof(ctx, funcutil.CallerName(1))

	panic("TODO")
}

func (p *Server) DeleteRoleByName(ctx context.Context, req *pbam.String) (*pbam.Bool, error) {
	logger.Infof(ctx, funcutil.CallerName(1))

	panic("TODO")
}

func (p *Server) GetRoleByName(ctx context.Context, req *pbam.String) (*pbam.Role, error) {
	logger.Infof(ctx, funcutil.CallerName(1))

	panic("TODO")
}

func (p *Server) GetRoleByXidList(ctx context.Context, req *pbam.XidList) (*pbam.RoleList, error) {
	logger.Infof(ctx, funcutil.CallerName(1))

	panic("TODO")
}

func (p *Server) ListRoles(ctx context.Context, req *pbam.NameFilter) (*pbam.RoleList, error) {
	logger.Infof(ctx, funcutil.CallerName(1))

	panic("TODO")
}

func (p *Server) CreateRoleXidBinding(ctx context.Context, req *pbam.RoleXidBindingList) (*pbam.Bool, error) {
	panic("TODO")
}

func (p *Server) DeleteRoleXidBindingByXid(ctx context.Context, req *pbam.XidList) (*pbam.Bool, error) {
	logger.Infof(ctx, funcutil.CallerName(1))

	panic("TODO")
}

func (p *Server) GetRoleXidBindingByRoleName(ctx context.Context, req *pbam.String) (*pbam.RoleXidBindingList, error) {
	logger.Infof(ctx, funcutil.CallerName(1))

	panic("TODO")
}
func (p *Server) GetRoleXidBindingByXidList(ctx context.Context, req *pbam.XidList) (*pbam.RoleXidBindingList, error) {
	logger.Infof(ctx, funcutil.CallerName(1))

	panic("TODO")
}
func (p *Server) ListRoleXidBindings(ctx context.Context, req *pbam.NameFilter) (*pbam.RoleXidBindingList, error) {
	logger.Infof(ctx, funcutil.CallerName(1))

	panic("TODO")
}
func (p *Server) ListRoleBindings(ctx context.Context, req *pbam.NameFilter) (*pbam.RoleXidBindingList, error) {
	logger.Infof(ctx, funcutil.CallerName(1))

	panic("TODO")
}

func (p *Server) CanDo(ctx context.Context, req *pbam.Action) (*pbam.Bool, error) {
	logger.Infof(ctx, funcutil.CallerName(1))

	panic("TODO")
}
