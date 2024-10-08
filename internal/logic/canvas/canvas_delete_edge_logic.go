package canvas

import (
	"context"
	"github.com/zeromicro/x/errors"
	"workflow/internal/logic"

	"workflow/internal/svc"
	"workflow/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type CanvasDeleteEdgeLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCanvasDeleteEdgeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CanvasDeleteEdgeLogic {
	return &CanvasDeleteEdgeLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CanvasDeleteEdgeLogic) CanvasDeleteEdge(req *types.CanvasDeleteEdgeRequest) (resp *types.CanvasDeleteEdgeResponse, err error) {
	err = l.svcCtx.EdgeModel.DeleteByEdgeIdAndWorkSpaceId(l.ctx, req.EdgeId, req.Id)
	if err != nil {
		return nil, errors.New(int(logic.SystemOrmError), "删除边失败")
	}
	return &types.CanvasDeleteEdgeResponse{EdgeId: req.EdgeId}, nil
}
