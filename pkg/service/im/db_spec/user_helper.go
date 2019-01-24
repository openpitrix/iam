// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package db_spec

func (p *User) AdjustForCreate() error {
	return nil
}
func (p *User) AdjustForUpdate() error {
	return nil
}

func (p *User) IsValidForCreate() error {
	return nil
}
func (p *User) IsValidForUpdate() error {
	return nil
}
