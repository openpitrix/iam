// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

// OpenPitrix Access Management service package.
package am

import (
	"context"
	"net/http"

	"google.golang.org/grpc"

	"openpitrix.io/iam/pkg/am/rbac"
	"openpitrix.io/iam/pkg/pb/am"
)

var (
	_ pbam.AccessManagerServer = (*Server)(nil)
)

type Server struct {
	rbac       rbac.Interface
	webServer  *http.Server
	grpcServer *grpc.Server
}

func OpenServer(dbtype, dbpath string) (*Server, error) {
	rbacManager, err := rbac.OpenDatabase(dbtype, dbpath)
	if err != nil {
		return nil, err
	}

	p := &Server{
		rbac: rbacManager,
	}
	return p, nil
}

func (p *Server) Close() error {
	var lastErr error

	if p.grpcServer != nil {
		p.grpcServer.Stop()
		p.grpcServer = nil
	}
	if p.webServer != nil {
		if err := p.webServer.Shutdown(context.Background()); err != nil {
			lastErr = err
		}
		p.webServer = nil
	}

	if err := p.rbac.Close(); err != nil {
		lastErr = err
	}

	if lastErr != nil {
		return lastErr
	}
	return nil
}
