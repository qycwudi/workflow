package logic

import (
	"context"

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
	// todo: add your logic here and delete this line

	return
}
