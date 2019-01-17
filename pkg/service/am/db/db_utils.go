// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package db

import (
	"fmt"
	"strings"
)

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
