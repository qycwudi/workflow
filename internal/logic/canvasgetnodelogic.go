package logic

import (
	"context"

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
	// todo: add your logic here and delete this line

	return
}
