package model

import "github.com/zeromicro/go-zero/core/stores/sqlx"

var _ WorkspaceModel = (*customWorkspaceModel)(nil)

type (
	// WorkspaceModel is an interface to be customized, add more methods here,
	// and implement the added methods in customWorkspaceModel.
	WorkspaceModel interface {
		workspaceModel
	}

	customWorkspaceModel struct {
		*defaultWorkspaceModel
	}
)

// NewWorkspaceModel returns a model for the database table.
func NewWorkspaceModel(conn sqlx.SqlConn) WorkspaceModel {
	return &customWorkspaceModel{
		defaultWorkspaceModel: newWorkspaceModel(conn),
	}
}
