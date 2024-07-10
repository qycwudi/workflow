package logic

import (
	"context"

	"gogogo/internal/svc"
	"gogogo/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type CommonGetVByKLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCommonGetVByKLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CommonGetVByKLogic {
	return &CommonGetVByKLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CommonGetVByKLogic) CommonGetVByK(req *types.GetVByKRequest) (resp *types.GetVByKResponse, err error) {
	// todo: add your logic here and delete this line

	return
}
