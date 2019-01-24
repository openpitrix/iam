// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package strutil

import (
	"strings"
)

func NewString(v string) *string {
	return &v
}

func SimplifyStringList(s []string) []string {
	b := s[:0]
	for _, x := range s {
		if x := strings.TrimSpace(x); x != "" {
			b = append(b, x)
		}
	}
	return b
}

// "\ta  b  c" => "a b c"
func SimplifyString(s string) string {
	s = strings.Replace(s, "\t", " ", -1)
	s = strings.Replace(s, "\n", " ", -1)
	s = strings.Replace(s, "\r", " ", -1)
	s = strings.TrimSpace(s)

	for {
		if sx := strings.Replace(s, "  ", " ", -1); sx == s {
			return s
		} else {
			s = sx
		}
	}
}
