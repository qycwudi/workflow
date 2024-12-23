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

type TagEditLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewTagEditLogic(ctx context.Context, svcCtx *svc.ServiceContext) *TagEditLogic {
	return &TagEditLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *TagEditLogic) TagEdit(req *types.TagEditRequest) (resp *types.TagEditResponse, err error) {
	// 编辑名称
	err = l.svcCtx.WorkSpaceTagModel.Update(l.ctx, &model.WorkspaceTag{
		Id:         req.Id,
		TagName:    req.Name,
		UpdateTime: time.Now(),
	})
	if err != nil {
		return nil, errors.New(int(logic.SystemOrmError), "修改标签失败")
	}
	resp = &types.TagEditResponse{
		Id:   req.Id,
		Name: req.Name,
	}
	return resp, nil
}
