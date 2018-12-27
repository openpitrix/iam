// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package service

import (
	"context"
	"database/sql"
	"fmt"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"openpitrix.io/iam/openpitrix/pkg/pb"
)

func (p *Database) CreateGroup(ctx context.Context, req *pb.CreateGroupRequest) (*pb.CreateGroupResponse, error) {
	//if err := req.Validate(); err != nil {
	//	return nil, err
	//}

	sql, values := pkgBuildSql_InsertInto(
		dbSpec.GroupTableName,
		req.GetValue(),
	)
	if len(values) == 0 {
		err := status.Errorf(codes.InvalidArgument, "empty field")
		return nil, err
	}
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

func (p *Database) DeleteGroups(ctx context.Context, req *pb.DeleteGroupsRequest) (*pb.DeleteGroupsResponse, error) {
	sql := pkgBuildSql_Delete(
		dbSpec.GroupTableName, dbSpec.GroupPrimaryKeyName,
		req.GroupId...,
	)

	_, err := p.DB.ExecContext(ctx, sql)
	if err != nil {
		return nil, err
	}

	reply := &pb.DeleteGroupsResponse{
		Head: &pb.ResponseHeader{
			UserId:     req.GetHead().GetUserId(),
			OwnerPath:  "", // TODO
			AccessPath: "", // TODO
		},
		GroupId: req.GroupId,
	}

	return reply, nil
}
func (p *Database) ModifyGroup(ctx context.Context, req *pb.ModifyGroupRequest) (*pb.ModifyGroupResponse, error) {
	sql, values := pkgBuildSql_Update(
		dbSpec.GroupTableName, req.GetValue(),
		dbSpec.GroupPrimaryKeyName,
	)

	_, err := p.DB.ExecContext(ctx, sql, values...)
	if err != nil {
		return nil, err
	}

	reply := &pb.ModifyGroupResponse{
		Head: &pb.ResponseHeader{
			UserId:     req.GetHead().GetUserId(),
			OwnerPath:  "", // TODO
			AccessPath: "", // TODO
		},
		GroupId: req.GetValue().GetGroupId(),
	}

	return reply, nil
}
func (p *Database) GetGroup(ctx context.Context, req *pb.GetGroupRequest) (*pb.GetGroupResponse, error) {
	var query = fmt.Sprintf(
		"SELECT * FROM %s WHERE %s=$1 LIMIT 1 OFFSET 0;",
		dbSpec.GroupTableName,
		dbSpec.GroupPrimaryKeyName,
	)

	rows, err := p.DB.QueryContext(ctx, query, req.GetGroupId())
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	if !rows.Next() {
		if err := rows.Err(); err != nil {
			return nil, err
		}
		return nil, sql.ErrNoRows
	}

	var msg = &pb.Group{}
	err = pkgSqlScanProtoMessge(rows, msg)
	if err != nil {
		return nil, err
	}

	reply := &pb.GetGroupResponse{
		Head: &pb.ResponseHeader{
			UserId:     req.GetHead().GetUserId(),
			OwnerPath:  "", // TODO
			AccessPath: "", // TODO
		},
		Value: msg,
	}

	return reply, nil
}
func (p *Database) DescribeGroups(ctx context.Context, req *pb.DescribeGroupsRequest) (*pb.DescribeGroupsResponse, error) {
	println("DescribeGroups 1")
	var searchWord = req.GetSearchWord()
	println("DescribeGroups 2")

	if searchWord == "" {
		println("DescribeGroups 3.1")

		return p._DescribeGroups_all(ctx, req)
	} else {
		println("DescribeGroups 3.2", searchWord)

		return p._DescribeGroups_bySearchWord(ctx, req)
	}
}

func (p *Database) _DescribeGroups_count(ctx context.Context, req *pb.DescribeGroupsRequest) (total int, err error) {
	var query = fmt.Sprintf("SELECT COUNT(*) FROM %s", dbSpec.GroupTableName)

	rows, err := p.DB.QueryContext(ctx, query)
	if err != nil {
		return 0, err
	}
	defer rows.Close()

	if rows.Next() {
		if err := rows.Scan(&total); err != nil {
			return 0, err
		}
	}

	return total, nil
}

func (p *Database) _DescribeGroups_all(ctx context.Context, req *pb.DescribeGroupsRequest) (*pb.DescribeGroupsResponse, error) {

	println("_DescribeGroups_all 1")

	total, err := p._DescribeGroups_count(ctx, req)
	if err != nil {
		return nil, err
	}
	println("_DescribeGroups_all 2")

	var query = fmt.Sprintf("SELECT * FROM %s", dbSpec.GroupTableName)
	if offset, limit := req.GetOffset(), req.GetLimit(); offset > 0 || limit > 0 {
		query += fmt.Sprintf(" LIMIT %d OFFSET %d;", limit, offset)
	} else {
		query += fmt.Sprintf(" LIMIT %d OFFSET %d;", 20, 0)
	}

	println("_DescribeGroups_all 3, query:", query)

	rows, err := p.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var groups []*pb.Group
	for i := 0; rows.Next(); i++ {
		println("_DescribeGroups_all 4:", i)

		var msg = &pb.Group{}
		if err := pkgSqlScanProtoMessge(rows, msg); err != nil {
			return nil, err
		}

		fmt.Println(msg)

		groups = append(groups, msg)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	reply := &pb.DescribeGroupsResponse{
		Head: &pb.ResponseHeader{
			UserId:     req.GetHead().GetUserId(),
			OwnerPath:  "", // TODO
			AccessPath: "", // TODO
		},
		Value:      groups,
		TotalCount: int32(total),
	}

	return reply, nil
}

func (p *Database) _DescribeGroups_bySearchWord(ctx context.Context, req *pb.DescribeGroupsRequest) (*pb.DescribeGroupsResponse, error) {
	var searchWord = req.GetSearchWord()

	if searchWord == "" {
		return p._DescribeGroups_all(ctx, req)
	}

	if !pkgSearchWordValid(searchWord) {
		return nil, fmt.Errorf("invalid search_word: %q", searchWord)
	}

	total, err := p._DescribeGroups_count(ctx, req)
	if err != nil {
		return nil, err
	}

	var searchWordFieldNames = pkgGetTableStringFieldNames(new(pb.Group))
	if len(searchWordFieldNames) == 0 {
		return p._DescribeGroups_all(ctx, req)
	}

	var query = fmt.Sprintf("SELECT * FROM %s", dbSpec.GroupTableName)
	for i, name := range searchWordFieldNames {
		if i == 0 {
			query += name + `WHERE LIKE %` + searchWord + `%`
		} else {
			query += `OR ` + name + ` LIKE %` + searchWord + `%`
		}
	}

	if offset, limit := req.GetOffset(), req.GetLimit(); offset > 0 || limit > 0 {
		query += fmt.Sprintf("LIMIT %d OFFSET %d;", limit, offset)
	} else {
		query += fmt.Sprintf("LIMIT %d OFFSET %d;", 20, 0)
	}

	rows, err := p.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var groups []*pb.Group
	for rows.Next() {
		var msg = &pb.Group{}
		if err := pkgSqlScanProtoMessge(rows, msg); err != nil {
			return nil, err
		}
		groups = append(groups, msg)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	reply := &pb.DescribeGroupsResponse{
		Head: &pb.ResponseHeader{
			UserId:     req.GetHead().GetUserId(),
			OwnerPath:  "", // TODO
			AccessPath: "", // TODO
		},
		Value:      groups,
		TotalCount: int32(total),
	}

	return reply, nil
}
