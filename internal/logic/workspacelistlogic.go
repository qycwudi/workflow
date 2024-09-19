package logic

import (
	"context"

	"gogogo/internal/svc"
	"gogogo/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type WorkSpaceListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewWorkSpaceListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *WorkSpaceListLogic {
	return &WorkSpaceListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *WorkSpaceListLogic) WorkSpaceList(req *types.WorkSpaceListRequest) (resp *types.WorkSpaceListResponse, err error) {
	// todo: add your logic here and delete this line

	return
}
