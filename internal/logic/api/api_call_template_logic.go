package api

import (
	"context"

	"github.com/bytedance/sonic"
	"github.com/zeromicro/go-zero/core/logx"

	"workflow/internal/model"
	"workflow/internal/svc"
	"workflow/internal/types"
)

type ApiCallTemplateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewApiCallTemplateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ApiCallTemplateLogic {
	return &ApiCallTemplateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ApiCallTemplateLogic) ApiCallTemplate(req *types.ApiCallTemplateRequest) (resp *types.ApiCallTemplateResponse, err error) {
	// 查询 api 参数
	api, err := l.svcCtx.ApiModel.FindOneByApiId(l.ctx, req.ApiId)
	if err != nil {
		return nil, err
	}
	url := l.svcCtx.Config.ApiUrl + "/" + req.ApiId
	// 查询 case 参数
	cases, err := l.svcCtx.CaseModel.FindOneByUid(l.ctx, req.CaseId)
	if err != nil {
		return nil, err
	}
	param := cases.Params
	// 查询 secret
	header := map[string]string{}
	secret, err := l.svcCtx.ApiSecretKeyModel.FindByApiId(l.ctx, api.ApiId)
	if err != nil && err != model.ErrNotFound {
		return nil, err
	}
	if secret == nil {
		header["Authorization"] = "Bearer " + "请创建密钥"
	} else {
		header["Authorization"] = "Bearer " + secret[0].SecretKey
	}
	headerJson, err := sonic.Marshal(header)
	if err != nil {
		return nil, err
	}

	resp = &types.ApiCallTemplateResponse{
		ApiId:  api.ApiId,
		Url:    url,
		Header: string(headerJson),
		Body:   param,
	}
	return resp, nil
}
