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

func (p *Database) CreateRole(ctx context.Context, req *pb.CreateRoleRequest) (*pb.CreateRoleResponse, error) {
	//if err := req.Validate(); err != nil {
	//	return nil, err
	//}

	sql, values := pkgBuildSql_InsertInto(
		dbSpec.RoleTableName,
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

	reply := &pb.CreateRoleResponse{
		Head: &pb.ResponseHeader{
			UserId:     req.GetHead().GetUserId(),
			OwnerPath:  "", // TODO
			AccessPath: "", // TODO
		},
		RoleId: req.GetValue().GetRoleId(),
	}

	return reply, nil
}

func (p *Database) DeleteRoles(ctx context.Context, req *pb.DeleteRolesRequest) (*pb.DeleteRolesResponse, error) {
	sql := pkgBuildSql_Delete(
		dbSpec.RoleTableName, dbSpec.RolePrimaryKeyName,
		req.RoleId...,
	)

	_, err := p.DB.ExecContext(ctx, sql)
	if err != nil {
		return nil, err
	}

	reply := &pb.DeleteRolesResponse{
		Head: &pb.ResponseHeader{
			UserId:     req.GetHead().GetUserId(),
			OwnerPath:  "", // TODO
			AccessPath: "", // TODO
		},
		RoleId: req.RoleId,
	}

	return reply, nil
}
func (p *Database) ModifyRole(ctx context.Context, req *pb.ModifyRoleRequest) (*pb.ModifyRoleResponse, error) {
	sql, values := pkgBuildSql_Update(
		dbSpec.RoleTableName, req.GetValue(),
		dbSpec.RolePrimaryKeyName,
	)

	_, err := p.DB.ExecContext(ctx, sql, values...)
	if err != nil {
		return nil, err
	}

	reply := &pb.ModifyRoleResponse{
		Head: &pb.ResponseHeader{
			UserId:     req.GetHead().GetUserId(),
			OwnerPath:  "", // TODO
			AccessPath: "", // TODO
		},
		RoleId: req.GetValue().GetRoleId(),
	}

	return reply, nil
}
func (p *Database) GetRole(ctx context.Context, req *pb.GetRoleRequest) (*pb.GetRoleResponse, error) {
	var query = fmt.Sprintf(
		"SELECT * FROM %s WHERE %s=$1 LIMIT 1 OFFSET 0;",
		dbSpec.RoleTableName,
		dbSpec.RolePrimaryKeyName,
	)

	rows, err := p.DB.QueryContext(ctx, query, req.GetRoleId())
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

	var msg = &pb.Role{}
	err = pkgSqlScanProtoMessge(rows, msg)
	if err != nil {
		return nil, err
	}

	reply := &pb.GetRoleResponse{
		Head: &pb.ResponseHeader{
			UserId:     req.GetHead().GetUserId(),
			OwnerPath:  "", // TODO
			AccessPath: "", // TODO
		},
		Value: msg,
	}

	return reply, nil
}
func (p *Database) DescribeRoles(ctx context.Context, req *pb.DescribeRolesRequest) (*pb.DescribeRolesResponse, error) {
	var searchWord = req.GetSearchWord()

	if searchWord == "" {
		return p._DescribeRoles_all(ctx, req)
	} else {
		return p._DescribeRoles_bySearchWord(ctx, req)
	}
}

func (p *Database) _DescribeRoles_count(ctx context.Context, req *pb.DescribeRolesRequest) (total int, err error) {
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

func (p *Database) _DescribeRoles_all(ctx context.Context, req *pb.DescribeRolesRequest) (*pb.DescribeRolesResponse, error) {
	total, err := p._DescribeRoles_count(ctx, req)
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

	var roles []*pb.Role
	for rows.Next() {
		var msg = &pb.Role{}
		if err := pkgSqlScanProtoMessge(rows, msg); err != nil {
			return nil, err
		}
		roles = append(roles, msg)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	reply := &pb.DescribeRolesResponse{
		Head: &pb.ResponseHeader{
			UserId:     req.GetHead().GetUserId(),
			OwnerPath:  "", // TODO
			AccessPath: "", // TODO
		},
		Value:      roles,
		TotalCount: int32(total),
	}

	return reply, nil
}

func (p *Database) _DescribeRoles_bySearchWord(ctx context.Context, req *pb.DescribeRolesRequest) (*pb.DescribeRolesResponse, error) {
	var searchWord = req.GetSearchWord()

	if searchWord == "" {
		return p._DescribeRoles_all(ctx, req)
	}

	if !pkgSearchWordValid(searchWord) {
		return nil, fmt.Errorf("invalid search_word: %q", searchWord)
	}

	total, err := p._DescribeRoles_count(ctx, req)
	if err != nil {
		return nil, err
	}

	var searchWordFieldNames = pkgGetTableStringFieldNames(new(pb.Group))
	if len(searchWordFieldNames) == 0 {
		return p._DescribeRoles_all(ctx, req)
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

	var roles []*pb.Role
	for rows.Next() {
		var msg = &pb.Role{}
		if err := pkgSqlScanProtoMessge(rows, msg); err != nil {
			return nil, err
		}
		roles = append(roles, msg)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	reply := &pb.DescribeRolesResponse{
		Head: &pb.ResponseHeader{
			UserId:     req.GetHead().GetUserId(),
			OwnerPath:  "", // TODO
			AccessPath: "", // TODO
		},
		Value:      roles,
		TotalCount: int32(total),
	}

	return reply, nil
}
