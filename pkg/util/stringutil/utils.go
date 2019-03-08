// Copyright 2019 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package stringutil

import (
	"regexp"
	"strings"
	"unicode/utf8"
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

func Contains(ss []string, s string) bool {
	for _, v := range ss {
		if v == s {
			return true
		}
	}
	return false
}

func Reverse(s string) string {
	size := len(s)
	buf := make([]byte, size)
	for start := 0; start < size; {
		r, n := utf8.DecodeRuneInString(s[start:])
		start += n
		utf8.EncodeRune(buf[size-start:], r)
	}
	return string(buf)
}

func Merge(left, right []string) []string {
	if len(left) == 0 {
		return right
	} else {
		var merge []string
		for _, s := range left {
			if Contains(right, s) {
				merge = append(merge, s)
			}
		}
		return merge
	}
}
