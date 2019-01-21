// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package db

import (
	"fmt"
	"sort"
	"strings"
)

func genLimitOffset(limit, offset int32) string {
	if offset < 0 {
		offset = 0
	}

	if limit <= 0 {
		limit = 20
	}
	if limit > 200 {
		limit = 200
	}

	return fmt.Sprintf("LIMIT %d OFFSET %d", limit, offset)
}

func genOrderBy(sortKey string, reverse bool) string {
	if sortKey == "" {
		return ""
	}
	if reverse {
		return "ORDER BY " + sortKey + " DESC"
	} else {
		return "ORDER BY " + sortKey + " ASC"
	}
}

func genWhereCondition(
	keyFileds map[string][]string,
	searchWordFieldNames []string,
	searchWord string,
) []string {
	m := make(map[string]string)

	// name IN(name1,name2,...)
	for name, values := range keyFileds {
		if len(values) > 0 {
			// GET /api/iam/im/users?uid=user1,user2,user3
			// uid[0] == "user1,user2,user3"
			if len(values) == 1 && strings.Contains(values[0], ",") {
				values = strings.Split(values[0], ",")
			}

			var b strings.Builder
			fmt.Fprintf(&b, "%s IN(", name)
			for i, v := range values {
				if i > 0 {
					fmt.Fprintf(&b, ",'%s'", v)
				} else {
					fmt.Fprintf(&b, "'%s'", v)
				}
			}
			fmt.Fprintf(&b, ")")
			m[name] = b.String()
		}
	}
	if len(m) == 0 && searchWord == "" {
		return nil // no where condition
	}

	if len(m) == 0 {
		var s = "1=0"
		for _, name := range searchWordFieldNames {
			s += fmt.Sprintf(
				" OR %s LIKE '%%%s%%'", name, searchWord,
			)

		}
		return []string{s}
	}

	// name LIKE '%search_word%'
	var conditionBySearchWord = `(1!=0 `
	if searchWord != "" {
		for _, name := range searchWordFieldNames {
			if _, exists := m[name]; !exists {
				conditionBySearchWord += fmt.Sprintf(
					"OR %s LIKE '%%%s%%'", name, searchWord,
				)
			}
		}
	}
	conditionBySearchWord += ")"

	var ss []string
	for _, v := range m {
		ss = append(ss, v)
	}
	sort.Strings(ss)
	ss = append(ss, conditionBySearchWord)
	return ss
}
