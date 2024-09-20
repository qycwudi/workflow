package model

import (
	"context"
	"fmt"
	"github.com/zeromicro/go-zero/core/logc"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"strings"
)

var _ WorkspaceTagMappingModel = (*customWorkspaceTagMappingModel)(nil)

type (
	// WorkspaceTagMappingModel is an interface to be customized, add more methods here,
	// and implement the added methods in customWorkspaceTagMappingModel.
	WorkspaceTagMappingModel interface {
		workspaceTagMappingModel
		FindByWorkSpaceId(ctx context.Context, workSpaceId []string) ([]*WorkspaceTagNameMapping, error)
		FindPageByTagId(ctx context.Context, current int, pageSize int, tagId []int64) ([]string, int64, error)
	}

	customWorkspaceTagMappingModel struct {
		*defaultWorkspaceTagMappingModel
	}

	WorkspaceTagNameMapping struct {
		Id          int64  `db:"id"`           // 主建
		TagId       int64  `db:"tag_id"`       // 标签ID
		TagName     string `db:"tag_name"`     // 标签Name
		WorkspaceId string `db:"workspace_id"` // 画布空间ID
	}
)

// NewWorkspaceTagMappingModel returns a model for the database table.
func NewWorkspaceTagMappingModel(conn sqlx.SqlConn) WorkspaceTagMappingModel {
	return &customWorkspaceTagMappingModel{
		defaultWorkspaceTagMappingModel: newWorkspaceTagMappingModel(conn),
	}
}

func (m *defaultWorkspaceTagMappingModel) FindByWorkSpaceId(ctx context.Context, workSpaceIds []string) ([]*WorkspaceTagNameMapping, error) {
	// 构建参数占位符，如 (?, ?, ?)
	placeholders := make([]string, len(workSpaceIds))
	for i := range workSpaceIds {
		placeholders[i] = "?"
	}
	inClause := strings.Join(placeholders, ", ")
	query := fmt.Sprintf("select a.id as id, a.tag_id as tag_id,a.workspace_id as workspace_id,b.tag_name as tag_name from `workspace_tag_mapping` as a join `workspace_tag` as b on  a.tag_id = b.id where a.`workspace_id` in (%s)", inClause)
	var resp []*WorkspaceTagNameMapping
	wsIds := make([]any, len(workSpaceIds))
	for i, v := range workSpaceIds {
		wsIds[i] = v
	}
	err := m.conn.QueryRowsCtx(ctx, &resp, query, wsIds...)
	switch err {
	case nil:
		return resp, nil
	default:
		logc.Infov(ctx, err)
		return nil, err
	}
}

func (c customWorkspaceTagMappingModel) FindPageByTagId(ctx context.Context, current int, pageSize int, tagId []int64) ([]string, int64, error) {
	placeholders := make([]string, len(tagId))
	for i := range tagId {
		placeholders[i] = "?"
	}
	inClause := strings.Join(placeholders, ", ")
	params := make([]any, len(tagId))
	for i, v := range tagId {
		params[i] = v
	}
	totalQuery := fmt.Sprintf("select count(*) from `workspace` as a join `workspace_tag_mapping` as b on a.workspace_id = b.workspace_id where b.tag_id in (%s)", inClause)
	var total int64
	_ = c.conn.QueryRowsCtx(ctx, &total, totalQuery, params...)

	query := fmt.Sprintf("select a.workspace_id from `workspace` as a join `workspace_tag_mapping` as b on a.workspace_id = b.workspace_id where b.tag_id in (%s) order by b.id desc LIMIT ?, ?", inClause)
	var resp []string
	params = append(params, (current-1)*pageSize)
	params = append(params, pageSize)
	err := c.conn.QueryRowsCtx(ctx, &resp, query, params...)
	switch err {
	case nil:
		return resp, total, nil
	default:
		logc.Infov(ctx, err)
		return nil, 0, err
	}
}
