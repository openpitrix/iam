// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package rbac

import (
	"github.com/bmatcuk/doublestar"

	"openpitrix.io/iam/pkg/pb/am"
)

func matchRule(method, namespace string, rule *pbam.ActionRule) bool {
	// check method
	if ok, _ := doublestar.Match(rule.MethodPattern, method); !ok {
		return false
	}

	// check verb
	var nsMatched bool
	for _, pattern := range rule.NamespacePattern {
		if ok, _ := doublestar.Match(pattern, namespace); !ok {
			nsMatched = true
			return false
		}
	}
	if !nsMatched {
		return false
	}

	// OK
	return true
}

func canDoAction(x *pbam.Action, rules []*pbam.ActionRule) bool {
	for _, rule := range rules {
		if matchRule(x.Method, x.Namespace, rule) {
			return true
		}
	}
	return false
}
