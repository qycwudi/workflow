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
	workspaceFieldNames          = builder.RawFieldNames(&Workspace{})
	workspaceRows                = strings.Join(workspaceFieldNames, ",")
	workspaceRowsExpectAutoSet   = strings.Join(stringx.Remove(workspaceFieldNames, "`id`"), ",")
	workspaceRowsWithPlaceHolder = strings.Join(stringx.Remove(workspaceFieldNames, "`id`"), "=?,") + "=?"
)

type (
	workspaceModel interface {
		Insert(ctx context.Context, data *Workspace) (sql.Result, error)
		FindOne(ctx context.Context, id int64) (*Workspace, error)
		FindOneByWorkspaceId(ctx context.Context, workspaceId string) (*Workspace, error)
		Update(ctx context.Context, data *Workspace) error
		Delete(ctx context.Context, id int64) error
	}

	defaultWorkspaceModel struct {
		conn  sqlx.SqlConn
		table string
	}

	Workspace struct {
		Id            int64          `db:"id"`             // 自增主建
		WorkspaceId   string         `db:"workspace_id"`   // 主建
		WorkspaceName string         `db:"workspace_name"` // 名称
		WorkspaceDesc sql.NullString `db:"workspace_desc"` // 描述
		WorkspaceType sql.NullString `db:"workspace_type"` // 类型flow|agent
		WorkspaceIcon sql.NullString `db:"workspace_icon"` // iconUrl
		CanvasConfig  sql.NullString `db:"canvas_config"`  // 前端画布配置
		CreateTime    time.Time      `db:"create_time"`    // 创建时间
		UpdateTime    time.Time      `db:"update_time"`    // 修改时间
		IsDelete      int64          `db:"is_delete"`      // 是否删除
	}
)

func newWorkspaceModel(conn sqlx.SqlConn) *defaultWorkspaceModel {
	return &defaultWorkspaceModel{
		conn:  conn,
		table: "`workspace`",
	}
}

func (m *defaultWorkspaceModel) withSession(session sqlx.Session) *defaultWorkspaceModel {
	return &defaultWorkspaceModel{
		conn:  sqlx.NewSqlConnFromSession(session),
		table: "`workspace`",
	}
}

func (m *defaultWorkspaceModel) Delete(ctx context.Context, id int64) error {
	query := fmt.Sprintf("delete from %s where `id` = ?", m.table)
	_, err := m.conn.ExecCtx(ctx, query, id)
	return err
}

func (m *defaultWorkspaceModel) FindOne(ctx context.Context, id int64) (*Workspace, error) {
	query := fmt.Sprintf("select %s from %s where `id` = ? limit 1", workspaceRows, m.table)
	var resp Workspace
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

func (m *defaultWorkspaceModel) FindOneByWorkspaceId(ctx context.Context, workspaceId string) (*Workspace, error) {
	var resp Workspace
	query := fmt.Sprintf("select %s from %s where `workspace_id` = ? limit 1", workspaceRows, m.table)
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

func (m *defaultWorkspaceModel) Insert(ctx context.Context, data *Workspace) (sql.Result, error) {
	query := fmt.Sprintf("insert into %s (%s) values (?, ?, ?, ?, ?, ?, ?, ?, ?)", m.table, workspaceRowsExpectAutoSet)
	ret, err := m.conn.ExecCtx(ctx, query, data.WorkspaceId, data.WorkspaceName, data.WorkspaceDesc, data.WorkspaceType, data.WorkspaceIcon, data.CanvasConfig, data.CreateTime, data.UpdateTime, data.IsDelete)
	return ret, err
}

func (m *defaultWorkspaceModel) Update(ctx context.Context, newData *Workspace) error {
	query := fmt.Sprintf("update %s set %s where `id` = ?", m.table, workspaceRowsWithPlaceHolder)
	_, err := m.conn.ExecCtx(ctx, query, newData.WorkspaceId, newData.WorkspaceName, newData.WorkspaceDesc, newData.WorkspaceType, newData.WorkspaceIcon, newData.CanvasConfig, newData.CreateTime, newData.UpdateTime, newData.IsDelete, newData.Id)
	return err
}

func (m *defaultWorkspaceModel) tableName() string {
	return m.table
}
