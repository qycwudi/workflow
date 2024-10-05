package model

import (
	"context"
	"fmt"
	"github.com/zeromicro/go-zero/core/logc"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ NodeModel = (*customNodeModel)(nil)

type (
	// NodeModel is an interface to be customized, add more methods here,
	// and implement the added methods in customNodeModel.
	NodeModel interface {
		nodeModel
		DeleteNodeByNodeIdAndWorkSpace(ctx context.Context, nodeId string, workSpaceId string) error
		FindOneByWorkSpace(ctx context.Context, workspaceId string) ([]*Node, error)
		UpdateByNodeId(ctx context.Context, newData *Node) error
	}

	customNodeModel struct {
		*defaultNodeModel
	}
)

func (c customNodeModel) UpdateByNodeId(ctx context.Context, newData *Node) error {
	query := fmt.Sprintf("update %s set %s where `node_id` = ?", c.table, nodeRowsWithPlaceHolder)
	_, err := c.conn.ExecCtx(ctx, query, newData.NodeId, newData.NodeType, newData.LabelConfig, newData.CustomConfig, newData.TaskConfig, newData.StyleConfig, newData.Position, newData.CreateTime, newData.UpdateTime, newData.NodeName, newData.Configuration, newData.WorkspaceId, newData.ModuleId, newData.NodeId)
	return err
}

func (c customNodeModel) FindOneByWorkSpace(ctx context.Context, workspaceId string) ([]*Node, error) {
	var resp []*Node
	query := fmt.Sprintf("select %s from %s where `workspace_id` = ?", nodeRows, c.table)
	err := c.conn.QueryRowsCtx(ctx, &resp, query, workspaceId)
	switch err {
	case nil:
		return resp, nil
	default:
		logc.Infov(ctx, err)
		return nil, err
	}
}

func (c customNodeModel) DeleteNodeByNodeIdAndWorkSpace(ctx context.Context, nodeId string, workSpaceId string) error {
	query := fmt.Sprintf("delete from %s where `node_id` = ? and workspace_id = ?", c.table)
	_, err := c.conn.ExecCtx(ctx, query, nodeId, workSpaceId)
	return err
}

// NewNodeModel returns a model for the database table.
func NewNodeModel(conn sqlx.SqlConn) NodeModel {
	return &customNodeModel{
		defaultNodeModel: newNodeModel(conn),
	}
}
