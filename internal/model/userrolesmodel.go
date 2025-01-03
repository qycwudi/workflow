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
