// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package service

import (
	"strings"

	"github.com/chai2010/spacestring"
	"github.com/fatih/structs"
	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/timestamp"
)

func pkgGetTableFiledNamesAndValues(v interface{}) (names []string, values []interface{}) {
	s := structs.New(v)
	for _, f := range s.Fields() {
		if !f.IsExported() || f.IsZero() {
			continue
		}
		if strings.HasPrefix(f.Name(), "XXX_") || f.Tag("json") == "-" {
			continue
		}

		var (
			db_field_name  = f.Name()
			db_field_value = f.Value()
		)

		if idx := strings.Index(f.Tag("json"), ","); idx > 0 {
			db_field_name = f.Tag("json")[:idx]
		}

		switch v := db_field_value.(type) {
		case string: // support space string
			if spacestring.IsSpace(v) {
				db_field_value = "" // clear field
			}
		case *timestamp.Timestamp:
			if v != nil {
				db_field_value, _ = ptypes.Timestamp(v)
			}
		}

		names = append(names, db_field_name)
		values = append(values, db_field_value)
	}
	return
}
