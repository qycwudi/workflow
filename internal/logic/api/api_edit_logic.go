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

type ApiEditLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewApiEditLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ApiEditLogic {
	return &ApiEditLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ApiEditLogic) ApiEdit(req *types.ApiEditRequest) (resp *types.ApiEditResponse, err error) {
	// 校验参数
	if req.ApiId == "" {
		return nil, errors.New(int(logic.SystemError), "apiId is required")
	}
	// 校验 apiName
	if req.ApiName == "" {
		return nil, errors.New(int(logic.SystemError), "apiName is required")
	}
	// 校验 apiDesc
	if req.ApiDesc == "" {
		return nil, errors.New(int(logic.SystemError), "apiDesc is required")
	}
	// 查询api
	api, err := l.svcCtx.ApiModel.FindOneByApiId(l.ctx, req.ApiId)
	if err != nil {
		return nil, errors.New(int(logic.SystemError), "apiId not found")
	}
	// 更新api
	api.ApiName = req.ApiName
	api.ApiDesc = req.ApiDesc
	tagJson, err := sonic.Marshal(req.Tag)
	if err != nil {
		return nil, errors.New(int(logic.SystemError), "tag conversion failed")
	}
	api.Tag = string(tagJson)

	err = l.svcCtx.ApiModel.Update(l.ctx, api)
	if err != nil {
		return nil, errors.New(int(logic.SystemError), "update api failed")
	}

	return &types.ApiEditResponse{
		ApiId: api.ApiId,
	}, nil
}
