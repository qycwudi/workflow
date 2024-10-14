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
	}

	customApiRecordModel struct {
		*defaultApiRecordModel
	}
)

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
