// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package db

import (
	"fmt"
	"strings"
)

func pkgBuildSql_InsertInto(tableName string, v interface{}) (sql string, values []interface{}) {
	names, values := pkgGetAllDBTableFiledNamesAndValues(v)
	if len(names) == 0 {
		return "", nil
	}

	tableFieldName := strings.Join(names, ",")

	taleFieldValue := func() string { // table field values
		// postgres
		// "$1",
		// "$1, $2"
		// "$1, $2, $3"
		//
		// mysql & sqlite3
		// ?, ?
		var b strings.Builder
		for i := 0; i < len(values); i++ {
			if i == 0 {
				fmt.Fprintf(&b, "?")
			} else {
				fmt.Fprintf(&b, ",?")
			}
		}
		return b.String()
	}()

	// postgres: VALUES ($1, $2)
	// mysql & sqlite3: VALUES (?, ?)
	// db.Exec("INSERT INTO place (country, telcode) VALUES ($1, $2)", "Singapore", "65")
	sql = fmt.Sprintf(
		"INSERT INTO %s (%s) VALUES (%s);",
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
		"DELETE FROM %s WHERE %s IN (%s);",
		tableName, primaryKeyName, primaryKeyValues,
	)
}

func pkgBuildSql_Update(
	tableName string, v interface{}, primaryKeyName string,
) (sql string, values []interface{}) {
	names, values := pkgGetNonZeroDBTableFiledNamesAndValues(v)
	if len(names) == 0 {
		return "", nil
	}

	var allFiledNames = []string{}
	var allFiledValues = []interface{}{}
	var primaryKeyValue interface{}

	for i := 0; i < len(names); i++ {
		if names[i] != primaryKeyName {
			allFiledNames = append(allFiledNames, names[i])
			allFiledValues = append(allFiledValues, values[i])
		} else {
			primaryKeyValue = values[i]
		}
	}
	allFiledNames = append(allFiledNames, primaryKeyName)
	allFiledValues = append(allFiledValues, primaryKeyValue)
	if len(allFiledNames) < 2 {
		return "", nil
	}

	// update user set user_name="user_name", position="position" where user_id="user_id"
	var b strings.Builder
	for i := 0; i < len(allFiledNames)-1; i++ {
		switch {
		case i == 0:
			fmt.Fprintf(&b, "%s = ?", allFiledNames[i])
		case i > 0:
			fmt.Fprintf(&b, ",%s = ?", allFiledNames[i])
		}
	}

	sql = fmt.Sprintf(
		"UPDATE %s SET %s WHERE %s = ?;",
		tableName, b.String(),
		primaryKeyName,
	)

	return sql, allFiledValues
}
