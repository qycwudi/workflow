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
	jobRecordFieldNames          = builder.RawFieldNames(&JobRecord{})
	jobRecordRows                = strings.Join(jobRecordFieldNames, ",")
	jobRecordRowsExpectAutoSet   = strings.Join(stringx.Remove(jobRecordFieldNames, "`id`"), ",")
	jobRecordRowsWithPlaceHolder = strings.Join(stringx.Remove(jobRecordFieldNames, "`id`"), "=?,") + "=?"
)

type (
	jobRecordModel interface {
		Insert(ctx context.Context, data *JobRecord) (sql.Result, error)
		FindOne(ctx context.Context, id int64) (*JobRecord, error)
		Update(ctx context.Context, data *JobRecord) error
		Delete(ctx context.Context, id int64) error
	}

	defaultJobRecordModel struct {
		conn  sqlx.SqlConn
		table string
	}

	JobRecord struct {
		Id       int64     `db:"id"`
		Status   string    `db:"status"`
		TraceId  string    `db:"trace_id"`
		Param    string    `db:"param"`  // 参数
		Result   string    `db:"result"` // 结果
		ExecTime time.Time `db:"exec_time"`
		JobId    string    `db:"job_id"`
		JobName  string    `db:"job_name"`
		ErrorMsg string    `db:"error_msg"`
	}
)

func newJobRecordModel(conn sqlx.SqlConn) *defaultJobRecordModel {
	return &defaultJobRecordModel{
		conn:  conn,
		table: "`job_record`",
	}
}

func (m *defaultJobRecordModel) withSession(session sqlx.Session) *defaultJobRecordModel {
	return &defaultJobRecordModel{
		conn:  sqlx.NewSqlConnFromSession(session),
		table: "`job_record`",
	}
}

func (m *defaultJobRecordModel) Delete(ctx context.Context, id int64) error {
	query := fmt.Sprintf("delete from %s where `id` = ?", m.table)
	_, err := m.conn.ExecCtx(ctx, query, id)
	return err
}

func (m *defaultJobRecordModel) FindOne(ctx context.Context, id int64) (*JobRecord, error) {
	query := fmt.Sprintf("select %s from %s where `id` = ? limit 1", jobRecordRows, m.table)
	var resp JobRecord
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

func (m *defaultJobRecordModel) Insert(ctx context.Context, data *JobRecord) (sql.Result, error) {
	query := fmt.Sprintf("insert into %s (%s) values (?, ?, ?, ?, ?, ?, ?, ?)", m.table, jobRecordRowsExpectAutoSet)
	ret, err := m.conn.ExecCtx(ctx, query, data.Status, data.TraceId, data.Param, data.Result, data.ExecTime, data.JobId, data.JobName, data.ErrorMsg)
	return ret, err
}

func (m *defaultJobRecordModel) Update(ctx context.Context, data *JobRecord) error {
	query := fmt.Sprintf("update %s set %s where `id` = ?", m.table, jobRecordRowsWithPlaceHolder)
	_, err := m.conn.ExecCtx(ctx, query, data.Status, data.TraceId, data.Param, data.Result, data.ExecTime, data.JobId, data.JobName, data.ErrorMsg, data.Id)
	return err
}

func (m *defaultJobRecordModel) tableName() string {
	return m.table
}
