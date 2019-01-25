// Copyright 2019 The OpenPitrix Authors. All rights reserved.
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

	"openpitrix.io/iam/pkg/service/am/config"
	"openpitrix.io/iam/pkg/service/am/db_spec"
	"openpitrix.io/logger"
)

type Database struct {
	cfg *config.Config
	*gorm.DB
}

type Options struct {
	SqlInitDB    []string
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

		// init hook
		if opt != nil && len(opt.SqlInitDB) > 0 {
			var lastInitErr error
			for _, sqlList := range opt.SqlInitDB {
				for _, sql := range strings.Split(sqlList, ";") {
					sql := strings.TrimSpace(sql)
					if sql == "" {
						continue
					}
					if strings.HasPrefix(sql, "/*") && strings.HasSuffix(sql, "*/") {
						continue
					}
					if _, err := db.Exec(sql); err != nil {
						lastInitErr = err
					}
				}
			}
			if lastInitErr != nil {
				logger.Warnf(nil, "%+v", lastInitErr)
			}
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

	// init hook
	if opt != nil && len(opt.SqlInitTable) > 0 {
		var lastInitErr error
		for _, sqlList := range opt.SqlInitTable {
			for _, sql := range strings.Split(sqlList, ";") {
				sql := strings.TrimSpace(sql)
				if sql == "" {
					continue
				}
				if strings.HasPrefix(sql, "/*") && strings.HasSuffix(sql, "*/") {
					continue
				}
				if err := p.DB.Exec(sql).Error; err != nil {
					lastInitErr = err
				}
			}
		}
		if lastInitErr != nil {
			logger.Warnf(nil, "%+v", lastInitErr)
		}
	}
	if opt != nil && len(opt.SqlInitData) > 0 {
		hasRecords := p.checkDbHasRecords()
		if hasRecords {
			logger.Infof(nil, "DB has records, skip opt.SqlInitData")
		} else {
			var lastInitErr error
			for _, sqlList := range opt.SqlInitData {
				for _, sql := range strings.Split(sqlList, ";") {
					sql := strings.TrimSpace(sql)
					if sql == "" {
						continue
					}
					if strings.HasPrefix(sql, "/*") && strings.HasSuffix(sql, "*/") {
						continue
					}
					if err := p.DB.Exec(sql).Error; err != nil {
						lastInitErr = err
					}
				}
			}
			if lastInitErr != nil {
				logger.Warnf(nil, "%+v", lastInitErr)
			}
		}
	}

	// greate tables
	{
		if !p.DB.HasTable(&db_spec.ModuleApi{}) {
			if err := p.DB.CreateTable(&db_spec.ModuleApi{}).Error; err != nil {
				logger.Warnf(nil, "%+v", err)
			}
		}
		if !p.DB.HasTable(&db_spec.Role{}) {
			if err := p.DB.CreateTable(&db_spec.Role{}).Error; err != nil {
				logger.Warnf(nil, "%+v", err)
			}
		}
		if !p.DB.HasTable(&db_spec.UserRoleBinding{}) {
			if err := p.DB.CreateTable(&db_spec.UserRoleBinding{}).Error; err != nil {
				logger.Warnf(nil, "%+v", err)
			}
		}
		if !p.DB.HasTable(&db_spec.RoleModuleBinding{}) {
			if err := p.DB.CreateTable(&db_spec.RoleModuleBinding{}).Error; err != nil {
				logger.Warnf(nil, "%+v", err)
			}
		}
		if !p.DB.HasTable(&db_spec.EnableAction{}) {
			if err := p.DB.CreateTable(&db_spec.EnableAction{}).Error; err != nil {
				logger.Warnf(nil, "%+v", err)
			}
		}
	}

	return p, nil
}

func (p *Database) checkDbHasRecords() bool {
	for _, name := range db_spec.TableNameList {
		var total int
		p.DB.Raw(fmt.Sprintf("select COUNT(*) from %s limit 1", name)).Count(&total)
		if err := p.DB.Error; err != nil {
			logger.Warnf(nil, "%+v", err)
		}
		if total > 0 {
			return true
		}
	}
	return false
}

func (p *Database) Close() error {
	return p.DB.Close()
}
