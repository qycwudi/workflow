package model

import (
	"context"
	"time"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ JobRecordModel = (*customJobRecordModel)(nil)

type (
	// JobRecordModel is an interface to be customized, add more methods here,
	// and implement the added methods in customJobRecordModel.
	JobRecordModel interface {
		jobRecordModel
		FindPage(ctx context.Context, current int, pageSize int, jobId string, startTime int64, endTime int64, status string) ([]*JobRecord, int64, error)
	}

	customJobRecordModel struct {
		*defaultJobRecordModel
	}
)

// NewJobRecordModel returns a model for the database table.
func NewJobRecordModel(conn sqlx.SqlConn) JobRecordModel {
	return &customJobRecordModel{
		defaultJobRecordModel: newJobRecordModel(conn),
	}
}

func (m *customJobRecordModel) FindPage(ctx context.Context, current int, pageSize int, jobId string, startTime int64, endTime int64, status string) ([]*JobRecord, int64, error) {
	var total int64
	var resp []*JobRecord

	var args []interface{}
	// 构建查询条件
	conditions := "1=1"
	if jobId != "" {
		conditions += " AND job_id = ?"
		args = append(args, jobId)
	}
	if startTime > 0 {
		startDate := time.Unix(startTime, 0)
		conditions += " AND create_time >= ?"
		args = append(args, startDate)
	}
	if endTime > 0 {
		endDate := time.Unix(endTime, 0)
		conditions += " AND create_time <= ?"
		args = append(args, endDate)
	}
	if status != "" {
		conditions += " AND status = ?"
		args = append(args, status)
	}

	// 查询总数
	query := "select count(*) from " + m.table + " WHERE " + conditions

	err := m.conn.QueryRowCtx(ctx, &total, query, args...)
	if err != nil {
		return nil, 0, err
	}

	// 分页查询数据
	if total > 0 {
		query = "select " + jobRecordRows + " from " + m.table + " WHERE " + conditions + " ORDER BY create_time DESC LIMIT ?,?"
		args = append(args, (current-1)*pageSize, pageSize)
		err = m.conn.QueryRowsCtx(ctx, &resp, query, args...)
		if err != nil {
			return nil, 0, err
		}
	}

	return resp, total, nil
}
