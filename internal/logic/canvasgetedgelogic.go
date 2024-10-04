package logic

import (
	"context"

	"workflow/internal/svc"
	"workflow/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type CanvasGetEdgeLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCanvasGetEdgeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CanvasGetEdgeLogic {
	return &CanvasGetEdgeLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CanvasGetEdgeLogic) CanvasGetEdge(req *types.CanvasGetEdgeRequest) (resp *types.ModuleListResponse, err error) {
	// todo: add your logic here and delete this line

	return
}
