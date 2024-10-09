package model

import "github.com/zeromicro/go-zero/core/stores/sqlx"

var _ SpaceRecordModel = (*customSpaceRecordModel)(nil)

type (
	// SpaceRecordModel is an interface to be customized, add more methods here,
	// and implement the added methods in customSpaceRecordModel.
	SpaceRecordModel interface {
		spaceRecordModel
	}

	customSpaceRecordModel struct {
		*defaultSpaceRecordModel
	}
)

// NewSpaceRecordModel returns a model for the database table.
func NewSpaceRecordModel(conn sqlx.SqlConn) SpaceRecordModel {
	return &customSpaceRecordModel{
		defaultSpaceRecordModel: newSpaceRecordModel(conn),
	}
}
