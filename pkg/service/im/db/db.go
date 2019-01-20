// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package db

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/fatih/structs"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"

	"openpitrix.io/iam/pkg/service/im/config"
	"openpitrix.io/iam/pkg/service/im/db_spec"
	"openpitrix.io/logger"
)

type Database struct {
	cfg *config.Config
	dbx *sqlx.DB
}

func OpenDatabase(cfg *config.Config) (*Database, error) {
	cfg = cfg.Clone()

	// init db
	func() {
		if strings.EqualFold(cfg.DB.Type, "mysql") {
			if !isValidDatabaseName(cfg.DB.Database) {
				logger.Warnf(nil, "invalid db name %s", cfg.DB.Database)
			}

			db, err := sql.Open("mysql", cfg.DB.GetHostUrl())
			if err != nil {
				logger.Warnf(nil, "%v", err)
			}
			defer db.Close()

			query := fmt.Sprintf(
				"CREATE DATABASE IF NOT EXISTS %s CHARACTER SET utf8 COLLATE utf8_general_ci;",
				cfg.DB.Database,
			)
			_, err = db.Exec(query)
			if err != nil {
				logger.Warnf(nil, "query = %s, err = %v", query, err)
			}
		}
	}()

	logger.Infof(nil, "DB config: begin")
	logger.Infof(nil, "\tType: %s", cfg.DB.Type)
	logger.Infof(nil, "\tHost: %s", cfg.DB.Host)
	logger.Infof(nil, "\tPort: %d", cfg.DB.Port)
	logger.Infof(nil, "\tUser: %s", cfg.DB.User)
	logger.Infof(nil, "\tDatabase: %s", cfg.DB.Database)
	logger.Infof(nil, "DB config: end")

	db, err := sqlx.Open(cfg.DB.Type, cfg.DB.GetUrlWithParseTime())
	if err != nil {
		return nil, err
	}

	p := &Database{
		cfg: cfg,
		dbx: db,
	}
	for _, v := range db_spec.TableMap {
		if !isValidDatabaseTableName(v.Name) {
			logger.Warnf(nil, "invalid table name %s", v.Name)
		}
		if _, err := p.dbx.Exec(v.Sql); err != nil {
			logger.Warnf(nil, "%s: %v", v.Name, err)
		}
	}

	if err := p.dbx.Ping(); err != nil {
		logger.Warnf(nil, "ping faild: %#v", err)
	} else {
		logger.Infof(nil, "ping ok")
	}

	logger.Infof(nil, "DB stats: begin")
	for _, f := range structs.Fields(p.dbx.Stats()) {
		logger.Infof(nil, "\t%s: %v", f.Name(), f.Value())
	}
	logger.Infof(nil, "DB stats: end")

	return p, nil
}

func (p *Database) Close() error {
	return p.dbx.Close()
}
