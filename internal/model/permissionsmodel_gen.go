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
	permissionsFieldNames          = builder.RawFieldNames(&Permissions{})
	permissionsRows                = strings.Join(permissionsFieldNames, ",")
	permissionsRowsExpectAutoSet   = strings.Join(stringx.Remove(permissionsFieldNames, "`id`"), ",")
	permissionsRowsWithPlaceHolder = strings.Join(stringx.Remove(permissionsFieldNames, "`id`"), "=?,") + "=?"
)

type (
	permissionsModel interface {
		Insert(ctx context.Context, data *Permissions) (sql.Result, error)
		FindOne(ctx context.Context, id int64) (*Permissions, error)
		FindOneByKey(ctx context.Context, key string) (*Permissions, error)
		Update(ctx context.Context, data *Permissions) error
		Delete(ctx context.Context, id int64) error
	}

	defaultPermissionsModel struct {
		conn  sqlx.SqlConn
		table string
	}

	Permissions struct {
		Id        int64          `db:"id"`
		Title     string         `db:"title"`      // 权限名称
		Key       string         `db:"key"`        // 权限编码
		Type      int64          `db:"type"`       // 类型 1:菜单 2:按钮 3:接口
		ParentKey string         `db:"parent_key"` // 父级权限编码
		Path      sql.NullString `db:"path"`       // 路径
		Method    sql.NullString `db:"method"`     // HTTP方法
		Sort      int64          `db:"sort"`       // 排序
		CreatedAt time.Time      `db:"created_at"`
		UpdatedAt time.Time      `db:"updated_at"`
	}
)

func newPermissionsModel(conn sqlx.SqlConn) *defaultPermissionsModel {
	return &defaultPermissionsModel{
		conn:  conn,
		table: "`permissions`",
	}
}

func (m *defaultPermissionsModel) withSession(session sqlx.Session) *defaultPermissionsModel {
	return &defaultPermissionsModel{
		conn:  sqlx.NewSqlConnFromSession(session),
		table: "`permissions`",
	}
}

func (m *defaultPermissionsModel) Delete(ctx context.Context, id int64) error {
	query := fmt.Sprintf("delete from %s where `id` = ?", m.table)
	_, err := m.conn.ExecCtx(ctx, query, id)
	return err
}

func (m *defaultPermissionsModel) FindOne(ctx context.Context, id int64) (*Permissions, error) {
	query := fmt.Sprintf("select %s from %s where `id` = ? limit 1", permissionsRows, m.table)
	var resp Permissions
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

func (m *defaultPermissionsModel) FindOneByKey(ctx context.Context, key string) (*Permissions, error) {
	var resp Permissions
	query := fmt.Sprintf("select %s from %s where `key` = ? limit 1", permissionsRows, m.table)
	err := m.conn.QueryRowCtx(ctx, &resp, query, key)
	switch err {
	case nil:
		return &resp, nil
	case sqlc.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

func (m *defaultPermissionsModel) Insert(ctx context.Context, data *Permissions) (sql.Result, error) {
	query := fmt.Sprintf("insert into %s (%s) values (?, ?, ?, ?, ?, ?, ?, ?, ?)", m.table, permissionsRowsExpectAutoSet)
	ret, err := m.conn.ExecCtx(ctx, query, data.Title, data.Key, data.Type, data.ParentKey, data.Path, data.Method, data.Sort, data.CreatedAt, data.UpdatedAt)
	return ret, err
}

func (m *defaultPermissionsModel) Update(ctx context.Context, newData *Permissions) error {
	query := fmt.Sprintf("update %s set %s where `id` = ?", m.table, permissionsRowsWithPlaceHolder)
	_, err := m.conn.ExecCtx(ctx, query, newData.Title, newData.Key, newData.Type, newData.ParentKey, newData.Path, newData.Method, newData.Sort, newData.CreatedAt, newData.UpdatedAt, newData.Id)
	return err
}

func (m *defaultPermissionsModel) tableName() string {
	return m.table
}
