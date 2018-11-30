// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package db

import (
	"encoding/json"

	"openpitrix.io/iam/pkg/pb/am"
	"openpitrix.io/logger"
)

func encodeRuleList(x []*pbam.Rule) string {
	if x == nil {
		x = []*pbam.Rule{}
	}

	data, err := json.Marshal(x)
	if err != nil {
		logger.Criticalf(nil, "%v+", err)
	}

	return string(data)
}

func decodeRuleList(jsonValue string) []*pbam.Rule {
	if jsonValue == "" {
		return []*pbam.Rule{}
	}

	var x []*pbam.Rule
	err := json.Unmarshal([]byte(jsonValue), &x)
	if err != nil {
		logger.Criticalf(nil, "%v+", err)
	}

	return x
}
