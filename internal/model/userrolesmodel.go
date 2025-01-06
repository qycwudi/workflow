package model

import (
	"context"
	"fmt"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ UserRolesModel = (*customUserRolesModel)(nil)

type (
	// UserRolesModel is an interface to be customized, add more methods here,
	// and implement the added methods in customUserRolesModel.
	UserRolesModel interface {
		userRolesModel
		FindOneByUserId(ctx context.Context, userId int64) (*UserRoles, error)
		FindByUserIds(ctx context.Context, userIds []int64) ([]*UserRoles, error)
	}

	customUserRolesModel struct {
		*defaultUserRolesModel
	}
)

// NewUserRolesModel returns a model for the database table.
func NewUserRolesModel(conn sqlx.SqlConn) UserRolesModel {
	return &customUserRolesModel{
		defaultUserRolesModel: newUserRolesModel(conn),
	}
}

func (m *customUserRolesModel) FindOneByUserId(ctx context.Context, userId int64) (*UserRoles, error) {
	query := fmt.Sprintf("SELECT %s FROM %s WHERE `user_id` = ? limit 1", userRolesRows, m.table)
	var resp UserRoles
	err := m.conn.QueryRowCtx(ctx, &resp, query, userId)
	switch err {
	case nil:
		return &resp, nil
	default:
		return nil, err
	}
}

func (m *customUserRolesModel) FindByUserIds(ctx context.Context, userIds []int64) ([]*UserRoles, error) {
	// 将 userIds 转换为逗号分隔的字符串
	var idStr string
	for i, id := range userIds {
		if i == 0 {
			idStr = fmt.Sprintf("%d", id)
		} else {
			idStr = fmt.Sprintf("%s,%d", idStr, id)
		}
	}
	query := fmt.Sprintf("SELECT %s FROM %s WHERE `user_id` IN (%s)", userRolesRows, m.table, idStr)
	var resp []*UserRoles
	err := m.conn.QueryRowsCtx(ctx, &resp, query)
	return resp, err
}
