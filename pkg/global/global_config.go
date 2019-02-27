// Copyright 2019 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package global

import (
	"sync"

	"github.com/google/gops/agent"

	"openpitrix.io/iam/pkg/config"
	"openpitrix.io/iam/pkg/db"
	"openpitrix.io/logger"
)

var global *Config
var globalMutex sync.RWMutex

func SetGlobal(config *config.Config) {
	globalMutex.Lock()
	global = NewConfig(config)
	globalMutex.Unlock()
}

func Global() *Config {
	globalMutex.RLock()
	defer globalMutex.RUnlock()
	return global
}

type Config struct {
	Config   *config.Config
	Database *db.Database
}

func NewConfig(config *config.Config) *Config {
	c := &Config{Config: config}
	c.openDatabase()

	if err := agent.Listen(agent.Options{
		ShutdownCleanup: true,
	}); err != nil {
		logger.Criticalf(nil, "failed to start gops agent")
	}

	return c
}

func (c *Config) openDatabase() {
	database, err := db.OpenDatabase(c.Config)
	if err != nil {
		logger.Criticalf(nil, "failed to connect database")
		panic(err)
	}
	c.Database = database
}
