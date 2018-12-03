// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

// OpenPitrix Access Management service config package.
package config

import (
	"fmt"

	"github.com/koding/multiconfig"

	"openpitrix.io/logger"
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
	Port        int    `default:"9115"`
	TlsEnabled  bool   `default:"false"`
	TlsCertFile string `default:"server.cert"`
	TlsKeyFile  string `default:"server.key"`
}

type LogConfig struct {
	Level      string `default:"info"` // debug, info, warn, error, fatal
	GrpcDetail bool
}

type MysqlConfig struct {
	Host     string `default:"openpitrix-db"`
	Port     int    `default:"3306"`
	User     string `default:"root"`
	Password string `default:"password"`
	Database string `default:"openpitrix"`
	Disable  bool   `default:"false"`
}

func (m *MysqlConfig) DbType() string {
	return "mysql"
}
func (m *MysqlConfig) GetUrl() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", m.User, m.Password, m.Host, m.Port, m.Database)
}

func LoadConf(path string) (*Config, error) {
	conf := new(Server)
	if err := multiconfig.NewWithPath(path).Load(conf); err != nil {
		return nil, err
	}
	return &Config{conf}, nil
}

func MustLoadConf(path string) *Config {
	conf := new(Server)
	if err := multiconfig.NewWithPath(path).Load(conf); err != nil {
		logger.Criticalf(nil, "%v", err)
		panic(err)
	}

	return &Config{conf}
}
