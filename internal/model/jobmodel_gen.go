// Code generated by goctl. DO NOT EDIT.

package model

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/zeromicro/go-zero/core/stores/builder"
	"github.com/zeromicro/go-zero/core/stores/sqlc"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"github.com/zeromicro/go-zero/core/stringx"
)

var (
	jobFieldNames          = builder.RawFieldNames(&Job{})
	jobRows                = strings.Join(jobFieldNames, ",")
	jobRowsExpectAutoSet   = strings.Join(stringx.Remove(jobFieldNames, "`id`"), ",")
	jobRowsWithPlaceHolder = strings.Join(stringx.Remove(jobFieldNames, "`id`"), "=?,") + "=?"
)

type (
	jobModel interface {
		Insert(ctx context.Context, data *Job) (sql.Result, error)
		FindOne(ctx context.Context, id int64) (*Job, error)
		FindOneByJobId(ctx context.Context, jobId string) (*Job, error)
		Update(ctx context.Context, data *Job) error
		Delete(ctx context.Context, id int64) error
	}

	defaultJobModel struct {
		conn  sqlx.SqlConn
		table string
	}

	Job struct {
		Id          int64     `db:"id"`
		WorkspaceId string    `db:"workspace_id"`
		JobId       string    `db:"job_id"`
		JobName     string    `db:"job_name"`
		JobCron     string    `db:"job_cron"`
		JobDesc     string    `db:"job_desc"`
		Dsl         string    `db:"dsl"`
		Status      string    `db:"status"`
		CreateTime  time.Time `db:"create_time"`
		UpdateTime  time.Time `db:"update_time"`
		HistoryId   int64     `db:"history_id"`
	}
)

func newJobModel(conn sqlx.SqlConn) *defaultJobModel {
	return &defaultJobModel{
		conn:  conn,
		table: "`job`",
	}
}

func (m *defaultJobModel) withSession(session sqlx.Session) *defaultJobModel {
	return &defaultJobModel{
		conn:  sqlx.NewSqlConnFromSession(session),
		table: "`job`",
	}
}

func (m *defaultJobModel) Delete(ctx context.Context, id int64) error {
	query := fmt.Sprintf("delete from %s where `id` = ?", m.table)
	_, err := m.conn.ExecCtx(ctx, query, id)
	return err
}

func (m *defaultJobModel) FindOne(ctx context.Context, id int64) (*Job, error) {
	query := fmt.Sprintf("select %s from %s where `id` = ? limit 1", jobRows, m.table)
	var resp Job
	err := m.conn.QueryRowCtx(ctx, &resp, query, id)
	switch err {
	case nil:
		return &resp, nil
	case sqlc.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

func (m *defaultJobModel) FindOneByJobId(ctx context.Context, jobId string) (*Job, error) {
	var resp Job
	query := fmt.Sprintf("select %s from %s where `job_id` = ? limit 1", jobRows, m.table)
	err := m.conn.QueryRowCtx(ctx, &resp, query, jobId)
	switch err {
	case nil:
		return &resp, nil
	case sqlc.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

func (m *defaultJobModel) Insert(ctx context.Context, data *Job) (sql.Result, error) {
	query := fmt.Sprintf("insert into %s (%s) values (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)", m.table, jobRowsExpectAutoSet)
	ret, err := m.conn.ExecCtx(ctx, query, data.WorkspaceId, data.JobId, data.JobName, data.JobCron, data.JobDesc, data.Dsl, data.Status, data.CreateTime, data.UpdateTime, data.HistoryId)
	return ret, err
}

func (m *defaultJobModel) Update(ctx context.Context, newData *Job) error {
	query := fmt.Sprintf("update %s set %s where `id` = ?", m.table, jobRowsWithPlaceHolder)
	_, err := m.conn.ExecCtx(ctx, query, newData.WorkspaceId, newData.JobId, newData.JobName, newData.JobCron, newData.JobDesc, newData.Dsl, newData.Status, newData.CreateTime, newData.UpdateTime, newData.HistoryId, newData.Id)
	return err
}

func (m *defaultJobModel) tableName() string {
	return m.table
}
