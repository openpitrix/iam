// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package service

import (
	"context"
	"fmt"
	"reflect"
	"strings"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"openpitrix.io/iam/openpitrix/pkg/pb"
)

func (p *Database) CreateGroup(ctx context.Context, req *pb.CreateGroupRequest) (*pb.CreateGroupResponse, error) {
	// TODO: check group valid

	names, values := pkgGetTableFiledNamesAndValues(req.GetValue())
	if len(names) == 0 {
		err := status.Errorf(codes.InvalidArgument, "empty field")
		return nil, err
	}

	// db.Exec("INSERT INTO place (country, telcode) VALUES ($1, $2)", "Singapore", "65")
	var sql = fmt.Sprintf(
		"INSERT INTO %s (%s) VALUES (%s)",
		func() string { // table name
			if t := reflect.TypeOf(req.GetValue()); t.Kind() == reflect.Ptr {
				return strings.ToLower(t.Elem().Name())
			} else {
				return strings.ToLower(t.Name())
			}
		}(),
		func() string { // table field names
			return strings.Join(names, ",")
		}(),
		func() string { // table field values
			// "$1",
			// "$1, $2"
			// "$1, $2, $3"
			var b strings.Builder
			for i := 0; i < len(values); i++ {
				if i == 0 {
					fmt.Fprintf(&b, "$%d", i)
				} else {
					fmt.Fprintf(&b, ",$%d", i)
				}
			}
			return b.String()
		}(),
	)

	_, err := p.DB.ExecContext(ctx, sql, values...)
	if err != nil {
		return nil, err
	}

	reply := &pb.CreateGroupResponse{
		Head: &pb.ResponseHeader{
			UserId:     req.GetHead().GetUserId(),
			OwnerPath:  "", // TODO
			AccessPath: "", // TODO
		},
		GroupId: req.GetValue().GetGroupId(),
	}

	return reply, nil
}

func (p *Database) DeleteGroups(context.Context, *pb.DeleteGroupsRequest) (*pb.DeleteGroupsResponse, error) {
	panic("TODO")
}
func (p *Database) ModifyGroup(context.Context, *pb.ModifyGroupRequest) (*pb.ModifyGroupResponse, error) {
	panic("TODO")
}
func (p *Database) GetGroup(context.Context, *pb.GetGroupRequest) (*pb.GetGroupResponse, error) {
	panic("TODO")
}
func (p *Database) DescribeGroups(context.Context, *pb.DescribeGroupsRequest) (*pb.DescribeGroupsResponse, error) {
	panic("TODO")
}
