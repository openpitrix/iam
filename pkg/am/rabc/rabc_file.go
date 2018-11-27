// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package rabc

import (
	"sync"

	"openpitrix.io/iam/pkg/pb/am"
)

type rabcFileServer struct {
	bindings []pbam.RoleBinding
	roles    []pbam.Role

	sync.Mutex
}

func openFileServer(jsonpath string) (Interface, error) {
	return new(rabcFileServer), nil
}

func (p *rabcFileServer) Close() error {
	return nil
}

func (p *rabcFileServer) CanDo(x pbam.Action) bool {
	panic("TODO")
}

func (p *rabcFileServer) AllRoles() []pbam.Role {
	panic("TODO")
}
func (p *rabcFileServer) AllRoleBindings() []pbam.RoleBinding {
	panic("TODO")
}

func (p *rabcFileServer) GetRoleByName(name string) (role pbam.Role, ok bool) {
	panic("TODO")
}
func (p *rabcFileServer) GetRoleByXid(xid []string) pbam.RoleList {
	panic("TODO")
}

func (p *rabcFileServer) CreateRole(role pbam.Role) error {
	panic("TODO")
}
func (p *rabcFileServer) ModifyRole(role pbam.Role) error {
	panic("TODO")
}
func (p *rabcFileServer) DeleteRole(name string) error {
	panic("TODO")
}

func (p *rabcFileServer) CreateRoleBinding(x []pbam.RoleBinding) error {
	panic("TODO")
}
func (p *rabcFileServer) DeleteRoleBinding(xid []string) error {
	panic("TODO")
}
