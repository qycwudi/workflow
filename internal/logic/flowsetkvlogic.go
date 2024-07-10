package logic

import (
	"context"

	"gogogo/internal/svc"
	"gogogo/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type FlowSetKvLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewFlowSetKvLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FlowSetKvLogic {
	return &FlowSetKvLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *FlowSetKvLogic) FlowSetKv(req *types.FlowRequest) (resp *types.FlowResponse, err error) {
	// todo: add your logic here and delete this line

	return
}
