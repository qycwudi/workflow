package logic

import (
	"context"

	"gogogo/internal/svc"
	"gogogo/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type WorkSpaceRemoveLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewWorkSpaceRemoveLogic(ctx context.Context, svcCtx *svc.ServiceContext) *WorkSpaceRemoveLogic {
	return &WorkSpaceRemoveLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *WorkSpaceRemoveLogic) WorkSpaceRemove(req *types.WorkRemoveRequest) (resp *types.WorkSpaceRemoveResponse, err error) {
	// todo: add your logic here and delete this line

	return
}
