package api

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/x/errors"

	"workflow/internal/logic"
	"workflow/internal/model"
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
	var apis []*model.Api
	var total int64
	if req.Id == "" {
		total, apis, err = l.svcCtx.ApiModel.FindAll(l.ctx, req.Current, req.PageSize)
	} else {
		total, apis, err = l.svcCtx.ApiModel.FindByWorkSpaceId(l.ctx, req.Id, req.Current, req.PageSize)
	}
	if err != nil {
		return nil, errors.New(int(logic.SystemOrmError), "查询API失败")
	}
	lists := make([]types.ApiPublishList, len(apis))
	for i, api := range apis {
		lists[i] = types.ApiPublishList{
			WorkSpaceId: api.WorkspaceId,
			ApiId:       api.ApiId,
			ApiName:     api.ApiName,
			ApiDesc:     api.ApiDesc,
			Status:      api.Status,
		}
	}

	resp = &types.ApiPublishListResponse{
		Current:  req.Current,
		PageSize: req.PageSize,
		Total:    total,
		List:     lists,
	}
	return resp, nil
}
