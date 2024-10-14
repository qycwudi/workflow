package model

import (
	"context"
	"fmt"
	"github.com/zeromicro/go-zero/core/stores/sqlc"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ ApiModel = (*customApiModel)(nil)

type (
	// ApiModel is an interface to be customized, add more methods here,
	// and implement the added methods in customApiModel.
	ApiModel interface {
		apiModel
		FindByName(ctx context.Context, name string) (*Api, error)
		FindByOn(ctx context.Context) ([]*Api, error)
	}

	customApiModel struct {
		*defaultApiModel
	}
)

func (c customApiModel) FindByOn(ctx context.Context) ([]*Api, error) {
	query := fmt.Sprintf("select %s from %s where `status` = ? ", apiRows, c.table)
	var resp []*Api
	err := c.conn.QueryRowsCtx(ctx, &resp, query, On)
	switch err {
	case nil:
		return resp, nil
	case sqlc.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

func (c customApiModel) FindByName(ctx context.Context, name string) (*Api, error) {
	query := fmt.Sprintf("select %s from %s where `api_name` = ? limit 1", apiRows, c.table)
	var resp Api
	err := c.conn.QueryRowCtx(ctx, &resp, query, name)
	switch err {
	case nil:
		return &resp, nil
	case sqlc.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

const (
	On  = "ON"
	Off = "OFF"
)

// NewApiModel returns a model for the database table.
func NewApiModel(conn sqlx.SqlConn) ApiModel {
	return &customApiModel{
		defaultApiModel: newApiModel(conn),
	}
}
