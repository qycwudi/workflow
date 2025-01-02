package model

import (
	"context"
	"fmt"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ RolePermissionsModel = (*customRolePermissionsModel)(nil)

type (
	// RolePermissionsModel is an interface to be customized, add more methods here,
	// and implement the added methods in customRolePermissionsModel.
	RolePermissionsModel interface {
		rolePermissionsModel
		GetBindPermission(ctx context.Context, roleId int64, permissionIds []int64) ([]int64, error)
	}

	customRolePermissionsModel struct {
		*defaultRolePermissionsModel
	}
)

// NewRolePermissionsModel returns a model for the database table.
func NewRolePermissionsModel(conn sqlx.SqlConn) RolePermissionsModel {
	return &customRolePermissionsModel{
		defaultRolePermissionsModel: newRolePermissionsModel(conn),
	}
}

// GetBindPermission 获取已绑定权限
func (s *defaultRolePermissionsModel) GetBindPermission(ctx context.Context, roleId int64, permissionIds []int64) ([]int64, error) {
	query := fmt.Sprintf("select permission_id from %s where `role_id` = ? and `permission_id` in (?)", s.table)

	var bindPermissions []int64
	err := s.conn.QueryRowsCtx(ctx, &bindPermissions, query, roleId, permissionIds)
	if err != nil {
		return nil, err
	}
	return bindPermissions, nil
}
