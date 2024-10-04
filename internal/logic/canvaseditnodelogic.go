package logic

import (
	"context"

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
	// todo: add your logic here and delete this line

	return
}
