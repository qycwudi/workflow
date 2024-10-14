package model

import (
	"context"
	"fmt"
	"github.com/zeromicro/go-zero/core/stores/sqlc"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ ApiSecretKeyModel = (*customApiSecretKeyModel)(nil)

type (
	// ApiSecretKeyModel is an interface to be customized, add more methods here,
	// and implement the added methods in customApiSecretKeyModel.
	ApiSecretKeyModel interface {
		apiSecretKeyModel
		FindByApiId(ctx context.Context, apiId string) ([]*ApiSecretKey, error)
	}

	customApiSecretKeyModel struct {
		*defaultApiSecretKeyModel
	}
)

func (c customApiSecretKeyModel) FindByApiId(ctx context.Context, apiId string) ([]*ApiSecretKey, error) {
	query := fmt.Sprintf("select %s from %s where `api_id` = ? ", apiSecretKeyRows, c.table)
	var resp []*ApiSecretKey
	err := c.conn.QueryRowsCtx(ctx, &resp, query, apiId)
	switch err {
	case nil:
		return resp, nil
	case sqlc.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

// NewApiSecretKeyModel returns a model for the database table.
func NewApiSecretKeyModel(conn sqlx.SqlConn) ApiSecretKeyModel {
	return &customApiSecretKeyModel{
		defaultApiSecretKeyModel: newApiSecretKeyModel(conn),
	}
}
