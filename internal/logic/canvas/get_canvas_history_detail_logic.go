package canvas

import (
	"context"

	"workflow/internal/svc"
	"workflow/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetCanvasHistoryDetailLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetCanvasHistoryDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetCanvasHistoryDetailLogic {
	return &GetCanvasHistoryDetailLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetCanvasHistoryDetailLogic) GetCanvasHistoryDetail(req *types.GetCanvasHistoryDetailReq) (resp *types.GetCanvasHistoryDetailResp, err error) {
	// todo: add your logic here and delete this line

	return
}
