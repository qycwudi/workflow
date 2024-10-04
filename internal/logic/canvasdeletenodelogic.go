package logic

import (
	"context"
	"github.com/zeromicro/x/errors"
	"workflow/internal/rolego"

	"workflow/internal/svc"
	"workflow/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type CanvasDeleteNodeLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCanvasDeleteNodeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CanvasDeleteNodeLogic {
	return &CanvasDeleteNodeLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CanvasDeleteNodeLogic) CanvasDeleteNode(req *types.CanvasDeleteNodeRequest) (resp *types.CanvasDeleteNodeResponse, err error) {
	// 起始节点不能删除
	edge, err := l.svcCtx.NodeModel.FindOneByNodeId(l.ctx, req.NodeId)
	if err != nil {
		return nil, errors.New(int(SystemOrmError), "查询节点失败")
	}
	if edge.NodeType == rolego.Start {
		return nil, errors.New(int(SystemError), "开始节点不允许删除")
	}

	err = l.svcCtx.NodeModel.DeleteNodeByNodeIdAndWorkSpace(l.ctx, req.NodeId, req.WorkSpaceId)
	if err != nil {
		return nil, errors.New(int(SystemOrmError), "删除节点失败")
	}
	return &types.CanvasDeleteNodeResponse{NodeId: req.NodeId}, nil
}