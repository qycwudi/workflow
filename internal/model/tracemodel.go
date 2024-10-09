package model

import "github.com/zeromicro/go-zero/core/stores/sqlx"

var _ TraceModel = (*customTraceModel)(nil)

type (
	// TraceModel is an interface to be customized, add more methods here,
	// and implement the added methods in customTraceModel.
	TraceModel interface {
		traceModel
	}

	customTraceModel struct {
		*defaultTraceModel
	}
)

// NewTraceModel returns a model for the database table.
func NewTraceModel(conn sqlx.SqlConn) TraceModel {
	return &customTraceModel{
		defaultTraceModel: newTraceModel(conn),
	}
}
