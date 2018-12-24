// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package service

import (
	"fmt"
	"reflect"
	"strings"
)

func pkgBuildSql_InsertInto(v interface{}) (sql string, values []interface{}) {
	names, values := pkgGetTableFiledNamesAndValues(v)
	if len(names) == 0 {
		return "", nil
	}

	tableName := func() string { // table name
		if t := reflect.TypeOf(v); t.Kind() == reflect.Ptr {
			return strings.ToLower(t.Elem().Name())
		} else {
			return strings.ToLower(t.Name())
		}
	}()

	tableFieldName := strings.Join(names, ",")

	taleFieldValue := func() string { // table field values
		// "$1",
		// "$1, $2"
		// "$1, $2, $3"
		var b strings.Builder
		for i := 0; i < len(values); i++ {
			if i == 0 {
				fmt.Fprintf(&b, "$%d", i)
			} else {
				fmt.Fprintf(&b, ",$%d", i)
			}
		}
		return b.String()
	}()

	// db.Exec("INSERT INTO place (country, telcode) VALUES ($1, $2)", "Singapore", "65")
	sql = fmt.Sprintf(
		"INSERT INTO %s (%s) VALUES (%s)",
		tableName,
		tableFieldName,
		taleFieldValue,
	)

	return
}
