// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package strutil

import (
	"regexp"
	"strings"
)

func NewString(v string) *string {
	return &v
}

func SimplifyStringList(s []string) []string {
	b := s[:0]
	for _, x := range s {
		if x := SimplifyString(x); x != "" {
			b = append(b, x)
		}
	}
	return b
}

var reMoreSpace = regexp.MustCompile(`\s+`)

// "\ta  b  c" => "a b c"
func SimplifyString(s string) string {
	return reMoreSpace.ReplaceAllString(strings.TrimSpace(s), " ")
}
