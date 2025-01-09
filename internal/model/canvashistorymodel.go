package model

import (
	"context"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ CanvasHistoryModel = (*customCanvasHistoryModel)(nil)

type (
	// CanvasHistoryModel is an interface to be customized, add more methods here,
	// and implement the added methods in customCanvasHistoryModel.
	CanvasHistoryModel interface {
		canvasHistoryModel
		FindAll(ctx context.Context, cond *CanvasHistory) ([]*CanvasHistory, error)
		FindPage(ctx context.Context, workspaceId string, name string, current int, pageSize int) ([]*CanvasHistory, int64, error)
		FindAllApiByWorkspaceId(ctx context.Context, workspaceId string, current int, pageSize int) ([]*CanvasHistory, int64, error)
	}

	customCanvasHistoryModel struct {
		*defaultCanvasHistoryModel
	}
)

// NewCanvasHistoryModel returns a model for the database table.
func NewCanvasHistoryModel(conn sqlx.SqlConn) CanvasHistoryModel {
	return &customCanvasHistoryModel{
		defaultCanvasHistoryModel: newCanvasHistoryModel(conn),
	}
}

func (m *customCanvasHistoryModel) FindAll(ctx context.Context, cond *CanvasHistory) ([]*CanvasHistory, error) {
	var resp []*CanvasHistory
	err := m.conn.QueryRowsCtx(ctx, &resp, "SELECT "+canvasHistoryRows+" FROM "+m.table+" WHERE workspace_id = ? order by id desc", cond.WorkspaceId)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (m *customCanvasHistoryModel) FindAllApiByWorkspaceId(ctx context.Context, workspaceId string, current int, pageSize int) ([]*CanvasHistory, int64, error) {
	var total int64
	err := m.conn.QueryRowCtx(ctx, &total, "SELECT count(*) FROM "+m.table+" WHERE workspace_id = ? and is_api = 1", workspaceId)
	if err != nil {
		return nil, 0, err
	}

	var resp []*CanvasHistory
	offset := (current - 1) * pageSize
	query := "SELECT " + canvasHistoryRows + " FROM " + m.table + " WHERE workspace_id = ? and is_api = 1 ORDER BY id DESC LIMIT ?,?"
	err = m.conn.QueryRowsCtx(ctx, &resp, query, workspaceId, offset, pageSize)
	if err != nil {
		return nil, 0, err
	}

	return resp, total, nil
}

func (m *customCanvasHistoryModel) FindPage(ctx context.Context, workspaceId string, name string, current int, pageSize int) ([]*CanvasHistory, int64, error) {
	var total int64
	err := m.conn.QueryRowCtx(ctx, &total, "SELECT count(*) FROM "+m.table+" WHERE workspace_id = ? and name like ?", workspaceId, "%"+name+"%")
	if err != nil {
		return nil, 0, err
	}

	var resp []*CanvasHistory
	offset := (current - 1) * pageSize
	query := "SELECT " + canvasHistoryRows + " FROM " + m.table + " WHERE workspace_id = ? and name like ? ORDER BY id DESC LIMIT ?,?"
	err = m.conn.QueryRowsCtx(ctx, &resp, query, workspaceId, "%"+name+"%", offset, pageSize)
	if err != nil {
		return nil, 0, err
	}
	return resp, total, nil
}
