package model

import (
	"context"
	"fmt"
	"strings"

	"github.com/zeromicro/go-zero/core/logc"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ UsersModel = (*customUsersModel)(nil)

type (
	// UsersModel is an interface to be customized, add more methods here,
	// and implement the added methods in customUsersModel.
	UsersModel interface {
		usersModel
		FindPage(ctx context.Context, username string, current int64, pageSize int64) (users []*Users, total int64, err error)
	}

	customUsersModel struct {
		*defaultUsersModel
	}
)

// NewUsersModel returns a model for the database table.
func NewUsersModel(conn sqlx.SqlConn) UsersModel {
	return &customUsersModel{
		defaultUsersModel: newUsersModel(conn),
	}
}

func (m *customUsersModel) FindPage(ctx context.Context, username string, current int64, pageSize int64) (users []*Users, total int64, err error) {
	// 初始化SQL语句和参数
	var queryBuilder strings.Builder
	var totalBuilder strings.Builder
	params := make([]interface{}, 0)

	_, _ = queryBuilder.WriteString(fmt.Sprintf("SELECT %s FROM %s ", usersRows, m.table))
	_, _ = totalBuilder.WriteString(fmt.Sprintf("SELECT count(*) FROM %s ", m.table))

	// 如果username不为空，则添加到查询条件
	if username != "" {
		_, _ = queryBuilder.WriteString(" WHERE `username` like CONCAT('%', ?, '%')")
		_, _ = totalBuilder.WriteString(" WHERE `username` like CONCAT('%', ?, '%')")
		params = append(params, username)
	}

	// 查询total
	err = m.conn.QueryRowCtx(ctx, &total, totalBuilder.String(), params...)
	if err != nil {
		logc.Infov(ctx, err)
		return nil, 0, err
	}

	// 添加排序
	_, _ = queryBuilder.WriteString(" order by id desc")
	// 添加分页条件
	_, _ = queryBuilder.WriteString(" LIMIT ?, ?")

	// 计算偏移量
	params = append(params, (current-1)*pageSize)
	params = append(params, pageSize)

	var resp []*Users
	err = m.conn.QueryRowsCtx(ctx, &resp, queryBuilder.String(), params...)
	switch err {
	case nil:
		return resp, total, nil
	default:
		logc.Infov(ctx, err)
		return nil, 0, err
	}
}
