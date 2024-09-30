package model

import "github.com/zeromicro/go-zero/core/stores/sqlx"

var _ EdgeModel = (*customEdgeModel)(nil)

type (
	// EdgeModel is an interface to be customized, add more methods here,
	// and implement the added methods in customEdgeModel.
	EdgeModel interface {
		edgeModel
	}

	customEdgeModel struct {
		*defaultEdgeModel
	}
)

// NewEdgeModel returns a model for the database table.
func NewEdgeModel(conn sqlx.SqlConn) EdgeModel {
	return &customEdgeModel{
		defaultEdgeModel: newEdgeModel(conn),
	}
}
