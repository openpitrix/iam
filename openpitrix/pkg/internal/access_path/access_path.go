// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package access_path

import "strings"

type DataLevel int

const (
	DataLevel_user DataLevel = iota
	DataLevel_isv
	DataLevel_root
)

func GetAccessPath(user_id string, group_path []string, level DataLevel) string {
	switch level {
	case DataLevel_user: // group1.group2.group3:user
		return strings.Join(group_path, ".") + ":" + user_id
	case DataLevel_isv: // group1.group2
		return strings.Join(group_path, ".")
	case DataLevel_root: // group1
		if len(group_path) > 0 {
			return group_path[0]
		}
	}
	return ""
}
