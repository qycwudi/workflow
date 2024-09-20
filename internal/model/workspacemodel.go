package model

import (
	"context"
	"fmt"
	"github.com/zeromicro/go-zero/core/logc"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"strings"
)

var _ WorkspaceModel = (*customWorkspaceModel)(nil)

type (
	// WorkspaceModel is an interface to be customized, add more methods here,
	// and implement the added methods in customWorkspaceModel.
	WorkspaceModel interface {
		workspaceModel
		FindPage(ctx context.Context, current, pageSize int, workSpaceType, workSpaceName string) ([]*Workspace, int64, error)
		FindInWorkSpaceId(ctx context.Context, workspaceId []string) ([]*Workspace, error)
	}

	customWorkspaceModel struct {
		*defaultWorkspaceModel
	}
)

func (c customWorkspaceModel) FindPage(ctx context.Context, current, pageSize int, workSpaceType, workSpaceName string) ([]*Workspace, int64, error) {
	// 初始化SQL语句和参数
	var queryBuilder strings.Builder
	var totalBuilder strings.Builder
	params := make([]interface{}, 0, 4) // 参数数量最多为4

	queryBuilder.WriteString(fmt.Sprintf("SELECT %s FROM %s WHERE is_delete = 0", workspaceRows, c.table))
	totalBuilder.WriteString(fmt.Sprintf("SELECT count(*) FROM %s WHERE is_delete = 0", c.table))

	// 如果workSpaceType不为空，则添加到查询条件
	if workSpaceType != "" {
		queryBuilder.WriteString(" AND `workspace_type` = ?")
		totalBuilder.WriteString(" AND `workspace_type` = ?")
		params = append(params, workSpaceType)
	}

	// 如果workSpaceName不为空，则添加到查询条件
	if workSpaceName != "" {
		queryBuilder.WriteString(" AND `workspace_name` like CONCAT('%', ?, '%')")
		totalBuilder.WriteString(" AND `workspace_name` like CONCAT('%', ?, '%')")
		params = append(params, workSpaceName)
	}
	// 查询total
	// 执行查询总数的SQL
	var total int64
	_ = c.conn.QueryRowsCtx(ctx, &total, totalBuilder.String(), params...)

	// 添加排序
	queryBuilder.WriteString(" order by id desc")
	// 添加分页条件
	queryBuilder.WriteString(" LIMIT ?, ?")

	// 计算偏移量
	params = append(params, (current-1)*pageSize)
	params = append(params, pageSize)

	var resp []*Workspace
	err := c.conn.QueryRowsCtx(ctx, &resp, queryBuilder.String(), params...)
	switch err {
	case nil:
		return resp, total, nil
	default:
		logc.Infov(ctx, err)
		return nil, 0, err
	}
}

func (c customWorkspaceModel) FindInWorkSpaceId(ctx context.Context, workspaceIds []string) ([]*Workspace, error) {
	placeholders := make([]string, len(workspaceIds))
	for i := range workspaceIds {
		placeholders[i] = "?"
	}
	inClause := strings.Join(placeholders, ", ")
	var resp []*Workspace
	query := fmt.Sprintf("select %s from %s where `workspace_id` in (%s)", workspaceRows, c.table, inClause)
	params := make([]any, len(workspaceIds))
	for i, v := range workspaceIds {
		params[i] = v
	}
	err := c.conn.QueryRowsCtx(ctx, &resp, query, params...)
	switch err {
	case nil:
		return resp, nil
	default:
		logc.Infov(ctx, err)
		return nil, err
	}
}

// NewWorkspaceModel returns a model for the database table.
func NewWorkspaceModel(conn sqlx.SqlConn) WorkspaceModel {
	return &customWorkspaceModel{
		defaultWorkspaceModel: newWorkspaceModel(conn),
	}
}
