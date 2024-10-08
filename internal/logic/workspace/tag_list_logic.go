package workspace

import (
	"context"
	"github.com/zeromicro/x/errors"
	"workflow/internal/logic"

	"workflow/internal/svc"
	"workflow/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type TagListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewTagListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *TagListLogic {
	return &TagListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *TagListLogic) TagList(req *types.TagListRequest) (resp *types.TagListResponse, err error) {
	tagList, err := l.svcCtx.WorkSpaceTagModel.FindAll(l.ctx)
	if err != nil {
		return nil, errors.New(int(logic.SystemOrmError), "查询标签列表失败")
	}
	tag := make([]types.TagEntity, len(tagList))
	for i, v := range tagList {
		tag[i] = types.TagEntity{
			Id:   int(v.Id),
			Name: v.TagName,
		}
	}
	resp = &types.TagListResponse{Tag: tag}
	return resp, nil
}
