package model

import "github.com/zeromicro/go-zero/core/stores/sqlx"

var _ ModuleRelationModel = (*customModuleRelationModel)(nil)

type (
	// ModuleRelationModel is an interface to be customized, add more methods here,
	// and implement the added methods in customModuleRelationModel.
	ModuleRelationModel interface {
		moduleRelationModel
	}

	customModuleRelationModel struct {
		*defaultModuleRelationModel
	}
)

// NewModuleRelationModel returns a model for the database table.
func NewModuleRelationModel(conn sqlx.SqlConn) ModuleRelationModel {
	return &customModuleRelationModel{
		defaultModuleRelationModel: newModuleRelationModel(conn),
	}
}
