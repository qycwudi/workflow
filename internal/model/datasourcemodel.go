package model

import (
	"context"
	"fmt"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ DatasourceModel = (*customDatasourceModel)(nil)

type (
	// DatasourceModel is an interface to be customized, add more methods here,
	// and implement the added methods in customDatasourceModel.
	DatasourceModel interface {
		datasourceModel
		FindDataSourcePageList(ctx context.Context, param PageListBuilder, current int64, pageSize int64) (int64, []*Datasource, error)
		FindBySwitch(ctx context.Context, switchStatus int64) ([]*Datasource, error)
		UpdateStatus(ctx context.Context, id int64, status string) error
	}

	customDatasourceModel struct {
		*defaultDatasourceModel
	}

	PageListBuilder struct {
		Type   string
		Status string
		Switch int64
	}
)

// NewDatasourceModel returns a model for the database table.
func NewDatasourceModel(conn sqlx.SqlConn) DatasourceModel {
	return &customDatasourceModel{
		defaultDatasourceModel: newDatasourceModel(conn),
	}
}

func (m *defaultDatasourceModel) FindDataSourcePageList(ctx context.Context, param PageListBuilder, current int64, pageSize int64) (int64, []*Datasource, error) {
	var count int64
	var list []*Datasource

	// 构建查询条件
	conditions := make([]string, 0)
	args := make([]interface{}, 0)

	if param.Type != "" {
		conditions = append(conditions, "type = ?")
		args = append(args, param.Type)
	}

	if param.Status != "" {
		conditions = append(conditions, "status = ?")
		args = append(args, param.Status)
	}

	if param.Switch >= 0 {
		conditions = append(conditions, "switch = ?")
		args = append(args, param.Switch)
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
	countQuery := "SELECT COUNT(*) FROM datasource " + whereClause
	err := m.conn.QueryRowCtx(ctx, &count, countQuery, args...)
	if err != nil {
		return 0, nil, err
	}

	// 查询列表
	offset := (current - 1) * pageSize
	query := fmt.Sprintf("SELECT %s FROM %s %s ORDER BY id DESC LIMIT ?, ?", datasourceRows, m.table, whereClause)
	args = append(args, offset, pageSize)

	err = m.conn.QueryRowsCtx(ctx, &list, query, args...)
	if err != nil {
		return 0, nil, err
	}
	if len(list) == 0 {
		return 0, nil, ErrNotFound
	}

	return count, list, nil
}

func (m *defaultDatasourceModel) FindBySwitch(ctx context.Context, switchStatus int64) ([]*Datasource, error) {
	var result []*Datasource
	err := m.conn.QueryRowsCtx(ctx, &result, "SELECT "+datasourceRows+" FROM "+m.table+" WHERE switch = ?", switchStatus)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (m *defaultDatasourceModel) UpdateStatus(ctx context.Context, id int64, status string) error {
	query := fmt.Sprintf("update %s set status = ? where id = ?", m.table)
	_, err := m.conn.ExecCtx(ctx, query, status, id)
	return err
}

const (
	DatasourceSwitchOn  = 1
	DatasourceSwitchOff = 0

	DatasourceStatusConnected = "connected"
	DatasourceStatusClosed    = "closed"
)
