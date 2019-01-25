// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package im

import (
	"openpitrix.io/logger"

	"openpitrix.io/iam/pkg/service/im/config"
	"openpitrix.io/iam/pkg/service/im/db"
)

type Server struct {
	cfg *config.Config
	db  *db.Database
}

func OpenServer(cfg *config.Config) (*Server, error) {
	cfg = cfg.Clone()

	db, err := db.OpenDatabase(cfg, nil)
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
	return p.db.Close()
}
