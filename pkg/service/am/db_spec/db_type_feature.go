// Copyright 2019 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package db_spec

type DBFeature struct {
	ModuleId    string `db:"module_id"`
	FeatureId   string `db:"feature_id"`
	FeatureName string `db:"feature_name"`
}
