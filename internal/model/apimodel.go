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
		FindAll(ctx context.Context, current int, pageSize int) (int64, []*Api, error)
		FindByWorkSpaceId(ctx context.Context, id string, current int, pageSize int) (int64, []*Api, error)
		UpdateStatusByApiId(ctx context.Context, apiId string, status string) error
	}

	customApiModel struct {
		*defaultApiModel
	}
)

func (c customApiModel) UpdateStatusByApiId(ctx context.Context, apiId string, status string) error {
	query := fmt.Sprintf("update %s set `status` = ? where `api_id` = ?", c.table)
	_, err := c.conn.ExecCtx(ctx, query, status, apiId)
	return err
}

func (c customApiModel) FindAll(ctx context.Context, current int, pageSize int) (int64, []*Api, error) {
	totalQuery := fmt.Sprintf("select count(*) from %s", c.table)
	var total int64
	_ = c.conn.QueryRowsCtx(ctx, &total, totalQuery)

	query := fmt.Sprintf("select %s from %s order by id desc limit ?, ?", apiRows, c.table)
	var resp []*Api
	err := c.conn.QueryRowsCtx(ctx, &resp, query, (current-1)*pageSize, pageSize)
	switch err {
	case nil:
		return total, resp, nil
	default:
		return 0, nil, err
	}
}

func (c customApiModel) FindByWorkSpaceId(ctx context.Context, id string, current int, pageSize int) (int64, []*Api, error) {
	totalQuery := fmt.Sprintf("select count(*) from %s where workspace_id = ?", c.table)
	var total int64
	_ = c.conn.QueryRowsCtx(ctx, &total, totalQuery, id)

	query := fmt.Sprintf("select %s from %s where workspace_id = ? order by id desc limit ?, ?", apiRows, c.table)
	var resp []*Api
	err := c.conn.QueryRowsCtx(ctx, &resp, query, id, (current-1)*pageSize, pageSize)
	switch err {
	case nil:
		return total, resp, nil
	default:
		return 0, nil, err
	}
}

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
