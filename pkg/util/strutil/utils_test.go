// Copyright 2019 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package strutil

import (
	"testing"

	. "openpitrix.io/iam/pkg/util/assert"
)

func TestSimplifyStringList(t *testing.T) {
	s0 := []string{"a", "", "c", "  ", " d "}
	s1 := SimplifyStringList(s0)

	Assert(t, len(s1) == 3)
	Assert(t, s1[0] == "a")
	Assert(t, s1[1] == "c")
	Assert(t, s1[2] == "d")
}

func TestSimplifyString(t *testing.T) {
	var tests = []struct{ s, expect string }{
		{s: "\ta  b  c", expect: "a b c"},
		{s: "a b c", expect: "a b c"},
		{s: "abc", expect: "abc"},
	}
	for _, v := range tests {
		got := SimplifyString(v.s)
		Assertf(t, got == v.expect, "expect = %q, got = %q", v.expect, got)
	}
}
