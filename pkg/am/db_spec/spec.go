// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

// OpenPitrix Access Management service database spec package.
package db_spec

var DbSchema = struct {
	DBName   string
	Engine   string
	Encoding string
}{
	DBName:   "openpitrix",
	Engine:   "InnoDB",
	Encoding: "UTF8",
}

var DbTables = map[string][]struct {
	FieldName  string
	FieldType  string
	PrimaryKey bool
	Default    string
	NotNull    bool
	GenIndex   bool
}{
	"role": {
		{
			FieldName:  "name",
			FieldType:  "VARCHAR(50)",
			PrimaryKey: true,
			NotNull:    true,
		},
		{
			FieldName: "rule",
			FieldType: "JSON",
		},
	},

	"role_binding": {
		{
			FieldName:  "name",
			FieldType:  "VARCHAR(50)",
			PrimaryKey: true,
			NotNull:    true,
		},
		{
			FieldName:  "xid",
			FieldType:  "VARCHAR(50)",
			PrimaryKey: true,
			NotNull:    true,
		},
	},
}
