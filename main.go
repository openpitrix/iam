// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

// Identity and Access Management System for OpenPitrix.
package main

import (
	"fmt"
	"io/ioutil"
	"net"
	"os"

	"github.com/urfave/cli"

	"openpitrix.io/iam/pkg/config"
	"openpitrix.io/iam/pkg/service/im"
	"openpitrix.io/iam/pkg/service/web"
	"openpitrix.io/iam/pkg/version"
	"openpitrix.io/logger"
)

var (
	appConfig *config.Config = nil
)

func main() {
	app := cli.NewApp()
	app.Name = "iam"
	app.Usage = "iam provides iam service."
	app.Version = version.GetVersionString()

	app.UsageText = `im [global options] command [options] [args...]

EXAMPLE:
   iam gen-config
   iam serve`

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   "config",
			Value:  "config.json",
			Usage:  "iam config file",
			EnvVar: "OPENPITRIX_IAM_CONFIG",
		},
		cli.StringFlag{
			Name:   "host",
			Value:  "openpitrix-iam2-service",
			EnvVar: "OPENPITRIX_IAM_HOST",
		},

		cli.StringFlag{
			Name:  "readme",
			Value: "_readme.md",
		},
	}

	app.Before = func(c *cli.Context) error {
		cfgpath := c.GlobalString("config")
		if _, err := os.Stat(cfgpath); os.IsNotExist(err) {
			appConfig = config.Default()
			ioutil.WriteFile(cfgpath, []byte(appConfig.JSONString()), 0666)
		} else {
			appConfig = config.MustLoad(c.GlobalString("config"))
		}

		logger.SetLevelByString(appConfig.LogLevel)
		return nil
	}

	app.Action = func(c *cli.Context) {
		serve(c)
	}

	app.Commands = []cli.Command{
		{
			Name:   "debug",
			Usage:  "debug app",
			Hidden: true,

			Action: func(c *cli.Context) {
				fmt.Println(nil)
				fmt.Println(version.GetVersion())
				return
			},
		},
		{
			Name:  "gen-config",
			Usage: "gen config file",
			Flags: []cli.Flag{
				cli.BoolFlag{
					Name:  "json",
					Usage: "json format (default)",
				},
				cli.BoolFlag{
					Name:  "toml",
					Usage: "toml format",
				},
				cli.BoolFlag{
					Name:  "yaml",
					Usage: "yaml format",
				},
			},

			Action: func(c *cli.Context) {
				switch {
				case c.Bool("json"):
					fmt.Println(config.Default().JSONString())
				case c.Bool("toml"):
					fmt.Println(config.Default().TOMLString())
				case c.Bool("yaml"):
					fmt.Println(config.Default().YAMLString())
				default:
					fmt.Println(config.Default().JSONString())
				}
				return
			},
		},

		{
			Name:  "serve",
			Usage: "run as service",
			Action: func(c *cli.Context) {
				serve(c)
			},
		},
	}

	app.CommandNotFound = func(ctx *cli.Context, command string) {
		fmt.Fprintf(ctx.App.Writer, "not found '%v'!\n", command)
	}

	app.Run(os.Args)
}

func serve(c *cli.Context) {
	imService, err := im.OpenServer(appConfig)
	if err != nil {
		logger.Criticalf(nil, "%v", err)
		os.Exit(1)
	}

	if !appConfig.TlsEnabled {
		logger.Infof(nil, version.GetVersionString())
		logger.Infof(nil, "IAM service http://%s:%d", appConfig.Host, appConfig.Port)

		err := web.ListenAndServe(
			fmt.Sprintf(":%d", appConfig.Port),
			[]web.GrpcServer{web.WithAccountManager(imService)},
			nil,
		)
		if err != nil {
			logger.Criticalf(nil, "%v", err)
			os.Exit(1)
		}
	} else {
		logger.Infof(nil, version.GetVersionString())
		logger.Infof(nil, "IAM service https://%s:%d", appConfig.Host, appConfig.Port)

		err := web.ListenAndServeTLS(
			fmt.Sprintf(":%d", appConfig.Port),
			appConfig.TlsCertFile, appConfig.TlsKeyFile,
			[]web.GrpcServer{web.WithAccountManager(imService)},
			nil,
		)
		if err != nil {
			logger.Criticalf(nil, "%v", err)
			os.Exit(1)
		}
	}
}

func getLocalIP() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return "127.0.0.1"
	}
	for _, address := range addrs {
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String()
			}
		}
	}
	return "127.0.0.1"
}
