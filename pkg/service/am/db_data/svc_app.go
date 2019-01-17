// Copyright 2019 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package db_data

func init() {
	AllServices = append(AllServices, appService)
}

var appService = Service{
	Name: "AppManager",
	MethodList: []ServiceMethod{
		{
			Name:  "",
			Url:   []string{},
			Verbs: []string{},
		},
	},
}
