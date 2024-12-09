package model

import (
	"context"
	"fmt"

	"github.com/zeromicro/go-zero/core/stores/sqlc"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ TraceModel = (*customTraceModel)(nil)

type (
	// TraceModel is an interface to be customized, add more methods here,
	// and implement the added methods in customTraceModel.
	TraceModel interface {
		traceModel
		UpdateByTraceIdAndNodeId(ctx context.Context, data *Trace) error
		FindByTraceId(ctx context.Context, id string) ([]*Trace, error)
	}

	customTraceModel struct {
		*defaultTraceModel
	}
)

func (c customTraceModel) FindByTraceId(ctx context.Context, id string) ([]*Trace, error) {
	query := fmt.Sprintf("select %s from %s where trace_id = ? order by id asc", traceRows, c.table)
	var resp []*Trace
	err := c.conn.QueryRowsCtx(ctx, &resp, query, id)
	switch err {
	case nil:
		if len(resp) == 0 {
			return nil, ErrNotFound
		}
		return resp, nil
	case sqlc.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

// NewTraceModel returns a model for the database table.
func NewTraceModel(conn sqlx.SqlConn) TraceModel {
	return &customTraceModel{
		defaultTraceModel: newTraceModel(conn),
	}
}

func (m *defaultTraceModel) UpdateByTraceIdAndNodeId(ctx context.Context, data *Trace) error {
	query := fmt.Sprintf("update %s set elapsed_time = ?,`output` = ?,status = ? where `trace_id` = ? and node_id = ?", m.table)
	_, err := m.conn.ExecCtx(ctx, query, data.ElapsedTime, data.Output, data.Status, data.TraceId, data.NodeId)
	return err
}
