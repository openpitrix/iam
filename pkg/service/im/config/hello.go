// Copyright 2019 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

// +build ignore

package main

import (
	"flag"
	"fmt"

	"openpitrix.io/iam/pkg/config"
)

var (
	flagConfFile = flag.String("conf", "", "config file")
)

func main() {
	flag.Parse()
	cfg, _ := config.Load(*flagConfFile)
	fmt.Println(cfg.JSONString())
}
