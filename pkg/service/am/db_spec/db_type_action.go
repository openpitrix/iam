// Copyright 2019 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package db_spec

type DBAction struct {
	ActionId    string `db:"action_id"`
	FeatureId   string `db:"feature_id"`
	Method      string `db:"method"`
	Description string `db:"description"`
	Url         string `db:"url"`
	UrlMethod   string `db:"url_method"`
}
