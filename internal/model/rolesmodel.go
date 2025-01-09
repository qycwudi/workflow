package model

import (
	"context"
	"fmt"
	"strings"

	"github.com/zeromicro/go-zero/core/logc"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ RolesModel = (*customRolesModel)(nil)

type (
	// RolesModel is an interface to be customized, add more methods here,
	// and implement the added methods in customRolesModel.
	RolesModel interface {
		rolesModel
		FindPage(ctx context.Context, name string, current, pageSize int64) ([]*Roles, int64, error)
		FindOneByUserId(ctx context.Context, userId int64) (*Roles, error)
		FindByRoleIds(ctx context.Context, roleIds []int64) ([]*Roles, error)
	}

	customRolesModel struct {
		*defaultRolesModel
	}
)

func (c customRolesModel) FindPage(ctx context.Context, name string, current, pageSize int64) ([]*Roles, int64, error) {
	// 初始化SQL语句和参数
	var queryBuilder strings.Builder
	var totalBuilder strings.Builder
	params := make([]interface{}, 0)

	_, _ = queryBuilder.WriteString(fmt.Sprintf("SELECT %s FROM %s ", rolesRows, c.table))
	_, _ = totalBuilder.WriteString(fmt.Sprintf("SELECT count(*) FROM %s ", c.table))

	// 如果name不为空，则添加到查询条件
	if name != "" {
		_, _ = queryBuilder.WriteString("WHERE `name` like CONCAT('%', ?, '%')")
		_, _ = totalBuilder.WriteString("WHERE `name` like CONCAT('%', ?, '%')")
		params = append(params, name)
	}

	// 查询total
	var total int64
	err := c.conn.QueryRowCtx(ctx, &total, totalBuilder.String(), params...)
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

	var resp []*Roles
	err = c.conn.QueryRowsCtx(ctx, &resp, queryBuilder.String(), params...)
	switch err {
	case nil:
		return resp, total, nil
	default:
		logc.Infov(ctx, err)
		return nil, 0, err
	}
}

// NewRolesModel returns a model for the database table.
func NewRolesModel(conn sqlx.SqlConn) RolesModel {
	return &customRolesModel{
		defaultRolesModel: newRolesModel(conn),
	}
}

func (c customRolesModel) FindOneByUserId(ctx context.Context, userId int64) (*Roles, error) {
	query := fmt.Sprintf("SELECT %s FROM %s WHERE `user_id` = ? limit 1", rolesRows, c.table)
	var resp Roles
	err := c.conn.QueryRowCtx(ctx, &resp, query, userId)
	if err != nil {
		logc.Infov(ctx, err)
		return nil, err
	}
	return &resp, nil
}
func (c customRolesModel) FindByRoleIds(ctx context.Context, roleIds []int64) ([]*Roles, error) {
	// 将 roleIds 转换为逗号分隔的字符串
	var idStr string
	for i, id := range roleIds {
		if i == 0 {
			idStr = fmt.Sprintf("%d", id)
		} else {
			idStr = fmt.Sprintf("%s,%d", idStr, id)
		}
	}
	query := fmt.Sprintf("SELECT %s FROM %s WHERE `id` IN (%s)", rolesRows, c.table, idStr)
	var resp []*Roles
	err := c.conn.QueryRowsCtx(ctx, &resp, query)
	return resp, err
}
