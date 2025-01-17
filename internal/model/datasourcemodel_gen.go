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
	datasourceFieldNames          = builder.RawFieldNames(&Datasource{})
	datasourceRows                = strings.Join(datasourceFieldNames, ",")
	datasourceRowsExpectAutoSet   = strings.Join(stringx.Remove(datasourceFieldNames, "`id`"), ",")
	datasourceRowsWithPlaceHolder = strings.Join(stringx.Remove(datasourceFieldNames, "`id`"), "=?,") + "=?"
)

type (
	datasourceModel interface {
		Insert(ctx context.Context, data *Datasource) (sql.Result, error)
		FindOne(ctx context.Context, id int64) (*Datasource, error)
		FindOneByName(ctx context.Context, name string) (*Datasource, error)
		Update(ctx context.Context, data *Datasource) error
		Delete(ctx context.Context, id int64) error
	}

	defaultDatasourceModel struct {
		conn  sqlx.SqlConn
		table string
	}

	Datasource struct {
		Id         int64     `db:"id"`
		Type       string    `db:"type"`
		Config     string    `db:"config"`
		Switch     int64     `db:"switch"`
		Hash       string    `db:"hash"`
		Status     string    `db:"status"`
		CreateTime time.Time `db:"create_time"`
		UpdateTime time.Time `db:"update_time"`
		Name       string    `db:"name"`
	}
)

func newDatasourceModel(conn sqlx.SqlConn) *defaultDatasourceModel {
	return &defaultDatasourceModel{
		conn:  conn,
		table: "`datasource`",
	}
}

func (m *defaultDatasourceModel) withSession(session sqlx.Session) *defaultDatasourceModel {
	return &defaultDatasourceModel{
		conn:  sqlx.NewSqlConnFromSession(session),
		table: "`datasource`",
	}
}

func (m *defaultDatasourceModel) Delete(ctx context.Context, id int64) error {
	query := fmt.Sprintf("delete from %s where `id` = ?", m.table)
	_, err := m.conn.ExecCtx(ctx, query, id)
	return err
}

func (m *defaultDatasourceModel) FindOne(ctx context.Context, id int64) (*Datasource, error) {
	query := fmt.Sprintf("select %s from %s where `id` = ? limit 1", datasourceRows, m.table)
	var resp Datasource
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

func (m *defaultDatasourceModel) FindOneByName(ctx context.Context, name string) (*Datasource, error) {
	var resp Datasource
	query := fmt.Sprintf("select %s from %s where `name` = ? limit 1", datasourceRows, m.table)
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

func (m *defaultDatasourceModel) Insert(ctx context.Context, data *Datasource) (sql.Result, error) {
	query := fmt.Sprintf("insert into %s (%s) values (?, ?, ?, ?, ?, ?, ?, ?)", m.table, datasourceRowsExpectAutoSet)
	ret, err := m.conn.ExecCtx(ctx, query, data.Type, data.Config, data.Switch, data.Hash, data.Status, data.CreateTime, data.UpdateTime, data.Name)
	return ret, err
}

func (m *defaultDatasourceModel) Update(ctx context.Context, newData *Datasource) error {
	query := fmt.Sprintf("update %s set %s where `id` = ?", m.table, datasourceRowsWithPlaceHolder)
	_, err := m.conn.ExecCtx(ctx, query, newData.Type, newData.Config, newData.Switch, newData.Hash, newData.Status, newData.CreateTime, newData.UpdateTime, newData.Name, newData.Id)
	return err
}

func (m *defaultDatasourceModel) tableName() string {
	return m.table
}
