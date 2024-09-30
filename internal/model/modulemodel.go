package model

import "github.com/zeromicro/go-zero/core/stores/sqlx"

var _ ModuleModel = (*customModuleModel)(nil)

type (
	// ModuleModel is an interface to be customized, add more methods here,
	// and implement the added methods in customModuleModel.
	ModuleModel interface {
		moduleModel
	}

	customModuleModel struct {
		*defaultModuleModel
	}
)

// NewModuleModel returns a model for the database table.
func NewModuleModel(conn sqlx.SqlConn) ModuleModel {
	return &customModuleModel{
		defaultModuleModel: newModuleModel(conn),
	}
}
