package model

import (
	"context"
	"fmt"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ PermissionsModel = (*customPermissionsModel)(nil)

type (
	// PermissionsModel is an interface to be customized, add more methods here,
	// and implement the added methods in customPermissionsModel.
	PermissionsModel interface {
		permissionsModel
		CheckPermission(ctx context.Context, userId int64, path string, method string) (bool, error)
		GetPermissionTree(ctx context.Context, parentId int64) ([]*Permissions, error)

		DeleteBindPermission(ctx context.Context, roleId int64, permissionId int64) error
	}

	customPermissionsModel struct {
		*defaultPermissionsModel
	}
)

// NewPermissionsModel returns a model for the database table.
func NewPermissionsModel(conn sqlx.SqlConn) PermissionsModel {
	return &customPermissionsModel{
		defaultPermissionsModel: newPermissionsModel(conn),
	}
}

// CheckPermission 检查用户是否有权限
func (s *defaultPermissionsModel) CheckPermission(ctx context.Context, userId int64, path string, method string) (bool, error) {
	query := `
        SELECT COUNT(1) FROM users u
        JOIN user_roles ur ON u.id = ur.user_id
        JOIN role_permissions rp ON ur.role_id = rp.role_id
        JOIN permissions p ON rp.permission_id = p.id
        WHERE u.id = ? AND p.path = ? AND p.method = ? AND u.status = 1
    `

	var count int
	err := s.conn.QueryRowCtx(ctx, &count, query, userId, path, method)
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

// GetPermissionTree 获取权限树
func (s *defaultPermissionsModel) GetPermissionTree(ctx context.Context, parentId int64) ([]*Permissions, error) {
	query := fmt.Sprintf("select %s from %s where parent_id = ? order by sort asc", permissionsRows, s.table)
	var permissions []*Permissions
	err := s.conn.QueryRowsCtx(ctx, &permissions, query, parentId)
	if err != nil {
		return nil, err
	}
	return permissions, nil
}

// DeleteBindPermission 删除角色与权限的绑定关系
func (s *defaultPermissionsModel) DeleteBindPermission(ctx context.Context, roleId int64, permissionId int64) error {
	query := fmt.Sprintf("delete from %s where `role_id` = ? and `permission_id` = ?", "`role_permissions`")
	_, err := s.conn.ExecCtx(ctx, query, roleId, permissionId)
	return err
}
