// Code generated by goctl. DO NOT EDIT.

package model

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/zeromicro/go-zero/core/stores/builder"
	"github.com/zeromicro/go-zero/core/stores/sqlc"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"github.com/zeromicro/go-zero/core/stringx"
)

var (
	edgeFieldNames          = builder.RawFieldNames(&Edge{})
	edgeRows                = strings.Join(edgeFieldNames, ",")
	edgeRowsExpectAutoSet   = strings.Join(stringx.Remove(edgeFieldNames, "`id`"), ",")
	edgeRowsWithPlaceHolder = strings.Join(stringx.Remove(edgeFieldNames, "`id`"), "=?,") + "=?"
)

type (
	edgeModel interface {
		Insert(ctx context.Context, data *Edge) (sql.Result, error)
		FindOne(ctx context.Context, id int64) (*Edge, error)
		FindOneByEdgeId(ctx context.Context, edgeId string) (*Edge, error)
		Update(ctx context.Context, data *Edge) error
		Delete(ctx context.Context, id int64) error
	}

	defaultEdgeModel struct {
		conn  sqlx.SqlConn
		table string
	}

	Edge struct {
		Id         int64  `db:"id"`
		EdgeId     string `db:"edge_id"`     // 边 ID
		EdgeType   string `db:"edge_type"`   // 边类型
		CustomData string `db:"custom_data"` // 自定义数据
		Source     string `db:"source"`      // 起点
		Target     string `db:"target"`      // 终点
		Style      string `db:"style"`       // 样式
		Route      string `db:"route"`       // 路由 True、False、Failure
	}
)

func newEdgeModel(conn sqlx.SqlConn) *defaultEdgeModel {
	return &defaultEdgeModel{
		conn:  conn,
		table: "`edge`",
	}
}

func (m *defaultEdgeModel) withSession(session sqlx.Session) *defaultEdgeModel {
	return &defaultEdgeModel{
		conn:  sqlx.NewSqlConnFromSession(session),
		table: "`edge`",
	}
}

func (m *defaultEdgeModel) Delete(ctx context.Context, id int64) error {
	query := fmt.Sprintf("delete from %s where `id` = ?", m.table)
	_, err := m.conn.ExecCtx(ctx, query, id)
	return err
}

func (m *defaultEdgeModel) FindOne(ctx context.Context, id int64) (*Edge, error) {
	query := fmt.Sprintf("select %s from %s where `id` = ? limit 1", edgeRows, m.table)
	var resp Edge
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

func (m *defaultEdgeModel) FindOneByEdgeId(ctx context.Context, edgeId string) (*Edge, error) {
	var resp Edge
	query := fmt.Sprintf("select %s from %s where `edge_id` = ? limit 1", edgeRows, m.table)
	err := m.conn.QueryRowCtx(ctx, &resp, query, edgeId)
	switch err {
	case nil:
		return &resp, nil
	case sqlc.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

func (m *defaultEdgeModel) Insert(ctx context.Context, data *Edge) (sql.Result, error) {
	query := fmt.Sprintf("insert into %s (%s) values (?, ?, ?, ?, ?, ?, ?)", m.table, edgeRowsExpectAutoSet)
	ret, err := m.conn.ExecCtx(ctx, query, data.EdgeId, data.EdgeType, data.CustomData, data.Source, data.Target, data.Style, data.Route)
	return ret, err
}

func (m *defaultEdgeModel) Update(ctx context.Context, newData *Edge) error {
	query := fmt.Sprintf("update %s set %s where `id` = ?", m.table, edgeRowsWithPlaceHolder)
	_, err := m.conn.ExecCtx(ctx, query, newData.EdgeId, newData.EdgeType, newData.CustomData, newData.Source, newData.Target, newData.Style, newData.Route, newData.Id)
	return err
}

func (m *defaultEdgeModel) tableName() string {
	return m.table
}
