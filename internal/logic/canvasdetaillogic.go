package logic

import (
	"context"
	"github.com/zeromicro/x/errors"

	"workflow/internal/svc"
	"workflow/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type CanvasDetailLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCanvasDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CanvasDetailLogic {
	return &CanvasDetailLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CanvasDetailLogic) CanvasDetail(req *types.CanvasDetailRequest) (resp *types.CanvasDetailResponse, err error) {
	// 查询点
	nodes, err := l.svcCtx.NodeModel.FindOneByWorkSpace(l.ctx, req.WorkSpaceId)
	if err != nil {
		return nil, errors.New(int(SystemOrmError), "查询节点失败")
	}

	// 查询边
	edges, err := l.svcCtx.EdgeModel.FindOneByWorkSpace(l.ctx, req.WorkSpaceId)
	if err != nil {
		return nil, errors.New(int(SystemOrmError), "查询边失败")
	}

	l.Infov(nodes)
	l.Infov("----------")
	l.Infov(edges)
	resp = &types.CanvasDetailResponse{
		Node: []types.CanvasNode{},
		Edge: []types.CanvasEdge{},
	}
	return resp, nil
}
