// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

// OpenPitrix Access Management service package.
package am

import (
	"openpitrix.io/iam/pkg/am/rbac"
	"openpitrix.io/iam/pkg/pb/am"
)

var (
	_ pbam.AccessManagerServer = (*Server)(nil)
)

type Server struct {
	rbac rbac.Interface
}

func NewManager(dbtype, dbpath string) *Server {
	return nil
}

func (p *Server) Close() error {
	panic("TODO")
}
