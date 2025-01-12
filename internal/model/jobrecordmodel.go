package model

import "github.com/zeromicro/go-zero/core/stores/sqlx"

var _ JobRecordModel = (*customJobRecordModel)(nil)

type (
	// JobRecordModel is an interface to be customized, add more methods here,
	// and implement the added methods in customJobRecordModel.
	JobRecordModel interface {
		jobRecordModel
	}

	customJobRecordModel struct {
		*defaultJobRecordModel
	}
)

// NewJobRecordModel returns a model for the database table.
func NewJobRecordModel(conn sqlx.SqlConn) JobRecordModel {
	return &customJobRecordModel{
		defaultJobRecordModel: newJobRecordModel(conn),
	}
}
