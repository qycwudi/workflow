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
		FindOneByNodeIdAndWorkspaceId(ctx context.Context, traceId, nodeId string) (*Trace, error)
		FindOneByNodeId(ctx context.Context, nodeId string) (*Trace, error)
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
	query := fmt.Sprintf("update %s set elapsed_time = ?,`output` = ?,status = ?,error_msg = ? where `trace_id` = ? and node_id = ?", m.table)
	_, err := m.conn.ExecCtx(ctx, query, data.ElapsedTime, data.Output, data.Status, data.ErrorMsg, data.TraceId, data.NodeId)
	return err
}

func (m *customTraceModel) FindOneByNodeIdAndWorkspaceId(ctx context.Context, traceId, nodeId string) (*Trace, error) {
	query := fmt.Sprintf("select %s from %s where workspace_id = ? and node_id = ? order by id desc limit 1", traceRows, m.table)
	var resp Trace
	err := m.conn.QueryRowCtx(ctx, &resp, query, traceId, nodeId)
	switch err {
	case nil:
		return &resp, nil
	case sqlc.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

func (m *customTraceModel) FindOneByNodeId(ctx context.Context, nodeId string) (*Trace, error) {
	query := fmt.Sprintf("select %s from %s where node_id = ? order by id desc limit 1", traceRows, m.table)
	var resp Trace
	err := m.conn.QueryRowCtx(ctx, &resp, query, nodeId)
	switch err {
	case nil:
		return &resp, nil
	case sqlc.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}
