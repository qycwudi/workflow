package model

import (
	"context"
	"fmt"
	"strings"

	"github.com/zeromicro/go-zero/core/stores/sqlc"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ JobModel = (*customJobModel)(nil)

const (
	JobStatusOn  = "ON"
	JobStatusOff = "OFF"
)

type (
	// JobModel is an interface to be customized, add more methods here,
	// and implement the added methods in customJobModel.
	JobModel interface {
		jobModel
		FindByName(ctx context.Context, name string) (*Job, error)
		FindByWorkspaceId(ctx context.Context, workspaceId string) (*Job, error)
		FindByOn(ctx context.Context) ([]*Job, error)
		FindByJobId(ctx context.Context, jobId string) (*Job, error)
		FindPage(ctx context.Context, jobName, workspaceId string, current, pageSize int) (int64, []*Job, error)
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
	query := fmt.Sprintf("SELECT %s FROM %s WHERE job_name = ?", jobRows, m.table)
	var resp Job
	err := m.conn.QueryRowCtx(ctx, &resp, query, name)
	switch err {
	case nil:
		return &resp, nil
	case sqlc.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

func (m *customJobModel) FindByWorkspaceId(ctx context.Context, workspaceId string) (*Job, error) {
	query := fmt.Sprintf("SELECT %s FROM %s WHERE workspace_id = ? limit 1", jobRows, m.table)
	var resp Job
	err := m.conn.QueryRowCtx(ctx, &resp, query, workspaceId)
	switch err {
	case nil:
		return &resp, nil
	case sqlc.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

func (m *customJobModel) FindByOn(ctx context.Context) ([]*Job, error) {
	query := fmt.Sprintf("SELECT %s FROM %s WHERE status = ?", jobRows, m.table)
	var resp []*Job
	err := m.conn.QueryRowsCtx(ctx, &resp, query, JobStatusOn)
	return resp, err
}

func (m *customJobModel) FindByJobId(ctx context.Context, jobId string) (*Job, error) {
	query := fmt.Sprintf("SELECT %s FROM %s WHERE job_id = ? limit 1", jobRows, m.table)
	var resp Job
	err := m.conn.QueryRowCtx(ctx, &resp, query, jobId)
	return &resp, err
}

func (m *customJobModel) FindPage(ctx context.Context, jobName, workspaceId string, current, pageSize int) (int64, []*Job, error) {
	var conditions []string
	var args []interface{}

	if jobName != "" {
		conditions = append(conditions, "job_name like ?")
		args = append(args, "%"+jobName+"%")
	}

	if workspaceId != "" {
		conditions = append(conditions, "workspace_id = ?")
		args = append(args, workspaceId)
	}

	var where string
	if len(conditions) > 0 {
		where = "WHERE " + strings.Join(conditions, " AND ")
	}

	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM %s %s", m.table, where)
	var total int64
	err := m.conn.QueryRowCtx(ctx, &total, countQuery, args...)
	if err != nil {
		return 0, nil, err
	}

	query := fmt.Sprintf("SELECT %s FROM %s %s ORDER BY create_time DESC LIMIT %d OFFSET %d",
		jobRows, m.table, where, pageSize, (current-1)*pageSize)

	var resp []*Job
	err = m.conn.QueryRowsCtx(ctx, &resp, query, args...)
	if err != nil {
		return 0, nil, err
	}

	return total, resp, nil
}
