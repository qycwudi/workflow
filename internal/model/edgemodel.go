package model

import (
	"context"
	"fmt"
	"github.com/zeromicro/go-zero/core/logc"
	"github.com/zeromicro/go-zero/core/stores/sqlc"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ EdgeModel = (*customEdgeModel)(nil)

type (
	// EdgeModel is an interface to be customized, add more methods here,
	// and implement the added methods in customEdgeModel.
	EdgeModel interface {
		edgeModel
		CheckEdge(ctx context.Context, source string, target string) (bool, error)
		DeleteByEdgeIdAndWorkSpaceId(ctx context.Context, edgeId string, workSpaceId string) error
		FindOneByWorkSpace(ctx context.Context, workspaceId string) ([]*Edge, error)
		UpdateByEdgeId(ctx context.Context, newData *Edge) error
	}

	customEdgeModel struct {
		*defaultEdgeModel
	}
)

func (m *defaultEdgeModel) UpdateByEdgeId(ctx context.Context, newData *Edge) error {
	query := fmt.Sprintf("update %s set %s where `edge_id` = ?", m.table, edgeRowsWithPlaceHolder)
	_, err := m.conn.ExecCtx(ctx, query, newData.EdgeId, newData.EdgeType, newData.CustomData, newData.Source, newData.Target, newData.Style, newData.Route, newData.WorkspaceId, newData.EdgeId)
	return err
}

func (c customEdgeModel) FindOneByWorkSpace(ctx context.Context, workspaceId string) ([]*Edge, error) {
	var resp []*Edge
	query := fmt.Sprintf("select %s from %s where `workspace_id` = ?", edgeRows, c.table)
	err := c.conn.QueryRowsCtx(ctx, &resp, query, workspaceId)
	switch err {
	case nil:
		return resp, nil
	default:
		logc.Infov(ctx, err)
		return nil, err
	}
}

func (c customEdgeModel) DeleteByEdgeIdAndWorkSpaceId(ctx context.Context, edgeId string, workSpaceId string) error {
	query := fmt.Sprintf("delete from %s where `edge_id` = ? and workspace_id = ?", c.table)
	_, err := c.conn.ExecCtx(ctx, query, edgeId, workSpaceId)
	return err
}

func (c customEdgeModel) CheckEdge(ctx context.Context, source string, target string) (bool, error) {
	var resp Edge
	query := fmt.Sprintf("select %s from %s where `source` = ? and `target` = ? limit 1", edgeRows, c.table)
	err := c.conn.QueryRowCtx(ctx, &resp, query, source, target)
	switch err {
	case nil:
		return false, nil
	case sqlc.ErrNotFound:
		return true, nil
	default:
		return false, err
	}
}

// NewEdgeModel returns a model for the database table.
func NewEdgeModel(conn sqlx.SqlConn) EdgeModel {
	return &customEdgeModel{
		defaultEdgeModel: newEdgeModel(conn),
	}
}
