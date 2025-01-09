package model

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ RolePermissionsModel = (*customRolePermissionsModel)(nil)

type (
	// RolePermissionsModel is an interface to be customized, add more methods here,
	// and implement the added methods in customRolePermissionsModel.
	RolePermissionsModel interface {
		rolePermissionsModel
		GetBindPermission(ctx context.Context, roleId int64, permissionIds []int64) ([]int64, error)
		GetRolePermission(ctx context.Context, roleId int64) ([]int64, error)
		DeleteRolePermission(ctx context.Context, roleId int64) error
		BatchBindPermission(ctx context.Context, roleId int64, permissionIds []int64) error
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
	// 将 permissionIds 转换为逗号分隔的字符串
	var idStr string
	for i, id := range permissionIds {
		if i == 0 {
			idStr = fmt.Sprintf("%d", id)
		} else {
			idStr = fmt.Sprintf("%s,%d", idStr, id)
		}
	}
	query := fmt.Sprintf("select permission_id from %s where `role_id` = ? and `permission_id` in (%s)", s.table, idStr)

	var bindPermissions []int64
	err := s.conn.QueryRowsCtx(ctx, &bindPermissions, query, roleId)
	if err != nil {
		return nil, err
	}
	return bindPermissions, nil
}

// GetRolePermission 获取角色权限
func (s *defaultRolePermissionsModel) GetRolePermission(ctx context.Context, roleId int64) ([]int64, error) {
	query := fmt.Sprintf("select permission_id from %s where `role_id` = ?", s.table)
	var rolePermissions []int64
	err := s.conn.QueryRowsCtx(ctx, &rolePermissions, query, roleId)
	if err != nil {
		return nil, err
	}
	return rolePermissions, nil
}

// DeleteRolePermission 删除角色权限
func (s *defaultRolePermissionsModel) DeleteRolePermission(ctx context.Context, roleId int64) error {
	query := fmt.Sprintf("delete from %s where `role_id` = ?", s.table)
	_, err := s.conn.ExecCtx(ctx, query, roleId)
	return err
}

// BatchBindPermission 批量绑定权限
func (s *defaultRolePermissionsModel) BatchBindPermission(ctx context.Context, roleId int64, permissionIds []int64) error {
	var values []string
	for _, permissionId := range permissionIds {
		values = append(values, fmt.Sprintf("(%d, %d, '%s')", roleId, permissionId, time.Now().Format("2006-01-02 15:04:05")))
	}
	query := fmt.Sprintf("insert into %s (`role_id`, `permission_id`, `created_at`) values %s", s.table, strings.Join(values, ","))
	_, err := s.conn.ExecCtx(ctx, query)
	return err
}
