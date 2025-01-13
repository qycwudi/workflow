package model

import (
	"context"
	"fmt"
	"strings"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ CanvasHistoryModel = (*customCanvasHistoryModel)(nil)

const (
	CanvasHistoryModeDraft = 0
	CanvasHistoryModeApi   = 1
	CanvasHistoryModeJob   = 2
)

type (
	// CanvasHistoryModel is an interface to be customized, add more methods here,
	// and implement the added methods in customCanvasHistoryModel.
	CanvasHistoryModel interface {
		canvasHistoryModel
		FindAll(ctx context.Context, cond *CanvasHistory) ([]*CanvasHistory, error)
		FindPage(ctx context.Context, workspaceId string, name string, mode int64, current int, pageSize int) ([]*CanvasHistory, int64, error)
		FindAllApiByWorkspaceId(ctx context.Context, workspaceId string, current int, pageSize int) ([]*CanvasHistory, int64, error)
		FindAllJobByWorkspaceId(ctx context.Context, workspaceId string, current int, pageSize int) ([]*CanvasHistory, int64, error)
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
	err := m.conn.QueryRowCtx(ctx, &total, "SELECT count(*) FROM "+m.table+" WHERE workspace_id = ? and mode = ?", workspaceId, CanvasHistoryModeApi)
	if err != nil {
		return nil, 0, err
	}

	var resp []*CanvasHistory
	offset := (current - 1) * pageSize
	query := "SELECT " + canvasHistoryRows + " FROM " + m.table + " WHERE workspace_id = ? and mode = ? ORDER BY id DESC LIMIT ?,?"
	err = m.conn.QueryRowsCtx(ctx, &resp, query, workspaceId, CanvasHistoryModeApi, offset, pageSize)
	if err != nil {
		return nil, 0, err
	}

	return resp, total, nil
}
func (m *customCanvasHistoryModel) FindPage(ctx context.Context, workspaceId string, name string, mode int64, current int, pageSize int) ([]*CanvasHistory, int64, error) {
	// 构建查询条件
	whereBuilder := strings.Builder{}
	whereBuilder.WriteString("WHERE 1=1")
	args := make([]interface{}, 0)

	if workspaceId != "" {
		whereBuilder.WriteString(" AND workspace_id = ?")
		args = append(args, workspaceId)
	}

	if name != "" {
		whereBuilder.WriteString(" AND name like ?")
		args = append(args, "%"+name+"%")
	}

	if mode > 0 {
		whereBuilder.WriteString(" AND mode = ?")
		args = append(args, mode)
	}

	// 查询总数
	countQuery := fmt.Sprintf("SELECT count(*) FROM %s %s", m.table, whereBuilder.String())
	var total int64
	err := m.conn.QueryRowCtx(ctx, &total, countQuery, args...)
	if err != nil {
		return nil, 0, err
	}

	// 分页查询数据
	offset := (current - 1) * pageSize
	query := fmt.Sprintf("SELECT %s FROM %s %s ORDER BY id DESC LIMIT ?,?",
		canvasHistoryRows, m.table, whereBuilder.String())
	args = append(args, offset, pageSize)

	var resp []*CanvasHistory
	err = m.conn.QueryRowsCtx(ctx, &resp, query, args...)
	if err != nil {
		return nil, 0, err
	}

	return resp, total, nil
}

func (m *customCanvasHistoryModel) FindAllJobByWorkspaceId(ctx context.Context, workspaceId string, current int, pageSize int) ([]*CanvasHistory, int64, error) {
	var total int64
	err := m.conn.QueryRowCtx(ctx, &total, "SELECT count(*) FROM "+m.table+" WHERE workspace_id = ? and mode = ?", workspaceId, CanvasHistoryModeJob)
	if err != nil {
		return nil, 0, err
	}

	var resp []*CanvasHistory
	offset := (current - 1) * pageSize
	query := "SELECT " + canvasHistoryRows + " FROM " + m.table + " WHERE workspace_id = ? and mode = ? ORDER BY id DESC LIMIT ?,?"
	err = m.conn.QueryRowsCtx(ctx, &resp, query, workspaceId, CanvasHistoryModeJob, offset, pageSize)
	if err != nil {
		return nil, 0, err
	}

	return resp, total, nil
}
