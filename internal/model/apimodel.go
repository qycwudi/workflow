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
		UpdateStatusByApiId(ctx context.Context, apiId string, status string) error
		Page(ctx context.Context, current, size int, apiId, name string) (*PageResponse[Api], error)
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

func (c customApiModel) Page(ctx context.Context, current, size int, id, name string) (*PageResponse[Api], error) {
	conditions := make([]string, 0)
	if id != "" {
		conditions = append(conditions, fmt.Sprintf("api_id = '%s'", id))
	}

	if name != "" {
		conditions = append(conditions, fmt.Sprintf("api_name LIKE '%s'", "%"+name+"%"))
	}

	resp, err := Paginate[Api](ctx, PageRequest{
		Current:    current,
		Size:       size,
		Table:      c.table,
		Conn:       c.conn,
		Conditions: conditions,
	})
	switch err {
	case nil:
		return resp, nil
	default:
		return nil, err
	}
}

func (c customApiModel) FindByOn(ctx context.Context) ([]*Api, error) {
	query := fmt.Sprintf("select %s from %s where `status` = ? ", apiRows, c.table)
	var resp []*Api
	err := c.conn.QueryRowsCtx(ctx, &resp, query, ApiStatusOn)
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
	ApiStatusOn  = "ON"
	ApiStatusOff = "OFF"
)

// NewApiModel returns a model for the database table.
func NewApiModel(conn sqlx.SqlConn) ApiModel {
	return &customApiModel{
		defaultApiModel: newApiModel(conn),
	}
}
