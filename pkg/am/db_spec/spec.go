// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

// OpenPitrix Access Management service database spec package.
package db_spec

type DbSchemaType struct {
	DBName   string
	Engine   string
	Encoding string
}

type DbTableSchemaType struct {
	FieldName     string
	FieldType     string
	FieldSize     int
	PrimaryKey    bool
	AutoIncrement bool
	Default       string
	NotNull       bool
	GenIndex      bool
}

var DbSchema = DbSchemaType{
	DBName:   "openpitrix",
	Engine:   "InnoDB",
	Encoding: "UTF8",
}

var DbTableSchema = map[string][]DbTableSchemaType{
	"role": {
		{
			FieldName:  "name",
			FieldType:  "VARCHAR",
			FieldSize:  50,
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
			FieldType:  "VARCHAR",
			FieldSize:  50,
			PrimaryKey: true,
			NotNull:    true,
		},
		{
			FieldName:  "xid",
			FieldType:  "VARCHAR",
			FieldSize:  50,
			PrimaryKey: true,
			NotNull:    true,
		},
	},
}
