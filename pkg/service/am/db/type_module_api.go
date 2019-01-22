// Copyright 2019 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package db

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"sort"
	"strings"

	pbam "openpitrix.io/iam/pkg/pb/am"
)

type ModuleApiInfoList []ModuleApiInfo

// keep same pbam.Action
type ModuleApiInfo struct {
	RoleId   string
	RoleName string
	Portal   string

	ModuleId            string
	ModuleName          string
	DataLevel           string
	IsFeatureAllChecked string

	FeatureId   string
	FeatureName string

	ActionId      string
	ActionName    string
	ActionEnabled string

	ApiId          string
	ApiMethod      string
	ApiDescription string

	Url       string
	UrlMethod string
}

func NewModuleApiInfoFromPB(p *pbam.Action) *ModuleApiInfo {
	if p == nil {
		return new(ModuleApiInfo)
	}

	var buf bytes.Buffer
	if err := gob.NewEncoder(&buf).Encode(p); err != nil {
		// return nil, err
	}

	var q = new(ModuleApiInfo)
	if err := gob.NewDecoder(bytes.NewBuffer(buf.Bytes())).Decode(q); err != nil {
		// return nil, err
	}

	return q
}

func (p *ModuleApiInfo) ToPB() *pbam.Action {
	q, err := p.ToProtoMessage()
	if err != nil {
		panic(err) // unreachable
	}
	return q
}

func (p *ModuleApiInfo) ToProtoMessage() (*pbam.Action, error) {
	var buf bytes.Buffer
	if err := gob.NewEncoder(&buf).Encode(p); err != nil {
		return nil, err
	}

	var q = new(pbam.Action)
	if err := gob.NewDecoder(bytes.NewBuffer(buf.Bytes())).Decode(q); err != nil {
		return nil, err
	}

	return q, nil
}

func NewModuleApiInfoListFromPB(m ...*pbam.RoleModule) (records ModuleApiInfoList) {
	for _, v0 := range m {
		for _, v1 := range v0.Module {
			for _, v2 := range v1.Feature {
				for _, v3 := range v2.Action {
					records = append(records, *NewModuleApiInfoFromPB(v3))
				}
			}
		}
	}
	return
}

func NewModuleApiInfoListFromPBMap(m map[string]*pbam.RoleModule) ModuleApiInfoList {
	var s []*pbam.RoleModule
	for _, x := range m {
		s = append(s, x)
	}
	return NewModuleApiInfoListFromPB(s...)
}

func (records ModuleApiInfoList) ToRoleModuleMap() map[string]*pbam.RoleModule {
	var (
		featureMap    = make(map[string]*pbam.Feature)
		moduleMap     = make(map[string]*pbam.Module)
		roleModuleMap = make(map[string]*pbam.RoleModule)
	)

	// action => feature map
	for _, v := range records {
		m := featureMap[v.FeatureId]
		if m == nil {
			m = new(pbam.Feature)
		}

		m.FeatureId = v.FeatureId
		m.FeatureName = v.FeatureName
		m.Action = append(m.Action, v.ToPB())
		if v.ActionEnabled == "1" || strings.EqualFold(v.ActionEnabled, "true") {
			m.CheckedActionId = append(m.CheckedActionId, v.ActionId)
		}

		featureMap[m.FeatureId] = m
	}
	for _, m := range featureMap {
		sort.Slice(m.Action, func(i, j int) bool {
			return m.Action[i].ActionId < m.Action[j].ActionId
		})
	}

	// feature map => module map
	for _, v := range featureMap {
		action := v.Action[0]

		m := moduleMap[action.ModuleId]
		if m == nil {
			m = new(pbam.Module)
		}

		m.ModuleId = action.ModuleId
		m.ModuleName = action.ModuleName
		m.Feature = append(m.Feature, v)
		m.DataLevel = action.DataLevel
		m.CheckAll = action.IsFeatureAllChecked == "1" || strings.EqualFold(action.IsFeatureAllChecked, "true")

		moduleMap[m.ModuleId] = m
	}
	for _, m := range moduleMap {
		sort.Slice(m.Feature, func(i, j int) bool {
			return m.Feature[i].FeatureId < m.Feature[j].FeatureId
		})
	}

	// module map => role module
	for _, v := range moduleMap {
		action := v.Feature[0].Action[0]

		m := roleModuleMap[action.RoleId]
		if m == nil {
			m = new(pbam.RoleModule)
		}

		//roleModuleMap
		m.RoleId = action.RoleId
		m.RoleName = action.RoleName
		m.Module = append(m.Module, v)

		roleModuleMap[action.RoleId] = m
	}

	// OK
	return roleModuleMap
}

func (p ModuleApiInfoList) JSONString() string {
	m := p.ToRoleModuleMap()

	data, _ := json.MarshalIndent(m, "", "\t")
	data = bytes.Replace(data, []byte("\n"), []byte("\r\n"), -1)
	return string(data)
}
