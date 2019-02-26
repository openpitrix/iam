// Copyright 2019 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package config

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/koding/multiconfig"
	yaml "gopkg.in/yaml.v2"

	"openpitrix.io/iam/pkg/util/jsonutil"
	"openpitrix.io/logger"
)

const EnvPrefix = "AM"

type Config struct {
	DB DBConfig

	AMHost      string `default:"am-service"`
	AMPort      int    `default:"9120"`
	IMHost      string `default:"im-service"`
	IMPort      int    `default:"9119"`
	TlsEnabled  bool   `default:"false"`
	TlsCertFile string `default:"server.cert"`
	TlsKeyFile  string `default:"server.key"`
	LogLevel    string `default:"DEBUG"`
}

type DBConfig struct {
	Type          string `default:"mysql"`
	Host          string `default:"am-db"`
	Port          int    `default:"3306"`
	User          string `default:"root"`
	Password      string `default:"password"`
	Database      string `default:"am"`
	LogModeEnable bool   `default:"false"`
}

func (m *DBConfig) GetHost() string {
	if m.Type == "sqlite3" {
		return m.Database
	}
	if m.Type == "mysql" {
		return fmt.Sprintf(
			"%s:%s@tcp(%s:%d)/",
			m.User, m.Password, m.Host, m.Port,
		)
	}
	return m.Database
}

func (m *DBConfig) GetUrl() string {
	if m.Type == "sqlite3" {
		return m.Database
	}
	if m.Type == "mysql" {
		return fmt.Sprintf(
			"%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=True&loc=Local",
			m.User, m.Password, m.Host, m.Port,
			m.Database,
		)
	}
	return m.Database
}

func Default() *Config {
	conf := new(Config)

	loader := newWithPath("")
	if err := loader.Load(conf); err != nil {
		logger.Criticalf(nil, "%v", err)
		panic(err)
	}
	return conf
}

func Load(path string) (*Config, error) {
	conf := new(Config)

	loader := newWithPath(path)
	if err := loader.Load(conf); err != nil {
		return nil, err
	}
	return conf, nil
}

func MustLoad(path string) *Config {
	conf := new(Config)

	loader := newWithPath(path)
	if err := loader.Load(conf); err != nil {
		logger.Criticalf(nil, "%s: %v", path, err)
		panic(err)
	}

	return conf
}

func (m *Config) JSONString() string {
	return jsonutil.ToString(m)
}

func (m *Config) TOMLString() string {
	buf := new(bytes.Buffer)
	if err := toml.NewEncoder(buf).Encode(m); err != nil {
		logger.Criticalf(nil, "%v", err)
		panic(err)
	}
	return buf.String()
}

func (m *Config) YAMLString() string {
	data, err := yaml.Marshal(m)
	if err != nil {
		logger.Criticalf(nil, "%v", err)
		panic(err)
	}
	return string(data)
}

func newWithPath(path string) *multiconfig.DefaultLoader {
	var loaders []multiconfig.Loader

	// Read default values defined via tag fields "default"
	loaders = append(loaders, &multiconfig.TagLoader{})

	// Choose what while is passed
	if strings.HasSuffix(path, ".toml") {
		loaders = append(loaders, &multiconfig.TOMLLoader{Path: path})
	}

	if strings.HasSuffix(path, ".json") {
		loaders = append(loaders, &multiconfig.JSONLoader{Path: path})
	}

	if strings.HasSuffix(path, ".yml") || strings.HasSuffix(path, ".yaml") {
		loaders = append(loaders, &multiconfig.YAMLLoader{Path: path})
	}

	loaders = append(loaders, &multiconfig.EnvironmentLoader{
		Prefix:    EnvPrefix,
		CamelCase: true,
	})

	loader := multiconfig.MultiLoader(loaders...)

	d := &multiconfig.DefaultLoader{}
	d.Loader = loader
	d.Validator = multiconfig.MultiValidator(&multiconfig.RequiredValidator{})
	return d
}
