package model

import (
	"context"
	"fmt"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ ApiRecordModel = (*customApiRecordModel)(nil)

type (
	// ApiRecordModel is an interface to be customized, add more methods here,
	// and implement the added methods in customApiRecordModel.
	ApiRecordModel interface {
		apiRecordModel
		UpdateStatusByTraceId(ctx context.Context, id string, status string) error
		FindByApiId(ctx context.Context, apiId string, current int, pageSize int) (int64, []*ApiRecord, error)
		FindByApiName(ctx context.Context, apiName string, current int, pageSize int) (int64, []*ApiRecord, error)
	}

	customApiRecordModel struct {
		*defaultApiRecordModel
	}
)

func (c customApiRecordModel) FindByApiId(ctx context.Context, apiId string, current int, pageSize int) (int64, []*ApiRecord, error) {
	totalQuery := fmt.Sprintf("select count(*) from %s where api_id = ?", c.table)
	var total int64
	_ = c.conn.QueryRowsCtx(ctx, &total, totalQuery, apiId)

	query := fmt.Sprintf("select %s from %s where api_id = ? order by id desc limit ?, ?", apiRows, c.table)
	var resp []*ApiRecord
	err := c.conn.QueryRowsCtx(ctx, &resp, query, apiId, (current-1)*pageSize, pageSize)
	switch err {
	case nil:
		return total, resp, nil
	default:
		return 0, nil, err
	}
}

func (c customApiRecordModel) FindByApiName(ctx context.Context, apiName string, current int, pageSize int) (int64, []*ApiRecord, error) {
	totalQuery := fmt.Sprintf("select count(*) from %s where api_name like CONCAT('%%', ?, '%%')", c.table)
	var total int64
	_ = c.conn.QueryRowsCtx(ctx, &total, totalQuery, apiName)

	query := fmt.Sprintf("select %s from %s where api_name like CONCAT('%%', ?, '%%') order by id desc limit ?, ?", apiRows, c.table)
	var resp []*ApiRecord
	err := c.conn.QueryRowsCtx(ctx, &resp, query, apiName, (current-1)*pageSize, pageSize)
	switch err {
	case nil:
		return total, resp, nil
	default:
		return 0, nil, err
	}
}

func (c customApiRecordModel) UpdateStatusByTraceId(ctx context.Context, id string, status string) error {
	query := fmt.Sprintf("update %s set status = ?  where `trace_id` = ?", c.table)
	_, err := c.conn.ExecCtx(ctx, query, status, id)
	return err
}

// NewApiRecordModel returns a model for the database table.
func NewApiRecordModel(conn sqlx.SqlConn) ApiRecordModel {
	return &customApiRecordModel{
		defaultApiRecordModel: newApiRecordModel(conn),
	}
}
