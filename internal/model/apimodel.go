package model

import "github.com/zeromicro/go-zero/core/stores/sqlx"

var _ ApiModel = (*customApiModel)(nil)

type (
	// ApiModel is an interface to be customized, add more methods here,
	// and implement the added methods in customApiModel.
	ApiModel interface {
		apiModel
	}

	customApiModel struct {
		*defaultApiModel
	}
)

// NewApiModel returns a model for the database table.
func NewApiModel(conn sqlx.SqlConn) ApiModel {
	return &customApiModel{
		defaultApiModel: newApiModel(conn),
	}
}
