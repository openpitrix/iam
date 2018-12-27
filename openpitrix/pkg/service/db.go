// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package service

import (
	"strings"

	"github.com/fatih/structs"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"

	"openpitrix.io/logger"
)

type Database struct {
	*sqlx.DB
}

func Open(dbtype, dbpath string) (*Database, error) {
	if dbtype == "mysql" {
		if !strings.Contains(dbpath, "parseTime=true") {
			dbpath += "?parseTime=true"
		}
	}

	db, err := sqlx.Open(dbtype, dbpath)
	if err != nil {
		return nil, err
	}

	p := &Database{DB: db}
	for _, v := range InitSqlList {
		if v.Name != "-" {
			if _, err := p.Exec(v.Sql); err != nil {
				logger.Warnf(nil, "%v", err)
			}
		}
	}

	if err := p.DB.Ping(); err != nil {
		logger.Warnf(nil, "ping faild: %#v", err)
	} else {
		logger.Infof(nil, "ping ok")
	}

	logger.Infof(nil, "DB stats: begin")
	for _, f := range structs.Fields(p.DB.Stats()) {
		logger.Infof(nil, "\t%s: %v", f.Name(), f.Value())
	}
	logger.Infof(nil, "DB stats: end")

	return p, nil
}

func (p *Database) Close() error {
	return p.DB.Close()
}
