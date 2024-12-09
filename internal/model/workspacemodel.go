package model

import (
	"context"
	"fmt"
	"strings"

	"github.com/zeromicro/go-zero/core/logc"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ WorkspaceModel = (*customWorkspaceModel)(nil)

type (
	// WorkspaceModel is an interface to be customized, add more methods here,
	// and implement the added methods in customWorkspaceModel.
	WorkspaceModel interface {
		workspaceModel
		FindPage(ctx context.Context, current, pageSize int, workSpaceType, workSpaceName string) ([]*Workspace, int64, error)
		FindInWorkSpaceId(ctx context.Context, workspaceId []string) ([]*Workspace, error)
		Remove(ctx context.Context, workSpaceId string) error
		UpdateByWorkspaceId(ctx context.Context, data *Workspace) error
	}

	customWorkspaceModel struct {
		*defaultWorkspaceModel
	}
)

func (c customWorkspaceModel) UpdateByWorkspaceId(ctx context.Context, data *Workspace) error {
	var setters []string
	var args []interface{}

	// 检查字段是否为空，如果不为空则添加到更新列表
	if data.WorkspaceName != "" {
		setters = append(setters, "`workspace_name` = ?")
		args = append(args, data.WorkspaceName)
	}
	if data.WorkspaceDesc.Valid {
		setters = append(setters, "`workspace_desc` = ?")
		args = append(args, data.WorkspaceDesc.String)
	}
	if data.WorkspaceType.Valid {
		setters = append(setters, "`workspace_type` = ?")
		args = append(args, data.WorkspaceType.String)
	}
	if data.WorkspaceIcon.Valid {
		setters = append(setters, "`workspace_icon` = ?")
		args = append(args, data.WorkspaceIcon.String)
	}
	if data.CanvasConfig.Valid {
		setters = append(setters, "`canvas_config` = ?")
		args = append(args, data.CanvasConfig.String)
	}
	if !data.UpdateTime.IsZero() {
		setters = append(setters, "`update_time` = ?")
		args = append(args, data.UpdateTime)
	}

	// 构建最终的SQL语句
	query := fmt.Sprintf("UPDATE %s SET %s WHERE `workspace_id` = ?", c.table, strings.Join(setters, ", "))
	args = append(args, data.WorkspaceId)
	_, err := c.conn.ExecCtx(ctx, query, args...)
	if err != nil {
		logc.Infov(ctx, err)
		return err
	}
	return nil
}

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
	err := c.conn.QueryRowCtx(ctx, &total, totalBuilder.String(), params...)
	if err != nil {
		logc.Infov(ctx, err)
		return nil, 0, err
	}

	// 添加排序
	queryBuilder.WriteString(" order by id desc")
	// 添加分页条件
	queryBuilder.WriteString(" LIMIT ?, ?")

	// 计算偏移量
	params = append(params, (current-1)*pageSize)
	params = append(params, pageSize)

	var resp []*Workspace
	err = c.conn.QueryRowsCtx(ctx, &resp, queryBuilder.String(), params...)
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

func (c customWorkspaceModel) Remove(ctx context.Context, workSpaceId string) error {
	// 把mapping关系也要删掉
	return c.conn.TransactCtx(ctx, func(ctx context.Context, session sqlx.Session) error {
		var err error
		update := fmt.Sprintf("update %s set is_delete = 1 where `workspace_id` = ?", c.table)
		_, err = session.ExecCtx(ctx, update, workSpaceId)
		if err != nil {
			logc.Infov(ctx, err)
			return err
		}
		query := "delete from workspace_tag_mapping where `workspace_id` = ?"
		_, err = session.ExecCtx(ctx, query, workSpaceId)
		if err != nil {
			logc.Infov(ctx, err)
			return err
		}
		return nil
	})
}

// NewWorkspaceModel returns a model for the database table.
func NewWorkspaceModel(conn sqlx.SqlConn) WorkspaceModel {
	return &customWorkspaceModel{
		defaultWorkspaceModel: newWorkspaceModel(conn),
	}
}
