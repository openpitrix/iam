// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package db

import (
	"fmt"
	"strings"
)

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
