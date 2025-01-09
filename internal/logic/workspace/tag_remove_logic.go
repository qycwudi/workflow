package workspace

import (
	"context"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/x/errors"

	"workflow/internal/logic"
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
	tag, err := l.svcCtx.WorkSpaceTagModel.FindOne(l.ctx, req.Id)
	if err != nil {
		return nil, errors.New(int(logic.SystemOrmError), "标签不存在")
	}
	if tag.IsDelete == 1 {
		return nil, errors.New(int(logic.SystemOrmError), "标签已删除")
	}
	// 逻辑删除
	tag.IsDelete = 1
	tag.UpdateTime = time.Now()
	err = l.svcCtx.WorkSpaceTagModel.Update(l.ctx, tag)
	if err != nil {
		return nil, errors.New(int(logic.SystemOrmError), "删除标签失败")
	}
	resp = &types.TagRemoveResponse{
		Id: req.Id,
	}

	return resp, nil
}
