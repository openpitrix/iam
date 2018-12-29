// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package service

import (
	"context"
	"fmt"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"openpitrix.io/iam/pkg/internal/funcutil"
	"openpitrix.io/iam/pkg/pb"
	"openpitrix.io/logger"
)

func (p *Database) CreateRole(ctx context.Context, req *pb.CreateRoleRequest) (*pb.CreateRoleResponse, error) {
	logger.Infof(ctx, funcutil.CallerName(1))

	var dbRole = pbRoleToDB(req.GetValue())

	if err := dbRole.ValidateForInsert(); err != nil {
		logger.Warnf(ctx, "%+v", err)
		return nil, err
	}

	sql, values := pkgBuildSql_InsertInto(
		dbSpec.RoleTableName,
		dbRole,
	)
	if len(values) == 0 {
		err := status.Errorf(codes.InvalidArgument, "empty field")
		logger.Warnf(ctx, "%v, %v", sql, values)
		logger.Warnf(ctx, "%+v", err)
		return nil, err
	}

	_, err := p.DB.ExecContext(ctx, sql, values...)
	if err != nil {
		logger.Warnf(ctx, "%v, %v", sql, values)
		logger.Warnf(ctx, "%+v", err)
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
	logger.Infof(ctx, funcutil.CallerName(1))

	sql := pkgBuildSql_Delete(
		dbSpec.RoleTableName, dbSpec.RolePrimaryKeyName,
		req.RoleId...,
	)

	_, err := p.DB.ExecContext(ctx, sql)
	if err != nil {
		logger.Warnf(ctx, "%v", sql)
		logger.Warnf(ctx, "%+v", err)
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
	logger.Infof(ctx, funcutil.CallerName(1))

	var dbRole = pbRoleToDB(req.GetValue())

	if err := dbRole.ValidateForUpdate(); err != nil {
		logger.Warnf(ctx, "%+v", err)
		return nil, err
	}

	sql, values := pkgBuildSql_Update(
		dbSpec.RoleTableName, dbRole,
		dbSpec.RolePrimaryKeyName,
	)

	_, err := p.DB.ExecContext(ctx, sql, values...)
	if err != nil {
		logger.Warnf(ctx, "%v, %v", sql, values)
		logger.Warnf(ctx, "%+v", err)
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
	logger.Infof(ctx, funcutil.CallerName(1))

	var query = fmt.Sprintf(
		"SELECT * FROM %s WHERE %s=? LIMIT 1 OFFSET 0;",
		dbSpec.RoleTableName,
		dbSpec.RolePrimaryKeyName,
	)

	var v = DBRole{}
	err := p.DB.GetContext(ctx, &v, query, req.GetRoleId())
	if err != nil {
		logger.Warnf(ctx, "%v", query)
		logger.Warnf(ctx, "%+v", err)
		return nil, err
	}

	reply := &pb.GetRoleResponse{
		Head: &pb.ResponseHeader{
			UserId:     req.GetHead().GetUserId(),
			OwnerPath:  "", // TODO
			AccessPath: "", // TODO
		},
		Value: v.ToPb(),
	}

	return reply, nil

}
func (p *Database) DescribeRoles(ctx context.Context, req *pb.DescribeRolesRequest) (*pb.DescribeRolesResponse, error) {
	logger.Infof(ctx, funcutil.CallerName(1))

	var searchWord = req.GetSearchWord()
	if searchWord == "" {
		return p._DescribeRoles_all(ctx, req)
	} else {
		return p._DescribeRoles_bySearchWord(ctx, req)
	}
}

func (p *Database) _DescribeRoles_count(ctx context.Context, req *pb.DescribeRolesRequest) (total int, err error) {
	logger.Infof(ctx, funcutil.CallerName(1))

	var query = fmt.Sprintf("SELECT COUNT(*) FROM %s", dbSpec.RoleTableName)

	rows, err := p.DB.QueryContext(ctx, query)
	if err != nil {
		logger.Warnf(ctx, "%v", query)
		logger.Warnf(ctx, "%+v", err)
		return 0, err
	}
	defer rows.Close()

	if rows.Next() {
		if err := rows.Scan(&total); err != nil {
			logger.Warnf(ctx, "%v", query)
			logger.Warnf(ctx, "%+v", err)
			return 0, err
		}
	}

	return total, nil
}

func (p *Database) _DescribeRoles_all(ctx context.Context, req *pb.DescribeRolesRequest) (*pb.DescribeRolesResponse, error) {
	logger.Infof(ctx, funcutil.CallerName(1))

	total, err := p._DescribeRoles_count(ctx, req)
	if err != nil {
		logger.Warnf(ctx, "%+v", err)
		return nil, err
	}

	var query = fmt.Sprintf("SELECT * FROM %s", dbSpec.RoleTableName)
	if offset, limit := req.GetOffset(), req.GetLimit(); offset > 0 || limit > 0 {
		query += fmt.Sprintf(" LIMIT %d OFFSET %d;", limit, offset)
	} else {
		query += fmt.Sprintf(" LIMIT %d OFFSET %d;", 20, 0)
	}

	var rows = []DBRole{}
	err = p.DB.SelectContext(ctx, &rows, query)
	if err != nil {
		logger.Warnf(ctx, "%v", query)
		logger.Warnf(ctx, "%+v", err)
		return nil, err
	}

	var sets []*pb.Role
	for _, v := range rows {
		sets = append(sets, v.ToPb())
	}

	reply := &pb.DescribeRolesResponse{
		Head: &pb.ResponseHeader{
			UserId:     req.GetHead().GetUserId(),
			OwnerPath:  "", // TODO
			AccessPath: "", // TODO
		},
		RoleSet:    sets,
		TotalCount: int32(total),
	}

	return reply, nil
}

func (p *Database) _DescribeRoles_bySearchWord(ctx context.Context, req *pb.DescribeRolesRequest) (*pb.DescribeRolesResponse, error) {
	logger.Infof(ctx, funcutil.CallerName(1))

	var searchWord = req.GetSearchWord()

	if searchWord == "" {
		return p._DescribeRoles_all(ctx, req)
	}

	if !pkgSearchWordValid(searchWord) {
		return nil, fmt.Errorf("invalid search_word: %q", searchWord)
	}

	var searchWordFieldNames = pkgGetDBTableStringFieldNames(new(DBRole))
	if len(searchWordFieldNames) == 0 {
		return p._DescribeRoles_all(ctx, req)
	}

	var (
		queryHeaer       = fmt.Sprintf("SELECT * FROM %s ", dbSpec.RoleTableName)
		queryCountHeader = fmt.Sprintf("SELECT COUNT(*) FROM %s ", dbSpec.RoleTableName)
		queryTail        string
	)

	for i, name := range searchWordFieldNames {
		if i == 0 {
			queryTail += " WHERE `" + name + "` LIKE '%" + searchWord + "%'"
		} else {
			queryTail += " OR `" + name + "` LIKE '%" + searchWord + "%'"
		}
	}

	if offset, limit := req.GetOffset(), req.GetLimit(); offset > 0 || limit > 0 {
		queryTail += fmt.Sprintf(" LIMIT %d OFFSET %d;", limit, offset)
	} else {
		queryTail += fmt.Sprintf(" LIMIT %d OFFSET %d;", 20, 0)
	}

	// total
	var total int
	{
		rows, err := p.DB.QueryContext(ctx, queryCountHeader+queryTail)
		if err != nil {
			logger.Warnf(ctx, "%v", queryCountHeader+queryTail)
			logger.Warnf(ctx, "%+v", err)
			return nil, err
		}
		defer rows.Close()

		if rows.Next() {
			if err := rows.Scan(&total); err != nil {
				logger.Warnf(ctx, "%v", queryCountHeader+queryTail)
				logger.Warnf(ctx, "%+v", err)
				return nil, err
			}
		}
	}

	var rows = []DBRole{}
	err := p.DB.SelectContext(ctx, &rows, queryHeaer+queryTail)
	if err != nil {
		logger.Warnf(ctx, "%v", queryCountHeader+queryTail)
		logger.Warnf(ctx, "%+v", err)
		return nil, err
	}

	var sets []*pb.Role
	for _, v := range rows {
		sets = append(sets, v.ToPb())
	}

	reply := &pb.DescribeRolesResponse{
		Head: &pb.ResponseHeader{
			UserId:     req.GetHead().GetUserId(),
			OwnerPath:  "", // TODO
			AccessPath: "", // TODO
		},
		RoleSet:    sets,
		TotalCount: int32(total),
	}

	return reply, nil
}
