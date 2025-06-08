package model

import (
	"context"
	"fmt"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ CaseModel = (*customCaseModel)(nil)

type (
	// CaseModel is an interface to be customized, add more methods here,
	// and implement the added methods in customCaseModel.
	CaseModel interface {
		caseModel
		withSession(session sqlx.Session) CaseModel
		FindByWorkspaceId(ctx context.Context, workspaceId string) ([]*Case, error)
	}

	customCaseModel struct {
		*defaultCaseModel
	}
)

// NewCaseModel returns a model for the database table.
func NewCaseModel(conn sqlx.SqlConn) CaseModel {
	return &customCaseModel{
		defaultCaseModel: newCaseModel(conn),
	}
}

func (m *customCaseModel) withSession(session sqlx.Session) CaseModel {
	return NewCaseModel(sqlx.NewSqlConnFromSession(session))
}

func (m *customCaseModel) FindByWorkspaceId(ctx context.Context, workspaceId string) ([]*Case, error) {
	query := fmt.Sprintf("SELECT %s FROM %s WHERE workspace_id = ? order by id desc", caseRows, m.table)
	var resp []*Case
	err := m.conn.QueryRowsCtx(ctx, &resp, query, workspaceId)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
