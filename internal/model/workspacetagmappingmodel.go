package model

import (
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ WorkspaceTagMappingModel = (*customWorkspaceTagMappingModel)(nil)

type (
	// WorkspaceTagMappingModel is an interface to be customized, add more methods here,
	// and implement the added methods in customWorkspaceTagMappingModel.
	WorkspaceTagMappingModel interface {
		workspaceTagMappingModel
	}

	customWorkspaceTagMappingModel struct {
		*defaultWorkspaceTagMappingModel
	}
)

// NewWorkspaceTagMappingModel returns a model for the database table.
func NewWorkspaceTagMappingModel(conn sqlx.SqlConn) WorkspaceTagMappingModel {
	return &customWorkspaceTagMappingModel{
		defaultWorkspaceTagMappingModel: newWorkspaceTagMappingModel(conn),
	}
}
