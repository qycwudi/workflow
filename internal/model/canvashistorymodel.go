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
