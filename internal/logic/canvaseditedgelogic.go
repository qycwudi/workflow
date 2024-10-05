package logic

import (
	"context"
	"github.com/rulego/rulego/utils/json"
	"github.com/zeromicro/x/errors"
	"workflow/internal/svc"
	"workflow/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type CanvasEditEdgeLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCanvasEditEdgeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CanvasEditEdgeLogic {
	return &CanvasEditEdgeLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CanvasEditEdgeLogic) CanvasEditEdge(req *types.CanvasEditEdgeRequest) (resp *types.CanvasEditEdgeResponse, err error) {
	// 检查
	if err = checkEdge(l.ctx, l.svcCtx, req.Source, req.Target); err != nil {
		return nil, err
	}

	// 查询边
	edge, err := l.svcCtx.EdgeModel.FindOneByEdgeId(l.ctx, req.EdgeId)
	if err != nil {
		return nil, errors.New(int(SystemOrmError), "查找边失败")
	}

	edgeCustomData := types.EdgeCustomData{
		SourcePoint: req.SourcePoint,
		TargetPoint: req.TargetPoint,
	}
	edgeCustomDataMar, _ := json.Marshal(edgeCustomData)
	edge.CustomData = string(edgeCustomDataMar)
	edge.Source = req.Source
	edge.Target = req.Target
	edge.Route = req.Route

	err = l.svcCtx.EdgeModel.UpdateByEdgeId(l.ctx, edge)
	if err != nil {
		return nil, errors.New(int(SystemOrmError), "修改边错误")
	}
	resp = &types.CanvasEditEdgeResponse{EdgeId: req.EdgeId}
	return resp, nil
}
