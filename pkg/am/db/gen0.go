// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

// +build ignore

package main

import (
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"strings"

	"openpitrix.io/iam/pkg/pb/am"
)

var (
	flagOutput = flag.String("outputx", "gen0.output.go", "set output file")
)

func main() {
	flag.Parse()

	var buf bytes.Buffer
	fmt.Fprintln(&buf, `// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

// Code generated. DO NOT EDIT.

package db

import (
	"github.com/golang/protobuf/proto"

	"openpitrix.io/iam/pkg/pb/am"
)
	`)

	fmt.Fprintln(&buf, "// Reference imports to suppress errors if they are not otherwise used.")
	fmt.Fprintln(&buf, "var _ pbam.DbSchema")
	fmt.Fprintln(&buf, "var _ proto.Message")
	fmt.Fprintln(&buf)

	fmt.Fprintln(&buf, "var DbSchemaTableList = []proto.Message{")
	for _, tname := range ReadTableList(new(pbam.DbSchema).GetTables()) {
		fmt.Fprintf(&buf, "\tnew(pbam.%s),\n", tname)
	}
	fmt.Fprintln(&buf, "}")
	fmt.Fprintln(&buf)

	err := ioutil.WriteFile(*flagOutput, buf.Bytes(), 0666)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("gen openpitrix.io/iam/pkg/am/db/%s ok\n", *flagOutput)
}

func ReadTableList(schema string) (tables []string) {
	for _, s := range strings.Split(strings.Trim(schema, ",; "), ",") {
		if s = strings.TrimSpace(s); s != "" {
			tables = append(tables, s)
		}
	}
	return
}
