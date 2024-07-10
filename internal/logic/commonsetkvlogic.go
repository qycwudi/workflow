package logic

import (
	"context"

	"gogogo/internal/svc"
	"gogogo/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type CommonSetKvLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCommonSetKvLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CommonSetKvLogic {
	return &CommonSetKvLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CommonSetKvLogic) CommonSetKv(req *types.SetKvRequest) (resp *types.SetKvResponse, err error) {
	l.Info("request:%+v", req)
	one, err := l.svcCtx.GogogoKvModel.FindByKey(l.ctx, "xuetu")
	response := types.SetKvResponse{
		Code:    0,
		Message: "SUCCESS:" + one.V,
	}
	return &response, nil
}
