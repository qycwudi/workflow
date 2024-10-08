package workspace

import (
	"context"
	"github.com/zeromicro/x/errors"
	"workflow/internal/logic"

	"workflow/internal/svc"
	"workflow/internal/types"

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
	err = l.svcCtx.WorkSpaceModel.Remove(l.ctx, req.WorkSpaceId)
	if err != nil {
		return nil, errors.New(int(logic.SystemOrmError), "空间删除失败")
	}
	return
}
