// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package im

import (
	"regexp"

	"github.com/fatih/structs"
)

func pkgGetDBTableFiledNames(v interface{}) (names []string) {
	for _, f := range structs.Fields(v) {
		if !f.IsExported() || f.Tag("db") == "-" {
			continue
		}

		var fieldName = f.Name()
		if s := f.Tag("db"); s != "" {
			fieldName = s
		}

		names = append(names, fieldName)
	}
	return
}
func pkgGetNonZeroDBTableFiledNamesAndValues(v interface{}) (names []string, values []interface{}) {
	for _, f := range structs.Fields(v) {
		if !f.IsExported() || f.IsZero() || f.Tag("db") == "-" {
			continue
		}

		var (
			fieldName  = f.Name()
			fieldValue = f.Value()
		)
		if s := f.Tag("db"); s != "" {
			fieldName = s
		}

		names = append(names, fieldName)
		values = append(values, fieldValue)
	}
	return
}
func pkgGetAllDBTableFiledNamesAndValues(v interface{}) (names []string, values []interface{}) {
	for _, f := range structs.Fields(v) {
		if !f.IsExported() || f.Tag("db") == "-" {
			continue
		}

		var (
			fieldName  = f.Name()
			fieldValue = f.Value()
		)
		if s := f.Tag("db"); s != "" {
			fieldName = s
		}

		names = append(names, fieldName)
		values = append(values, fieldValue)
	}
	return
}

func pkgGetDBTableStringFieldNames(v interface{}) (names []string) {
	for _, f := range structs.Fields(v) {
		if !f.IsExported() || f.Tag("db") == "-" {
			continue
		}

		var (
			fieldName  = f.Name()
			fieldValue = f.Value()
		)
		if s := f.Tag("db"); s != "" {
			fieldName = s
		}

		if _, ok := fieldValue.(string); ok {
			names = append(names, fieldName)
		}
	}
	return
}

var reSearchWord = regexp.MustCompile(`^[a-z0-9-_]+$`)

func pkgSearchWordValid(s string) bool {
	return reSearchWord.MatchString(s)
}
