// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package service

import (
	"fmt"
	"strings"
)

func pkgBuildSql_InsertInto(tableName string, v interface{}) (sql string, values []interface{}) {
	names, values := pkgGetTableFiledNamesAndValues(v)
	if len(names) == 0 {
		return "", nil
	}

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

func pkgBuildSql_Delete(tableName, primaryKeyName string, key ...string) (sql string) {
	primaryKeyValues := func() string {
		var b strings.Builder
		for i := 0; i < len(key); i++ {
			if i == 0 {
				fmt.Fprintf(&b, `"%s"`, key[i])
			} else {
				fmt.Fprintf(&b, `,"%s"`, key[i])
			}
		}
		return b.String()
	}()

	// delete * from group where group_id in ("group1","group2")
	return fmt.Sprintf(
		"DELETE * FROM %s WHERE %s IN (%s)",
		tableName, primaryKeyName, primaryKeyValues,
	)
}
