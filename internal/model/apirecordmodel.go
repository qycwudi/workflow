package model

import "github.com/zeromicro/go-zero/core/stores/sqlx"

var _ ApiRecordModel = (*customApiRecordModel)(nil)

type (
	// ApiRecordModel is an interface to be customized, add more methods here,
	// and implement the added methods in customApiRecordModel.
	ApiRecordModel interface {
		apiRecordModel
	}

	customApiRecordModel struct {
		*defaultApiRecordModel
	}
)

// NewApiRecordModel returns a model for the database table.
func NewApiRecordModel(conn sqlx.SqlConn) ApiRecordModel {
	return &customApiRecordModel{
		defaultApiRecordModel: newApiRecordModel(conn),
	}
}
