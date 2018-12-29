// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

// OpenPitrix IAM service app.
package main

import (
	"fmt"
	"io/ioutil"
	"net"
	"os"

	"github.com/urfave/cli"

	"openpitrix.io/iam/openpitrix/pkg/config"
	"openpitrix.io/iam/openpitrix/pkg/service"
	"openpitrix.io/iam/openpitrix/pkg/version"
	"openpitrix.io/logger"
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
			Value:  "localhost",
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
			data := config.Default().JSONString()
			ioutil.WriteFile(cfgpath, []byte(data), 0666)
		}

		cfg := config.MustLoad(c.GlobalString("config"))
		logger.SetLevelByString(cfg.LogLevel)
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
	cfg := config.MustLoad(c.GlobalString("config"))

	host := c.GlobalString("host")
	if host == "" || host == "localhost" {
		host = getLocalIP()
	}

	readme := c.GlobalString("readme")

	if !cfg.TlsEnabled {
		logger.Infof(nil, version.GetVersionString())
		logger.Infof(nil, "IAM service http://%s:%d", host, cfg.Port)

		server, err := service.OpenServer(cfg.DB.Type, cfg.DB.GetUrl(), readme)
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
		logger.Infof(nil, version.GetVersionString())
		logger.Infof(nil, "IAM service https://%s:%d", host, cfg.Port)

		server, err := service.OpenServer(cfg.DB.Type, cfg.DB.GetUrl(), readme)
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
