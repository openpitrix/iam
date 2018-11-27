// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

// +build ignore

package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"sort"
	"strings"

	"openpitrix.io/iam/pkg/am/db_spec"
	"openpitrix.io/iam/pkg/internal/camel"
)

var (
	flagOutput = flag.String("outputx", "gen.output.go", "set output file")
)

func main() {
	flag.Parse()

	var buf bytes.Buffer
	fmt.Fprintln(&buf, `// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

// Code generated. DO NOT EDIT.

package db
`)

	for _, name := range GetTableNameList() {
		GenTableByName(&buf, name)
	}

	err := ioutil.WriteFile(*flagOutput, buf.Bytes(), 0666)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("gen openpitrix.io/iam/pkg/am/db/%s ok\n", *flagOutput)
}

func GetTableNameList() (names []string) {
	for k, _ := range db_spec.DbTableSchema {
		names = append(names, k)
	}
	sort.Strings(names)
	return
}

func GenTableByName(w io.Writer, name string) {
	fmt.Fprintf(w, "type %s struct {\n", camel.CamelCase(name))
	defer fmt.Fprintln(w, "}\n")

	for _, v := range db_spec.DbTableSchema[name] {
		fmt.Fprintf(w, "\t%s %s `%s`\n", camel.CamelCase(v.FieldName),
			goTypeString(v), gorpTagString(v),
		)
	}
}

func goTypeString(v db_spec.DbTableSchemaType) string {
	// int
	if strings.HasPrefix(v.FieldType, "INT") {
		return "int64"
	}

	// string
	if strings.HasPrefix(v.FieldType, "CHAR") {
		return "string"
	}
	if strings.HasPrefix(v.FieldType, "VARCHAR") {
		return "string"
	}
	if strings.HasPrefix(v.FieldType, "TEXT") {
		return "string"
	}
	if strings.HasPrefix(v.FieldType, "JSON") {
		return "string"
	}

	// timestamp
	if strings.HasPrefix(v.FieldType, "TIMESTAMP") {
		return "time.Time"
	}

	// unknown
	return "???"
}

func gorpTagString(v db_spec.DbTableSchemaType) string {
	var tag string = v.FieldName

	if v.FieldType == "VARCHAR" {
		tag += fmt.Sprintf(", size:%d", v.FieldSize)
	}

	if v.PrimaryKey {
		tag += ", primarykey"
	}
	if v.AutoIncrement {
		tag += ", autoincrement"
	}

	return fmt.Sprintf(`db:"%s"`, tag)
}
