// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package rabc

import (
	"sync"

	"openpitrix.io/iam/pkg/pb/am"
)

type rabcDbServer struct {
	bindings []pbam.RoleBinding
	roles    []pbam.Role

	sync.Mutex
}

func openDatabase(dbpath string, opt *DBOptions) (Interface, error) {
	return new(rabcDbServer), nil
}

func (p *rabcDbServer) Close() error {
	return nil
}

func (p *rabcDbServer) CanDo(x pbam.Action) bool {
	panic("TODO")
}

func (p *rabcDbServer) AllRoles() []pbam.Role {
	panic("TODO")
}
func (p *rabcDbServer) AllRoleBindings() []pbam.RoleBinding {
	panic("TODO")
}

func (p *rabcDbServer) GetRoleByName(name string) (role pbam.Role, ok bool) {
	panic("TODO")
}
func (p *rabcDbServer) GetRoleByXid(xid []string) pbam.RoleList {
	panic("TODO")
}

func (p *rabcDbServer) CreateRole(role pbam.Role) error {
	panic("TODO")
}
func (p *rabcDbServer) ModifyRole(role pbam.Role) error {
	panic("TODO")
}
func (p *rabcDbServer) DeleteRole(name string) error {
	panic("TODO")
}

func (p *rabcDbServer) CreateRoleBinding(x []pbam.RoleBinding) error {
	panic("TODO")
}
func (p *rabcDbServer) DeleteRoleBinding(xid []string) error {
	panic("TODO")
}
