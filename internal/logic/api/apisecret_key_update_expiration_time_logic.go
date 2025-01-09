package api

import (
	"context"
	"fmt"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/x/errors"

	"workflow/internal/cache"
	"workflow/internal/logic"
	"workflow/internal/svc"
	"workflow/internal/types"
)

type ApisecretKeyUpdateExpirationTimeLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewApisecretKeyUpdateExpirationTimeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ApisecretKeyUpdateExpirationTimeLogic {
	return &ApisecretKeyUpdateExpirationTimeLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ApisecretKeyUpdateExpirationTimeLogic) ApisecretKeyUpdateExpirationTime(req *types.ApiSecretKeyUpdateExpirationTimeRequest) (resp *types.ApiSecretKeyUpdateExpirationTimeResponse, err error) {

	// 检查过期时间
	if req.ExpirationTime <= time.Now().UnixMilli() {
		return nil, errors.New(int(logic.SystemOrmError), "过期时间不能小于当前时间")
	}
	// 查询api信息
	api, err := l.svcCtx.ApiSecretKeyModel.FindOneBySecretKey(l.ctx, req.SecretKey)
	if err != nil {
		return nil, errors.New(int(logic.SystemOrmError), "查询 API 信息失败")
	}

	// 修改过期时间
	expirationTime := time.UnixMilli(req.ExpirationTime)
	err = l.svcCtx.ApiSecretKeyModel.UpdateExpirationTime(l.ctx, req.SecretKey, expirationTime)
	if err != nil {
		return nil, errors.New(int(logic.SystemOrmError), "修改 API Secret Key 过期时间失败")
	}
	// 删除redis中的secretKey信息
	err = cache.Redis.Del(l.ctx, fmt.Sprintf(cache.ApiSecretKeyRedisKey, api.ApiId, api.SecretKey))
	if err != nil {
		logx.Errorf("delete redis api info error: %s", err)
		return nil, errors.New(int(logic.SystemError), "删除缓存中的API Secret Key 信息失败")
	}
	resp = &types.ApiSecretKeyUpdateExpirationTimeResponse{
		ExpirationTime: expirationTime.Format(time.DateTime),
	}
	return resp, nil
}
