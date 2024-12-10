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
	locksFieldNames          = builder.RawFieldNames(&Locks{})
	locksRows                = strings.Join(locksFieldNames, ",")
	locksRowsExpectAutoSet   = strings.Join(stringx.Remove(locksFieldNames, "`id`"), ",")
	locksRowsWithPlaceHolder = strings.Join(stringx.Remove(locksFieldNames, "`id`"), "=?,") + "=?"
)

type (
	locksModel interface {
		Insert(ctx context.Context, data *Locks) (sql.Result, error)
		FindOne(ctx context.Context, id int64) (*Locks, error)
		Update(ctx context.Context, data *Locks) error
		Delete(ctx context.Context, id int64) error
	}

	defaultLocksModel struct {
		conn  sqlx.SqlConn
		table string
	}

	Locks struct {
		LockName    string    `db:"lock_name"`    // 锁名称
		IsLocked    int64     `db:"is_locked"`    // 锁是否被持有
		HeldBy      string    `db:"held_by"`      // 锁持有者
		LockedTime  time.Time `db:"locked_time"`  // 锁开始持有时间
		Timeout     int64     `db:"timeout"`      // 锁超时时间（秒）
		UpdatedTime time.Time `db:"updated_time"` // 锁更新时间
		Id          int64     `db:"id"`
	}
)

func newLocksModel(conn sqlx.SqlConn) *defaultLocksModel {
	return &defaultLocksModel{
		conn:  conn,
		table: "`locks`",
	}
}

func (m *defaultLocksModel) withSession(session sqlx.Session) *defaultLocksModel {
	return &defaultLocksModel{
		conn:  sqlx.NewSqlConnFromSession(session),
		table: "`locks`",
	}
}

func (m *defaultLocksModel) Delete(ctx context.Context, id int64) error {
	query := fmt.Sprintf("delete from %s where `id` = ?", m.table)
	_, err := m.conn.ExecCtx(ctx, query, id)
	return err
}

func (m *defaultLocksModel) FindOne(ctx context.Context, id int64) (*Locks, error) {
	query := fmt.Sprintf("select %s from %s where `id` = ? limit 1", locksRows, m.table)
	var resp Locks
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

func (m *defaultLocksModel) Insert(ctx context.Context, data *Locks) (sql.Result, error) {
	query := fmt.Sprintf("insert into %s (%s) values (?, ?, ?, ?, ?, ?)", m.table, locksRowsExpectAutoSet)
	ret, err := m.conn.ExecCtx(ctx, query, data.LockName, data.IsLocked, data.HeldBy, data.LockedTime, data.Timeout, data.UpdatedTime)
	return ret, err
}

func (m *defaultLocksModel) Update(ctx context.Context, data *Locks) error {
	query := fmt.Sprintf("update %s set %s where `id` = ?", m.table, locksRowsWithPlaceHolder)
	_, err := m.conn.ExecCtx(ctx, query, data.LockName, data.IsLocked, data.HeldBy, data.LockedTime, data.Timeout, data.UpdatedTime, data.Id)
	return err
}

func (m *defaultLocksModel) tableName() string {
	return m.table
}
