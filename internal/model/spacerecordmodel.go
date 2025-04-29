package model

import (
	"context"
	"fmt"
	"strings"

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
		FindByWorkspaceIdPage(ctx context.Context, id string, page, pageSize int64) (int64, []*SpaceRecord, error)
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

// FindByWorkspaceIdPage 实现了 SpaceRecordModel。
func (c *customSpaceRecordModel) FindByWorkspaceIdPage(ctx context.Context, workspaceId string, page int64, pageSize int64) (int64, []*SpaceRecord, error) {
	// 初始化SQL语句和参数
	var queryBuilder strings.Builder
	var totalBuilder strings.Builder
	params := make([]interface{}, 0, 4) // 参数数量最多为4

	queryBuilder.WriteString(fmt.Sprintf("SELECT %s FROM %s WHERE workspace_id = ?", spaceRecordRows, c.table))
	totalBuilder.WriteString(fmt.Sprintf("SELECT count(*) FROM %s WHERE workspace_id = ?", c.table))

	// 添加排序
	queryBuilder.WriteString(" order by id desc")
	// 添加分页条件
	queryBuilder.WriteString(" LIMIT ?, ?")

	// 计算偏移量
	params = append(params, workspaceId)
	params = append(params, (page-1)*pageSize)
	params = append(params, pageSize)

	// 查询总数
	var total int64
	err := c.conn.QueryRowCtx(ctx, &total, totalBuilder.String(), params[0])
	if err != nil {
		return 0, nil, err
	}

	var resp []*SpaceRecord
	err = c.conn.QueryRowsCtx(ctx, &resp, queryBuilder.String(), params...)
	switch err {
	case nil:
		return total, resp, nil
	case sqlc.ErrNotFound:
		return 0, nil, ErrNotFound
	default:
		return 0, nil, err
	}
}

func (c customSpaceRecordModel) FindHistory(ctx context.Context, id string) ([]*SpaceRecord, error) {
	query := fmt.Sprintf("select %s from %s where workspace_id = ? order by id desc", spaceRecordRows, c.table)
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
