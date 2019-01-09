// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package db

import (
	"context"

	"github.com/golang/protobuf/ptypes"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"openpitrix.io/iam/pkg/internal/funcutil"
	"openpitrix.io/iam/pkg/pb/im"
	"openpitrix.io/iam/pkg/service/im/db_spec"
	"openpitrix.io/logger"
)

func (p *Database) CreateGroup(ctx context.Context, req *pbim.Group) (*pbim.Group, error) {
	logger.Infof(ctx, funcutil.CallerName(1))

	if req == nil {
		err := status.Errorf(codes.InvalidArgument, "empty field")
		logger.Warnf(ctx, "%+v", err)
		return nil, err
	}
	if req != nil {
		if req.Gid == "" {
			req.Gid = genGid()
		}

		if isZeroTimestamp(req.CreateTime) {
			req.CreateTime = ptypes.TimestampNow()
		}
		if isZeroTimestamp(req.UpdateTime) {
			req.UpdateTime = ptypes.TimestampNow()
		}
		if isZeroTimestamp(req.StatusTime) {
			req.StatusTime = ptypes.TimestampNow()
		}
	}

	if err := req.Validate(); err != nil {
		logger.Warnf(ctx, "%+v", err)
		return nil, err
	}

	// TODO: check group_path valid

	var dbGroup = db_spec.PBGroupToDB(req)
	if err := dbGroup.ValidateForInsert(); err != nil {
		logger.Warnf(ctx, "%+v", err)
		return nil, err
	}

	sql, values := pkgBuildSql_InsertInto(
		db_spec.DBSpec.UserGroupTableName,
		dbGroup,
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

	return req, nil
}

func (p *Database) DeleteGroups(ctx context.Context, req *pbim.GroupIdList) (*pbim.Empty, error) {
	logger.Infof(ctx, funcutil.CallerName(1))

	if req == nil || len(req.Gid) == 0 || !isValidIds(req.Gid...) {
		err := status.Errorf(codes.InvalidArgument, "empty field")
		logger.Warnf(ctx, "%+v", err)
		return nil, err
	}

	sql := pkgBuildSql_Delete(
		db_spec.DBSpec.UserGroupTableName,
		db_spec.DBSpec.UserGroupPrimaryKeyName,
		req.Gid...,
	)

	_, err := p.DB.ExecContext(ctx, sql)
	if err != nil {
		logger.Warnf(ctx, "%v", sql)
		logger.Warnf(ctx, "%+v", err)
		return nil, err
	}

	reply := &pbim.Empty{}
	return reply, nil
}

/*

func (p *Database) ModifyGroup(ctx context.Context, req *pb.ModifyGroupRequest) (*pb.ModifyGroupResponse, error) {
	logger.Infof(ctx, funcutil.CallerName(1))

	var dbGroup = pbGroupToDB(req.GetValue())

	if err := dbGroup.ValidateForUpdate(); err != nil {
		logger.Warnf(ctx, "%+v", err)
		return nil, err
	}

	sql, values := pkgBuildSql_Update(
		dbSpec.GroupTableName, dbGroup,
		dbSpec.GroupPrimaryKeyName,
	)

	_, err := p.DB.ExecContext(ctx, sql, values...)
	if err != nil {
		logger.Warnf(ctx, "%v, %v", sql, values)
		logger.Warnf(ctx, "%+v", err)
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
	logger.Infof(ctx, funcutil.CallerName(1))

	var query = fmt.Sprintf(
		"SELECT * FROM %s WHERE %s=? LIMIT 1 OFFSET 0;",
		dbSpec.GroupTableName,
		dbSpec.GroupPrimaryKeyName,
	)

	var v = DBGroup{}
	err := p.DB.GetContext(ctx, &v, query, req.GetGroupId())
	if err != nil {
		logger.Warnf(ctx, "%v", query)
		logger.Warnf(ctx, "%+v", err)
		return nil, err
	}

	reply := &pb.GetGroupResponse{
		Head: &pb.ResponseHeader{
			UserId:     req.GetHead().GetUserId(),
			OwnerPath:  "", // TODO
			AccessPath: "", // TODO
		},
		Value: v.ToPb(),
	}

	return reply, nil
}
func (p *Database) DescribeGroups(ctx context.Context, req *pb.DescribeGroupsRequest) (*pb.DescribeGroupsResponse, error) {
	logger.Infof(ctx, funcutil.CallerName(1))

	var searchWord = req.GetSearchWord()
	if searchWord == "" {
		return p._DescribeGroups_all(ctx, req)
	} else {
		return p._DescribeGroups_bySearchWord(ctx, req)
	}
}

func (p *Database) _DescribeGroups_count(ctx context.Context, req *pb.DescribeGroupsRequest) (total int, err error) {
	logger.Infof(ctx, funcutil.CallerName(1))

	var query = fmt.Sprintf("SELECT COUNT(*) FROM %s", dbSpec.GroupTableName)

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

func (p *Database) _DescribeGroups_all(ctx context.Context, req *pb.DescribeGroupsRequest) (*pb.DescribeGroupsResponse, error) {
	logger.Infof(ctx, funcutil.CallerName(1))

	total, err := p._DescribeGroups_count(ctx, req)
	if err != nil {
		logger.Warnf(ctx, "%+v", err)
		return nil, err
	}

	var query = fmt.Sprintf("SELECT * FROM %s", dbSpec.GroupTableName)
	if offset, limit := req.GetOffset(), req.GetLimit(); offset > 0 || limit > 0 {
		query += fmt.Sprintf(" LIMIT %d OFFSET %d;", limit, offset)
	} else {
		query += fmt.Sprintf(" LIMIT %d OFFSET %d;", 20, 0)
	}

	var rows = []DBGroup{}
	err = p.DB.SelectContext(ctx, &rows, query)
	if err != nil {
		logger.Warnf(ctx, "%v", query)
		logger.Warnf(ctx, "%+v", err)
		return nil, err
	}

	var sets []*pb.Group
	for _, v := range rows {
		sets = append(sets, v.ToPb())
	}

	reply := &pb.DescribeGroupsResponse{
		Head: &pb.ResponseHeader{
			UserId:     req.GetHead().GetUserId(),
			OwnerPath:  "", // TODO
			AccessPath: "", // TODO
		},
		GroupSet:   sets,
		TotalCount: int32(total),
	}

	return reply, nil
}

func (p *Database) _DescribeGroups_bySearchWord(ctx context.Context, req *pb.DescribeGroupsRequest) (*pb.DescribeGroupsResponse, error) {
	logger.Infof(ctx, funcutil.CallerName(1))

	var searchWord = req.GetSearchWord()

	if searchWord == "" {
		return p._DescribeGroups_all(ctx, req)
	}

	if !pkgSearchWordValid(searchWord) {
		return nil, fmt.Errorf("invalid search_word: %q", searchWord)
	}

	var searchWordFieldNames = pkgGetDBTableStringFieldNames(new(DBGroup))
	if len(searchWordFieldNames) == 0 {
		return p._DescribeGroups_all(ctx, req)
	}

	var (
		queryHeaer       = fmt.Sprintf("SELECT * FROM %s ", dbSpec.GroupTableName)
		queryCountHeader = fmt.Sprintf("SELECT COUNT(*) FROM %s ", dbSpec.GroupTableName)
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

	var rows = []DBGroup{}
	err := p.DB.SelectContext(ctx, &rows, queryHeaer+queryTail)
	if err != nil {
		logger.Warnf(ctx, "%v", queryCountHeader+queryTail)
		logger.Warnf(ctx, "%+v", err)
		return nil, err
	}

	var sets []*pb.Group
	for _, v := range rows {
		sets = append(sets, v.ToPb())
	}

	reply := &pb.DescribeGroupsResponse{
		Head: &pb.ResponseHeader{
			UserId:     req.GetHead().GetUserId(),
			OwnerPath:  "", // TODO
			AccessPath: "", // TODO
		},
		GroupSet:   sets,
		TotalCount: int32(total),
	}

	return reply, nil
}

*/
