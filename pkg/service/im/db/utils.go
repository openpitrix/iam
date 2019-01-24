// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package db

/*

import (
	"crypto/rand"
	"regexp"
	"strings"
	"unicode"

	"github.com/golang/protobuf/ptypes/timestamp"

	"openpitrix.io/iam/pkg/internal/base58"
)

func newString(v string) *string {
	return &v
}

func genId(prefix string, maxLen int) string {
	if prefix == "" {
		prefix = "xid-"
	}
	if maxLen <= 0 {
		maxLen = 12
	}

	if maxLen <= len(prefix) {
		maxLen += len(prefix)
	}

	buf := make([]byte, maxLen-len(prefix))
	rand.Read(buf)
	s := string(base58.EncodeBase58(buf))
	return prefix + s[:len(buf)]
}

func isValidIds(ids ...string) bool {
	var re = regexp.MustCompile(`^[a-zA-Z0-9_-]*$`)
	for _, id := range ids {
		if !re.MatchString(id) {
			return false
		}
	}
	return true
}

func isValidGroupPath(s string) bool {
	var re = regexp.MustCompile(`^[a-zA-Z0-9_.-]{2,255}$`)
	return re.MatchString(s)
}
func isValidSearchWord(name string) bool {
	var re = regexp.MustCompile(`^[a-zA-Z0-9_-]*$`)
	return re.MatchString(name)
}
func isValidSortKey(name string) bool {
	var re = regexp.MustCompile(`^[a-zA-Z0-9_-]*$`)
	return re.MatchString(name)
}
func isValidEmails(emails ...string) bool {
	for _, v := range emails {
		if v == "" {
			return false
		}
		if idx := strings.IndexByte(v, '@'); idx <= 0 || idx >= len(v) {
			return false
		}
		if strings.Count(v, "@") != 1 {
			return false
		}
		for _, c := range v {
			if unicode.IsSpace(c) {
				return false
			}
		}
	}
	return true
}
func isValidPhoneNumbers(phoneNumbers ...string) bool {
	var re = regexp.MustCompile(`^[0-9]+$`)
	for _, id := range phoneNumbers {
		if !re.MatchString(id) {
			return false
		}
	}
	return true
}

func isZeroTimestamp(x *timestamp.Timestamp) bool {
	if x == nil {
		return true
	}
	if x.Seconds == 0 && x.Nanos == 0 {
		return true
	}
	return false
}

func simplifyStringList(s []string) []string {
	b := s[:0]
	for _, x := range s {
		if x := strings.TrimSpace(x); x != "" {
			b = append(b, x)
		}
	}
	return b
}

// "\ta  b  c" => "a b c"
func simplifyString(s string) string {
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

*/
