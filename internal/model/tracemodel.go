package model

import (
	"context"
	"fmt"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ TraceModel = (*customTraceModel)(nil)

type (
	// TraceModel is an interface to be customized, add more methods here,
	// and implement the added methods in customTraceModel.
	TraceModel interface {
		traceModel
		UpdateByTraceIdAndNodeId(ctx context.Context, data *Trace) error
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

func (m *defaultTraceModel) UpdateByTraceIdAndNodeId(ctx context.Context, data *Trace) error {
	query := fmt.Sprintf("update %s set %s where `trace_id` = ? and node_id = ?", m.table, traceRowsWithPlaceHolder)
	_, err := m.conn.ExecCtx(ctx, query, data.WorkspaceId, data.TraceId, data.Input, data.Logic, data.Output, data.Step, data.NodeId, data.NodeName, data.Status, data.ElapsedTime, data.StartTime, data.TraceId, data.NodeId)
	return err
}
