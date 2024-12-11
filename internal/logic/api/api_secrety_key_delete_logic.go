package api

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/x/errors"

	"workflow/internal/logic"
	"workflow/internal/svc"
	"workflow/internal/types"
)

type ApiSecretyKeyDeleteLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewApiSecretyKeyDeleteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ApiSecretyKeyDeleteLogic {
	return &ApiSecretyKeyDeleteLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ApiSecretyKeyDeleteLogic) ApiSecretyKeyDelete(req *types.ApiSecretyKeyDeleteRequest) (resp *types.ApiSecretyKeyDeleteResponse, err error) {
	// 检查API是否存在
	_, err = l.svcCtx.ApiModel.FindOneByApiId(l.ctx, req.ApiId)
	if err != nil {
		return nil, errors.New(int(logic.SystemOrmError), "API 不存在")
	}

	// 逻辑删除
	err = l.svcCtx.ApiSecretKeyModel.LogicalDelete(l.ctx, req.ApiId, req.SecretyKey)
	if err != nil {
		return nil, errors.New(int(logic.SystemOrmError), "删除 API Secret Key 记录失败")
	}
	resp = &types.ApiSecretyKeyDeleteResponse{
		ApiId:      req.ApiId,
		SecretyKey: req.SecretyKey,
	}
	return resp, nil
}
