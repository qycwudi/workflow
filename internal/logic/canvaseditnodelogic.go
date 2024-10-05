package logic

import (
	"context"
	"github.com/rulego/rulego/utils/json"
	"github.com/zeromicro/x/errors"
	"time"
	"workflow/internal/svc"
	"workflow/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type CanvasEditNodeLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCanvasEditNodeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CanvasEditNodeLogic {
	return &CanvasEditNodeLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CanvasEditNodeLogic) CanvasEditNode(req *types.CanvasEditNodeRequest) (resp *types.CanvasEditNodeResponse, err error) {

	node, err := l.svcCtx.NodeModel.FindOneByNodeId(l.ctx, req.NodeId)
	if err != nil {
		return nil, errors.New(int(SystemOrmError), "节点不存在")
	}
	// 修改
	position := types.NodePosition{
		X: req.Position.X,
		Y: req.Position.Y,
	}
	positionMar, _ := json.Marshal(position)
	node.Position = string(positionMar)

	if len(req.ModuleConfig) > 0 {
		moduleConfigMarshal, _ := json.Marshal(req.ModuleConfig)
		node.Configuration = string(moduleConfigMarshal)
	}
	node.UpdateTime = time.Now()

	err = l.svcCtx.NodeModel.UpdateByNodeId(l.ctx, node)
	if err != nil {
		return nil, errors.New(int(SystemOrmError), "修改节点错误")
	}

	resp = &types.CanvasEditNodeResponse{NodeId: req.NodeId}
	return resp, nil
}
