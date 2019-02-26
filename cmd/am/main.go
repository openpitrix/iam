// Copyright 2019 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package main

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/urfave/cli"

	"openpitrix.io/iam/pkg/config"
	"openpitrix.io/iam/pkg/service/am"
	"openpitrix.io/iam/pkg/version"
	"openpitrix.io/logger"
)

var (
	appConfig *config.Config = nil
)

func main() {
	app := cli.NewApp()
	app.Name = "am"
	app.Usage = "provide am service."
	app.Version = version.GetVersionString()

	app.UsageText = `am [global options] command [options] [args...]

EXAMPLE:
   am gen-config
   am serve`

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   "config",
			Value:  "config-am.json",
			Usage:  "am config file",
			EnvVar: "am_CONFIG",
		},
		cli.StringFlag{
			Name:   "host",
			Value:  "am-service",
			EnvVar: "am_HOST",
		},
	}

	app.Before = func(c *cli.Context) error {
		cfgPath := c.GlobalString("config")
		if _, err := os.Stat(cfgPath); os.IsNotExist(err) {
			appConfig = config.Default()
			ioutil.WriteFile(cfgPath, []byte(appConfig.JSONString()), 0666)
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
	am.Serve(appConfig)
}
