// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package validator

import (
	"regexp"
	"strings"
	"unicode"
)

func IsValidId(ids ...string) bool {
	var re = regexp.MustCompile(`^[a-zA-Z0-9_-]*$`)
	for _, id := range ids {
		if !re.MatchString(id) {
			return false
		}
	}
	return true
}

func IsValidName(name ...string) bool {
	for _, v := range name {
		if s := strings.TrimSpace(v); s == "" || s != v {
			return false
		}
	}
	return true
}

func IsValidStatus(status ...string) bool {
	var re = regexp.MustCompile(`^[a-zA-Z0-9_-]*$`)
	for _, v := range status {
		if !re.MatchString(v) {
			return false
		}
	}
	return true
}

func IsValidEmail(emails ...string) bool {
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

func IsValidPhoneNumber(phoneNumbers ...string) bool {
	var re = regexp.MustCompile(`^[0-9]+$`)
	for _, id := range phoneNumbers {
		if !re.MatchString(id) {
			return false
		}
	}
	return true
}

func IsValidGroupPath(s string) bool {
	var re = regexp.MustCompile(`^[a-zA-Z0-9_.-]{2,255}$`)
	return re.MatchString(s)
}

func IsValidSearchWord(name string) bool {
	var re = regexp.MustCompile(`^[a-zA-Z0-9_-]*$`)
	return re.MatchString(name)
}

func IsValidSortKey(name string) bool {
	var re = regexp.MustCompile(`^[a-zA-Z0-9_-]*$`)
	return re.MatchString(name)
}
