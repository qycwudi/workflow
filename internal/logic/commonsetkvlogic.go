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
	// todo: add your logic here and delete this line
	l.Info("request:%+v", req)
	response := types.SetKvResponse{
		Code:    0,
		Message: "SUCCESS",
	}
	return &response, nil
}
