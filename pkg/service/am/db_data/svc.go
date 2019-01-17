// Copyright 2019 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package db_data

type Service struct {
	Name       string
	MethodList []ServiceMethod
}

type ServiceMethod struct {
	Name  string
	Url   []string
	Verbs []string
}

var AllServices = []Service{}
