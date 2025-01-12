package model

import (
	"context"
	"fmt"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ JobModel = (*customJobModel)(nil)

const (
	JobStatusOn  = "on"
	JobStatusOff = "off"
)

type (
	// JobModel is an interface to be customized, add more methods here,
	// and implement the added methods in customJobModel.
	JobModel interface {
		jobModel
		FindByName(ctx context.Context, name string) (*Job, error)
		FindByWorkspaceId(ctx context.Context, workspaceId string) (*Job, error)
		FindByOn(ctx context.Context) ([]*Job, error)
	}

	customJobModel struct {
		*defaultJobModel
	}
)

// NewJobModel returns a model for the database table.
func NewJobModel(conn sqlx.SqlConn) JobModel {
	return &customJobModel{
		defaultJobModel: newJobModel(conn),
	}
}

func (m *customJobModel) FindByName(ctx context.Context, name string) (*Job, error) {
	query := fmt.Sprintf("SELECT %s FROM %s WHERE name = ?", jobRows, m.table)
	var resp Job
	err := m.conn.QueryRowCtx(ctx, &resp, query, name)
	switch err {
	case nil:
		return &resp, nil
	default:
		return nil, err
	}
}

func (m *customJobModel) FindByWorkspaceId(ctx context.Context, workspaceId string) (*Job, error) {
	query := fmt.Sprintf("SELECT %s FROM %s WHERE workspace_id = ? limit 1", jobRows, m.table)
	var resp Job
	err := m.conn.QueryRowCtx(ctx, &resp, query, workspaceId)
	return &resp, err
}

func (m *customJobModel) FindByOn(ctx context.Context) ([]*Job, error) {
	query := fmt.Sprintf("SELECT %s FROM %s WHERE status = ?", jobRows, m.table)
	var resp []*Job
	err := m.conn.QueryRowsCtx(ctx, &resp, query, JobStatusOn)
	return resp, err
}
