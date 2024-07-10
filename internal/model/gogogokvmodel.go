package model

import (
	"context"
	"fmt"
	"github.com/zeromicro/go-zero/core/stores/sqlc"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ GogogoKvModel = (*customGogogoKvModel)(nil)

type (
	// GogogoKvModel is an interface to be customized, add more methods here,
	// and implement the added methods in customGogogoKvModel.
	GogogoKvModel interface {
		gogogoKvModel
		FindByKey(ctx context.Context, id string) (*GogogoKv, error)
	}

	customGogogoKvModel struct {
		*defaultGogogoKvModel
	}
)

func (c customGogogoKvModel) FindByKey(ctx context.Context, key string) (*GogogoKv, error) {
	query := fmt.Sprintf("select %s from %s where `k` = ? limit 1", gogogoKvRows, c.table)
	var resp GogogoKv
	err := c.conn.QueryRowCtx(ctx, &resp, query, key)
	switch err {
	case nil:
		return &resp, nil
	case sqlc.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

// NewGogogoKvModel returns a model for the database table.
func NewGogogoKvModel(conn sqlx.SqlConn) GogogoKvModel {
	return &customGogogoKvModel{
		defaultGogogoKvModel: newGogogoKvModel(conn),
	}
}
