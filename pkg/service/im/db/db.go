// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package db

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	_ "github.com/jinzhu/gorm/dialects/sqlite"

	"openpitrix.io/iam/pkg/service/im/config"
	"openpitrix.io/logger"
)

type Database struct {
	cfg *config.Config
	*gorm.DB
}

type Options struct {
	SqlInitTable []string
	SqlInitData  []string
}

func OpenDatabase(cfg *config.Config, opt *Options) (*Database, error) {
	cfg = cfg.Clone()

	// create db if not exists
	if strings.EqualFold(cfg.DB.Type, "mysql") {
		db, err := sql.Open("mysql", cfg.DB.GetHost())
		if err != nil {
			logger.Warnf(nil, "%v", err)
		}
		defer db.Close()

		query := fmt.Sprintf(
			"CREATE DATABASE IF NOT EXISTS %s DEFAULT CHARACTER SET utf8mb4 DEFAULT COLLATE utf8mb4_unicode_ci;",
			cfg.DB.Database,
		)

		_, err = db.Exec(query)
		if err != nil {
			logger.Warnf(nil, "query = %s, err = %v", query, err)
		}
	}

	logger.Infof(nil, "DB config: begin")
	logger.Infof(nil, "\tType: %s", cfg.DB.Type)
	logger.Infof(nil, "\tHost: %s", cfg.DB.Host)
	logger.Infof(nil, "\tPort: %d", cfg.DB.Port)
	logger.Infof(nil, "\tUser: %s", cfg.DB.User)
	logger.Infof(nil, "\tDatabase: %s", cfg.DB.Database)
	logger.Infof(nil, "DB config: end")

	var p = &Database{cfg: cfg}
	var err error

	p.DB, err = gorm.Open(cfg.DB.Type, cfg.DB.GetUrl())
	if err != nil {
		return nil, err
	}

	p.DB.SingularTable(true)
	p.DB.Set("gorm:table_options", "ENGINE=InnoDB DEFAULT CHARSET=utf8")

	// init hook
	if opt != nil && len(opt.SqlInitTable) > 0 {
		for _, sql := range opt.SqlInitTable {
			if err := p.DB.Exec(sql).Error; err != nil {
				logger.Warnf(nil, "%+v", err)
			}
		}
	}
	if opt != nil && len(opt.SqlInitData) > 0 {
		for _, sql := range opt.SqlInitData {
			if err := p.DB.Exec(sql).Error; err != nil {
				logger.Warnf(nil, "%+v", err)
			}
		}
	}

	// greate tables
	{
		if !p.DB.HasTable(&User{}) {
			if err := p.DB.CreateTable(&User{}).Error; err != nil {
				logger.Warnf(nil, "%+v", err)
			}
		}
		if !p.DB.HasTable(&UserGroup{}) {
			if err := p.DB.CreateTable(&UserGroup{}).Error; err != nil {
				logger.Warnf(nil, "%+v", err)
			}
		}
		if !p.DB.HasTable(&UserGroupBinding{}) {
			if err := p.DB.CreateTable(&UserGroupBinding{}).Error; err != nil {
				logger.Warnf(nil, "%+v", err)
			}
		}
	}

	return p, nil
}

func (p *Database) Close() error {
	return p.DB.Close()
}
