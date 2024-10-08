package canvas

import (
	"context"
	"github.com/zeromicro/x/errors"
	"workflow/internal/logic"

	"workflow/internal/svc"
	"workflow/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type CanvasGetNodeLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCanvasGetNodeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CanvasGetNodeLogic {
	return &CanvasGetNodeLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CanvasGetNodeLogic) CanvasGetNode(req *types.CanvasGetNodeRequest) (resp *types.CanvasGetNodeResponse, err error) {
	node, err := l.svcCtx.NodeModel.FindOneByNodeId(l.ctx, req.NodeId)
	if err != nil {
		return nil, errors.New(int(logic.SystemOrmError), "节点不存在")
	}
	resp = &types.CanvasGetNodeResponse{NodeId: node.NodeId}
	// todo 补充
	return resp, nil
}
