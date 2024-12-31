package model

import (
	"context"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ PermissionsModel = (*customPermissionsModel)(nil)

type (
	// PermissionsModel is an interface to be customized, add more methods here,
	// and implement the added methods in customPermissionsModel.
	PermissionsModel interface {
		permissionsModel
		CheckPermission(ctx context.Context, userId int64, path string, method string) (bool, error)
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
