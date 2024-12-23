package workspace

import (
	"context"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/x/errors"

	"workflow/internal/logic"
	"workflow/internal/model"
	"workflow/internal/svc"
	"workflow/internal/types"
)

type TagRemoveLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewTagRemoveLogic(ctx context.Context, svcCtx *svc.ServiceContext) *TagRemoveLogic {
	return &TagRemoveLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *TagRemoveLogic) TagRemove(req *types.TagRemoveRequest) (resp *types.TagRemoveResponse, err error) {
	// 逻辑删除
	err = l.svcCtx.WorkSpaceTagModel.Update(l.ctx, &model.WorkspaceTag{
		Id:         req.Id,
		IsDelete:   1,
		UpdateTime: time.Now(),
	})
	if err != nil {
		return nil, errors.New(int(logic.SystemOrmError), "删除标签失败")
	}
	resp = &types.TagRemoveResponse{
		Id: req.Id,
	}

	return resp, nil
}
