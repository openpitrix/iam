// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

// OpenPitrix IAM service app.
package main

import (
	"fmt"
	"os"

	"github.com/urfave/cli"

	"openpitrix.io/iam/openpitrix/pkg/config"
	"openpitrix.io/iam/openpitrix/pkg/service"
	"openpitrix.io/iam/openpitrix/pkg/version"
	"openpitrix.io/logger"
)

func Main() {
	app := cli.NewApp()
	app.Name = "iam"
	app.Usage = "iam provides iam service."
	app.Version = version.GetVersionString()

	app.UsageText = `im [global options] command [options] [args...]

EXAMPLE:
   iam gen-config
   iam info
   iam list
   iam ping
   iam getv key
   iam serve
   iam tour`

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   "config",
			Value:  "iam-config.json",
			Usage:  "iam config file",
			EnvVar: "OPENPITRIX_IAM_CONFIG",
		},
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
			Name:  "serve",
			Usage: "run as service",
			Action: func(c *cli.Context) {
				cfg := config.MustLoad(c.GlobalString("config"))
				if cfg.TlsEnabled {
					server, err := service.OpenServer(cfg.Mysql.DbType(), cfg.Mysql.GetUrl())
					if err != nil {
						logger.Criticalf(nil, "%v", err)
						os.Exit(1)
					}
					err = server.ListenAndServe(fmt.Sprintf(":%d", cfg.Port))
					if err != nil {
						logger.Criticalf(nil, "%v", err)
						os.Exit(1)
					}
				} else {
					server, err := service.OpenServer(cfg.Mysql.DbType(), cfg.Mysql.GetUrl())
					if err != nil {
						logger.Criticalf(nil, "%v", err)
						os.Exit(1)
					}
					err = server.ListenAndServeTLS(
						fmt.Sprintf(":%d", cfg.Port),
						cfg.TlsCertFile,
						cfg.TlsKeyFile,
					)
					if err != nil {
						logger.Criticalf(nil, "%v", err)
						os.Exit(1)
					}
				}
			},
		},
	}

	app.CommandNotFound = func(ctx *cli.Context, command string) {
		fmt.Fprintf(ctx.App.Writer, "not found '%v'!\n", command)
	}

	app.Run(os.Args)
}
