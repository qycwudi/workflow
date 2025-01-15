package model

import (
	"context"
	"fmt"
	"time"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ ApiRecordModel = (*customApiRecordModel)(nil)

type (
	// ApiRecordModel is an interface to be customized, add more methods here,
	// and implement the added methods in customApiRecordModel.
	ApiRecordModel interface {
		apiRecordModel
		UpdateStatusAndResultByTraceId(ctx context.Context, id string, status string, result string, errMsg string) error
		UpdateStatusByTraceId(ctx context.Context, id string, status string, errMsg string) error
		FindByApiId(ctx context.Context, apiId string, startTime int64, endTime int64, request string, response string, current int, pageSize int) (int64, []*ApiRecord, error)
		FindByApiName(ctx context.Context, apiName string, current int, pageSize int) (int64, []*ApiRecord, error)
		GetApiCallStatistics(ctx context.Context, apiId string, startTime int64, endTime int64) (*ApiCallStatistics, error)
	}

	customApiRecordModel struct {
		*defaultApiRecordModel
	}
	ApiCallStatistics struct {
		XAxis []string `json:"xAxis"`
		YAxis []int64  `json:"yAxis"`
	}
	ApiCallStatPoint struct {
		TimeMinute string `json:"timeMinute"`
		CallCount  int64  `json:"callCount"`
	}
)

func (c customApiRecordModel) FindByApiId(ctx context.Context, apiId string, startTime int64, endTime int64, request string, response string, current int, pageSize int) (int64, []*ApiRecord, error) {
	var total int64
	var resp []*ApiRecord

	// 构建查询条件
	conditions := make([]string, 0)
	args := make([]interface{}, 0)

	conditions = append(conditions, "api_id = ?")
	args = append(args, apiId)

	if startTime != 0 {
		conditions = append(conditions, "call_time >= ?")
		args = append(args, time.UnixMilli(startTime))
	}

	if endTime != 0 {
		conditions = append(conditions, "call_time <= ?")
		args = append(args, time.UnixMilli(endTime))
	}

	if request != "" {
		conditions = append(conditions, "param LIKE ?")
		args = append(args, "%"+request+"%")
	}

	if response != "" {
		conditions = append(conditions, "extend LIKE ?")
		args = append(args, "%"+response+"%")
	}

	// 构建WHERE子句
	whereClause := ""
	if len(conditions) > 0 {
		whereClause = "WHERE " + conditions[0]
		for i := 1; i < len(conditions); i++ {
			whereClause += " AND " + conditions[i]
		}
	}

	// 查询总数
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM %s %s", c.table, whereClause)
	err := c.conn.QueryRowCtx(ctx, &total, countQuery, args...)
	if err != nil {
		return 0, nil, err
	}

	// 查询列表
	offset := (current - 1) * pageSize
	query := fmt.Sprintf("SELECT %s FROM %s %s ORDER BY id DESC LIMIT ?, ?", apiRecordRows, c.table, whereClause)
	queryArgs := append(args, offset, pageSize)

	err = c.conn.QueryRowsCtx(ctx, &resp, query, queryArgs...)
	if err != nil {
		return 0, nil, err
	}

	return total, resp, nil
}

func (c customApiRecordModel) FindByApiName(ctx context.Context, apiName string, current int, pageSize int) (int64, []*ApiRecord, error) {
	totalQuery := fmt.Sprintf("select count(*) from %s where api_name like CONCAT('%%', ?, '%%')", c.table)
	var total int64
	_ = c.conn.QueryRowsCtx(ctx, &total, totalQuery, apiName)

	query := fmt.Sprintf("select %s from %s where api_name like CONCAT('%%', ?, '%%') order by id desc limit ?, ?", apiRecordRows, c.table)
	var resp []*ApiRecord
	err := c.conn.QueryRowsCtx(ctx, &resp, query, apiName, (current-1)*pageSize, pageSize)
	switch err {
	case nil:
		return total, resp, nil
	default:
		return 0, nil, err
	}
}

func (c customApiRecordModel) UpdateStatusAndResultByTraceId(ctx context.Context, id string, status string, result string, errMsg string) error {
	query := fmt.Sprintf("update %s set status = ?, extend = ?, error_msg = ? where `trace_id` = ?", c.table)
	_, err := c.conn.ExecCtx(ctx, query, status, result, errMsg, id)
	return err
}

func (c customApiRecordModel) UpdateStatusByTraceId(ctx context.Context, id string, status string, errMsg string) error {
	query := fmt.Sprintf("update %s set status = ?, error_msg = ? where `trace_id` = ?", c.table)
	_, err := c.conn.ExecCtx(ctx, query, status, errMsg, id)
	return err
}

func (c customApiRecordModel) GetApiCallStatistics(ctx context.Context, apiId string, startTime int64, endTime int64) (*ApiCallStatistics, error) {
	query := fmt.Sprintf(`SELECT 
		DATE_FORMAT(call_time, '%%Y-%%m-%%d %%H:%%i:00') as time_minute,
		COUNT(*) as call_count 
		FROM %s 
		WHERE api_id = ? 
		AND call_time BETWEEN ? AND ?
		GROUP BY DATE_FORMAT(call_time, '%%Y-%%m-%%d %%H:%%i:00')
		ORDER BY time_minute`, c.table)

	var stats []*ApiCallStatPoint
	err := c.conn.QueryRowsCtx(ctx, &stats, query, apiId, time.UnixMilli(startTime), time.UnixMilli(endTime))
	if err != nil {
		return nil, err
	}
	// 构建xAxis
	xAxis := make([]string, 0)
	for _, stat := range stats {
		xAxis = append(xAxis, stat.TimeMinute)
	}
	// 构建yAxis
	yAxis := make([]int64, 0)
	for _, stat := range stats {
		yAxis = append(yAxis, stat.CallCount)
	}

	return &ApiCallStatistics{
		XAxis: xAxis,
		YAxis: yAxis,
	}, nil
}

// NewApiRecordModel returns a model for the database table.
func NewApiRecordModel(conn sqlx.SqlConn) ApiRecordModel {
	return &customApiRecordModel{
		defaultApiRecordModel: newApiRecordModel(conn),
	}
}
