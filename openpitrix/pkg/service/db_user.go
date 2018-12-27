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

func (p *Database) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
	//if err := req.Validate(); err != nil {
	//	return nil, err
	//}

	sql, values := pkgBuildSql_InsertInto(
		dbSpec.UserTableName,
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

	reply := &pb.CreateUserResponse{
		Head: &pb.ResponseHeader{
			UserId:     req.GetHead().GetUserId(),
			OwnerPath:  "", // TODO
			AccessPath: "", // TODO
		},
		UserId: req.GetValue().GetUserId(),
	}

	return reply, nil
}
func (p *Database) DeleteUsers(ctx context.Context, req *pb.DeleteUsersRequest) (*pb.DeleteUsersResponse, error) {
	sql := pkgBuildSql_Delete(
		dbSpec.UserTableName, dbSpec.UserPrimaryKeyName,
		req.UserId...,
	)

	_, err := p.DB.ExecContext(ctx, sql)
	if err != nil {
		return nil, err
	}

	reply := &pb.DeleteUsersResponse{
		Head: &pb.ResponseHeader{
			UserId:     req.GetHead().GetUserId(),
			OwnerPath:  "", // TODO
			AccessPath: "", // TODO
		},
		UserId: req.UserId,
	}

	return reply, nil
}
func (p *Database) ModifyUser(ctx context.Context, req *pb.ModifyUserRequest) (*pb.ModifyUserResponse, error) {
	sql, values := pkgBuildSql_Update(
		dbSpec.UserTableName, req.GetValue(),
		dbSpec.UserPrimaryKeyName,
	)

	_, err := p.DB.ExecContext(ctx, sql, values...)
	if err != nil {
		return nil, err
	}

	reply := &pb.ModifyUserResponse{
		Head: &pb.ResponseHeader{
			UserId:     req.GetHead().GetUserId(),
			OwnerPath:  "", // TODO
			AccessPath: "", // TODO
		},
		UserId: req.GetValue().GetUserId(),
	}

	return reply, nil
}
func (p *Database) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.GetUserResponse, error) {
	var query = fmt.Sprintf(
		"SELECT * FROM %s WHERE %s=$1 LIMIT 1 OFFSET 0;",
		dbSpec.GroupTableName,
		dbSpec.GroupPrimaryKeyName,
	)

	rows, err := p.DB.QueryContext(ctx, query, req.GetUserId())
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

	var msg = &pb.User{}
	err = pkgSqlScanProtoMessge(rows, msg)
	if err != nil {
		return nil, err
	}

	reply := &pb.GetUserResponse{
		Head: &pb.ResponseHeader{
			UserId:     req.GetHead().GetUserId(),
			OwnerPath:  "", // TODO
			AccessPath: "", // TODO
		},
		Value: msg,
	}

	return reply, nil
}
func (p *Database) DescribeUsers(ctx context.Context, req *pb.DescribeUsersRequest) (*pb.DescribeUsersResponse, error) {
	var searchWord = req.GetSearchWord()

	if searchWord == "" {
		return p._DescribeUsers_all(ctx, req)
	} else {
		return p._DescribeUsers_bySearchWord(ctx, req)
	}
}

func (p *Database) _DescribeUsers_count(ctx context.Context, req *pb.DescribeUsersRequest) (total int, err error) {
	var query = fmt.Sprintf("SELECT COUNT(*) FROM %s", dbSpec.GroupTableName)

	rows, err := p.DB.QueryContext(ctx, query)
	if err != nil {
		return 0, err
	}
	defer rows.Close()

	if err := rows.Scan(&total); err != nil {
		return 0, err
	}

	return total, nil
}

func (p *Database) _DescribeUsers_all(ctx context.Context, req *pb.DescribeUsersRequest) (*pb.DescribeUsersResponse, error) {
	total, err := p._DescribeUsers_count(ctx, req)
	if err != nil {
		return nil, err
	}

	var query = fmt.Sprintf("SELECT * FROM %s", dbSpec.GroupTableName)
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

	var users []*pb.User
	for rows.Next() {
		var msg = &pb.User{}
		if err := pkgSqlScanProtoMessge(rows, msg); err != nil {
			return nil, err
		}

		msg.Password = ""
		msg.OldPassword = ""

		users = append(users, msg)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	reply := &pb.DescribeUsersResponse{
		Head: &pb.ResponseHeader{
			UserId:     req.GetHead().GetUserId(),
			OwnerPath:  "", // TODO
			AccessPath: "", // TODO
		},
		Value:      users,
		TotalCount: int32(total),
	}

	return reply, nil
}

func (p *Database) _DescribeUsers_bySearchWord(ctx context.Context, req *pb.DescribeUsersRequest) (*pb.DescribeUsersResponse, error) {
	var searchWord = req.GetSearchWord()

	if searchWord == "" {
		return p._DescribeUsers_all(ctx, req)
	}

	if !pkgSearchWordValid(searchWord) {
		return nil, fmt.Errorf("invalid search_word: %q", searchWord)
	}

	total, err := p._DescribeUsers_count(ctx, req)
	if err != nil {
		return nil, err
	}

	var searchWordFieldNames = pkgGetTableStringFieldNames(new(pb.Group))
	if len(searchWordFieldNames) == 0 {
		return p._DescribeUsers_all(ctx, req)
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

	var users []*pb.User
	for rows.Next() {
		var msg = &pb.User{}
		if err := pkgSqlScanProtoMessge(rows, msg); err != nil {
			return nil, err
		}

		msg.Password = ""
		msg.OldPassword = ""

		users = append(users, msg)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	reply := &pb.DescribeUsersResponse{
		Head: &pb.ResponseHeader{
			UserId:     req.GetHead().GetUserId(),
			OwnerPath:  "", // TODO
			AccessPath: "", // TODO
		},
		Value:      users,
		TotalCount: int32(total),
	}

	return reply, nil
}
