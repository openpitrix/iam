// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package service

import (
	"database/sql"
	"strings"
	"time"

	"github.com/chai2010/spacestring"
	"github.com/fatih/structs"
	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/timestamp"
)

func pkgSqlScanProtoMessge(rows *sql.Rows, msg proto.Message) error {
	cols, _ := rows.Columns()

	// Create a slice of interface{}'s to represent each column,
	// and a second slice to contain pointers to each item in the columns slice.
	columns := make([]interface{}, len(cols))
	columnPointers := make([]interface{}, len(cols))
	for i, _ := range columns {
		columnPointers[i] = &columns[i]
	}

	// Scan the result into the column pointers...
	if err := rows.Scan(columnPointers...); err != nil {
		return err
	}

	// fill struct
	structs := structs.New(msg)
	for i, colName := range cols {
		if f, ok := structs.FieldOk(colName); ok {
			switch f.Value().(type) {
			case int:
				if v, ok := columns[i].(int); ok {
					f.Set(v)
				}
			case int32:
				if v, ok := columns[i].(int32); ok {
					f.Set(v)
				}
			case int64:
				if v, ok := columns[i].(int64); ok {
					f.Set(v)
				}
			case string:
				if v, ok := columns[i].(string); ok {
					f.Set(v)
				}
			case *timestamp.Timestamp:
				if t, ok := columns[i].(time.Time); ok {
					if t, err := ptypes.TimestampProto(t); err == nil {
						f.Set(&t)
					}
				}
			}
		}
	}

	return nil
}

func pkgSqlScanMap(rows *sql.Rows) (map[string]interface{}, error) {
	cols, _ := rows.Columns()

	// Create a slice of interface{}'s to represent each column,
	// and a second slice to contain pointers to each item in the columns slice.
	columns := make([]interface{}, len(cols))
	columnPointers := make([]interface{}, len(cols))
	for i, _ := range columns {
		columnPointers[i] = &columns[i]
	}

	// Scan the result into the column pointers...
	if err := rows.Scan(columnPointers...); err != nil {
		return nil, err
	}

	// Create our map, and retrieve the value for each column from the pointers slice,
	// storing it in the map with the name of the column as the key.
	m := make(map[string]interface{})
	for i, colName := range cols {
		val := columnPointers[i].(*interface{})
		m[colName] = *val
	}

	return m, nil
}

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
