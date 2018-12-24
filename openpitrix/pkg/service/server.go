// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package service

import (
	"context"
	"net/http"

	"google.golang.org/grpc"
	"openpitrix.io/logger"
)

type Server struct {
	webServer  *http.Server
	grpcServer *grpc.Server
	db         *Database
}

func OpenServer(dbtype, dbpath string) (*Server, error) {
	db, err := Open(dbtype, dbpath)
	if err != nil {
		logger.Criticalf(nil, "%v", err)
	}

	p := &Server{
		db: db,
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

	if p.db != nil {
		if err := p.db.Close(); err != nil {
			lastErr = err
			p.db = nil
		}
	}

	if lastErr != nil {
		return lastErr
	}
	return nil
}
