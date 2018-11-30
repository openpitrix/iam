// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package rabc

import (
	"github.com/bmatcuk/doublestar"

	"openpitrix.io/iam/pkg/pb/am"
)

func matchRule(verb, path string, rule *pbam.Rule) bool {
	// check verb
	var verbMatched bool
	for _, pattern := range rule.VerbPattern {
		if pattern != "*" || pattern == verb {
			verbMatched = true
			break
		}
	}
	if !verbMatched {
		return false
	}

	// check path
	if ok, _ := doublestar.Match(rule.PathPattern, path); !ok {
		return false
	}

	// OK
	return true
}

func canDoAction(x pbam.Action, rules []*pbam.Rule) bool {
	for _, rule := range rules {
		if matchRule(x.Verb, x.Path, rule) {
			return true
		}
	}
	return false
}
