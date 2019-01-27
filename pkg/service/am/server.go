// Copyright 2019 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package am

import (
	"openpitrix.io/logger"

	"openpitrix.io/iam/pkg/service/am/config"
	"openpitrix.io/iam/pkg/service/am/db"
	"openpitrix.io/iam/pkg/service/am/db/static"
)

type Server struct {
	cfg *config.Config
	db  *db.Database
}

func OpenServer(cfg *config.Config) (*Server, error) {
	cfg = cfg.Clone()

	dbInitOpt := &db.Options{}
	dbInitOpt.SqlInitDB = append(dbInitOpt.SqlInitDB, static.Files["V0_1__init.sql"])

	dbInitOpt.SqlInitDB = append(dbInitOpt.SqlInitDB, cfg.DB.InitDB.SqlInitDB...)
	dbInitOpt.SqlInitTable = append(dbInitOpt.SqlInitTable, cfg.DB.InitDB.SqlInitTable...)
	dbInitOpt.SqlInitData = append(dbInitOpt.SqlInitData, cfg.DB.InitDB.SqlInitData...)

	db, err := db.OpenDatabase(cfg, dbInitOpt)
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
