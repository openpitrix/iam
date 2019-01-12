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
			m[name] = fmt.Sprintf(
				"%s IN(%s)", name, strings.Join(values, ","),
			)
		}
	}

	// name LIKE '%search_word%'
	if searchWord != "" {
		for _, name := range searchWordFieldNames {
			if _, exists := m[name]; !exists {
				m[name] = fmt.Sprintf(
					"%s LIKE '%%%s%%'", name, searchWord,
				)
			}
		}
	}
	if len(m) == 0 {
		return nil // no where condition
	}

	var ss []string
	for _, v := range m {
		ss = append(ss, v)
	}
	sort.Strings(ss)
	return ss
}
