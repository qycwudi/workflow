package canvas

import (
	"context"
	"github.com/zeromicro/x/errors"
	"workflow/internal/logic"

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

func (l *CanvasGetEdgeLogic) CanvasGetEdge(req *types.CanvasGetEdgeRequest) (resp *types.CanvasGetEdgeResponse, err error) {
	edge, err := l.svcCtx.EdgeModel.FindOneByEdgeId(l.ctx, req.EdgeId)
	if err != nil {
		return nil, errors.New(int(logic.SystemOrmError), "边不存在")
	}
	resp = &types.CanvasGetEdgeResponse{EdgeId: edge.EdgeId}
	// todo 补充
	return resp, nil
}
