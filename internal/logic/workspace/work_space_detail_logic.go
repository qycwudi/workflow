package workspace

import (
	"context"

	"workflow/internal/svc"
	"workflow/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type WorkSpaceDetailLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewWorkSpaceDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *WorkSpaceDetailLogic {
	return &WorkSpaceDetailLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *WorkSpaceDetailLogic) WorkSpaceDetail(req *types.WorkSpaceDetailRequest) (resp *types.WorkSpaceDetailResponse, err error) {
	// todo: add your logic here and delete this line

	return
}
