// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

// OpenPitrix Access Management service config package.
package config

import (
	"fmt"

	"github.com/koding/multiconfig"
)

type Config struct {
	*Server
}

type Server struct {
	Log   LogConfig
	Mysql MysqlConfig
	AM    AMConfig
}

type AMConfig struct {
	Port     string `default:"9115"`
	RuleFile string `default:"rabc.json"`
}

type LogConfig struct {
	Level      string `default:"info"` // debug, info, warn, error, fatal
	GrpcDetail bool
}

type MysqlConfig struct {
	Host     string `default:"openpitrix-db"`
	Port     string `default:"3306"`
	User     string `default:"root"`
	Password string `default:"password"`
	Database string `default:"openpitrix"`
	Disable  bool   `default:"false"`
}

func (m *MysqlConfig) GetUrl() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", m.User, m.Password, m.Host, m.Port, m.Database)
}

func LoadConf(path string) (*Config, error) {
	conf := new(Server)
	if err := multiconfig.NewWithPath(path).Load(conf); err != nil {
		return nil, err
	}
	return &Config{conf}, nil
}
