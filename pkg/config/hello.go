// Copyright 2019 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

// +build ignore

package main

import (
	"fmt"

	"openpitrix.io/iam/pkg/config"
)

func main() {
	fmt.Println(config.Default().JSONString())
}
