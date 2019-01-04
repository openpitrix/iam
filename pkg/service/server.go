// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package service

import (
	"context"
	"net/http"

	"google.golang.org/grpc"
	"openpitrix.io/logger"

	"openpitrix.io/iam/pkg/config"
)

type Server struct {
	cfg        *config.Config
	webServer  *http.Server
	grpcServer *grpc.Server
	db         *Database
}

func OpenServer(cfg *config.Config) (*Server, error) {
	cfg = cfg.Clone()

	db, err := OpenDatabase(cfg)
	if err != nil {
		logger.Criticalf(nil, "%v", err)
	}

	p := &Server{
		cfg: cfg,
		db:  db,
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
