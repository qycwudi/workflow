package workspace

import (
	"context"

	"workflow/internal/svc"
	"workflow/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type WorkSpaceEditTagLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewWorkSpaceEditTagLogic(ctx context.Context, svcCtx *svc.ServiceContext) *WorkSpaceEditTagLogic {
	return &WorkSpaceEditTagLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *WorkSpaceEditTagLogic) WorkSpaceEditTag(req *types.WorkSpaceEditTagRequest) (resp *types.WorkSpaceEditTagResponse, err error) {
	// todo: add your logic here and delete this line

	return
}
