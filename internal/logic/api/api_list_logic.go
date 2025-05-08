package api

import (
	"context"

	"github.com/bytedance/sonic"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/x/errors"

	"workflow/internal/logic"
	"workflow/internal/svc"
	"workflow/internal/types"
)

type ApiListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewApiListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ApiListLogic {
	return &ApiListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ApiListLogic) ApiList(req *types.ApiPublishListRequest) (resp *types.ApiPublishListResponse, err error) {
	pagin, err := l.svcCtx.ApiModel.Page(l.ctx, req.Current, req.PageSize, req.Id, req.Name, req.Tag)
	if err != nil {
		return nil, errors.New(int(logic.SystemOrmError), "查询API失败")
	}
	lists := make([]types.ApiPublishList, len(pagin.List))
	for i, api := range pagin.List {
		var tag []string
		err := sonic.Unmarshal([]byte(api.Tag), &tag)
		if err != nil {
			return nil, errors.New(int(logic.SystemError), "标签转换失败")
		}
		lists[i] = types.ApiPublishList{
			WorkSpaceId: api.WorkspaceId,
			ApiId:       api.ApiId,
			ApiName:     api.ApiName,
			Tag:         tag,
			ApiDesc:     api.ApiDesc,
			Status:      api.Status,
			PublishTime: api.CreateTime.Format("2006-01-02 15:04:05"),
		}
	}

	resp = &types.ApiPublishListResponse{
		Current:  req.Current,
		PageSize: req.PageSize,
		Total:    pagin.Total,
		List:     lists,
	}
	return resp, nil
}
