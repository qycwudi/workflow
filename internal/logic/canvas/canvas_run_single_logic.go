package canvas

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"workflow/internal/svc"
	"workflow/internal/types"
)

type CanvasRunSingleLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCanvasRunSingleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CanvasRunSingleLogic {
	return &CanvasRunSingleLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CanvasRunSingleLogic) CanvasRunSingle(req *types.CanvasRunSingleRequest) (resp *types.CanvasRunSingleResponse, err error) {
	// todo: add your logic here and delete this line
	
	return
}
