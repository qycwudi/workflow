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
	apiRecordFieldNames          = builder.RawFieldNames(&ApiRecord{})
	apiRecordRows                = strings.Join(apiRecordFieldNames, ",")
	apiRecordRowsExpectAutoSet   = strings.Join(stringx.Remove(apiRecordFieldNames, "`id`"), ",")
	apiRecordRowsWithPlaceHolder = strings.Join(stringx.Remove(apiRecordFieldNames, "`id`"), "=?,") + "=?"
)

type (
	apiRecordModel interface {
		Insert(ctx context.Context, data *ApiRecord) (sql.Result, error)
		FindOne(ctx context.Context, id int64) (*ApiRecord, error)
		Update(ctx context.Context, data *ApiRecord) error
		Delete(ctx context.Context, id int64) error
	}

	defaultApiRecordModel struct {
		conn  sqlx.SqlConn
		table string
	}

	ApiRecord struct {
		Id       int64     `db:"id"`
		Status   string    `db:"status"`
		TraceId  string    `db:"trace_id"`
		Param    string    `db:"param"`  // 参数
		Extend   string    `db:"extend"` // 扩展
		CallTime time.Time `db:"call_time"`
	}
)

func newApiRecordModel(conn sqlx.SqlConn) *defaultApiRecordModel {
	return &defaultApiRecordModel{
		conn:  conn,
		table: "`api_record`",
	}
}

func (m *defaultApiRecordModel) withSession(session sqlx.Session) *defaultApiRecordModel {
	return &defaultApiRecordModel{
		conn:  sqlx.NewSqlConnFromSession(session),
		table: "`api_record`",
	}
}

func (m *defaultApiRecordModel) Delete(ctx context.Context, id int64) error {
	query := fmt.Sprintf("delete from %s where `id` = ?", m.table)
	_, err := m.conn.ExecCtx(ctx, query, id)
	return err
}

func (m *defaultApiRecordModel) FindOne(ctx context.Context, id int64) (*ApiRecord, error) {
	query := fmt.Sprintf("select %s from %s where `id` = ? limit 1", apiRecordRows, m.table)
	var resp ApiRecord
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

func (m *defaultApiRecordModel) Insert(ctx context.Context, data *ApiRecord) (sql.Result, error) {
	query := fmt.Sprintf("insert into %s (%s) values (?, ?, ?, ?, ?)", m.table, apiRecordRowsExpectAutoSet)
	ret, err := m.conn.ExecCtx(ctx, query, data.Status, data.TraceId, data.Param, data.Extend, data.CallTime)
	return ret, err
}

func (m *defaultApiRecordModel) Update(ctx context.Context, data *ApiRecord) error {
	query := fmt.Sprintf("update %s set %s where `id` = ?", m.table, apiRecordRowsWithPlaceHolder)
	_, err := m.conn.ExecCtx(ctx, query, data.Status, data.TraceId, data.Param, data.Extend, data.CallTime, data.Id)
	return err
}

func (m *defaultApiRecordModel) tableName() string {
	return m.table
}