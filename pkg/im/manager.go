// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

// OpenPitrix Identity Management service package.
package im

import (
	"context"

	"openpitrix.io/iam/pkg/pb/im"
)

var (
	_ pbim.AccountManagerServer = (*AccountManager)(nil)
)

type AccountManager struct {
	p int
}

func NewAccountManager() *AccountManager {
	return nil
}

func (p *AccountManager) CreateUser(ctx context.Context, req *pbim.User) (*pbim.User, error) {
	panic("TODO")
}

func (p *AccountManager) CreateGroup(ctx context.Context, req *pbim.Group) (*pbim.Group, error) {
	panic("TODO")
}

func (p *AccountManager) GetUser(ctx context.Context, req *pbim.Id) (*pbim.User, error) {
	panic("TODO")
}

func (p *AccountManager) GetGroup(ctx context.Context, req *pbim.Id) (*pbim.Group, error) {
	panic("TODO")
}

func (p *AccountManager) ListUesrs(ctx context.Context, req *pbim.Range) (*pbim.ListUesrsResponse, error) {
	panic("TODO")
}

func (p *AccountManager) ListGroups(ctx context.Context, req *pbim.Range) (*pbim.ListGroupsResponse, error) {
	panic("TODO")
}

func (p *AccountManager) ModifyUser(ctx context.Context, req *pbim.ModifyUsersRequest) (*pbim.User, error) {
	panic("TODO")
}

func (p *AccountManager) ModifyGroup(ctx context.Context, req *pbim.ModifyGroupsRequest) (*pbim.Group, error) {
	panic("TODO")
}

func (p *AccountManager) ComparePassword(ctx context.Context, req *pbim.Password) (*pbim.Bool, error) {
	panic("TODO")
}

func (p *AccountManager) ModifyPassword(ctx context.Context, req *pbim.Password) (*pbim.Bool, error) {
	panic("TODO")
}

func (p *AccountManager) DeleteUsers(ctx context.Context, req *pbim.IdList) (*pbim.Bool, error) {
	panic("TODO")
}

func (p *AccountManager) DeleteGroups(ctx context.Context, req *pbim.IdList) (*pbim.Bool, error) {
	panic("TODO")
}
