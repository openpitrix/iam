// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

// copy from https://godoc.org/github.com/golang/protobuf/protoc-gen-go/generator

package camel

import (
	"testing"
)

func TestCamelCase(t *testing.T) {
	tests := []struct {
		in, want string
	}{
		{"one", "One"},
		{"one_two", "OneTwo"},
		{"_my_field_name_2", "XMyFieldName_2"},
		{"_my_field_name_2_ext", "XMyFieldName_2Ext"},
		{"Something_Capped", "Something_Capped"},
		{"my_Name", "My_Name"},
		{"OneTwo", "OneTwo"},
		{"_", "X"},
		{"_a_", "XA_"},
	}
	for _, tc := range tests {
		if got := CamelCase(tc.in); got != tc.want {
			t.Errorf("CamelCase(%q) = %q, want %q", tc.in, got, tc.want)
		}
	}
}
