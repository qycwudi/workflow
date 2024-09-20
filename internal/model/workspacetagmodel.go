package model

import (
	"context"
	"fmt"
	"github.com/zeromicro/go-zero/core/stores/sqlc"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ WorkspaceTagModel = (*customWorkspaceTagModel)(nil)

type (
	// WorkspaceTagModel is an interface to be customized, add more methods here,
	// and implement the added methods in customWorkspaceTagModel.
	WorkspaceTagModel interface {
		workspaceTagModel
		FindOneByName(ctx context.Context, tagName string) (*WorkspaceTag, error)
	}

	customWorkspaceTagModel struct {
		*defaultWorkspaceTagModel
	}
)

func (c customWorkspaceTagModel) FindOneByName(ctx context.Context, tagName string) (*WorkspaceTag, error) {
	query := fmt.Sprintf("select %s from %s where `tag_name` = ? limit 1", workspaceTagRows, c.table)
	var resp WorkspaceTag
	err := c.conn.QueryRowCtx(ctx, &resp, query, tagName)
	switch err {
	case nil:
		return &resp, nil
	case sqlc.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

// NewWorkspaceTagModel returns a model for the database table.
func NewWorkspaceTagModel(conn sqlx.SqlConn) WorkspaceTagModel {
	return &customWorkspaceTagModel{
		defaultWorkspaceTagModel: newWorkspaceTagModel(conn),
	}
}