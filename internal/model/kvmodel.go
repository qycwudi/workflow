package model

import (
	"context"
	"fmt"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ KvModel = (*customKvModel)(nil)

type (
	// KvModel is an interface to be customized, add more methods here,
	// and implement the added methods in customKvModel.
	KvModel interface {
		kvModel
		FindAll(ctx context.Context, key string, current int64, pageSize int64) ([]*Kv, int64, error)
	}

	customKvModel struct {
		*defaultKvModel
	}
)

// NewKvModel returns a model for the database table.
func NewKvModel(conn sqlx.SqlConn) KvModel {
	return &customKvModel{
		defaultKvModel: newKvModel(conn),
	}
}

func (m *customKvModel) FindAll(ctx context.Context, key string, current int64, pageSize int64) ([]*Kv, int64, error) {
	var total int64
	countQuery := fmt.Sprintf("select count(*) from %s where `key` like CONCAT('%%', ?, '%%')", m.table)
	err := m.conn.QueryRowCtx(ctx, &total, countQuery, key)
	if err != nil {
		return nil, 0, err
	}

	query := fmt.Sprintf("select %s from %s where `key` like CONCAT('%%', ?, '%%') order by id desc limit ?, ?", kvRows, m.table)
	var resp []*Kv
	err = m.conn.QueryRowsCtx(ctx, &resp, query, key, (current-1)*pageSize, pageSize)
	if err != nil {
		return nil, total, err
	}
	return resp, total, nil
}
