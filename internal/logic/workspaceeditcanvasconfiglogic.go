package logic

import (
	"context"

	"gogogo/internal/svc"
	"gogogo/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type WorkSpaceEditCanvasConfigLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewWorkSpaceEditCanvasConfigLogic(ctx context.Context, svcCtx *svc.ServiceContext) *WorkSpaceEditCanvasConfigLogic {
	return &WorkSpaceEditCanvasConfigLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *WorkSpaceEditCanvasConfigLogic) WorkSpaceEditCanvasConfig(req *types.WorkSpaceUploadCanvasConfigTagRequest) (resp *types.WorkSpaceUploadCanvasConfigTagResponse, err error) {
	// todo: add your logic here and delete this line

	return
}
