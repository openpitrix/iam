// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

// OpenPitrix Identity Management service app.
package main

import (
	"context"
	"fmt"
	"os"

	"github.com/urfave/cli"

	"openpitrix.io/iam/pkg/am/client"
	"openpitrix.io/iam/pkg/am/config"
	"openpitrix.io/iam/pkg/am/service"
	"openpitrix.io/iam/pkg/pb/am"
	"openpitrix.io/iam/pkg/version"
	"openpitrix.io/logger"
)

func main() {
	app := cli.NewApp()
	app.Name = "am"
	app.Usage = "am provides am service."
	app.Version = version.GetVersionString()

	app.UsageText = `am [global options] command [options] [args...]

EXAMPLE:
   am info
   am can-do

   am list-role
   am list-role-binding

   am create-role
   am modifu-role
   am delete-role

   am create-binding
   am delete-binding

   am load-rbac
   am save-rbac

   am serve`

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   "config",
			Value:  "config.toml",
			Usage:  "am config file",
			EnvVar: "OPENPITRIX_AM_CONFIG",
		},
	}

	app.Commands = []cli.Command{
		{
			Name:   "debug",
			Usage:  "debug app",
			Hidden: true,

			Action: func(c *cli.Context) {
				fmt.Println(version.GetVersion())
				return
			},
		},

		{
			Name:  "info",
			Usage: "show config",
			Flags: []cli.Flag{
				cli.BoolFlag{
					Name:  "json",
					Usage: "json format",
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
				cfg := config.MustLoad(c.GlobalString("config"))
				switch {
				case c.Bool("json"):
					fmt.Println(cfg.JSONString())
				case c.Bool("toml"):
					fmt.Println(cfg.TOMLString())
				case c.Bool("yaml"):
					fmt.Println(cfg.YAMLString())
				default:
					fmt.Println(cfg.JSONString())
				}
				return
			},
		},

		{
			Name:      "can-do",
			Usage:     "can do action",
			ArgsUsage: "Serve.Method /path/to/home",
			Flags: []cli.Flag{
				cli.IntFlag{
					Name:  "port",
					Usage: "set host port",
					Value: 0,
				},
				cli.StringSliceFlag{
					Name:  "xid",
					Usage: "xid list",
					Value: &cli.StringSlice{"user1", "group1"},
				},
				cli.StringSliceFlag{
					Name:  "role",
					Usage: "role name",
					Value: &cli.StringSlice{"role1"},
				},
			},
			Action: func(c *cli.Context) {
				if c.NArg() < 2 {
					fmt.Println("args: missing Verb or Path")
					return
				}
				cfg := config.MustLoad(c.GlobalString("config"))

				port := c.Int("port")
				if port <= 0 {
					port = cfg.Port
				}

				client, conn, err := client.DialService(c.String("host"), port)
				if err != nil {
					logger.Criticalf(nil, "%v", err)
					os.Exit(1)
				}
				defer conn.Close()

				reply, err := client.CanDo(context.Background(), &pbam.Action{
					RoleName:  c.StringSlice("role"),
					Xid:       c.StringSlice("xid"),
					Method:    c.Args().First(),
					Namespace: c.Args().Get(1),
				})
				if err != nil {
					logger.Criticalf(nil, "%v", err)
					os.Exit(1)
				}

				fmt.Println(reply.GetValue())
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
