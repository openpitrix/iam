// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package db

import gorp "gopkg.in/gorp.v2"

type gorpHooker interface {
	gorp.HasPostGet

	gorp.HasPreInsert
	gorp.HasPreUpdate
	gorp.HasPreDelete
}
