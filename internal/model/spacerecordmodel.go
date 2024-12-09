package model

import (
	"context"
	"fmt"

	"github.com/zeromicro/go-zero/core/stores/sqlc"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ SpaceRecordModel = (*customSpaceRecordModel)(nil)

type (
	// SpaceRecordModel is an interface to be customized, add more methods here,
	// and implement the added methods in customSpaceRecordModel.
	SpaceRecordModel interface {
		spaceRecordModel
		FindAll(ctx context.Context, id string) ([]*SpaceRecord, error)
		FindHistory(ctx context.Context, id string) ([]*SpaceRecord, error)
		FindBySerialNumber(ctx context.Context, id string) (*SpaceRecord, error)
		UpdateStatusBySid(ctx context.Context, sid string, status string, duration int64) error
	}

	customSpaceRecordModel struct {
		*defaultSpaceRecordModel
	}
)

func (c customSpaceRecordModel) UpdateStatusBySid(ctx context.Context, sid string, status string, duration int64) error {
	query := fmt.Sprintf("update %s set status = ?, duration = ?  where `serial_number` = ?", c.table)
	_, err := c.conn.ExecCtx(ctx, query, status, duration, sid)
	return err
}

func (c customSpaceRecordModel) FindAll(ctx context.Context, id string) ([]*SpaceRecord, error) {
	query := fmt.Sprintf("select %s from %s where workspace_id = ?", spaceRecordRows, c.table)
	var resp []*SpaceRecord
	err := c.conn.QueryRowsCtx(ctx, &resp, query, id)
	switch err {
	case nil:
		return resp, nil
	case sqlc.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

func (c customSpaceRecordModel) FindHistory(ctx context.Context, id string) ([]*SpaceRecord, error) {
	query := fmt.Sprintf("select %s from %s where workspace_id = ? order by id desc", spaceRecordRows, c.table)
	var resp []*SpaceRecord
	err := c.conn.QueryRowsCtx(ctx, &resp, query, id)
	return resp, err
}

func (c customSpaceRecordModel) FindBySerialNumber(ctx context.Context, id string) (*SpaceRecord, error) {
	query := fmt.Sprintf("select %s from %s where serial_number = ?", spaceRecordRows, c.table)
	var resp SpaceRecord
	err := c.conn.QueryRowCtx(ctx, &resp, query, id)
	return &resp, err
}

// NewSpaceRecordModel returns a model for the database table.
func NewSpaceRecordModel(conn sqlx.SqlConn) SpaceRecordModel {
	return &customSpaceRecordModel{
		defaultSpaceRecordModel: newSpaceRecordModel(conn),
	}
}
