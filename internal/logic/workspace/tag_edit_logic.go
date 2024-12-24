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
	if req.Name == "" {
		return nil, errors.New(int(logic.SystemOrmError), "标签名称不能为空")
	}
	// 编辑名称
	tag, err := l.svcCtx.WorkSpaceTagModel.FindOne(l.ctx, req.Id)
	if err != nil {
		return nil, errors.New(int(logic.SystemOrmError), "标签不存在")
	}
	if tag.TagName == req.Name {
		return nil, errors.New(int(logic.SystemOrmError), "标签名称相同")
	}
	tag.TagName = req.Name
	tag.UpdateTime = time.Now()
	err = l.svcCtx.WorkSpaceTagModel.Update(l.ctx, tag)
	if err != nil {
		return nil, errors.New(int(logic.SystemOrmError), "修改标签失败")
	}
	resp = &types.TagEditResponse{
		Id:   req.Id,
		Name: req.Name,
	}
	return resp, nil
}
