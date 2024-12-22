package api

import (
	"context"
	"fmt"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/x/errors"

	"workflow/internal/cache"
	"workflow/internal/logic"
	"workflow/internal/svc"
	"workflow/internal/types"
)

type ApisecretKeyDeleteLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewApisecretKeyDeleteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ApisecretKeyDeleteLogic {
	return &ApisecretKeyDeleteLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ApisecretKeyDeleteLogic) ApisecretKeyDelete(req *types.ApiSecretKeyDeleteRequest) (resp *types.ApiSecretKeyDeleteResponse, err error) {

	// 查询api信息
	api, err := l.svcCtx.ApiSecretKeyModel.FindOneBySecretKey(l.ctx, req.SecretKey)
	if err != nil {
		return nil, errors.New(int(logic.SystemOrmError), "查询 API 信息失败")
	}

	// 逻辑删除
	err = l.svcCtx.ApiSecretKeyModel.LogicalDelete(l.ctx, req.SecretKey)
	if err != nil {
		return nil, errors.New(int(logic.SystemOrmError), "删除 API Secret Key 记录失败")
	}

	// 删除redis中的secretKey信息
	err = cache.Redis.Del(l.ctx, fmt.Sprintf(cache.ApiSecretKeyRedisKey, api.ApiId, api.SecretKey))
	if err != nil {
		logx.Errorf("delete redis api info error: %s", err)
		return nil, errors.New(int(logic.SystemError), "删除缓存中的API Secret Key 信息失败")
	}

	resp = &types.ApiSecretKeyDeleteResponse{
		SecretKey: req.SecretKey,
	}
	return resp, nil
}
