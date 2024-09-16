package logic

import (
	"context"

	"gogogo/internal/svc"
	"gogogo/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type WorkSpaceNewLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewWorkSpaceNewLogic(ctx context.Context, svcCtx *svc.ServiceContext) *WorkSpaceNewLogic {
	return &WorkSpaceNewLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *WorkSpaceNewLogic) WorkSpaceNew(req *types.WorkSpaceRequest) (resp *types.WorkSpaceNewResponse, err error) {
	// todo: add your logic here and delete this line

	return
}
